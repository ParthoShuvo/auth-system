package route

import (
	"net/http"
	"time"

	log "github.com/parthoshuvo/authsvc/log4u"
	"github.com/parthoshuvo/authsvc/resource"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

// Builder holds all routes.
type Builder struct {
	allowCors  bool
	pr         resource.Protector
	router     *mux.Router
	serverName string
	isLogDebug bool
}

// NewRouteBuilder creates a route builder.
func NewRouteBuilder(allowCors bool, pr resource.Protector, serverName string, isLogDebug bool) *Builder {
	return &Builder{allowCors, pr, mux.NewRouter().StrictSlash(true), serverName, isLogDebug}
}

// SubrouteBuilder creates a subroute builder.
func (rb *Builder) SubrouteBuilder(pathPrefix string) *Builder {
	return rb.partialClone(rb.router.PathPrefix(pathPrefix).Subrouter())
}

// AddSafe adds a protected route.
func (rb *Builder) AddSafe(action resource.Action, method, path string, handlerFunc http.HandlerFunc) *mux.Route {
	handler := rb.generalHandler(rb.corsHandler(rb.pr.Protect(action, rb.performanceLogger(handlerFunc, action))))
	return rb.add(action, method, path, handler)
}

// Add a route.
func (rb *Builder) Add(action resource.Action, method, path string, handlerFunc http.HandlerFunc) *mux.Route {
	handler := rb.generalHandler(rb.corsHandler(rb.performanceLogger(handlerFunc, action)))
	return rb.add(action, method, path, handler)
}

// add a route.
func (rb *Builder) add(action resource.Action, method, path string, handler http.Handler) *mux.Route {
	return rb.router.Methods(method).Path(path).Name(action.String()).Handler(handler)
}

// Router fetch the configured router.
func (rb *Builder) Router() *mux.Router {
	return rb.router
}

func (rb *Builder) corsHandler(handler http.Handler) http.Handler {
	if rb.allowCors {
		return cors.New(cors.Options{
			AllowedHeaders:   []string{"authorization"},
			AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
			AllowCredentials: true}).Handler(handler)
	}
	return handler
}

func (rb *Builder) generalHandler(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Server", rb.serverName)
		inner.ServeHTTP(w, r)
	})
}

func (rb *Builder) performanceLogger(inner http.HandlerFunc, action resource.Action) http.Handler {
	if rb.isLogDebug {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			inner.ServeHTTP(w, r)
			log.Debugf("%s %s %s %s", r.Method, r.RequestURI, action, time.Since(start))
		})
	}
	return inner
}

func (rb *Builder) partialClone(router *mux.Router) *Builder {
	return &Builder{rb.allowCors, rb.pr, router, rb.serverName, rb.isLogDebug}
}
