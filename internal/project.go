package internal

import (
	"fmt"
	"os"
	"os/exec"

	constants "github.com/luigimorel/gogen/consants"
)

type ProjectGenerator struct {
	RouterGenerator *RouterGenerator
}

type WebProjectConfig struct {
	ProjectName       string
	ModuleName        string
	Router            string
	FrontendFramework string
	Runtime           string
	UseTypeScript     bool
	UseTailwind       bool
	UseDocker         bool
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

	if err := os.WriteFile("main.go", []byte(mainContent), 0600); err != nil {
		return err
	}

	if err := pg.InitGitRepository(projectName, constants.CLITemplate); err != nil {
		fmt.Printf("Warning: failed to initialize git repository: %v\n", err)
	}

	if err := pg.CreateAirFile(".", constants.CLITemplate); err != nil {
		fmt.Printf("Warning: failed to create .air.toml file: %v\n", err)
	}

	if err := pg.CreateGitignoreFile(constants.CLITemplate, "."); err != nil {
		fmt.Printf("Warning: failed to create .gitignore file: %v\n", err)
	}

	cmd := exec.Command("go", "mod", "tidy")
	return cmd.Run()
}

func (pg *ProjectGenerator) CreateWebProject(projectName, moduleName, router, frontendFramework, runtime string, useTypeScript, useTailwind, useDocker bool) error {
	config := &WebProjectConfig{
		ProjectName:       projectName,
		ModuleName:        moduleName,
		Router:            router,
		FrontendFramework: frontendFramework,
		Runtime:           runtime,
		UseTypeScript:     useTypeScript,
		UseTailwind:       useTailwind,
		UseDocker:         useDocker,
	}

	return pg.CreateWebProjectWithConfig(config)
}

func (pg *ProjectGenerator) CreateWebProjectWithConfig(config *WebProjectConfig) error {
	dm, err := NewDirectoryManager()
	if err != nil {
		return fmt.Errorf("failed to initialize directory manager: %w", err)
	}

	if err := pg.setupAPIProject(config, dm); err != nil {
		return fmt.Errorf("failed to setup API project: %w", err)
	}

	if err := pg.setupFrontendProject(config, dm); err != nil {
		return fmt.Errorf("failed to setup frontend project: %w", err)
	}

	return nil
}

func (pg *ProjectGenerator) setupAPIProject(config *WebProjectConfig, dm *DirectoryManager) error {
	if err := pg.createAPIProjectInDir(constants.APIDir, config.ProjectName, config.ModuleName, config.Router); err != nil {
		return fmt.Errorf("failed to create API project: %w", err)
	}

	if err := dm.ChangeToDir(constants.APIDir); err != nil {
		return fmt.Errorf("failed to change to API directory: %w", err)
	}
	defer func() {
		if err := dm.RootDir(); err != nil {
			fmt.Printf("Warning: failed to change root directory: %v\n", err)
		}
	}()

	if config.UseDocker {
		if err := pg.CreateDockerfile(".", constants.APIDir, config.Runtime); err != nil {
			return fmt.Errorf("failed to create Docker files for API: %w", err)
		}

		if err := pg.CreateDockerComposeFile(".."); err != nil {
			return fmt.Errorf("failed to create docker-compose files: %w", err)
		}
	}

	if err := pg.createConfigFiles(); err != nil {
		return fmt.Errorf("failed to create config files: %w", err)
	}

	return nil
}

func (pg *ProjectGenerator) setupFrontendProject(config *WebProjectConfig, dm *DirectoryManager) error {
	if err := dm.RootDir(); err != nil {
		return fmt.Errorf("failed to change to root directory before frontend setup: %w", err)
	}

	if err := pg.CreateFrontendProject(config.FrontendFramework, constants.FrontendDir, config.UseTypeScript, config.Runtime, config.UseTailwind); err != nil {
		return fmt.Errorf("failed to create frontend project: %w", err)
	}

	if err := dm.ChangeToDir(constants.FrontendDir); err != nil {
		return fmt.Errorf("failed to change to frontend directory: %w", err)
	}

	pg.createFrontendConfigFiles(config)

	defer func() {
		if restoreErr := dm.RootDir(); restoreErr != nil {
			fmt.Printf("Warning: failed to change to root directory: %v\n", restoreErr)
		}
	}()

	return nil
}

func (pg *ProjectGenerator) createFrontendConfigFiles(config *WebProjectConfig) {
	if err := pg.CreateEnvFile(constants.FrontendDir, "."); err != nil {
		fmt.Printf("Warning: failed to create env file: %v\n", err)
	}

	if err := pg.CreateEnvConfig(".", config.FrontendFramework, config.UseTypeScript); err != nil {
		fmt.Printf("Warning: failed to create env config file: %v\n", err)
	}

	if config.UseDocker {
		if err := pg.CreateDockerfile(".", constants.FrontendDir, config.Runtime); err != nil {
			fmt.Printf("Warning: failed to create Docker files for frontend: %v\n", err)
		}
	}

	if err := pg.RemoveGitRepository("."); err != nil {
		fmt.Printf("Warning: failed to remove git repository from frontend: %v\n", err)
	}
}

func (pg *ProjectGenerator) createConfigFiles() error {
	if err := pg.CreateEnvFile(constants.APIDir, "."); err != nil {
		return fmt.Errorf("warning: failed to create env file in api: %v", err)
	}

	if err := pg.CreateGitignoreFile(constants.APIDir, "."); err != nil {
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

	if err := os.MkdirAll(cmdWebDir, 0750); err != nil {
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

	default:
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

	if err := os.WriteFile(mainGoPath, []byte(mainContent), 0600); err != nil {
		return err
	}

	if err := os.WriteFile(routesPath, []byte(routesContent), 0600); err != nil {
		return err
	}

	baseModuleName := pg.setModuleName(moduleName, projectName)
	apiModContent := fmt.Sprintf("module %s\n\ngo 1.21\n", baseModuleName)
	if err := os.WriteFile(goModPath, []byte(apiModContent), 0600); err != nil {
		return err
	}

	if err := pg.CreateEnvFile(constants.APIDir, baseDir); err != nil {
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
			if chdirErr := os.Chdir(originalDir); chdirErr != nil {
				return fmt.Errorf("failed to tidy go.mod and to change back to original directory: %w, %w", err, chdirErr)
			}
		}
		return fmt.Errorf("failed to tidy go.mod: %w", err)
	}

	if baseDir != "." {
		if err := os.Chdir(originalDir); err != nil {
			return fmt.Errorf("failed to change to original directory: %w", err)
		}
	}

	if err := pg.CreateAirFile(".", constants.WebTemplate); err != nil {
		fmt.Printf("Warning: failed to create .air.toml file: %v\n", err)
	}

	return nil
}
