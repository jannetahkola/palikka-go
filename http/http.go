package http

import (
	"bufio"
	"context"
	"encoding/json"
	"log/slog"
	"net"
	"net/http"
	"os"
	palikka "palikka-go"
)

const (
	HeaderContentType = "Content-Type"
	HeaderLocation    = "Location"
	MediaTypeJSON     = "application/json"
)

const (
	ErrMsgInvalidRequest = "invalid request"
	ErrMsgUnauthorized   = "unauthorized"
	ErrMsgForbidden      = "forbidden"
	ErrMsgNotFound       = "not found"
	ErrMsgInternal       = "internal server error"
)

// responseWrapper wraps a ResponseWriter to expose any response code written to it
// with WriteHeader.
//
// Supports http.ResponseController by implementing Unwrap.
// Supports http.Hijacker by implementing Hijack (for WebSockets).
//
// See https://pkg.go.dev/net/http#ResponseController
type responseWrapper struct {
	http.ResponseWriter
	code int
}

func (rw *responseWrapper) WriteHeader(code int) {
	rw.code = code
	rw.ResponseWriter.WriteHeader(code)
}

// Unwrap returns the original ResponseWriter.
func (rw *responseWrapper) Unwrap() http.ResponseWriter {
	return rw.ResponseWriter
}

// Hijack calls the original ResponseWriter's Hijack func.
func (rw *responseWrapper) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return rw.ResponseWriter.(http.Hijacker).Hijack()
}

func authorizeAll() []string {
	return []string{}
}

func authorizeAllOf(perms ...string) []string {
	if len(perms) == 0 {
		panic("no permissions for authorizeAllOf()")
	}
	return perms
}

// **********
// request context
// **********

type requestContextKey string

const (
	requestContextKeyMethod requestContextKey = "method"
	requestContextKeyPath   requestContextKey = "path"
)

func copyContextWithRequestAttributes(r *http.Request) *http.Request {
	r = r.WithContext(context.WithValue(r.Context(), requestContextKeyMethod, r.Method))
	r = r.WithContext(context.WithValue(r.Context(), requestContextKeyPath, r.URL.Path))
	return r
}

func copyContextWithLogger(r *http.Request) *http.Request {
	l := slog.New(&logContextHandler{Handler: slog.NewJSONHandler(os.Stderr, nil)})
	r = r.WithContext(palikka.CopyContextWithLogger(r.Context(), l))
	return r
}

// **********
// logging
// **********

type logContextHandler struct {
	slog.Handler
}

func (h *logContextHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.Handler.Enabled(ctx, level)
}

func (h *logContextHandler) Handle(ctx context.Context, r slog.Record) error {
	r.AddAttrs(slog.Any("method", ctx.Value(requestContextKeyMethod)))
	r.AddAttrs(slog.Any("path", ctx.Value(requestContextKeyPath)))

	if session := palikka.SessionFromContext(ctx); session != nil {
		r.AddAttrs(slog.Any("session", session.ID))
		r.AddAttrs(slog.Any("user", session.UserID))
	}

	// todo do we need user here?
	//if user := palikka.UserFromContext(ctx); user != nil {
	//	r.AddAttrs(slog.Any("user", user.ID))
	//}

	return h.Handler.Handle(ctx, r)
}

func (h *logContextHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h.Handler.WithAttrs(attrs)
}

func (h *logContextHandler) WithGroup(name string) slog.Handler {
	return h.Handler.WithGroup(name)
}

func logInfo(r *http.Request, message string, args ...any) {
	l := palikka.LoggerFromContext(r.Context())
	l.InfoContext(r.Context(), message, args...)
}

func logError(r *http.Request, message string, args ...any) {
	l := palikka.LoggerFromContext(r.Context())
	l.ErrorContext(r.Context(), message, args...)
}

func (s *Server) respondError(w http.ResponseWriter, r *http.Request, err error) {
	code, message := palikka.ErrorCode(err), palikka.ErrorMessage(err)

	// log the full error
	if code == palikka.ErrInternal {
		logError(r, "[http]", "error", err)
	} else {
		logInfo(r, "[http]", "error", err)
	}

	// todo test
	statusCode := errorStatusCode(code)
	if statusCode == http.StatusUnauthorized {
		http.SetCookie(w, s.CookieService.NewExpired())
	}

	w.Header().Set(HeaderContentType, MediaTypeJSON)
	w.WriteHeader(statusCode)

	body := &errorResponse{Error: message}
	if encodeError := json.NewEncoder(w).Encode(body); encodeError != nil {
		logError(r, "[http] failed to write error response", "error", encodeError)
	}
}

// **********
// error handling
// **********

type errorResponse struct {
	Error string `json:"error"`
}

var codes = map[string]int{
	palikka.ErrInvalid:      http.StatusBadRequest,
	palikka.ErrUnauthorized: http.StatusUnauthorized,
	palikka.ErrForbidden:    http.StatusForbidden,
	palikka.ErrNotFound:     http.StatusNotFound,
	palikka.ErrInternal:     http.StatusInternalServerError,
}

func errorStatusCode(code string) int {
	if v, ok := codes[code]; ok {
		return v
	}
	return http.StatusInternalServerError
}
