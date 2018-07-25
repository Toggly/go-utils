package utils

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/middleware"
	uuid "github.com/satori/go.uuid"
)

type contextKey int

const (
	// XRequestID gets name of header
	XRequestID = "X-Toggly-Request-Id"
	// XServiceName gets service name in header
	XServiceName = "X-Toggly-Service-Name"
	// XServiceVersion gets version
	XServiceVersion = "X-Toggly-Service-Version"
	// ContextReqIDKey gets key for request id
	ContextReqIDKey contextKey = iota
	// ContextVersionKey gets key for version
	ContextVersionKey contextKey = iota
)

// VersionCtx adds api version to context
func VersionCtx(version string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r = r.WithContext(context.WithValue(r.Context(), ContextVersionKey, version))
			next.ServeHTTP(w, r)
		})
	}
}

// RequestIDCtx adds request id to context
func RequestIDCtx(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		rid := r.Header.Get(http.CanonicalHeaderKey(XRequestID))
		if rid == "" {
			rid = uuid.NewV4().String()
		}
		ctx := r.Context()
		ctx = context.WithValue(ctx, ContextReqIDKey, rid)
		ctx = context.WithValue(ctx, middleware.RequestIDKey, rid)

		// set header to response
		w.Header().Set(XRequestID, rid)
		next.ServeHTTP(w, r.WithContext(ctx))

	}
	return http.HandlerFunc(fn)
}

// ServiceInfo adds service information to the response header
func ServiceInfo(name string, version string) func(http.Handler) http.Handler {
	f := func(h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set(http.CanonicalHeaderKey(XServiceName), name)
			w.Header().Set(http.CanonicalHeaderKey(XServiceVersion), version)
			h.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
	return f
}
