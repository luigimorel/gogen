package internal

import (
	"fmt"
	"os"
	"os/exec"
)

type TailwindConfig struct {
	Framework     string
	Runtime       string
	DirName       string
	UseTypeScript bool
}

func NewTailwindConfig(framework, runtime, dirName string) *TailwindConfig {
	useTypeScript := false
	if _, err := os.Stat(dirName + "/tsconfig.json"); err == nil {
		useTypeScript = true
	}

	return &TailwindConfig{
		Framework:     framework,
		Runtime:       runtime,
		DirName:       dirName,
		UseTypeScript: useTypeScript,
	}
}

func (tc *TailwindConfig) InstallTailwindCSS() error {
	originalDir, _ := os.Getwd()
	defer func() {
		if err := os.Chdir(originalDir); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to return to original directory: %v\n", err)
		}
	}()

	if err := os.Chdir(tc.DirName); err != nil {
		return fmt.Errorf("failed to change to frontend directory: %w", err)
	}

	fmt.Println("Installing Tailwind CSS...")

	if err := tc.tailwindLibInstall(tc.Framework, tc.Runtime); err != nil {
		return fmt.Errorf("failed to install Tailwind packages: %w", err)
	}

	if err := tc.updateConfigFile(tc.Framework); err != nil {
		return fmt.Errorf("failed to update config files: %w", err)
	}

	if err := tc.updateStylesFile(tc.Framework); err != nil {
		return fmt.Errorf("failed to update style sheets file: %w", err)
	}

	fmt.Println("✅ Tailwind CSS configured successfully!")
	return nil
}
func (tc *TailwindConfig) tailwindLibInstall(framework, runtime string) error {
	if runtime == "node" {
		runtime = "npm"
	}

	switch framework {
	case "react", "vue", "svelte", "solidjs":
		args := append([]string{"add"}, "tailwindcss", "@tailwindcss/vite")
		return exec.Command(runtime, args...).Run()

	case "angular":
		if runtime == "bun" {
			return exec.Command(runtime, "add", "tailwindcss", "@tailwindcss/postcss", "postcss", "--force").Run()
		}
		return exec.Command("npm", "install", "tailwindcss", "@tailwindcss/postcss", "postcss", "--force").Run()

	default:
		return fmt.Errorf("unsupported framework: %s", framework)
	}
}

func (tc *TailwindConfig) updateConfigFile(framework string) error {
	configExt := ".js"
	if tc.UseTypeScript && framework != "angular" {
		configExt = ".ts"
	}

	switch framework {
	case "react":
		viteConfig := `import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import tailwindcss from '@tailwindcss/vite'

export default defineConfig({
  plugins: [
  react(),
  tailwindcss(),
  ],
})
`
		return os.WriteFile("vite.config"+configExt, []byte(viteConfig), 0644)
	case "vue":
		viteConfig := `import { fileURLToPath, URL } from "node:url";

import vue from "@vitejs/plugin-vue";
import vueJsx from "@vitejs/plugin-vue-jsx";
import { defineConfig } from "vite";
import vueDevTools from "vite-plugin-vue-devtools";
import tailwindcss from "@tailwindcss/vite";

// https://vite.dev/config/
export default defineConfig({
  plugins: [vue(), vueJsx(), vueDevTools(), tailwindcss()],
  resolve: {
    alias: {
      "@": fileURLToPath(new URL("./src", import.meta.url)),
    },
  },
});

`
		return os.WriteFile("vite.config"+configExt, []byte(viteConfig), 0644)
	case "svelte":
		viteConfig := `import tailwindcss from '@tailwindcss/vite';
import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';

export default defineConfig({
	plugins: [tailwindcss(), sveltekit()]
});
`
		return os.WriteFile("vite.config"+configExt, []byte(viteConfig), 0644)
	case "solidjs":
		viteConfig := `import { defineConfig } from 'vite';
import solidPlugin from 'vite-plugin-solid';
import tailwindcss from '@tailwindcss/vite';

export default defineConfig({
  plugins: [
	tailwindcss(),
	solidPlugin(),
  ],
  server: {
	port: 3000,
  },
  build: {
	target: 'esnext',
  },
});
`
		return os.WriteFile("vite.config"+configExt, []byte(viteConfig), 0644)
	case "angular":
		postcssConfig := `{  "plugins": {    "@tailwindcss/postcss": {}  }}`
		return os.WriteFile(".postcssrc.json", []byte(postcssConfig), 0644)
	}

	return nil
}

func (tc *TailwindConfig) updateStylesFile(framework string) error {
	var cssFile string
	switch framework {
	case "react":
		cssFile = "src/index.css"
	case "vue":
		cssFile = "src/assets/main.css"
	case "svelte":
		cssFile = "src/main.css"
	case "solidjs":
		cssFile = "src/index.css"
	case "angular":
		cssFile = "src/styles.css"
	default:
		return nil
	}

	existingContent, err := os.ReadFile(cssFile)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to read CSS file: %w", err)
	}

	tailwindImport := "@import \"tailwindcss\";\n\n"

	var newContent []byte
	if len(existingContent) > 0 {
		newContent = []byte(tailwindImport + string(existingContent))
	} else {
		newContent = []byte(tailwindImport)
	}

	return os.WriteFile(cssFile, newContent, 0644)
}
