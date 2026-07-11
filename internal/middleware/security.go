package middleware

import (
	"log"
	"net/http"
)

// SecurityHeaders adds security headers to HTTP responses
func SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Content Security Policy
		w.Header().Set("Content-Security-Policy",
			"default-src 'self'; "+
			"script-src 'self' 'unsafe-inline'; "+
			"style-src 'self' 'unsafe-inline' https://fonts.googleapis.com; "+
			"font-src 'self' https://fonts.gstatic.com; "+
			"img-src 'self' data: https:; "+
			"connect-src 'self'; "+
			"frame-ancestors 'none'; "+
			"form-action 'self';")

		// X-Frame-Options
		w.Header().Set("X-Frame-Options", "DENY")

		// X-Content-Type-Options
		w.Header().Set("X-Content-Type-Options", "nosniff")

		// Referrer Policy
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

		// X-XSS-Protection (legacy, but still useful)
		w.Header().Set("X-XSS-Protection", "1; mode=block")

		// Strict-Transport-Security (only in HTTPS)
		if r.URL.Scheme == "https" {
			w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
		}

		// Permissions Policy
		w.Header().Set("Permissions-Policy",
			"geolocation=(), "+
			"microphone=(), "+
			"camera=(), "+
			"payment=(), "+
			"usb=(), "+
			"magnetometer=(), "+
			"gyroscope=(), "+
			"accelerometer=()")

		// Remove server information
		w.Header().Set("Server", "")

		next.ServeHTTP(w, r)
	})
}

// CORS provides CORS headers for API endpoints
func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		w.Header().Set("Access-Control-Max-Age", "86400")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// RateLimiting provides basic rate limiting (simplified for demo)
func RateLimiting(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// In production, you'd want to use a proper rate limiting solution
		// like a token bucket algorithm with Redis or similar

		// For now, we'll just add a simple delay for form submissions
		if r.Method == "POST" && r.URL.Path == "/contact" {
			// Add a small delay to prevent brute force attacks
			// In production, you'd want proper rate limiting
		}

		next.ServeHTTP(w, r)
	})
}

// RequestLogging logs HTTP requests
func RequestLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Log basic request information
		// In production, you'd want structured logging
		log.Printf("%s %s %s", r.Method, r.URL.Path, r.RemoteAddr)

		next.ServeHTTP(w, r)
	})
}
