package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"palikka-go/internal/logs"
	"palikka-go/internal/session"
)

// Middleware is a function type that wraps a [http.Handler].
// An implementation may set up pre and/or postprocessing before
// handing the request off to the original handler.
type Middleware func(http.Handler) http.Handler

// NewMiddlewareStack creates a new stack of [Middleware]
// where each is called in the same order they
// were passed in.
func NewMiddlewareStack(mws ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(mws) - 1; i >= 0; i-- {
			mw := mws[i]
			next = mw(next)
		}
		return next
	}
}

type LoggerMiddleware interface {
	LogRequests(next http.Handler) http.Handler
}

type loggerMiddleware struct {
}

func NewLoggerMiddleware() LoggerMiddleware {
	return &loggerMiddleware{}
}

func (mw *loggerMiddleware) LogRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		r = logs.AddRequestAttributes(r)

		// read headers
		var headers string
		for k, v := range r.Header {
			for _, value := range v {
				headers = fmt.Sprintf("%s%s%s\n", headers, k, value)
			}
		}

		slog.InfoContext(r.Context(), "--> in", slog.Any("headers", headers))

		next.ServeHTTP(w, r)

		slog.InfoContext(r.Context(), "<-- out")
	})
}

type SessionMiddleware interface {
	EnsureAuth(next http.Handler) http.Handler
}

type sessionMiddleware struct {
	sessionStore session.Store
}

func NewSessionMiddleware(sessionStore session.Store) SessionMiddleware {
	return &sessionMiddleware{
		sessionStore: sessionStore,
	}
}

func (mw *sessionMiddleware) EnsureAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := mw.sessionStore.Load(r)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// add session attributes to context
		ctx := logs.ContextKeySessionID().AddToContext(r.Context(), session.ID)
		r = r.WithContext(ctx)

		slog.InfoContext(r.Context(), "request authenticated successfully")

		next.ServeHTTP(w, r)
	})
}
