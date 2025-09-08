# gogen

A fast and simple CLI tool for generating Go project boilerplates.

## Features

- **Quick Project Setup** - Generate Go projects with proper structure in seconds
- **Multiple Templates** - Support for CLI, web, and API project templates
- **Auto Configuration** - Automatically initializes Go modules and dependencies
- **Cross Platform** - Works on Linux, macOS, and Windows
- **Zero Configuration** - Works out of the box with sensible defaults

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

After building or downloading the binary:

```bash
./gogen install
```

This will automatically detect your system and install gogen to your PATH.

## Usage

### Create a New Project

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

Generate an API server (explicit):

```bash
gogen new --name my-api --template api
```

Specify a custom module name:

```bash
gogen new --name my-project --module github.com/username/my-project
```

### Available Templates

- **api** (default) - REST API server with JSON responses
- **cli** - CLI application using urfave/cli/v2
- **web** - HTTP web server with basic routing

### Command Options

```bash
gogen new --help
  --name, -n     Project name (required)
  --module, -m   Go module path (default: project name)
  --template, -t Project template (cli, web, api) (default: api)
```

## Development

### Prerequisites

- Go 1.21.13 or later
- Make (optional, for build automation)

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
