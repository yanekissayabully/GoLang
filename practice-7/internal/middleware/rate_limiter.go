package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type Limiter struct {
	mu    sync.Mutex
	store map[string][]time.Time
	limit int
	win   int
}

func NewRateLimiter(limit, window int) *Limiter {
	return &Limiter{
		store: make(map[string][]time.Time),
		limit: limit,
		win:   window,
	}
}

func (l *Limiter) allow(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	
	now := time.Now()
	cut := now.Add(-time.Duration(l.win) * time.Second)
	
	times := l.store[key]
	valid := []time.Time{}
	for _, t := range times {
		if t.After(cut) {
			valid = append(valid, t)
		}
	}
	
	if len(valid) >= l.limit {
		return false
	}
	
	l.store[key] = append(valid, now)
	return true
}

func RateLimitMiddleware(l *Limiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.ClientIP()
		if id, exists := c.Get("userID"); exists {
			key = "user:" + id.(string)
		}
		
		if !l.allow(key) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "too many requests"})
			return
		}
		c.Next()
	}
}