package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"time"
)

func IsRetryable(resp *http.Response, err error) bool {
	if err != nil {
		return true
	}
	if resp == nil {
		return false
	}
	if resp.StatusCode == http.StatusTooManyRequests ||
		resp.StatusCode == http.StatusInternalServerError ||
		resp.StatusCode == http.StatusBadGateway ||
		resp.StatusCode == http.StatusServiceUnavailable ||
		resp.StatusCode == http.StatusGatewayTimeout {
		return true
	}
	return false
}

func CalculateBackoff(attempt int) time.Duration {
	baseDelay := 500 * time.Millisecond
	maxDelay := 10 * time.Second
	backoff := baseDelay * time.Duration(math.Pow(2, float64(attempt)))
	if backoff > maxDelay {
		backoff = maxDelay
	}
	jitter := time.Duration(rand.Int63n(int64(backoff)))
	return jitter
}

type PaymentClient struct {
	client     *http.Client
	maxRetries int
}

func NewPaymentClient(maxRetries int) *PaymentClient {
	return &PaymentClient{
		client:     &http.Client{Timeout: 5 * time.Second},
		maxRetries: maxRetries,
	}
}

func (c *PaymentClient) ExecutePayment(ctx context.Context, url string) error {
	var lastErr error
	var lastResp *http.Response

	for attempt := 0; attempt <= c.maxRetries; attempt++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		req, err := http.NewRequestWithContext(ctx, "POST", url, nil)
		if err != nil {
			return err
		}

		lastResp, lastErr = c.client.Do(req)
		
		if lastErr == nil && lastResp.StatusCode == http.StatusOK {
			if lastResp.Body != nil {
				lastResp.Body.Close()
			}
			fmt.Printf("Attempt %d: Success!\n", attempt+1)
			return nil
		}

		if IsRetryable(lastResp, lastErr) {
			if lastResp != nil && lastResp.Body != nil {
				lastResp.Body.Close()
			}
			
			if attempt == c.maxRetries {
				break
			}
			
			waitTime := CalculateBackoff(attempt)
			fmt.Printf("Attempt %d failed: waiting %v...\n", attempt+1, waitTime)
			
			timer := time.NewTimer(waitTime)
			select {
			case <-ctx.Done():
				timer.Stop()
				return ctx.Err()
			case <-timer.C:
			}
		} else {
			if lastResp != nil && lastResp.Body != nil {
				lastResp.Body.Close()
			}
			return fmt.Errorf("non-retryable error: %v", lastErr)
		}
	}
	
	if lastResp != nil && lastResp.Body != nil {
		lastResp.Body.Close()
	}
	return fmt.Errorf("max retries exceeded: %v", lastErr)
}

func main1() {
	rand.Seed(time.Now().UnixNano())
	
	requestCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		if requestCount < 4 {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(`{"error": "service unavailable"}`))
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})
	}))
	defer server.Close()
	
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	client := NewPaymentClient(5)
	err := client.ExecutePayment(ctx, server.URL)
	if err != nil {
		fmt.Printf("Payment failed: %v\n", err)
	} else {
		fmt.Println("Payment completed successfully")
	}
}