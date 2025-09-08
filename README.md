# gogen

A fast and simple CLI tool for generating Go project boilerplates.

## Features

- **Quick Project Setup** - Generate Go projects with proper structure in seconds
- **Multiple Templates** - Support for CLI, web, and API project templates
- **Router Integration** - Built-in support for popular Go routers (Chi, Gorilla Mux, HttpRouter, stdlib)
- **Frontend Integration** - Add modern frontend frameworks (React, Vue, Svelte, SolidJS, Angular) with TypeScript support
- **Auto Configuration** - Automatically initializes Go modules and dependencies
- **Self-Installing** - Install gogen to your system PATH with a single command
- **Cross Platform** - Works on Linux, macOS, and Windows
- **Zero Configuration** - Works out of the box with sensible defaults
- **Extensible Architecture** - Add routers and frontends to existing projects

## Installation

### Quick Install

```bash
# Build from source
git clone https://github.com/luigimorel/gogen.git
cd gogen
make build
```

### Using Go Install

```bash
go install github.com/luigimorel/gogen@latest
```

### Using the CLI (Self-Install)

After building or downloading the binary, gogen can install itself to your system PATH:

```bash
# Auto-detect best installation method for your system
./gogen install

# Force reinstallation if already installed
./gogen install --force

# Choose specific installation method
./gogen install --method binary  # Direct binary installation
./gogen install --method nix     # Nix package manager (planned)
./gogen install --method brew    # Homebrew (planned)
```

The install command supports:

- **Linux/macOS**: Installs to `~/.local/bin` with PATH configuration help
- **Windows**: Installs to `%USERPROFILE%\AppData\Local\gogen` with PATH setup instructions
- **Auto-detection**: Automatically chooses the best method for your system

## Usage

### Global Commands

```bash
gogen --help                # Show all available commands
gogen --version             # Show version information
```

### Create a New Project

The `new` command creates a new Go project with proper structure and initialization.

#### Basic Usage

Generate an API server (default):

```bash
gogen new --name my-project
```

Generate a CLI application:

```bash
gogen new --name my-cli --template cli
```

Generate a web server:

```bash
gogen new --name my-web-app --template web
```

#### Advanced Usage

Create a web project with React frontend:

```bash
gogen new --name my-fullstack --template web --frontend react
```

Create a web project with Vue.js and TypeScript:

```bash
gogen new --name my-vue-app --template web --frontend vue --ts
```

Create an API with Chi router:

```bash
gogen new --name my-api --template api --router chi
```

Specify custom module name and directory:

```bash
gogen new --name my-project --module github.com/username/my-project --dir custom-dir
```

#### Available Templates

- **api** (default) - REST API server with JSON responses
- **cli** - CLI application using urfave/cli/v2
- **web** - HTTP web server with optional frontend integration

#### Available Routers

- **stdlib** (default) - Go standard library http.ServeMux
- **chi** - Chi lightweight router with middleware support
- **gorilla** - Gorilla Mux with advanced routing features
- **httprouter** - High-performance HttpRouter

#### Available Frontend Frameworks

- **react** - React with Vite build tool
- **vue** - Vue.js with Vite
- **svelte** - Svelte with Vite
- **solidjs** - SolidJS with Vite
- **angular** - Angular with Angular CLI

#### New Command Options

```bash
gogen new --help
  --name, -n         Project name (required)
  --module, -m       Go module path (default: project name)
  --template, -t     Project template (cli, web, api) (default: "api")
  --router, -r       Router type for API/web projects (stdlib, chi, gorilla, httprouter) (default: "stdlib")
  --frontend, --fe   Frontend framework for web projects (react, vue, svelte, solidjs, angular)
  --dir              Directory name for the project (default: project name)
  --ts               Use TypeScript for frontend projects (only applicable with --frontend)
```

### Install gogen to System PATH

The `install` command automatically installs gogen to your system PATH for easy access from anywhere.

```bash
# Auto-detect best installation method
gogen install

# Force reinstallation
gogen install --force

# Specify installation method
gogen install --method binary
gogen install --method nix
gogen install --method brew
```

#### Install Command Options

```bash
gogen install --help
  --method, -m  Installation method (auto, binary, nix, brew) (default: "auto")
  --force, -f   Force reinstall even if already installed
```

#### Supported Installation Methods

- **auto** - Automatically detects the best method for your system
- **binary** - Direct binary installation to ~/.local/bin (Linux/macOS) or %USERPROFILE%\AppData\Local\gogen (Windows)
- **nix** - Nix package manager (planned)
- **brew** - Homebrew package manager (planned)

### Add Router to Existing Project

The `router` command adds a router to your existing Go project and updates your main.go file.

```bash
# Add Chi router (lightweight with middleware)
gogen router chi

# Add Gorilla Mux router
gogen router gorilla

# Add HttpRouter (high performance)
gogen router httprouter

# Add router without updating main.go
gogen router chi --update=false
```

#### Router Command Options

```bash
gogen router --help
  <router-type>     Router type: chi, gorilla, httprouter, or stdlib (required)
  --update, -u      Update main.go with router implementation (default: true)
```

#### Router Features

- **stdlib** - Go standard library http.ServeMux with pattern matching
- **chi** - Lightweight router with built-in middleware (Logger, Recoverer, RequestID)
- **gorilla** - Full-featured router with path variables and advanced matching
- **httprouter** - Ultra-fast router with zero memory allocation and path parameters

### Add Frontend to Existing Project

The `frontend` command adds a frontend framework to your existing project.

```bash
# Add React frontend
gogen frontend react

# Add Vue.js with TypeScript
gogen frontend vue --typescript

# Add Svelte in custom directory
gogen frontend svelte --dir client

# Add Angular frontend
gogen frontend angular
```

#### Frontend Command Options

```bash
gogen frontend --help
  <framework-type>       Frontend framework: react, vue, svelte, solidjs, or angular (required)
  --dir, -d              Directory name for the frontend project (default: "frontend")
  --typescript, --ts     Use TypeScript (where supported)
```

#### Frontend Framework Details

- **react** - React 18+ with Vite, hot reloading, and modern tooling
- **vue** - Vue 3 with Composition API, Vite, and TypeScript support
- **svelte** - Svelte with SvelteKit and Vite integration
- **solidjs** - SolidJS with fine-grained reactivity and Vite
- **angular** - Angular with CLI, TypeScript, and modern build tools

## Examples

### Full-Stack Web Application

Create a complete full-stack application with Go backend and React frontend:

```bash
# Create web project with React frontend
gogen new --name my-app --template web --frontend react --ts

cd my-app

# Start the API server (in one terminal)
cd api
go run main.go

# Start the frontend dev server (in another terminal)
cd frontend
npm run dev
```

### Microservice API

Create a high-performance API service with Chi router:

```bash
# Create API project with Chi router
gogen new --name my-service --template api --router chi --module github.com/company/my-service

cd my-service
go run main.go

# Test the API
curl http://localhost:8080/api/hello
curl http://localhost:8080/api/health
```

### CLI Application

Create a CLI tool:

```bash
gogen new --name my-cli --template cli --module github.com/company/my-cli

cd my-cli
go run main.go --help
```

### Add Features to Existing Project

Add a router to an existing Go project:

```bash
cd existing-go-project
gogen router gorilla
```

Add a frontend to an existing web project:

```bash
cd existing-web-project
gogen frontend vue --typescript --dir client
```

## Quick Reference

### Commands

| Command          | Description             | Example                                                   |
| ---------------- | ----------------------- | --------------------------------------------------------- |
| `gogen new`      | Create a new project    | `gogen new --name my-app --template web --frontend react` |
| `gogen install`  | Install gogen to PATH   | `gogen install --force`                                   |
| `gogen router`   | Add router to project   | `gogen router chi`                                        |
| `gogen frontend` | Add frontend to project | `gogen frontend vue --ts`                                 |

### Templates

| Template | Description     | Use Case                                       |
| -------- | --------------- | ---------------------------------------------- |
| `api`    | REST API server | Microservices, APIs, backend services          |
| `web`    | Web server      | Full-stack applications, server-rendered sites |
| `cli`    | CLI application | Command-line tools, utilities                  |

### Routers

| Router       | Description                 | Best For                      |
| ------------ | --------------------------- | ----------------------------- |
| `stdlib`     | Go standard library         | Simple applications, learning |
| `chi`        | Lightweight with middleware | Most web applications         |
| `gorilla`    | Full-featured router        | Complex routing requirements  |
| `httprouter` | High-performance            | High-throughput APIs          |

### Frontend Frameworks

| Framework | Description      | TypeScript | Build Tool  |
| --------- | ---------------- | ---------- | ----------- |
| `react`   | React 18+        | ✅         | Vite        |
| `vue`     | Vue 3            | ✅         | Vite        |
| `svelte`  | Svelte/SvelteKit | ✅         | Vite        |
| `solidjs` | SolidJS          | ✅         | Vite        |
| `angular` | Angular          | ✅         | Angular CLI |

## Development

### Prerequisites

- Go 1.21.13 or later
- Make (optional, for build automation)
- Node.js and npm (required for frontend features)

### Building from Source

```bash
git clone https://github.com/luigimorel/gogen.git
cd gogen
make build
```

### Development Workflow

```bash
# Install dependencies
make deps

# Run in development mode with hot reloading
make dev

# Run tests
make test

# Format and lint code
make check

# Build for all platforms
make build-all
```

### Testing Generated Projects

gogen includes a test project structure in `bin/testweb/` that demonstrates a full-stack web application with both API and frontend components. This serves as both a testing ground and a reference implementation.

### Project Structure

```text
gogen/
├── main.go              # Entry point
├── cmd/                 # CLI commands
│   ├── init.go         # App initialization
│   ├── new.go          # Project creation command
│   ├── install.go      # Self-installation command
│   ├── router.go       # Router management command
│   └── frontend.go     # Frontend integration command
├── internal/           # Internal packages
│   ├── project.go      # Project generation logic
│   ├── create_frontend.go # Frontend creation
│   ├── env-file.go     # Environment file handling
│   └── git-init.go     # Git initialization
└── bin/testweb/       # Test/example project
    ├── api/           # Go API server
    └── frontend/      # React frontend
```

### Available Make Commands

- `make build` - Build the application
- `make run` - Build and run the application
- `make dev` - Start development server with hot reloading
- `make test` - Run tests
- `make fmt` - Format code
- `make lint` - Lint code
- `make vet` - Vet code
- `make check` - Run all code quality checks
- `make clean` - Clean build artifacts
- `make build-all` - Build for multiple platforms

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Run tests and checks (`make check test`)
5. Commit your changes (`git commit -m 'Add amazing feature'`)
6. Push to the branch (`git push origin feature/amazing-feature`)
7. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](./LICENSE) file for details.
