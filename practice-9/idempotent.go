package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"time"
	"github.com/go-redis/redis/v8"
)

type CachedResponse struct {
	StatusCode int
	Body       []byte
	Completed  bool
}

type IdempotencyStore interface {
	StartProcessing(ctx context.Context, key string) (bool, error)
	Get(ctx context.Context, key string) (*CachedResponse, error)
	Finish(ctx context.Context, key string, statusCode int, body []byte) error
}

type RedisStore struct {
	client *redis.Client
	ttl    time.Duration
}

func NewRedisStore(addr string, ttl time.Duration) *RedisStore {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	return &RedisStore{
		client: client,
		ttl:    ttl,
	}
}

func (s *RedisStore) StartProcessing(ctx context.Context, key string) (bool, error) {
	val, err := s.client.SetNX(ctx, key, "processing", 5*time.Minute).Result()
	return val, err
}

func (s *RedisStore) Get(ctx context.Context, key string) (*CachedResponse, error) {
	val, err := s.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	
	if val == "processing" {
		return &CachedResponse{Completed: false}, nil
	}
	
	var resp CachedResponse
	if err := json.Unmarshal([]byte(val), &resp); err != nil {
		return nil, err
	}
	resp.Completed = true
	return &resp, nil
}

func (s *RedisStore) Finish(ctx context.Context, key string, statusCode int, body []byte) error {
	resp := CachedResponse{
		StatusCode: statusCode,
		Body:       body,
		Completed:  true,
	}
	data, err := json.Marshal(resp)
	if err != nil {
		return err
	}
	return s.client.Set(ctx, key, data, s.ttl).Err()
}

type IdempotencyMiddleware struct {
	store IdempotencyStore
}

func NewIdempotencyMiddleware(store IdempotencyStore) *IdempotencyMiddleware {
	return &IdempotencyMiddleware{store: store}
}

func (m *IdempotencyMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.Header.Get("Idempotency-Key")
		if key == "" {
			http.Error(w, "Idempotency-Key header is required", http.StatusBadRequest)
			return
		}
		
		ctx := r.Context()
		
		cached, err := m.store.Get(ctx, key)
		if err != nil {
			http.Error(w, "Storage error", http.StatusInternalServerError)
			return
		}
		
		if cached != nil && cached.Completed {
			w.WriteHeader(cached.StatusCode)
			w.Write(cached.Body)
			return
		}
		
		if cached != nil && !cached.Completed {
			http.Error(w, "Duplicate request in progress", http.StatusConflict)
			return
		}
		
		started, err := m.store.StartProcessing(ctx, key)
		if err != nil {
			http.Error(w, "Storage error", http.StatusInternalServerError)
			return
		}
		
		if !started {
			cached, _ := m.store.Get(ctx, key)
			if cached != nil && cached.Completed {
				w.WriteHeader(cached.StatusCode)
				w.Write(cached.Body)
				return
			}
			http.Error(w, "Duplicate request in progress", http.StatusConflict)
			return
		}
		
		recorder := httptest.NewRecorder()
		next.ServeHTTP(recorder, r)
		
		m.store.Finish(ctx, key, recorder.Code, recorder.Body.Bytes())
		
		for k, vals := range recorder.Header() {
			for _, v := range vals {
				w.Header().Add(k, v)
			}
		}
		w.WriteHeader(recorder.Code)
		w.Write(recorder.Body.Bytes())
	})
}

type PaymentHandler struct{}

func (h *PaymentHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	time.Sleep(2 * time.Second)
	
	response := map[string]interface{}{
		"status":         "paid",
		"amount":         1000,
		"transaction_id": "uuid-1234-5678",
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func main2() {
	store := NewRedisStore("localhost:6379", 24*time.Hour)
	middleware := NewIdempotencyMiddleware(store)
	
	paymentHandler := &PaymentHandler{}
	
	mux := http.NewServeMux()
	mux.Handle("/pay", middleware.Handler(paymentHandler))
	
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	
	go func() {
		fmt.Println("Server starting on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Server error: %v\n", err)
		}
	}()
	
	time.Sleep(1 * time.Second)
	
	testConcurrentRequests()
	
	select {}
}

func testConcurrentRequests() {
	client := &http.Client{Timeout: 10 * time.Second}
	url := "http://localhost:8080/pay"
	idempotencyKey := "test-key-123"
	
	var wg sync.WaitGroup
	results := make(chan int, 10)
	
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			
			req, _ := http.NewRequest("POST", url, nil)
			req.Header.Set("Idempotency-Key", idempotencyKey)
			
			resp, err := client.Do(req)
			if err != nil {
				fmt.Printf("Goroutine %d: error - %v\n", id, err)
				results <- 0
				return
			}
			defer resp.Body.Close()
			
			fmt.Printf("Goroutine %d: status %d\n", id, resp.StatusCode)
			results <- resp.StatusCode
		}(i)
	}
	
	wg.Wait()
	close(results)
	
	successCount := 0
	conflictCount := 0
	for code := range results {
		if code == 200 {
			successCount++
		} else if code == 409 {
			conflictCount++
		}
	}
	
	fmt.Printf("\nResults: %d successful, %d conflicts\n", successCount, conflictCount)
}