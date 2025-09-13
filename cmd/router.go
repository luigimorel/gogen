package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/urfave/cli/v2"
)

// Router type constants
const (
	RouterStdlib     = "stdlib"
	RouterChi        = "chi"
	RouterGorilla    = "gorilla"
	RouterHttpRouter = "httprouter"
)

type Router struct {
	Type       string
	UpdateMain bool
}

func NewRouter(routerType string, updateMain bool) *Router {
	return &Router{
		Type:       routerType,
		UpdateMain: updateMain,
	}
}

func RouterCommand() *cli.Command {
	return &cli.Command{
		Name:      "router",
		Usage:     "Add a router to your Go project",
		ArgsUsage: "<router-type>",
		Description: `Add a router to your existing Go project.
This command will add the selected router dependency and update your main.go file with the router setup.

Supported routers:
- chi: Chi lightweight router
- gorilla: Gorilla Mux router
- stdlib: Plain Go standard library
- httprouter: HttpRouter high performance router

Usage:
  gogen router chi
  gogen router gorilla
  gogen router httprouter
  gogen router stdlib`,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "update",
				Aliases: []string{"u"},
				Usage:   "Update main.go with router implementation",
				Value:   true,
			},
		},
		Action: func(c *cli.Context) error {
			if c.NArg() == 0 {
				return fmt.Errorf("router type is required. Usage: gogen router <router-type>")
			}

			routerType := c.Args().Get(0)
			updateMain := c.Bool("update")

			router := NewRouter(routerType, updateMain)
			return router.execute()
		},
	}
}

func (r *Router) execute() error {
	if err := r.validateProject(); err != nil {
		return err
	}

	if err := r.installDependency(); err != nil {
		return fmt.Errorf("failed to install router dependency: %w", err)
	}

	if r.UpdateMain {
		if err := r.updateMainFile(); err != nil {
			return fmt.Errorf("failed to update main.go: %w", err)
		}
	}

	r.printInstructions()

	return nil
}

func (r *Router) validateProject() error {
	if _, err := os.Stat("go.mod"); err != nil {
		return fmt.Errorf("no go.mod found - please run this command in a Go project directory")
	}
	return nil
}

func (r *Router) installDependency() error {
	var dependency string

	switch r.Type {
	case RouterStdlib:
		fmt.Println("Using standard library http.ServeMux - no additional dependency needed")
		return nil
	case RouterChi:
		dependency = "github.com/go-chi/chi/v5"
	case RouterGorilla:
		dependency = "github.com/gorilla/mux"
	case RouterHttpRouter:
		dependency = "github.com/julienschmidt/httprouter"
	default:
		return fmt.Errorf("unsupported router type: %s", r.Type)
	}

	if dependency != "" {
		fmt.Printf("Installing %s...\n", dependency)
		cmd := exec.Command("go", "get", dependency)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to install dependency %s: %w", dependency, err)
		}
	}

	return nil
}

func (r *Router) updateMainFile() error {
	mainContent, err := os.ReadFile("main.go")
	if err != nil {
		return fmt.Errorf("failed to read main.go: %w", err)
	}

	newContent := r.generateMainContent(string(mainContent))

	backupPath := "main.go.backup"
	if err := os.WriteFile(backupPath, mainContent, 0600); err != nil {
		fmt.Printf("Warning: failed to create backup at %s: %v\n", backupPath, err)
	}
	if err := os.WriteFile("main.go", []byte(newContent), 0600); err != nil {
		return fmt.Errorf("failed to write updated main.go: %w", err)
	}

	return nil
}

func (r *Router) generateMainContent(existingContent string) string {
	switch r.Type {
	case RouterStdlib:
		return r.generateServeMuxContent()
	case RouterChi:
		return r.generateChiContent()
	case RouterGorilla:
		return r.generateGorillaContent()
	case RouterHttpRouter:
		return r.generateHttpRouterContent()
	default:
		return existingContent
	}
}

func (r *Router) printInstructions() {
	fmt.Println("\nNext steps:")
	fmt.Println("   go run main.go")
	fmt.Println("   curl http://localhost:8080/api/hello")
	fmt.Println("   curl http://localhost:8080/api/health")

	switch r.Type {
	case RouterChi:
		fmt.Println("\nChi router features:")
		fmt.Println("   - Built-in middleware (Logger, Recoverer, RequestID)")
		fmt.Println("   - Route groups and subrouting")
		fmt.Println("   - Fast and lightweight")
	case RouterGorilla:
		fmt.Println("\nGorilla Mux features:")
		fmt.Println("   - Path variables: r.HandleFunc(\"/users/{id}\", handler)")
		fmt.Println("   - Query parameter matching")
		fmt.Println("   - Host and scheme matching")
	case RouterHttpRouter:
		fmt.Println("\nHttpRouter features:")
		fmt.Println("   - Extremely fast performance")
		fmt.Println("   - Path parameters: router.GET(\"/users/:id\", handler)")
		fmt.Println("   - Zero memory allocation")
	case RouterStdlib:
		fmt.Println("\nhttp.ServeMux features:")
		fmt.Println("   - Part of Go standard library")
		fmt.Println("   - Simple and reliable")
		fmt.Println("   - Pattern matching with wildcards")
	}
}

func (r *Router) generateServeMuxContent() string {
	return `package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Response struct {
	Message string ` + "`json:\"message\"`" + `
	Router  string ` + "`json:\"router\"`" + `
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		response := Response{
			Message: "Hello from http.ServeMux",
			Router:  "http.ServeMux",
		}
		json.NewEncoder(w).Encode(response)
	})

	mux.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := Response{
			Message: "API is healthy",
			Router:  "http.ServeMux",
		}
		json.NewEncoder(w).Encode(response)
	})

	port := ":8080"
	fmt.Printf("Starting API server with http.ServeMux on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, mux))
}
`
}

func (r *Router) generateChiContent() string {
	return `package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Response struct {
	Message string ` + "`json:\"message\"`" + `
	Router  string ` + "`json:\"router\"`" + `
}

func main() {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)

	r.Route("/api", func(r chi.Router) {
		r.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			response := Response{
				Message: "Hello from Chi router",
				Router:  "Chi",
			}
			json.NewEncoder(w).Encode(response)
		})

		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := Response{
				Message: "API is healthy",
				Router:  "Chi",
			}
			json.NewEncoder(w).Encode(response)
		})
	})

	port := ":8080"
	fmt.Printf("Starting API server with Chi router on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, r))
}
`
}

func (r *Router) generateGorillaContent() string {
	return `package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Response struct {
	Message string ` + "`json:\"message\"`" + `
	Router  string ` + "`json:\"router\"`" + `
}

func main() {
	r := mux.NewRouter()

	// API routes
	api := r.PathPrefix("/api").Subrouter()
	
	api.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		response := Response{
			Message: "Hello from Gorilla Mux",
			Router:  "Gorilla Mux",
		}
		json.NewEncoder(w).Encode(response)
	}).Methods("GET")

	api.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := Response{
			Message: "API is healthy",
			Router:  "Gorilla Mux",
		}
		json.NewEncoder(w).Encode(response)
	}).Methods("GET")

	port := ":8080"
	fmt.Printf("Starting API server with Gorilla Mux on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, r))
}
`
}

func (r *Router) generateHttpRouterContent() string {
	return `package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Response struct {
	Message string ` + "`json:\"message\"`" + `
	Router  string ` + "`json:\"router\"`" + `
}

func main() {
	router := httprouter.New()

	router.GET("/api/hello", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.Header().Set("Content-Type", "application/json")
		response := Response{
			Message: "Hello from HttpRouter",
			Router:  "HttpRouter",
		}
		json.NewEncoder(w).Encode(response)
	})

	router.GET("/api/health", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := Response{
			Message: "API is healthy",
			Router:  "HttpRouter",
		}
		json.NewEncoder(w).Encode(response)
	})

	port := ":8080"
	fmt.Printf("Starting API server with HttpRouter on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, router))
}
`
}
