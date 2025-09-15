package stdlibtemplate_test

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"testing"

	"github.com/luigimorel/gogen/internal/stdlibtemplate"
)

func init() {
	// Find project root by looking for go.mod
	for {
		_, err := os.Stat("go.mod")
		if os.IsNotExist(err) {
			// go.mod not found, try parent directory
			if err := os.Chdir(".."); err != nil {
				log.Fatalf("init: failed to change directory: %v", err)
			}
			// Check if we are at the root of the filesystem
			cwd, _ := os.Getwd()
			if cwd == "/" || cwd == filepath.Dir(cwd) {
				log.Fatalf("init: go.mod not found in any parent directory")
			}
			continue
		}
		if err != nil {
			log.Fatalf("init: error checking for go.mod: %v", err)
		}

		return
	}
}

func TestCreateRouterSetup(t *testing.T) {
	// Create a temporary directory for the test to run in.
	// This is safer than creating files in the project root.
	tempDir := t.TempDir()

	// Get the current working directory to return to it after the test.
	originalWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current working directory: %v", err)
	}
	// Change to the temporary directory.
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("failed to change to temp directory: %v", err)
	}
	// Use t.Cleanup to ensure we change back to the original directory.
	t.Cleanup(func() {
		if err := os.Chdir(originalWD); err != nil {
			t.Errorf("failed to change back to original directory: %v", err)
		}
	})

	// Run the function to be tested. This will create the 'router' dir inside the tempDir.
	if err := stdlibtemplate.CreateRouterSetup(); err != nil {
		t.Fatalf("CreateRouterSetup() failed: %v", err)
	}

	// Define the expected file structure based on the provided tree.
	expectedFiles := []string{
		"router/handler.go",
		"router/ocsp.go",
		"router/redirects.go",
		"router/router.go",
		"router/serve.go",
	}
	sort.Strings(expectedFiles)

	var actualFiles []string
	err = filepath.Walk("router", func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			// Normalize path separators for consistent comparison.
			actualFiles = append(actualFiles, filepath.ToSlash(path))
		}
		return nil
	})

	if err != nil {
		t.Fatalf("failed to walk created 'router' directory: %v", err)
	}
	sort.Strings(actualFiles)

	// Compare the actual file list with the expected file list.
	if !reflect.DeepEqual(expectedFiles, actualFiles) {
		t.Errorf(`file structure mismatch:
Expected: %v
Actual:   %v`, expectedFiles, actualFiles)
	}
}
