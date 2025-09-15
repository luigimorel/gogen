package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"

	"github.com/luigimorel/gogen/internal/stdlibtemplate"
	"github.com/urfave/cli/v2"
)

// Router type constants
const (
	RouterStdlib     = "stdlib"
	RouterChi        = "chi"
	RouterGorilla    = "gorilla"
	RouterHTTPRouter = "httprouter"
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
				return errors.New("router type is required. Usage: gogen router <router-type>")
			}

			routerType := c.Args().Get(0)
			updateMain := c.Bool("update")

			router := NewRouter(routerType, updateMain)
			return router.execute()
		},
	}
}

func (r *Router) execute() error {
	// A main.go is not required for a valid project, ex packages, creeate it in updatemainFile()
	// if err := r.validateProject(); err != nil {
	// 	return err
	// }

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
	case RouterHTTPRouter:
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
	var mainFile *os.File
	var err error

	mainFile, err = os.Open("main.go")
	switch {
	case os.IsNotExist(err):
		_, err := os.Create("main.go")
		if err != nil {
			return fmt.Errorf("failed to create main.go: %w", err)
		}
	case err != nil:
		return fmt.Errorf("failed to check if main.go exists: %w", err)
	default:

		var backupFileName = "main.go.backup"
		_, err := os.Stat(backupFileName)
		switch {
		case os.IsNotExist(err):
			// keep the name
		case err != nil:
			return fmt.Errorf("failed to check if backup file exists: %w", err)
		default:
			// If stat does not return a error, the file exist, hence create a timestamped backup file
			backupFileName = "main.go.backup" + time.Now().Format("20060102150405")

		}

		backupFile, err := os.Create(backupFileName)
		switch {
		case err != nil:
			return fmt.Errorf("failed to create backup file: %w", err)
		default:
			defer backupFile.Close()
			_, err = io.Copy(backupFile, mainFile)
			if err != nil {
				return fmt.Errorf("failed to copy main.go to backup file: %w", err)
			}
		}
	}

	defer mainFile.Close()

	mainContent := &bytes.Buffer{}

	_, err = io.Copy(mainContent, mainFile)
	if err != nil {
		return fmt.Errorf("failed to read main.go: %w", err)
	}

	_, err = mainContent.WriteString(r.generateMainContent(mainContent.String()))
	if err != nil {
		return fmt.Errorf("failed to write updated main.go: %w", err)
	}

	_ = mainFile.Truncate(0)

	_, err = mainFile.Seek(0, 0)
	if err != nil {
		return fmt.Errorf("failed to seek to beginning of main.go: %w", err)
	}
	_, err = io.Copy(mainFile, mainContent)
	if err != nil {
		return fmt.Errorf("failed to write updated main.go: %w", err)
	}

	return nil
}

func (r *Router) generateMainContent(existingContent string) string {
	switch r.Type {
	case RouterStdlib:
		// ugly addition as the command currently only support updating a main.go file
		err := stdlibtemplate.CreateRouterSetup()
		if err != nil {
			fmt.Printf("Warning: failed to create router setup: %v", err)
			return "" // not ideal
		}

		return r.generateServeMuxContent()
	case RouterChi:
		return r.generateChiContent()
	case RouterGorilla:
		return r.generateGorillaContent()
	case RouterHTTPRouter:
		return r.generateHTTPRouterContent()
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
	case RouterHTTPRouter:
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
	"fmt"
	"log"
	"net/http"

	"<your-module-name>/router"
)

func main() {
	// Example of adding a handler to the router
	router.Router.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello from the stdlib router!")
	})

	addr := ":8080"
	fmt.Printf("Starting server on http://localhost%s\n", addr)
	log.Fatal(router.Serve(addr))
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

func (r *Router) generateHTTPRouterContent() string {
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
