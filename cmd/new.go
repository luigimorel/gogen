package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/urfave/cli/v2"
)

// NewCommand creates the new project command for the CLI
func NewCommand() *cli.Command {
	return &cli.Command{
		Name:  "new",
		Usage: "Create a new Go project",
		Description: `Create a new Go project with proper structure and initialization.
This command will create a new directory, initialize a Go module, and create a new api project`,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "name",
				Aliases:  []string{"n"},
				Usage:    "Project name",
				Required: true,
			},
			&cli.StringFlag{
				Name:    "module",
				Aliases: []string{"m"},
				Usage:   "Go module path (default: project name)",
			},
			&cli.StringFlag{
				Name:    "template",
				Aliases: []string{"t"},
				Usage:   "Project template (cli, web, api)",
				Value:   "api",
			},
			&cli.StringFlag{
				Name:    "router",
				Aliases: []string{"r"},
				Usage:   "Router type for API/web projects (stdlib, chi, gorilla, httprouter)",
				Value:   "stdlib",
			},
			&cli.StringFlag{
				Name:    "frontend",
				Aliases: []string{"fe"},
				Usage:   "Frontend framework for web projects (react, vue, svelte, solidjs, angular)",
			},
			&cli.StringFlag{
				Name:  "dir",
				Usage: "Directory name for the project (default: project name)",
			},
			&cli.BoolFlag{
				Name:  "ts",
				Usage: "Use TypeScript for frontend projects (only applicable with --frontend)",
				Value: false,
			},
		},
		Action: func(c *cli.Context) error {
			projectName := c.String("name")
			moduleName := c.String("module")
			template := c.String("template")
			router := c.String("router")
			frontend := c.String("frontend")
			projectDir := c.String("dir")
			useTypeScript := c.Bool("ts")

			if frontend != "" && template != "web" {
				return fmt.Errorf("frontend flag is only applicable when template is 'web'")
			}

			if useTypeScript && frontend == "" {
				return fmt.Errorf("TypeScript flag is only applicable when frontend is specified")
			}

			if moduleName == "" {
				moduleName = projectName
			}

			if projectDir == "" {
				projectDir = projectName
			}

			fmt.Printf("Creating new Go project '%s'...\n", projectName)

			if err := os.Mkdir(projectDir, 0755); err != nil {
				return fmt.Errorf("failed to create project directory: %w", err)
			}

			originalDir, _ := os.Getwd()
			if err := os.Chdir(projectDir); err != nil {
				return fmt.Errorf("failed to change to project directory: %w", err)
			}
			defer func() {
				if err := os.Chdir(originalDir); err != nil {
					fmt.Fprintf(os.Stderr, "Warning: failed to change back to original directory: %v\n", err)
				}
			}()

			fmt.Printf("Initializing Go module '%s'...\n", moduleName)

			if template != "web" {
				cmd := exec.Command("go", "mod", "init", "github.com/"+moduleName)
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				if err := cmd.Run(); err != nil {
					return fmt.Errorf("failed to initialize go module: %w", err)
				}
			}

			if err := createProjectFiles(projectName, moduleName, template, router, frontend, useTypeScript); err != nil {
				return fmt.Errorf("failed to create project files: %w", err)
			}

			fmt.Printf("Project '%s' created successfully\n", projectName)
			fmt.Printf("Location: %s\n", filepath.Join(originalDir, projectDir))
			fmt.Printf("Template: %s\n", template)
			fmt.Printf("Router: %s\n", router)
			if frontend != "" {
				fmt.Printf("Frontend: %s\n", frontend)
				if useTypeScript {
					fmt.Println("TypeScript: enabled")
				}
			}
			fmt.Println("\nNext steps:")
			fmt.Printf("   cd %s\n", projectDir)
			if template == "web" {
				fmt.Println("   cd api")
				fmt.Println("   go mod tidy")
				fmt.Println("   go run main.go")
				if frontend != "" {
					fmt.Println("\n   # In another terminal:")
					fmt.Println("   cd frontend")
					fmt.Println("   npm install")
					fmt.Println("   npm run dev")
				}
			} else {
				fmt.Println("   go mod tidy")
				fmt.Println("   go run main.go")
			}

			return nil
		},
	}
}

func createProjectFiles(projectName, moduleName, template, router, frontend string, useTypeScript bool) error {
	switch template {
	case "cli":
		return createCLIProject(projectName, moduleName)
	case "web":
		return createWebProject(projectName, moduleName, router, frontend, useTypeScript)
	case "api":
		return createAPIProject(projectName, moduleName, router)
	default:
		return fmt.Errorf("unsupported template: %s", template)
	}
}

func createCLIProject(projectName, moduleName string) error {
	mainContent := fmt.Sprintf(`package main

import (
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

	cmd := exec.Command("go", "get", "github.com/urfave/cli/v2")
	return cmd.Run()
}

func createWebProject(projectName, moduleName, router, frontend string, useTypeScript bool) error {
	if err := os.MkdirAll("api/cmd/web", 0755); err != nil {
		return fmt.Errorf("failed to create api/cmd/web directory: %w", err)
	}

	var mainContent string
	var routesContent string

	switch router {
	case "chi":
		mainContent = `package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/` + moduleName + `/api/cmd/web"
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
	"fmt"
	"log"
	"net/http"

	"github.com/` + moduleName + `/api/cmd/web"
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
	"fmt"
	"log"
	"net/http"

	"github.com/` + moduleName + `/api/cmd/web"
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
	"fmt"
	"log"
	"net/http"

	"github.com/` + moduleName + `/api/cmd/web"
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

	apiModContent := fmt.Sprintf("module github.com/%s/api\n\ngo 1.21\n", moduleName)
	if err := os.WriteFile("api/go.mod", []byte(apiModContent), 0644); err != nil {
		return err
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
		if err := createFrontendProject(frontend, "frontend", useTypeScript); err != nil {
			return fmt.Errorf("failed to create frontend project: %w", err)
		}
	}

	return nil
}

func createAPIProject(projectName, moduleName, router string) error {
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
	"fmt"
	"log"
	"net/http"

	"github.com/` + moduleName + `/cmd/api"
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
	"fmt"
	"log"
	"net/http"

	"github.com/` + moduleName + `/cmd/api"
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
	"fmt"
	"log"
	"net/http"

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
	"fmt"
	"log"
	"net/http"

	"github.com/` + moduleName + `/cmd/api"
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

	cmd := exec.Command("go", "mod", "tidy")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to tidy go.mod: %w", err)
	}

	return nil
}
