// Package api provides middleware for the AI Agent Society HTTP server.
//
// =============================================================================
// PURPOSE:
// =============================================================================
// This file contains HTTP middleware functions for cross-cutting concerns
// like logging, authentication, rate limiting, and request tracking.
//
// =============================================================================
// MIDDLEWARE IMPLEMENTATIONS:
// =============================================================================
//
// func RequestIDMiddleware() gin.HandlerFunc
//   - Generate unique request ID for each incoming request
//   - Add to context and response headers (X-Request-ID)
//   - Enable request tracing across logs
//
// func LoggerMiddleware(logger *slog.Logger) gin.HandlerFunc
//   - Log request method, path, status code, latency
//   - Include request ID for correlation
//   - Log request body for POST/PUT (with size limits)
//
// func CORSMiddleware() gin.HandlerFunc
//   - Allow cross-origin requests from dashboard frontend
//   - Configure allowed methods: GET, POST, PUT, DELETE, OPTIONS
//   - Configure allowed headers: Content-Type, Authorization
//
// func RateLimitMiddleware(rps int) gin.HandlerFunc
//   - Protect API from abuse
//   - Token bucket algorithm per IP address
//   - Return 429 Too Many Requests when exceeded
//
// func RecoveryMiddleware(logger *slog.Logger) gin.HandlerFunc
//   - Catch panics in handlers
//   - Log stack trace with request context
//   - Return 500 Internal Server Error
//
// =============================================================================
// ERROR HANDLING:
// =============================================================================
//
// type APIError struct {
//     Code    string `json:"code"`
//     Message string `json:"message"`
//     Details any    `json:"details,omitempty"`
// }
//
// func ErrorHandler() gin.HandlerFunc
//   - Centralized error response formatting
//   - Map internal errors to appropriate HTTP status codes
//   - Sanitize error messages for production (hide internal details)

package api
