package router

import (
	"context"
	"github.com/go-chi/chi"
	"net/http"
)

// ======================================
// Router wrapper
// ======================================
type Router struct {
	chi chi.Router
}

// ======================================
// Create new Router
// ======================================
func NewRouter() *Router {
	return &Router{chi.NewRouter()}
}

// ======================================
// Default methods
// ======================================
func (r *Router) Route(pattern string, fn func(*Router)) {
	r.chi.Route(pattern, func(c chi.Router) {
		fn(&Router{c})
	})
}

func (r *Router) Get(pattern string, fn apiHandler) {
	r.chi.Get(pattern, handler(fn))
}
func (r *Router) Post(pattern string, fn apiHandler) {
	r.chi.Post(pattern, handler(fn))
}
func (r *Router) Put(pattern string, fn apiHandler) {
	r.chi.Put(pattern, handler(fn))
}
func (r *Router) Delete(pattern string, fn apiHandler) {
	r.chi.Delete(pattern, handler(fn))
}

// ======================================
// Use custom middleware
// ======================================
func (r *Router) With(fn middlewareHandler) *Router {
	c := r.chi.With(middleware(fn))
	return &Router{c}
}

func (r *Router) WithBypass(fn func(next http.Handler) http.Handler) *Router {
	c := r.chi.With(fn)
	return &Router{c}
}

func (r *Router) Use(fn middlewareHandler) {
	r.chi.Use(middleware(fn))
}
func (r *Router) UseBypass(fn func(next http.Handler) http.Handler) {
	r.chi.Use(fn)
}

// ======================================
// Serve
// ======================================
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.chi.ServeHTTP(w, req)
}

// ======================================
// Custom error handler
// ======================================
type apiHandler func(w http.ResponseWriter, r *http.Request) error

func handler(fn apiHandler) http.HandlerFunc {
	return fn.serve
}

func (h apiHandler) serve(w http.ResponseWriter, r *http.Request) {
	if err := h(w, r); err != nil {
		handleError(err, w, r)
	}
}

// ======================================
// Custom middleware handler
// ======================================
type middlewareHandler func(w http.ResponseWriter, r *http.Request) (context.Context, error)

func (m middlewareHandler) handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.serve(next, w, r)
	})
}

func (m middlewareHandler) serve(next http.Handler, w http.ResponseWriter, r *http.Request) {
	ctx, err := m(w, r)
	if err != nil {
		handleError(err, w, r)
		return
	}
	if ctx != nil {
		r = r.WithContext(ctx)
	}
	next.ServeHTTP(w, r)
}

func middleware(fn middlewareHandler) func(http.Handler) http.Handler {
	return fn.handler
}
