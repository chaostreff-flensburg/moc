package router

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"runtime/debug"

	session "github.com/chaostreff-flensburg/moc/context"
)

// HTTPError is an error with a message and an HTTP status code.
type HTTPError struct {
	Object          string      `json:"object"`
	Code            int         `json:"code"`
	Message         string      `json:"msg"`
	Json            interface{} `json:"json"`
	InternalError   error       `json:"-"`
	InternalMessage string      `json:"-"`
	ErrorID         string      `json:"error_id,omitempty"`
}

// ======================================
// Return bad request error
// ======================================
func BadRequestError(fmtString string, args ...interface{}) *HTTPError {
	return httpError(http.StatusBadRequest, fmtString, args...)
}

// ======================================
// Return internal server error
// ======================================
func InternalServerError(fmtString string, args ...interface{}) *HTTPError {
	return httpError(http.StatusInternalServerError, fmtString, args...)
}

// ======================================
// Return not found error
// ======================================
func NotFoundError(fmtString string, args ...interface{}) *HTTPError {
	return httpError(http.StatusNotFound, fmtString, args...)
}

// ======================================
// Return unauthorized error
// ======================================
func UnauthorizedError(fmtString string, args ...interface{}) *HTTPError {
	return httpError(http.StatusUnauthorized, fmtString, args...)
}

// ======================================
// Return unavailable service error
// ======================================
func UnavailableServiceError(fmtString string, args ...interface{}) *HTTPError {
	return httpError(http.StatusServiceUnavailable, fmtString, args...)
}

// ======================================
// Return error as string
// ======================================
func (e *HTTPError) Error() string {
	if e.InternalMessage != "" {
		return e.InternalMessage
	}
	return fmt.Sprintf("%d: %s", e.Code, e.Message)
}

// ======================================
// Cause returns the root cause error
// ======================================
func (e *HTTPError) Cause() error {
	if e.InternalError != nil {
		return e.InternalError
	}
	return e
}

// ======================================
// WithJsonError
// ======================================
func (e *HTTPError) WithJsonError(json interface{}) *HTTPError {
	e.Json = json
	return e
}

// ======================================
// WithInternalError adds internal error information to the error
// ======================================
func (e *HTTPError) WithInternalError(err error) *HTTPError {
	e.InternalError = err
	return e
}

// ======================================
// WithInternalMessage adds internal message information to the error
// ======================================
func (e *HTTPError) WithInternalMessage(fmtString string, args ...interface{}) *HTTPError {
	e.InternalMessage = fmt.Sprintf(fmtString, args...)
	return e
}

func httpError(code int, fmtString string, args ...interface{}) *HTTPError {
	return &HTTPError{
		Object:  "error",
		Code:    code,
		Message: fmt.Sprintf(fmtString, args...),
	}
}

// ======================================
// Recoverer is a middleware that recovers from panics, logs the panic (and a
// backtrace), and returns a HTTP 500 (Internal Server Error) status if
// possible. Recoverer prints a request ID if one is provided.
// ======================================
func Recoverer(w http.ResponseWriter, r *http.Request) (context.Context, error) {
	ctx := r.Context()

	defer func() {
		if rvr := recover(); rvr != nil {

			log := session.GetLogger(ctx)
			if log != nil {
				log.Panic(rvr, debug.Stack())
			} else {
				fmt.Fprintf(os.Stderr, "Panic: %+v\n", rvr)
				debug.PrintStack()
			}

			se := &HTTPError{
				Object:  "error",
				Code:    http.StatusInternalServerError,
				Message: http.StatusText(http.StatusInternalServerError),
			}
			handleError(se, w, r)
		}
	}()

	return nil, nil
}

// ======================================
// Handle Errors
// ======================================
func handleError(err error, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	log := session.GetLogger(ctx)
	errorID := session.GetRequestID(r.Context())
	switch e := err.(type) {
	case *HTTPError:
		if e.Code >= http.StatusInternalServerError {
			e.ErrorID = *errorID
			// this will get us the stack trace too
			log.WithError(e.Cause()).Error(e.Error())
		} else {
			log.WithError(e.Cause()).Info(e.Error())
		}
		if jsonErr := SendJSON(w, e.Code, e); jsonErr != nil {
			handleError(jsonErr, w, r)
		}
	default:
		log.WithError(e).Errorf("Unhandled server error: %s", e.Error())
		// hide real error details from response to prevent info leaks
		w.WriteHeader(http.StatusInternalServerError)
		if _, writeErr := w.Write([]byte(`{object:"error",code":500,"msg":"Internal server error","error_id":"` + *errorID + `"}`)); writeErr != nil {
			log.WithError(writeErr).Error("Error writing generic error message")
		}
	}
}
