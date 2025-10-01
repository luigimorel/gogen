# gogen

A fast and simple CLI tool for generating Go project boilerplates.

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
```

The install command supports:

- **Linux/macOS**: Installs to `~/.local/bin` with PATH configuration help
- **Windows**: Installs to `%USERPROFILE%\AppData\Local\gogen` with PATH setup instructions
- **Auto-detection**: Automatically chooses the best method for your system

## Usage

### Global Commands and Flags

```bash
gogen --help                # Show all available commands
gogen -h                    # Short form of help
gogen --version             # Show version information
gogen -v                    # Short form of version
```

### Global Flags (Available for all commands)

| Flag        | Short | Description               | Default |
| ----------- | ----- | ------------------------- | ------- |
| `--help`    | `-h`  | Show help for the command | false   |
| `--version` | `-v`  | Show version information  | false   |
| `--verbose` |       | Enable verbose output     | false   |

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

Create a fullstack project with Chi router, docker, bun runtime and React for the frontend:

```bash
gogen new --name my-app --docker --template web --router chi --ts
```

Specify custom module name and directory:

```bash
gogen new --name my-project --module github.com/username/my-project --dir custom-dir
```

#### New Command - All Flags

```bash
gogen new --help
```

| Flag           | Short  | Description                                    | Default      |
| -------------- | ------ | ---------------------------------------------- | ------------ |
| `--name`       | `-n`   | Project name                                   |              |
| `--module`     | `-m`   | Go module path                                 | project name |
| `--template`   | `-t`   | Project template (api, web, cli)               | "api"        |
| `--router`     | `-r`   | Router type (stdlib, chi, gorilla, httprouter) | "stdlib"     |
| `--frontend`   | `--fe` | Frontend framework (react, vue, svelte, etc.)  |              |
| `--dir`        | `-d`   | Directory name for the project                 | project name |
| `--typescript` | `--ts` | Use TypeScript for frontend projects           | false        |
| `--docker`     |        | Create dockerfiles and dockercompose           | false        |

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

### Install gogen to System PATH

The `install` command automatically installs gogen to your system PATH for easy access from anywhere.

```bash
# Auto-detect best installation method
gogen install

# Force reinstallation
gogen install --force

# Specify installation method
gogen install --method binary
```

#### Router Features

- **stdlib** - Go standard library http.ServeMux with pattern matching
- **chi** - Lightweight router with built-in middleware (Logger, Recoverer, RequestID)
- **gorilla** - Full-featured router with path variables and advanced matching
- **httprouter** - Ultra-fast router with zero memory allocation and path parameters

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
# Create web project with React frontend and all features
gogen new --name my-app --template web --frontend react --typescript \
  --router chi --cors --auth --logging --tailwind --testing

cd my-app

# Start the API server (in one terminal)
cd api
go run main.go

# Start the frontend dev server (in another terminal)
cd frontend
npm run dev
```

### Microservice API

Create a high-performance API service with comprehensive features:

```bash
# Create API project with full microservice setup
gogen new --name my-service --template microservice --router gin \
  --module github.com/company/my-service --middleware --cors \
  --metrics --rate-limit --swagger --docker

cd my-service
go run main.go

# Test the API
curl http://localhost:8080/
curl http://localhost:8080/health
```

### CLI Application

Create a comprehensive CLI tool:

```bash
gogen new --name my-cli --template cli --module github.com/company/my-cli

cd my-cli
go run main.go --help
```

## Quick Reference

### Commands

| Command         | Description           | Example                                 |
| --------------- | --------------------- | --------------------------------------- |
| `gogen new`     | Create a new project  | `gogen new -n my-app -t web --fe react` |
| `gogen install` | Install gogen to PATH | `gogen install --force`                 |

### Templates

| Template | Description     | Use Case                          |
| -------- | --------------- | --------------------------------- |
| `api`    | REST API server | Microservices, APIs, backends     |
| `web`    | Web server      | Full-stack applications, websites |
| `cli`    | CLI application | Command-line tools, utilities     |

### Routers

| Router       | Description                 | Best For                      |
| ------------ | --------------------------- | ----------------------------- |
| `stdlib`     | Go standard library         | Simple applications, learning |
| `chi`        | Lightweight with middleware | Most web applications         |
| `gorilla`    | Full-featured router        | Complex routing requirements  |
| `httprouter` | High-performance            | High-throughput APIs          |

### Frontend Frameworks

| Framework | Description      | TypeScript | Build Tool  | State Management |
| --------- | ---------------- | ---------- | ----------- | ---------------- |
| `react`   | React 18+        | ✅         | Vite        | Redux, Zustand   |
| `vue`     | Vue 3            | ✅         | Vite        | Pinia, Vuex      |
| `svelte`  | Svelte/SvelteKit | ✅         | Vite        | Svelte stores    |
| `solidjs` | SolidJS          | ✅         | Vite        | Built-in stores  |
| `angular` | Angular          | ✅         | Angular CLI | NgRx, Services   |

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
# Run in development mode with hot reloading
make dev

# Run tests
make test

# Format and lint code
make check

# Build for all platforms
make build-all
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
