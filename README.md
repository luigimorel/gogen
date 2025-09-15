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

### Global Commands and Flags

```bash
gogen --help                # Show all available commands
gogen -h                    # Short form of help
gogen --version             # Show version information
gogen -v                    # Short form of version
```

### Global Flags (Available for all commands)

| Flag         | Short | Description                   | Default |
| ------------ | ----- | ----------------------------- | ------- |
| `--help`     | `-h`  | Show help for the command     | false   |
| `--version`  | `-v`  | Show version information      | false   |
| `--verbose`  |       | Enable verbose output         | false   |
| `--quiet`    | `-q`  | Suppress non-essential output | false   |
| `--no-color` |       | Disable colored output        | false   |
| `--config`   | `-c`  | Path to configuration file    |         |

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

#### New Command - All Flags

```bash
gogen new --help
```

| Flag            | Short  | Description                                    | Default      |
| --------------- | ------ | ---------------------------------------------- | ------------ |
| `--name`        | `-n`   | Project name                                   |              |
| `--module`      | `-m`   | Go module path                                 | project name |
| `--template`    | `-t`   | Project template (api, web, cli)               | "api"        |
| `--router`      | `-r`   | Router type (stdlib, chi, gorilla, httprouter) | "stdlib"     |
| `--frontend`    | `--fe` | Frontend framework (react, vue, svelte, etc.)  |              |
| `--dir`         | `-d`   | Directory name for the project                 | project name |
| `--typescript`  | `--ts` | Use TypeScript for frontend projects           | false        |
| `--git`         | `-g`   | Initialize Git repository                      | true         |
| `--no-git`      |        | Skip Git repository initialization             | false        |
| `--mod-init`    |        | Initialize Go module                           | true         |
| `--no-mod-init` |        | Skip Go module initialization                  | false        |
| `--force`       | `-f`   | Overwrite existing directory                   | false        |
| `--dry-run`     |        | Show what would be created without executing   | false        |
| `--output`      | `-o`   | Output format (text, json)                     | "text"       |
| `--go-version`  |        | Go version to use in go.mod                    | "1.21"       |
| `--license`     | `-l`   | License type (mit, apache, gpl3, bsd)          | "mit"        |
| `--author`      | `-a`   | Author name for license                        | git config   |
| `--description` |        | Project description                            |              |
| `--private`     |        | Mark as private repository                     | false        |
| `--tags`        |        | Comma-separated list of project tags           |              |

#### Available Templates

- **api** (default) - REST API server with JSON responses
- **cli** - CLI application using urfave/cli/v2
- **web** - HTTP web server with optional frontend integration
- **microservice** - Microservice with Docker and health checks
- **grpc** - gRPC service with protobuf definitions
- **worker** - Background worker with queue processing

#### Available Routers

- **stdlib** (default) - Go standard library http.ServeMux
- **chi** - Chi lightweight router with middleware support
- **gorilla** - Gorilla Mux with advanced routing features
- **httprouter** - High-performance HttpRouter
- **gin** - Gin web framework with performance focus
- **echo** - Echo web framework with middleware
- **fiber** - Fiber web framework (Express-like)

#### Available Frontend Frameworks

- **react** - React with Vite build tool
- **vue** - Vue.js with Vite
- **svelte** - Svelte with Vite
- **solidjs** - SolidJS with Vite
- **angular** - Angular with Angular CLI
- **nextjs** - Next.js with React
- **nuxtjs** - Nuxt.js with Vue
- **vanilla** - Vanilla JavaScript/TypeScript

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

#### Install Command - All Flags

```bash
gogen install --help
```

| Flag        | Short | Description                                   | Default |
| ----------- | ----- | --------------------------------------------- | ------- |
| `--method`  | `-m`  | Installation method (auto, binary, nix, brew) | "auto"  |
| `--force`   | `-f`  | Force reinstall even if already installed     | false   |
| `--path`    | `-p`  | Custom installation path                      |         |
| `--symlink` | `-s`  | Create symlink instead of copying binary      | false   |
| `--shell`   |       | Shell type for PATH setup (bash, zsh, fish)   | auto    |
| `--no-path` |       | Skip PATH configuration                       | false   |
| `--backup`  | `-b`  | Backup existing installation                  | true    |
| `--verify`  |       | Verify installation after completion          | true    |

#### Supported Installation Methods

- **auto** - Automatically detects the best method for your system
- **binary** - Direct binary installation to ~/.local/bin (Linux/macOS) or %USERPROFILE%\AppData\Local\gogen (Windows)
- **nix** - Nix package manager
- **brew** - Homebrew package manager
- **snap** - Snap package manager (Linux)
- **scoop** - Scoop package manager (Windows)
- **winget** - Windows Package Manager

### Add Router to Existing Project

The `router` command adds a router to your existing Go project and updates your main.go file.

```bash
# Add Chi router (lightweight with middleware)
gogen router chi

# Add Gorilla Mux router
gogen router gorilla --update=false
```

#### Router Command - All Flags

```bash
gogen router --help
```

| Flag            | Short | Description                          | Default |
| --------------- | ----- | ------------------------------------ | ------- |
| `<router-type>` |       | Router type (chi, gorilla, etc.)     |         |
| `--update`      | `-u`  | Update main.go with router impl      | true    |
| `--middleware`  | `-w`  | Include common middleware            | true    |
| `--cors`        |       | Add CORS middleware                  | false   |
| `--auth`        |       | Add authentication middleware        | false   |
| `--logging`     |       | Add request logging middleware       | true    |
| `--recovery`    |       | Add panic recovery middleware        | true    |
| `--rate-limit`  |       | Add rate limiting middleware         | false   |
| `--compression` |       | Add response compression middleware  | false   |
| `--timeout`     |       | Add request timeout middleware       | false   |
| `--swagger`     |       | Add Swagger documentation support    | false   |
| `--metrics`     |       | Add Prometheus metrics middleware    | false   |
| `--force`       | `-f`  | Force overwrite existing router      | false   |
| `--backup`      | `-b`  | Backup existing files before changes | true    |

#### Router Features

- **stdlib** - Go standard library http.ServeMux with pattern matching
- **chi** - Lightweight router with built-in middleware (Logger, Recoverer, RequestID)
- **gorilla** - Full-featured router with path variables and advanced matching
- **httprouter** - Ultra-fast router with zero memory allocation and path parameters
- **gin** - High-performance web framework with JSON binding
- **echo** - Minimalist web framework with built-in middleware
- **fiber** - Express-inspired web framework built on Fasthttp

### Add Frontend to Existing Project

The `frontend` command adds a frontend framework to your existing project.

```bash
# Add React frontend
gogen frontend react

# Add Vue.js with TypeScript
gogen frontend vue --typescript
```

#### Frontend Command - All Flags

```bash
gogen frontend --help
```

| Flag                | Short  | Description                              | Default    |
| ------------------- | ------ | ---------------------------------------- | ---------- |
| `<framework>`       |        | Frontend framework type                  |            |
| `--dir`             | `-d`   | Directory name for frontend project      | "frontend" |
| `--typescript`      | `--ts` | Use TypeScript                           | false      |
| `--eslint`          |        | Include ESLint configuration             | true       |
| `--prettier`        |        | Include Prettier configuration           | true       |
| `--testing`         |        | Include testing framework setup          | false      |
| `--pwa`             |        | Configure as Progressive Web App         | false      |
| `--docker`          |        | Include Docker configuration             | false      |
| `--tailwind`        |        | Include Tailwind CSS                     | false      |
| `--sass`            |        | Include Sass/SCSS support                | false      |
| `--router`          |        | Include routing configuration            | true       |
| `--state`           |        | State management (redux, zustand, pinia) |            |
| `--api-base`        |        | Base URL for API calls                   | "/api"     |
| `--port`            | `-p`   | Development server port                  | 3000       |
| `--proxy`           |        | Proxy API calls to backend               | true       |
| `--proxy-target`    |        | Backend server URL for proxy             | ":8080"    |
| `--force`           | `-f`   | Overwrite existing frontend directory    | false      |
| `--template`        |        | Use specific project template            |            |
| `--package-manager` |        | Package manager (npm, yarn, pnpm)        | "npm"      |

#### Frontend Framework Details

- **react** - React 18+ with Vite, hot reloading, and modern tooling
- **vue** - Vue 3 with Composition API, Vite, and TypeScript support
- **svelte** - Svelte with SvelteKit and Vite integration
- **solidjs** - SolidJS with fine-grained reactivity and Vite
- **angular** - Angular with CLI, TypeScript, and modern build tools
- **nextjs** - Next.js with React and built-in optimizations
- **nuxtjs** - Nuxt.js with Vue and server-side rendering
- **vanilla** - Plain JavaScript/TypeScript with modern tooling

### Configuration Management

The `config` command manages gogen configuration settings.

```bash
# Show current configuration
gogen config show

# Set configuration value
gogen config set author.name "John Doe"
```

#### Config Command - All Flags

```bash
gogen config --help
```

| Flag       | Short | Description                      | Default |
| ---------- | ----- | -------------------------------- | ------- |
| `<action>` |       | Action (show, set, get, reset)   |         |
| `--global` | `-g`  | Use global configuration         | false   |
| `--local`  | `-l`  | Use local project configuration  | false   |
| `--file`   | `-f`  | Configuration file path          |         |
| `--format` |       | Output format (json, yaml, toml) | "yaml"  |

### Template Management

The `template` command manages custom project templates.

```bash
# List available templates
gogen template list

# Create custom template
gogen template create my-template --from api
```

#### Template Command - All Flags

```bash
gogen template --help
```

| Flag         | Short | Description                         | Default |
| ------------ | ----- | ----------------------------------- | ------- |
| `<action>`   |       | Action (list, create, delete, show) |         |
| `--from`     | `-f`  | Base template to extend             |         |
| `--path`     | `-p`  | Template directory path             |         |
| `--force`    |       | Force overwrite existing template   | false   |
| `--validate` | `-v`  | Validate template before creation   | true    |

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
curl http://localhost:8080/api/hello
curl http://localhost:8080/api/health
curl http://localhost:8080/metrics
```

### CLI Application

Create a comprehensive CLI tool:

```bash
gogen new --name my-cli --template cli --module github.com/company/my-cli \
  --author "John Doe" --license mit --description "My awesome CLI tool"

cd my-cli
go run main.go --help
```

### Add Features to Existing Project

Add a comprehensive router setup:

```bash
cd existing-go-project
gogen router chi --cors --auth --logging --metrics --swagger
```

Add a full-featured frontend:

```bash
cd existing-web-project
gogen frontend react --typescript --tailwind --testing --pwa \
  --state redux --dir client --port 3001
```

## Quick Reference

### Commands

| Command          | Description             | Example                                   |
| ---------------- | ----------------------- | ----------------------------------------- |
| `gogen new`      | Create a new project    | `gogen new -n my-app -t web --fe react`   |
| `gogen install`  | Install gogen to PATH   | `gogen install --force`                   |
| `gogen router`   | Add router to project   | `gogen router chi --cors --auth`          |
| `gogen frontend` | Add frontend to project | `gogen frontend vue --ts --tailwind`      |
| `gogen config`   | Manage configuration    | `gogen config set author.name "John"`     |
| `gogen template` | Manage custom templates | `gogen template create my-api --from api` |

### Templates

| Template       | Description              | Use Case                          |
| -------------- | ------------------------ | --------------------------------- |
| `api`          | REST API server          | Microservices, APIs, backends     |
| `web`          | Web server               | Full-stack applications, websites |
| `cli`          | CLI application          | Command-line tools, utilities     |
| `microservice` | Microservice with Docker | Containerized services            |
| `grpc`         | gRPC service             | High-performance RPC services     |
| `worker`       | Background worker        | Queue processing, batch jobs      |

### Routers

| Router       | Description                    | Best For                      |
| ------------ | ------------------------------ | ----------------------------- |
| `stdlib`     | Go standard library            | Simple applications, learning |
| `chi`        | Lightweight with middleware    | Most web applications         |
| `gorilla`    | Full-featured router           | Complex routing requirements  |
| `httprouter` | High-performance               | High-throughput APIs          |
| `gin`        | High-performance web framework | APIs with JSON processing     |
| `echo`       | Minimalist web framework       | Lightweight web services      |
| `fiber`      | Express-inspired framework     | High-performance web apps     |

### Frontend Frameworks

| Framework | Description      | TypeScript | Build Tool  | State Management |
| --------- | ---------------- | ---------- | ----------- | ---------------- |
| `react`   | React 18+        | ✅         | Vite        | Redux, Zustand   |
| `vue`     | Vue 3            | ✅         | Vite        | Pinia, Vuex      |
| `svelte`  | Svelte/SvelteKit | ✅         | Vite        | Svelte stores    |
| `solidjs` | SolidJS          | ✅         | Vite        | Built-in stores  |
| `angular` | Angular          | ✅         | Angular CLI | NgRx, Services   |
| `nextjs`  | Next.js          | ✅         | Next.js     | Redux, Zustand   |
| `nuxtjs`  | Nuxt.js          | ✅         | Nuxt        | Pinia, Vuex      |
| `vanilla` | Vanilla JS/TS    | ✅         | Vite        | Custom           |

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
│   ├── frontend.go     # Frontend integration command
│   ├── config.go       # Configuration management
│   └── template.go     # Template management
├── internal/           # Internal packages
│   ├── project.go      # Project generation logic
│   ├── create_frontend.go # Frontend creation
│   ├── env-file.go     # Environment file handling
│   ├── git-init.go     # Git initialization
│   ├── config/         # Configuration management
│   └── templates/      # Template definitions
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
