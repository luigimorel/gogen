package internal

type RouterGenerator struct{}

func NewRouterGenerator() *RouterGenerator {
	return &RouterGenerator{}
}

func (rg *RouterGenerator) generateStdlibContent() string {
	return `package web

import (
	"fmt"
	"net/http"
)

func SetupRoutes() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/health", healthHandler)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello from web server")
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "OK")
}`
}

func (rg *RouterGenerator) generateChiContent() string {
	return `package web

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func SetupRoutes() *chi.Mux {
	r := chi.NewRouter()
	
	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Routes
	r.Get("/", homeHandler)
	r.Get("/health", healthHandler)

	return r
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello from web server")
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "OK")
}`
}

func (rg *RouterGenerator) generateGorillaContent() string {
	return `package web

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func SetupRoutes() *mux.Router {
	r := mux.NewRouter()

	// Routes
	r.HandleFunc("/", homeHandler).Methods("GET")
	r.HandleFunc("/health", healthHandler).Methods("GET")

	return r
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello from web server")
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "OK")
}`
}

func (rg *RouterGenerator) generateHttpRouterContent() string {
	return `package web

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func SetupRoutes() *httprouter.Router {
	router := httprouter.New()

	// Routes
	router.GET("/", homeHandler)
	router.GET("/health", healthHandler)

	return router
}

func homeHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprintf(w, "Hello from web server")
}

// healthHandler handles the health check endpoint
func healthHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "OK")
}`
}
