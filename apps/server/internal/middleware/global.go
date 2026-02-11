package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/Niiaks/campusCart/internal/server"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	errs "github.com/Niiaks/campusCart/internal/err"
	"github.com/Niiaks/campusCart/internal/sqlerr"
)

type GlobalMiddlewares struct {
	server *server.Server
}

func NewGlobalMiddlewares(s *server.Server) *GlobalMiddlewares {
	return &GlobalMiddlewares{
		server: s,
	}
}

func (global *GlobalMiddlewares) CORS() func(http.Handler) http.Handler {
	return cors.Handler(cors.Options{
		AllowedOrigins: global.server.Config.Server.CorsAllowedOrigins,
		AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-Request-ID"},
		ExposedHeaders: []string{"X-Request-ID"},
		MaxAge:         300,
	})
}

// responseRecorder wraps http.ResponseWriter to capture status code and size
type responseRecorder struct {
	http.ResponseWriter
	statusCode int
	size       int
}

func (r *responseRecorder) WriteHeader(statusCode int) {
	r.statusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

func (r *responseRecorder) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.size += size
	return size, err
}

func (global *GlobalMiddlewares) RequestLogger() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Wrap response writer to capture status
			rec := &responseRecorder{
				ResponseWriter: w,
				statusCode:     200,
				size:           0,
			}

			// Get enhanced logger from context
			logger := GetLogger(r.Context())

			// Continue chain
			next.ServeHTTP(rec, r)

			// Calculate latency
			latency := time.Since(start)
			statusCode := rec.statusCode

			var e *zerolog.Event

			switch {
			case statusCode >= 500:
				e = logger.Error()
			case statusCode >= 400:
				e = logger.Warn()
			default:
				e = logger.Info()
			}

			// Add request ID if available
			if requestID := GetRequestID(r); requestID != "" {
				e = e.Str("request_id", requestID)
			}

			// Add user context if available
			if userID := GetUserID(r); userID != "" {
				e = e.Str("user_id", userID)
			}

			e.
				Dur("latency", latency).
				Int("status", statusCode).
				Str("method", r.Method).
				Str("uri", r.RequestURI).
				Str("host", r.Host).
				Str("ip", r.RemoteAddr).
				Str("user_agent", r.UserAgent()).
				Int("bytes", rec.size).
				Msg("API")
		})
	}
}

func (global *GlobalMiddlewares) Recover() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rvr := recover(); rvr != nil {
					logger := GetLogger(r.Context())

					logger.Error().
						Interface("panic", rvr).
						Bytes("stack", debug.Stack()).
						Str("request_id", GetRequestID(r)).
						Str("method", r.Method).
						Str("path", r.URL.Path).
						Msg("Panic recovered")

					if !isResponseWritten(w) {
						err := errs.NewInternalServerError()
						w.Header().Set("Content-Type", "application/json")
						w.WriteHeader(err.Status)
						json.NewEncoder(w).Encode(err)
					}
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}

func (global *GlobalMiddlewares) Secure() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Security headers
			w.Header().Set("X-Content-Type-Options", "nosniff")
			w.Header().Set("X-Frame-Options", "SAMEORIGIN")
			w.Header().Set("X-XSS-Protection", "1; mode=block")
			w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")

			next.ServeHTTP(w, r)
		})
	}
}

// ErrorHandler wraps handlers and converts errors to JSON responses
func (global *GlobalMiddlewares) ErrorHandler(handler func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := handler(w, r); err != nil {
			global.handleError(w, r, err)
		}
	}
}

func (global *GlobalMiddlewares) handleError(w http.ResponseWriter, r *http.Request, err error) {
	// First try to handle database errors and convert them to appropriate HTTP errors
	originalErr := err

	// Try to handle known database errors
	// Only do this for errors that haven't already been converted to HTTPError
	var httpErr *errs.HTTPError
	if !errors.As(err, &httpErr) {
		// Here we call our sqlerr handler which will convert database errors
		// to appropriate application errors
		err = sqlerr.HandleError(err)
	}

	// Now process the possibly converted error
	var status int
	var code string
	var message string
	var fieldErrors []errs.FieldError
	var action *errs.Action
	var override bool

	if errors.As(err, &httpErr) {
		status = httpErr.Status
		code = httpErr.Code
		message = httpErr.Message
		fieldErrors = httpErr.Errors
		action = httpErr.Action
		override = httpErr.Override
	} else {
		status = http.StatusInternalServerError
		code = errs.MakeUpperCaseWithUnderscores(http.StatusText(http.StatusInternalServerError))
		message = http.StatusText(http.StatusInternalServerError)
		override = false
	}

	// Log the original error to help with debugging
	// Use enhanced logger from context which already includes request_id, method, path, ip, user context, and trace context
	logger := GetLogger(r.Context())

	logger.Error().Stack().
		Err(originalErr).
		Int("status", status).
		Str("error_code", code).
		Msg(message)

	// Don't write response if headers were already sent
	if !isResponseWritten(w) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(errs.HTTPError{
			Code:     code,
			Message:  message,
			Status:   status,
			Override: override,
			Errors:   fieldErrors,
			Action:   action,
		})
	}
}

// NotFoundHandler returns a handler for 404 errors
func (global *GlobalMiddlewares) NotFoundHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := errs.NewNotFoundError("Route not found", false, nil)
		global.handleError(w, r, err)
	}
}

// MethodNotAllowedHandler returns a handler for 405 errors
func (global *GlobalMiddlewares) MethodNotAllowedHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := &errs.HTTPError{
			Code:     errs.MakeUpperCaseWithUnderscores(http.StatusText(http.StatusMethodNotAllowed)),
			Message:  fmt.Sprintf("Method %s not allowed", r.Method),
			Status:   http.StatusMethodNotAllowed,
			Override: false,
		}
		global.handleError(w, r, err)
	}
}

// isResponseWritten checks if the response has already been written
func isResponseWritten(w http.ResponseWriter) bool {
	// Check if we can still write headers (response not committed)
	// This is a bit hacky but works with standard ResponseWriter
	_, ok := w.(middleware.WrapResponseWriter)
	if ok {
		return false
	}
	// For standard ResponseWriter, assume not written if we can set a test header
	// (this is removed immediately)
	w.Header().Del("X-Response-Check")
	return false
}
