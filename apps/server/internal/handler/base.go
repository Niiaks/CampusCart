package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Niiaks/campusCart/internal/middleware"
	"github.com/Niiaks/campusCart/internal/server"
	"github.com/Niiaks/campusCart/internal/validation"
	nrpkgerrors "github.com/newrelic/go-agent/v3/integrations/nrpkgerrors"
	"github.com/newrelic/go-agent/v3/newrelic"
)

// Handler provides base functionality for all handlers
type Handler struct {
	server *server.Server
}

func NewHandler(server *server.Server) Handler {
	return Handler{server: server}
}

func (h Handler) JSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(payload)
	if err != nil {
		return
	}
}

// HandlerFunc represents a typed handler function that processes a request and returns a response
type HandlerFunc[Req validation.Validatable, Res any] func(w http.ResponseWriter, r *http.Request, req Req) (Res, error)

// HandlerFuncNoContent represents a typed handler function that processes a request without returning content
type HandlerFuncNoContent[Req validation.Validatable] func(w http.ResponseWriter, r *http.Request, req Req) error

// ResponseHandler defines the interface for handling different response types
type ResponseHandler interface {
	Handle(w http.ResponseWriter, result interface{}) error
	GetOperation() string
	AddAttributes(txn *newrelic.Transaction, result interface{})
}

// JSONResponseHandler handles JSON responses
type JSONResponseHandler struct {
	status int
}

func (h JSONResponseHandler) Handle(w http.ResponseWriter, result interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(h.status)
	return json.NewEncoder(w).Encode(result)
}

func (h JSONResponseHandler) GetOperation() string {
	return "handler"
}

func (h JSONResponseHandler) AddAttributes(txn *newrelic.Transaction, result interface{}) {
	// http.status_code is already set by tracing middleware
}

// NoContentResponseHandler handles no-content responses
type NoContentResponseHandler struct {
	status int
}

func (h NoContentResponseHandler) Handle(w http.ResponseWriter, result interface{}) error {
	w.WriteHeader(h.status)
	return nil
}

func (h NoContentResponseHandler) GetOperation() string {
	return "handler_no_content"
}

func (h NoContentResponseHandler) AddAttributes(txn *newrelic.Transaction, result interface{}) {
	// http.status_code is already set by tracing middleware
}

// FileResponseHandler handles file responses
type FileResponseHandler struct {
	status      int
	filename    string
	contentType string
}

func (h FileResponseHandler) Handle(w http.ResponseWriter, result interface{}) error {
	data, ok := result.([]byte)
	if !ok {
		return fmt.Errorf("expected []byte for file response, got %T", result)
	}
	w.Header().Set("Content-Type", h.contentType)
	w.Header().Set("Content-Disposition", "attachment; filename="+h.filename)
	w.WriteHeader(h.status)
	_, err := w.Write(data)
	return err
}

func (h FileResponseHandler) GetOperation() string {
	return "handler_file"
}

func (h FileResponseHandler) AddAttributes(txn *newrelic.Transaction, result interface{}) {
	if txn != nil {
		// http.status_code is already set by tracing middleware
		txn.AddAttribute("file.name", h.filename)
		txn.AddAttribute("file.content_type", h.contentType)
		if data, ok := result.([]byte); ok {
			txn.AddAttribute("file.size_bytes", len(data))
		}
	}
}

// handleRequest is the unified handler function that eliminates code duplication
func handleRequest[Req validation.Validatable](
	w http.ResponseWriter,
	r *http.Request,
	req Req,
	handler func(w http.ResponseWriter, r *http.Request, req Req) (interface{}, error),
	responseHandler ResponseHandler,
) error {
	start := time.Now()
	method := r.Method
	path := r.URL.Path

	// Get New Relic transaction from context
	txn := newrelic.FromContext(r.Context())
	if txn != nil {
		txn.AddAttribute("handler.name", path)
		// http.method and http.route are already set by nrhttp middleware
		responseHandler.AddAttributes(txn, nil)
	}

	// Get context-enhanced logger
	loggerBuilder := middleware.GetLogger(r.Context()).With().
		Str("operation", responseHandler.GetOperation()).
		Str("method", method).
		Str("path", path)

	// Add file-specific fields to logger if it's a file handler
	if fileHandler, ok := responseHandler.(FileResponseHandler); ok {
		loggerBuilder = loggerBuilder.
			Str("filename", fileHandler.filename).
			Str("content_type", fileHandler.contentType)
	}

	logger := loggerBuilder.Logger()

	// user.id is already set by tracing middleware

	logger.Info().Msg("handling request")

	// Validation with observability
	validationStart := time.Now()
	if err := validation.BindAndValidate(r, req); err != nil {
		validationDuration := time.Since(validationStart)

		logger.Error().
			Err(err).
			Dur("validation_duration", validationDuration).
			Msg("request validation failed")

		if txn != nil {
			txn.NoticeError(nrpkgerrors.Wrap(err))
			txn.AddAttribute("validation.status", "failed")
			txn.AddAttribute("validation.duration_ms", validationDuration.Milliseconds())
		}
		return err
	}

	validationDuration := time.Since(validationStart)
	if txn != nil {
		txn.AddAttribute("validation.status", "success")
		txn.AddAttribute("validation.duration_ms", validationDuration.Milliseconds())
	}

	logger.Debug().
		Dur("validation_duration", validationDuration).
		Msg("request validation successful")

	// Execute handler with observability
	handlerStart := time.Now()
	result, err := handler(w, r, req)
	handlerDuration := time.Since(handlerStart)

	if err != nil {
		totalDuration := time.Since(start)

		logger.Error().
			Err(err).
			Dur("handler_duration", handlerDuration).
			Dur("total_duration", totalDuration).
			Msg("handler execution failed")

		if txn != nil {
			txn.NoticeError(nrpkgerrors.Wrap(err))
			txn.AddAttribute("handler.status", "error")
			txn.AddAttribute("handler.duration_ms", handlerDuration.Milliseconds())
			txn.AddAttribute("total.duration_ms", totalDuration.Milliseconds())
		}
		return err
	}

	totalDuration := time.Since(start)

	// Record success metrics and tracing
	if txn != nil {
		txn.AddAttribute("handler.status", "success")
		txn.AddAttribute("handler.duration_ms", handlerDuration.Milliseconds())
		txn.AddAttribute("total.duration_ms", totalDuration.Milliseconds())
		responseHandler.AddAttributes(txn, result)
	}

	logger.Info().
		Dur("handler_duration", handlerDuration).
		Dur("validation_duration", validationDuration).
		Dur("total_duration", totalDuration).
		Msg("request completed successfully")

	return responseHandler.Handle(w, result)
}

// Handle wraps a handler with validation, error handling, logging, metrics, and tracing
func Handle[Req validation.Validatable, Res any](
	h Handler,
	handler HandlerFunc[Req, Res],
	status int,
	reqFactory func() Req,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := reqFactory() // Create fresh instance per request to avoid concurrency issues
		err := handleRequest(w, r, req, func(w http.ResponseWriter, r *http.Request, req Req) (interface{}, error) {
			return handler(w, r, req)
		}, JSONResponseHandler{status: status})
		if err != nil {
			// Error is handled by global error handler middleware
			panic(err)
		}
	}
}

// HandleFile wraps a handler that returns a file with validation, error handling, logging, metrics, and tracing
func HandleFile[Req validation.Validatable](
	h Handler,
	handler HandlerFunc[Req, []byte],
	status int,
	reqFactory func() Req,
	filename string,
	contentType string,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := reqFactory() // Create fresh instance per request to avoid concurrency issues
		err := handleRequest(w, r, req, func(w http.ResponseWriter, r *http.Request, req Req) (interface{}, error) {
			return handler(w, r, req)
		}, FileResponseHandler{
			status:      status,
			filename:    filename,
			contentType: contentType,
		})
		if err != nil {
			// Error is handled by global error handler middleware
			panic(err)
		}
	}
}

// HandleNoContent wraps a handler with validation, error handling, logging, metrics, and tracing for endpoints that don't return content
func HandleNoContent[Req validation.Validatable](
	h Handler,
	handler HandlerFuncNoContent[Req],
	status int,
	reqFactory func() Req,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := reqFactory() // Create fresh instance per request to avoid concurrency issues
		err := handleRequest(w, r, req, func(w http.ResponseWriter, r *http.Request, req Req) (interface{}, error) {
			err := handler(w, r, req)
			return nil, err
		}, NoContentResponseHandler{status: status})
		if err != nil {
			// Error is handled by global error handler middleware
			panic(err)
		}
	}
}
