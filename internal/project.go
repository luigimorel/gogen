package internal

import (
	"fmt"
	"os"
	"os/exec"
)

type ProjectGenerator struct{}

func NewProjectGenerator() *ProjectGenerator {
	return &ProjectGenerator{}
}

func (pg *ProjectGenerator) CreateCLIProject(projectName, moduleName string) error {
	mainContent := fmt.Sprintf(`package main

import (
    "fmt"
    "log"
    "os"

    "github.com/urfave/cli/v2"
)

func main() {
    app := &cli.App{
        Name:  "%s",
        Usage: "A CLI application built with gogen",
        Action: func(c *cli.Context) error {
            return cli.ShowAppHelp(c)
        },
        Commands: []*cli.Command{
            {
                Name:    "greet",
                Aliases: []string{"g"},
                Usage:   "Greet someone",
                Flags: []cli.Flag{
                    &cli.StringFlag{
                        Name:  "name",
                        Value: "World",
                        Usage: "Name to greet",
                    },
                },
                Action: func(c *cli.Context) error {
                    name := c.String("name")
                    fmt.Printf("Hello %%s\n", name)
                    return nil
                },
            },
        },
    }

    if err := app.Run(os.Args); err != nil {
        log.Fatal(err)
    }
}
`, projectName)

	if err := os.WriteFile("main.go", []byte(mainContent), 0644); err != nil {
		return err
	}

	if err := pg.InitGitRepository(projectName, "cli"); err != nil {
		fmt.Printf("Warning: failed to initialize git repository: %v\n", err)
	}

	cmd := exec.Command("go", "mod", "tidy")
	return cmd.Run()
}

func (pg *ProjectGenerator) setModuleName(moduleName, projectName string) string {
	if moduleName == "" {
		return "github.com/" + projectName
	}
	return moduleName
}

func (pg *ProjectGenerator) setDefaultPackages() string {
	return `"fmt"
	"log"
	"net/http"`
}

func (pg *ProjectGenerator) CreateWebProject(projectName, moduleName, router, frontend string, useTypeScript bool) error {
	if err := os.MkdirAll("api/cmd/web", 0755); err != nil {
		return fmt.Errorf("failed to create api/cmd/web directory: %w", err)
	}

	var mainContent string
	var routesContent string

	switch router {
	case "chi":
		mainContent = `package main

import (
	` + pg.setDefaultPackages() + `

	"` + pg.setModuleName(moduleName, projectName) + `/api/cmd/web"
)

func main() {
	r := web.SetupRoutes()

	port := ":8080"
	fmt.Printf("Starting ` + projectName + ` web server on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, r))
}
`
		routesContent = `package web

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// SetupRoutes configures and returns the router with all routes
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

// homeHandler handles the home page
func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello from ` + projectName + ` web server")
}

// healthHandler handles the health check endpoint
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "OK")
}
`

	case "gorilla":
		mainContent = `package main

import (
	` + pg.setDefaultPackages() + `

	"` + pg.setModuleName(moduleName, projectName) + `/api/cmd/web"
)

func main() {
	r := web.SetupRoutes()

	port := ":8080"
	fmt.Printf("Starting ` + projectName + ` web server on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, r))
}
`
		routesContent = `package web

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// SetupRoutes configures and returns the router with all routes
func SetupRoutes() *mux.Router {
	r := mux.NewRouter()

	// Routes
	r.HandleFunc("/", homeHandler).Methods("GET")
	r.HandleFunc("/health", healthHandler).Methods("GET")

	return r
}

// homeHandler handles the home page
func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello from ` + projectName + ` web server")
}

// healthHandler handles the health check endpoint
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "OK")
}
`

	case "httprouter":
		mainContent = `package main

import (
	` + pg.setDefaultPackages() + `

	"` + pg.setModuleName(moduleName, projectName) + `/api/cmd/web"
)

func main() {
	router := web.SetupRoutes()

	port := ":8080"
	fmt.Printf("Starting ` + projectName + ` web server on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, router))
}
`
		routesContent = `package web

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// SetupRoutes configures and returns the router with all routes
func SetupRoutes() *httprouter.Router {
	router := httprouter.New()

	// Routes
	router.GET("/", homeHandler)
	router.GET("/health", healthHandler)

	return router
}

// homeHandler handles the home page
func homeHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprintf(w, "Hello from ` + projectName + ` web server")
}

// healthHandler handles the health check endpoint
func healthHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "OK")
}
`

	default: // stdlib
		mainContent = `package main

import (
	` + pg.setDefaultPackages() + `

	"` + pg.setModuleName(moduleName, projectName) + `/api/cmd/web"
)

func main() {
	web.SetupRoutes()

	port := ":8080"
	fmt.Printf("Starting ` + projectName + ` web server on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
`
		routesContent = `package web

import (
	"fmt"
	"net/http"
)

// SetupRoutes configures all HTTP routes
func SetupRoutes() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/health", healthHandler)
}

// homeHandler handles the home page
func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello from ` + projectName + ` web server")
}

// healthHandler handles the health check endpoint
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "OK")
}
`
	}

	if err := os.WriteFile("api/main.go", []byte(mainContent), 0644); err != nil {
		return err
	}

	if err := os.WriteFile("api/cmd/web/routes.go", []byte(routesContent), 0644); err != nil {
		return err
	}

	baseModuleName := pg.setModuleName(moduleName, projectName)
	apiModContent := fmt.Sprintf("module %s/api\n\ngo 1.21\n", baseModuleName)
	if err := os.WriteFile("api/go.mod", []byte(apiModContent), 0644); err != nil {
		return err
	}

	if err := pg.InitGitRepository(projectName, "web"); err != nil {
		fmt.Printf("Warning: failed to initialize git repository: %v\n", err)
	}

	originalDir, _ := os.Getwd()
	if err := os.Chdir("api"); err != nil {
		return fmt.Errorf("failed to change to api directory: %w", err)
	}

	cmd := exec.Command("go", "mod", "tidy")
	if err := cmd.Run(); err != nil {
		os.Chdir(originalDir)
		return fmt.Errorf("failed to tidy go.mod file: %w", err)
	}

	if err := os.Chdir(originalDir); err != nil {
		return fmt.Errorf("failed to change back to original directory: %w", err)
	}

	if frontend != "" {
		if err := pg.CreateFrontendProject(frontend, "frontend", useTypeScript); err != nil {
			return fmt.Errorf("failed to create frontend project: %w", err)
		}

		if err := pg.CreateEnvFile("frontend"); err != nil {
			fmt.Printf("Warning: failed to create env file: %v\n", err)
		}

		if err := pg.CreateGitignoreFile("frontend", "frontend"); err != nil {
			fmt.Printf("Warning: failed to create .gitignore file in frontend: %v\n", err)
		}
	}

	if err := pg.CreateEnvFile("api"); err != nil {
		fmt.Printf("Warning: failed to create env file in api: %v\n", err)
	}

	if err := pg.CreateGitignoreFile("api", "api"); err != nil {
		fmt.Printf("Warning: failed to create .gitignore file in api: %v\n", err)
	}

	return nil
}

// CreateAPIProject creates an API-only project with the specified parameters
func (pg *ProjectGenerator) CreateAPIProject(projectName, moduleName, router string) error {
	if err := os.MkdirAll("cmd/api", 0755); err != nil {
		return fmt.Errorf("failed to create cmd/api directory: %w", err)
	}

	var mainContent string
	var routesContent string

	responseStruct := `type Response struct {
	Message string ` + "`" + `json:"message"` + "`" + `
	Service string ` + "`" + `json:"service"` + "`" + `
}`

	switch router {
	case "chi":
		mainContent = `package main

import (
	` + pg.setDefaultPackages() + `

	"` + pg.setModuleName(moduleName, projectName) + `/cmd/api"
)

func main() {
	r := api.SetupRoutes()

	port := ":8080"
	fmt.Printf("Starting ` + projectName + ` API server on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, r))
}
`
		routesContent = fmt.Sprintf(`package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

%s

// SetupRoutes configures and returns the router with all routes
func SetupRoutes() *chi.Mux {
	r := chi.NewRouter()
	
	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.SetHeader("Content-Type", "application/json"))

	// Routes
	r.Get("/api/hello", helloHandler)
	r.Get("/api/health", healthHandler)

	return r
}

// helloHandler handles the hello endpoint
func helloHandler(w http.ResponseWriter, r *http.Request) {
	response := Response{
		Message: "Hello from %s API",
		Service: "%s",
	}
	json.NewEncoder(w).Encode(response)
}

// healthHandler handles the health check endpoint
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	response := Response{
		Message: "API is healthy",
		Service: "%s",
	}
	json.NewEncoder(w).Encode(response)
}
`, responseStruct, projectName, projectName, projectName)

	case "gorilla":
		mainContent = `package main

import (
	` + pg.setDefaultPackages() + `

	"` + pg.setModuleName(moduleName, projectName) + `/cmd/api"
)

func main() {
	r := api.SetupRoutes()

	port := ":8080"
	fmt.Printf("Starting ` + projectName + ` API server on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, r))
}
`
		routesContent = fmt.Sprintf(`package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

%s

// SetupRoutes configures and returns the router with all routes
func SetupRoutes() *mux.Router {
	r := mux.NewRouter()

	// Routes
	r.HandleFunc("/api/hello", helloHandler).Methods("GET")
	r.HandleFunc("/api/health", healthHandler).Methods("GET")

	return r
}

// helloHandler handles the hello endpoint
func helloHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := Response{
		Message: "Hello from %s API",
		Service: "%s",
	}
	json.NewEncoder(w).Encode(response)
}

// healthHandler handles the health check endpoint
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := Response{
		Message: "API is healthy",
		Service: "%s",
	}
	json.NewEncoder(w).Encode(response)
}
`, responseStruct, projectName, projectName, projectName)

	case "httprouter":
		mainContent = `package main

import (
	` + pg.setDefaultPackages() + `

	"github.com/` + moduleName + `/cmd/api"
)

func main() {
	router := api.SetupRoutes()

	port := ":8080"
	fmt.Printf("Starting ` + projectName + ` API server on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, router))
}
`
		routesContent = fmt.Sprintf(`package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

%s

// SetupRoutes configures and returns the router with all routes
func SetupRoutes() *httprouter.Router {
	router := httprouter.New()

	// Routes
	router.GET("/api/hello", helloHandler)
	router.GET("/api/health", healthHandler)

	return router
}

// helloHandler handles the hello endpoint
func helloHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	response := Response{
		Message: "Hello from %s API",
		Service: "%s",
	}
	json.NewEncoder(w).Encode(response)
}

// healthHandler handles the health check endpoint
func healthHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := Response{
		Message: "API is healthy",
		Service: "%s",
	}
	json.NewEncoder(w).Encode(response)
}
`, responseStruct, projectName, projectName, projectName)

	default: // stdlib
		mainContent = `package main

import (
	` + pg.setDefaultPackages() + `

	"` + pg.setModuleName(moduleName, projectName) + `/cmd/api"
)

func main() {
	api.SetupRoutes()

	port := ":8080"
	fmt.Printf("Starting ` + projectName + ` API server on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
`
		routesContent = fmt.Sprintf(`package api

import (
	"encoding/json"
	"net/http"
)

%s

// SetupRoutes configures all HTTP routes
func SetupRoutes() {
	http.HandleFunc("/api/hello", helloHandler)
	http.HandleFunc("/api/health", healthHandler)
}

// helloHandler handles the hello endpoint
func helloHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := Response{
		Message: "Hello from %s API",
		Service: "%s",
	}
	json.NewEncoder(w).Encode(response)
}

// healthHandler handles the health check endpoint
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := Response{
		Message: "API is healthy",
		Service: "%s",
	}
	json.NewEncoder(w).Encode(response)
}
`, responseStruct, projectName, projectName, projectName)
	}

	if err := os.WriteFile("main.go", []byte(mainContent), 0644); err != nil {
		return err
	}

	if err := os.WriteFile("cmd/api/routes.go", []byte(routesContent), 0644); err != nil {
		return err
	}

	if err := pg.CreateEnvFile("."); err != nil {
		fmt.Printf("Warning: failed to create env file: %v\n", err)
	}

	if err := pg.InitGitRepository(projectName, "api"); err != nil {
		fmt.Printf("Warning: failed to initialize git repository: %v\n", err)
	}

	cmd := exec.Command("go", "mod", "tidy")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to tidy go.mod: %w", err)
	}

	return nil
}
