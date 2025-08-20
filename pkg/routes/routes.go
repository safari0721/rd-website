package routes

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/parthvinchhi/rd-website/pkg/handlers"
)

func SetupRouter(authHandler *handlers.AuthHandler) *gin.Engine {
	router := gin.Default()
	// Limit request body size to mitigate DoS via large payloads
	router.Use(func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 1<<20) // 1MB
		c.Next()
	})

	router.POST("/signup", authHandler.Signup)

	// Simple in-memory rate-limit for login endpoint
	loginLimiter := newIPRateLimiter(5, time.Minute)
	router.POST("/login", func(c *gin.Context) {
		if !loginLimiter.Allow(c.ClientIP()) {
			c.JSON(429, gin.H{"error": "too many requests"})
			return
		}
		authHandler.Login(c)
	})

	return router
}

// in-memory IP rate limiter (fixed window)
type ipRateLimiter struct {
	mu        sync.Mutex
	allowance map[string]int
	resetAt   time.Time
	limit     int
	window    time.Duration
}

func newIPRateLimiter(limit int, window time.Duration) *ipRateLimiter {
	return &ipRateLimiter{
		allowance: make(map[string]int),
		resetAt:   time.Now().Add(window),
		limit:     limit,
		window:    window,
	}
}

func (l *ipRateLimiter) Allow(ip string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	if time.Now().After(l.resetAt) {
		l.allowance = make(map[string]int)
		l.resetAt = time.Now().Add(l.window)
	}
	l.allowance[ip]++
	return l.allowance[ip] <= l.limit
}