package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/mlaitechio/vagais/internal/models"
	"github.com/mlaitechio/vagais/internal/services"
)

// AuthMiddleware validates JWT tokens
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Bearer token required"})
			c.Abort()
			return
		}

		user, err := services.AuthServiceInstance.GetUserFromToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Set("user", user)
		c.Next()
	}
}

// OptionalAuthMiddleware validates JWT tokens if present
func OptionalAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.Next()
			return
		}

		user, err := services.AuthServiceInstance.GetUserFromToken(tokenString)
		if err == nil {
			c.Set("user", user)
		}

		c.Next()
	}
}

// RoleMiddleware checks if user has required role
func RoleMiddleware(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userInterface, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			c.Abort()
			return
		}

		user := userInterface.(*models.User)
		hasRole := false
		for _, role := range roles {
			if user.Role == role {
				hasRole = true
				break
			}
		}

		if !hasRole {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// CORS middleware for cross-origin requests
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		c.Header("Access-Control-Expose-Headers", "Content-Length")
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// RateLimiter middleware for API rate limiting
func RateLimiter() gin.HandlerFunc {
	// Simple in-memory rate limiter
	// In production, use Redis for distributed rate limiting
	clients := make(map[string]*rateLimiter)

	return func(c *gin.Context) {
		clientIP := c.ClientIP()

		limiter, exists := clients[clientIP]
		if !exists {
			limiter = &rateLimiter{
				requests: make([]time.Time, 0),
				limit:    100, // requests per minute
				window:   time.Minute,
			}
			clients[clientIP] = limiter
		}

		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit exceeded"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// SecurityHeaders adds security headers
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		c.Header("Content-Security-Policy", "default-src 'self'")
		c.Next()
	}
}

// RequestID adds a unique request ID
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}
		c.Header("X-Request-ID", requestID)
		c.Set("request_id", requestID)
		c.Next()
	}
}

// DomainBlockMiddleware blocks specific domains from registration
func DomainBlockMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "POST" && c.Request.URL.Path == "/v1/auth/register" {
			var req services.RegisterRequest
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
				c.Abort()
				return
			}

			// Check blocked domains
			email := strings.ToLower(req.Email)
			for _, blockedDomain := range services.AuthServiceInstance.GetConfig().Security.BlockedDomains {
				if strings.Contains(email, "@"+blockedDomain) {
					c.JSON(http.StatusForbidden, gin.H{"error": "Domain not allowed"})
					c.Abort()
					return
				}
			}
		}

		c.Next()
	}
}

// LicenseCheckMiddleware checks license status (optional feature)
func LicenseCheckMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userInterface, exists := c.Get("user")
		if !exists {
			c.Next()
			return
		}

		user := userInterface.(*models.User)

		// Skip license check for single user orgs
		if user.OrganizationID == nil {
			c.Next()
			return
		}

		// Check license if organization exists
		if user.OrganizationID != nil {
			if !services.LicenseServiceInstance.IsAvailable() {
				c.Next()
				return
			}

			isValid, _ := models.GetLicenseStatus(user.OrganizationID, services.LicenseServiceInstance.GetDB())
			if !isValid {
				c.JSON(http.StatusForbidden, gin.H{"error": "License required or expired"})
				c.Abort()
				return
			}
		}

		c.Next()
	}
}

// PaymentCheckMiddleware checks payment status (optional feature)
func PaymentCheckMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userInterface, exists := c.Get("user")
		if !exists {
			c.Next()
			return
		}

		user := userInterface.(*models.User)

		// Check payment status
		if !services.PaymentServiceInstance.IsAvailable() {
			c.Next()
			return
		}

		isValid, _ := models.GetPaymentStatus(user.ID, services.PaymentServiceInstance.GetDB())
		if !isValid {
			c.JSON(http.StatusForbidden, gin.H{"error": "Payment required"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// rateLimiter implements a simple rate limiter
type rateLimiter struct {
	requests []time.Time
	limit    int
	window   time.Duration
}

func (rl *rateLimiter) Allow() bool {
	now := time.Now()

	// Remove old requests outside the window
	var validRequests []time.Time
	for _, req := range rl.requests {
		if now.Sub(req) <= rl.window {
			validRequests = append(validRequests, req)
		}
	}
	rl.requests = validRequests

	// Check if we're under the limit
	if len(rl.requests) < rl.limit {
		rl.requests = append(rl.requests, now)
		return true
	}

	return false
}
