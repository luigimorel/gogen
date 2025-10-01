package internal

import (
	"fmt"
	"os"
	"path/filepath"

	constants "github.com/luigimorel/gogen/consants"
)

var runTimeContent = `
FROM oven/bun:1-alpine AS builder

WORKDIR /app

COPY package.json bun.lock ./

RUN bun install --frozen-lockfile

COPY . .

RUN bun run build

FROM oven/bun:1-alpine

RUN addgroup -g 1001 -S user

RUN adduser -S user -u 1001 -G user

WORKDIR /app

RUN chown -R user:user /app

USER user

COPY --chown=user:user package.json bun.lock ./

RUN bun install --frozen-lockfile

COPY --from=builder --chown=user:user /app/dist ./dist

EXPOSE 4173

CMD ["bun", "run", "preview", "--host", "0.0.0.0", "--port", "4173"]
`

func (pg *ProjectGenerator) CreateDockerfile(dirName, dirType, runtime string) error {
	var dockerContent string
	var dockerIgnoreContent string

	if dirType == constants.APIDir {
		dockerContent = `FROM golang:1.21-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest

RUN apk --no-cache add ca-certificates

RUN addgroup -g 1001 -S appgroup && \
    adduser -S appuser -u 1001 -G appgroup

WORKDIR /root/

COPY --from=builder /app/main .

COPY --from=builder /app/.env* ./

RUN chown -R appuser:appgroup /root

USER appuser

EXPOSE 8080

CMD ["./main"]
`
		dockerIgnoreContent = `# Binaries
*.exe
*.exe~
*.dll
*.so
*.dylib
main
tmp/*
*.test
*.out
go.work
.vscode/
.idea/
*.swp
*.swo
*~
.DS_Store
.DS_Store?
._*
.Spotlight-V100
.Trashes
ehthumbs.db
Thumbs.db
*.log
logs/
tmp/
*.tmp
*.temp
.git/
.gitignore
*.md
README*
Dockerfile*
docker-compose*
.dockerignore
`
	} else {
		if runtime == "bun" {
			dockerContent = runTimeContent
		} else {
			dockerContent = `FROM node:22-alpine AS builder

WORKDIR /app

COPY package*.json ./

RUN npm ci

COPY . .

RUN npm run build

FROM node:22-alpine

RUN addgroup -g 1001 -S user

RUN adduser -S user -u 1001 -G user

WORKDIR /app

RUN chown -R user:user /app

USER user

COPY --chown=user:user package*.json ./

RUN npm ci --only=production && npm install vite

COPY --from=builder --chown=user:user /app/dist ./dist

EXPOSE 4173

# Vite preview port is 4173. You can change to nginx if needed.
CMD ["npm", "run", "preview", "--", "--host", "0.0.0.0", "--port", "4173"]
`
		}
		dockerIgnoreContent = `node_modules
		# Dependencies
node_modules/
npm-debug.log*
yarn-debug.log*
yarn-error.log*
pnpm-debug.log*

# Production build
dist/
build/

# Runtime data
pids
*.pid
*.seed
*.pid.lock

# Coverage directory used by tools like istanbul
coverage/
*.lcov

# IDE files
.vscode/
.idea/
*.swp
*.swo
*~

# OS files
.DS_Store
.DS_Store?
._*
.Spotlight-V100
.Trashes
ehthumbs.db
Thumbs.db

# Logs
*.log
logs/

# Environment variables
.env
.env.local
.env.development.local
.env.test.local
.env.production.local

# Git
.git/
.gitignore

# Documentation
*.md
README*

# Docker files
Dockerfile*
docker-compose*
.dockerignore`
	}

	dockerfilePath := filepath.Join(dirName, "Dockerfile")

	dockerignorePath := filepath.Join(dirName, ".dockerignore")

	if err := os.WriteFile(dockerfilePath, []byte(dockerContent), 0600); err != nil {
		return fmt.Errorf("failed to create Dockerfile: %w", err)
	}

	if err := os.WriteFile(dockerignorePath, []byte(dockerIgnoreContent), 0600); err != nil {
		return fmt.Errorf("failed to create .dockerignore: %w", err)
	}

	return nil
}

func (pg *ProjectGenerator) CreateDockerComposeFile(dirName string) error {
	var dockerComposeContent string
	var dockerComposeOverrideContent string

	dockerComposeContent = `services:
  api:
    build:
      context: ./api
      dockerfile: Dockerfile
    container_name: audits-api
    ports:
      - "8080:8080"
    environment:
      - PORT=8080
      - ENV=production
    volumes:
      - ./api/.env:/root/.env:ro
    networks:
      - default
    restart: unless-stopped
    healthcheck:
      test:
        [
          "CMD",
          "wget",
          "--quiet",
          "--tries=1",
          "--spider",
          "http://localhost:8080/health",
        ]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s

  frontend:
    build:
      context: ./frontend
      dockerfile: Dockerfile
    container_name: audits-frontend
    ports:
      - "4173:4173"
    depends_on:
      api:
        condition: service_healthy
    networks:
      - default
    restart: unless-stopped
    healthcheck:
      test:
        [
          "CMD",
          "wget",
          "--quiet",
          "--tries=1",
          "--spider",
          "http://localhost:4173",
        ]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s

networks:
  default:
    driver: bridge

volumes:
  api-data:
    driver: local
`

	dockerComposeOverrideContent = `services:
  api:
    build:
      target: builder
    volumes:
      - ./api:/app
      - /app/tmp
    environment:
      - ENV=development
      - GO_ENV=development
    command: ["sh", "-c", "go mod download && go run main.go"]
    ports:
      - "8080:8080"
      - "2345:2345" #Delve debugger

  frontend:
    build:
      target: builder
    volumes:
      - ./frontend:/app
      - /app/node_modules
      - /app/dist
    environment:
      - NODE_ENV=development
    command: ["npm", "run", "dev", "--", "--host", "0.0.0.0"]
    ports:
      - "5173:5173"

`

	dockerComposefilePath := dirName + "/docker-compose.yml"
	if err := os.WriteFile(dockerComposefilePath, []byte(dockerComposeContent), 0600); err != nil {
		return fmt.Errorf("failed to create Dockerfile: %w", err)
	}

	composeOverrideFilePath := dirName + "/docker-compose.override.yml"
	if err := os.WriteFile(composeOverrideFilePath, []byte(dockerComposeOverrideContent), 0600); err != nil {
		return fmt.Errorf("failed to create docker-compose.override.yml: %w", err)
	}

	return nil
}
