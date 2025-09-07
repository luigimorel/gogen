package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/urfave/cli/v2"
)

// RouterCommand creates the router command for the CLI
func RouterCommand() *cli.Command {
	return &cli.Command{
		Name:  "router",
		Usage: "Add a router to your Go project",
		Description: `Add a router to your existing Go project.
This command will add the selected router dependency and update your main.go file with the router setup.

Supported routers:
- chi: Chi lightweight router
- gorilla: Gorilla Mux router
- stdlib: Plain Go standard library,
- httprouter: HttpRouter high performance router`,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "type",
				Aliases: []string{"t"},
				Usage:   "Router type (stdlib, chi, gorilla, httprouter)",
				Value:   "stdlib",
			},
			&cli.BoolFlag{
				Name:    "update",
				Aliases: []string{"u"},
				Usage:   "Update main.go with router implementation",
				Value:   true,
			},
		},
		Action: func(c *cli.Context) error {
			routerType := c.String("type")
			updateMain := c.Bool("update")

			fmt.Printf("Adding %s router to your project...\n", getRouterDisplayName(routerType))

			if err := validateProject(); err != nil {
				return err
			}

			if err := installRouterDependency(routerType); err != nil {
				return fmt.Errorf("failed to install router dependency: %w", err)
			}

			if updateMain {
				if err := updateMainWithRouter(routerType); err != nil {
					return fmt.Errorf("failed to update main.go: %w", err)
				}
			}

			fmt.Printf("Successfully added %s router to your project!\n", getRouterDisplayName(routerType))
			printRouterInstructions(routerType)

			return nil
		},
	}
}

func getRouterDisplayName(routerType string) string {
	switch routerType {
	case "stdlib":
		return "http.ServeMux"
	case "chi":
		return "Chi"
	case "gorilla":
		return "Gorilla Mux"
	case "httprouter":
		return "HttpRouter"
	default:
		return routerType
	}
}

func validateProject() error {
	if _, err := os.Stat("go.mod"); err != nil {
		return fmt.Errorf("no go.mod found - please run this command in a Go project directory")
	}
	return nil
}

func installRouterDependency(routerType string) error {
	var dependency string

	switch routerType {
	case "stdlib":
		fmt.Println("Using standard library http.ServeMux - no additional dependency needed")
		return nil
	case "chi":
		dependency = "github.com/go-chi/chi/v5"
	case "gorilla":
		dependency = "github.com/gorilla/mux"
	case "httprouter":
		dependency = "github.com/julienschmidt/httprouter"
	default:
		return fmt.Errorf("unsupported router type: %s", routerType)
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

func updateMainWithRouter(routerType string) error {
	mainContent, err := os.ReadFile("main.go")
	if err != nil {
		return fmt.Errorf("failed to read main.go: %w", err)
	}

	newContent := generateRouterMainContent(routerType, string(mainContent))

	backupPath := "main.go.backup"
	if err := os.WriteFile(backupPath, mainContent, 0644); err != nil {
		fmt.Printf("Warning: failed to create backup at %s: %v\n", backupPath, err)
	}
	if err := os.WriteFile("main.go", []byte(newContent), 0644); err != nil {
		return fmt.Errorf("failed to write updated main.go: %w", err)
	}

	return nil
}

func generateRouterMainContent(routerType, existingContent string) string {
	switch routerType {
	case "stdlib":
		return generateServeMuxContent()
	case "chi":
		return generateChiContent()
	case "gorilla":
		return generateGorillaContent()
	case "httprouter":
		return generateHttpRouterContent()
	default:
		return existingContent
	}
}

func generateServeMuxContent() string {
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

func generateChiContent() string {
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

func generateGorillaContent() string {
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

func generateHttpRouterContent() string {
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

func printRouterInstructions(routerType string) {
	fmt.Println("\nNext steps:")
	fmt.Println("   go run main.go")
	fmt.Println("   curl http://localhost:8080/api/hello")
	fmt.Println("   curl http://localhost:8080/api/health")

	switch routerType {
	case "chi":
		fmt.Println("\nChi router features:")
		fmt.Println("   - Built-in middleware (Logger, Recoverer, RequestID)")
		fmt.Println("   - Route groups and subrouting")
		fmt.Println("   - Fast and lightweight")
	case "gorilla":
		fmt.Println("\nGorilla Mux features:")
		fmt.Println("   - Path variables: r.HandleFunc(\"/users/{id}\", handler)")
		fmt.Println("   - Query parameter matching")
		fmt.Println("   - Host and scheme matching")
	case "httprouter":
		fmt.Println("\nHttpRouter features:")
		fmt.Println("   - Extremely fast performance")
		fmt.Println("   - Path parameters: router.GET(\"/users/:id\", handler)")
		fmt.Println("   - Zero memory allocation")
	case "stdlib":
		fmt.Println("\nhttp.ServeMux features:")
		fmt.Println("   - Part of Go standard library")
		fmt.Println("   - Simple and reliable")
		fmt.Println("   - Pattern matching with wildcards")
	}
}
