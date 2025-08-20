package routes

import (
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/parthvinchhi/rd-website/pkg/handlers"
)

func SetupRouter(authHandler *handlers.AuthHandler) *gin.Engine {
	router := gin.Default()

	// CORS configuration
	corsCfg := cors.DefaultConfig()
	allowed := os.Getenv("CORS_ALLOWED_ORIGINS")
	if allowed == "" {
		corsCfg.AllowOrigins = []string{"http://localhost:5173", "http://127.0.0.1:5173"}
	} else {
		corsCfg.AllowOrigins = strings.Split(allowed, ",")
	}
	corsCfg.AllowCredentials = true
	corsCfg.AddAllowHeaders("Content-Type", "Authorization")
	router.Use(cors.New(corsCfg))

	// Limit request body size to mitigate DoS via large payloads
	router.Use(func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 1<<20) // 1MB
		c.Next()
	})

	api := router.Group("/api")
	{
		api.POST("/signup", authHandler.Signup)

		// Simple in-memory rate-limit for login endpoint
		loginLimiter := newIPRateLimiter(5, time.Minute)
		api.POST("/login", func(c *gin.Context) {
			if !loginLimiter.Allow(c.ClientIP()) {
				c.JSON(http.StatusTooManyRequests, gin.H{"error": "too many requests"})
				return
			}
			authHandler.Login(c)
		})
	}

	// Static file serving for frontend (assumes build output in web/dist)
	router.Static("/assets", "web/dist/assets")

	// SPA fallback for non-API routes
	router.NoRoute(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/api") {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.File("web/dist/index.html")
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
