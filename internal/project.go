package internal

import (
	"fmt"
	"os"
	"os/exec"
)

type ProjectGenerator struct {
	RouterGenerator *RouterGenerator
}

func NewProjectGenerator() *ProjectGenerator {
	return &ProjectGenerator{
		RouterGenerator: NewRouterGenerator(),
	}
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
	"net/http"
	"os"

	"github.com/joho/godotenv"`
}

func (pg *ProjectGenerator) generateMainContent(moduleName, projectName, routerType string) string {
	var routerSetup, serverStart string

	if routerType == "stdlib" {
		routerSetup = "web.SetupRoutes()"
		serverStart = "log.Fatal(http.ListenAndServe(port, nil))"
	} else {
		routerSetup = "r := web.SetupRoutes()"
		serverStart = "log.Fatal(http.ListenAndServe(port, r))"
	}

	return `package main

import (
	` + pg.setDefaultPackages() + `
	"` + pg.setModuleName(moduleName, projectName) + `/cmd/web"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	` + routerSetup + `

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	port = ":" + port

	fmt.Printf("Starting web server on http://localhost%s\n", port)
	` + serverStart + `
}`
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

func (pg *ProjectGenerator) CreateWebProject(projectName, moduleName, router, frontendFramework string, useTypeScript bool, runtime string) error {
	if err := pg.createAPIProjectInDir("api", projectName, moduleName, router); err != nil {
		return fmt.Errorf("failed to create API project: %w", err)
	}

	originalDir, _ := os.Getwd()
	if err := os.Chdir("api"); err != nil {
		return fmt.Errorf("failed to change to api directory: %w", err)
	}

	if err := pg.createConfigFiles(); err != nil {
		os.Chdir(originalDir)
		return fmt.Errorf("failed to create config files: %w", err)
	}

	if err := os.Chdir(originalDir); err != nil {
		return fmt.Errorf("failed to change back to original directory: %w", err)
	}

	if err := pg.CreateFrontendProject(frontendFramework, "frontend", useTypeScript, runtime); err != nil {
		return fmt.Errorf("failed to create frontend project: %w", err)
	}

	if err := os.Chdir("frontend"); err != nil {
		return fmt.Errorf("failed to change to frontend directory: %w", err)
	}

	if err := pg.CreateEnvFile("frontend", "."); err != nil {
		fmt.Printf("Warning: failed to create env file: %v\n", err)
	}

	if err := pg.RemoveGitRepository("."); err != nil {
		fmt.Printf("Warning: failed to remove git repository from frontend: %v\n", err)
	} else {
		fmt.Println("Removed git repository from frontend")
	}

	return nil
}

func (pg *ProjectGenerator) createConfigFiles() error {
	if err := pg.CreateEnvFile("api", "."); err != nil {
		return fmt.Errorf("warning: failed to create env file in api: %v", err)
	}

	if err := pg.CreateGitignoreFile("api", "."); err != nil {
		return fmt.Errorf("warning: failed to create .gitignore file in api: %v", err)
	}
	return nil
}

func (pg *ProjectGenerator) CreateAPIProject(projectName, moduleName, router string) error {
	return pg.createAPIProjectInDir(".", projectName, moduleName, router)
}

func (pg *ProjectGenerator) createAPIProjectInDir(baseDir, projectName, moduleName, router string) error {
	cmdWebDir := fmt.Sprintf("%s/cmd/web", baseDir)
	if baseDir == "." {
		cmdWebDir = "cmd/web"
	}

	if err := os.MkdirAll(cmdWebDir, 0755); err != nil {
		return fmt.Errorf("failed to create cmd/web directory: %w", err)
	}

	var mainContent string
	var routesContent string

	switch router {
	case "chi":
		mainContent = pg.generateMainContent(moduleName, projectName, "chi")
		routesContent = pg.RouterGenerator.generateChiContent()

	case "gorilla":
		mainContent = pg.generateMainContent(moduleName, projectName, "gorilla")
		routesContent = pg.RouterGenerator.generateGorillaContent()

	case "httprouter":
		mainContent = pg.generateMainContent(moduleName, projectName, "httprouter")
		routesContent = pg.RouterGenerator.generateHttpRouterContent()

	default: // stdlib
		mainContent = pg.generateMainContent(moduleName, projectName, "stdlib")
		routesContent = pg.RouterGenerator.generateStdlibContent()
	}

	mainGoPath := fmt.Sprintf("%s/main.go", baseDir)
	routesPath := fmt.Sprintf("%s/cmd/web/routes.go", baseDir)
	goModPath := fmt.Sprintf("%s/go.mod", baseDir)

	if baseDir == "." {
		mainGoPath = "main.go"
		routesPath = "cmd/web/routes.go"
		goModPath = "go.mod"
	}

	if err := os.WriteFile(mainGoPath, []byte(mainContent), 0644); err != nil {
		return err
	}

	if err := os.WriteFile(routesPath, []byte(routesContent), 0644); err != nil {
		return err
	}

	baseModuleName := pg.setModuleName(moduleName, projectName)
	apiModContent := fmt.Sprintf("module %s\n\ngo 1.21\n", baseModuleName)
	if err := os.WriteFile(goModPath, []byte(apiModContent), 0644); err != nil {
		return err
	}

	if err := pg.CreateEnvFile("api", baseDir); err != nil {
		fmt.Printf("Warning: failed to create env file: %v\n", err)
	}

	if err := pg.InitGitRepository(projectName, baseDir); err != nil {
		fmt.Printf("Warning: failed to initialize git repository: %v\n", err)
	}

	originalDir, _ := os.Getwd()
	if baseDir != "." {
		if err := os.Chdir(baseDir); err != nil {
			return fmt.Errorf("failed to change to %s directory: %w", baseDir, err)
		}
	}

	cmd := exec.Command("go", "mod", "tidy")
	if err := cmd.Run(); err != nil {
		if baseDir != "." {
			os.Chdir(originalDir)
		}
		return fmt.Errorf("failed to tidy go.mod: %w", err)
	}

	if baseDir != "." {
		if err := os.Chdir(originalDir); err != nil {
			return fmt.Errorf("failed to change to original directory: %w", err)
		}
	}

	return nil
}
