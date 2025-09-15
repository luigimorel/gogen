	package gentstypes

	import (
		"errors"
		"go/ast"
		"go/parser"
		"go/token"
		"os"
		"path/filepath"
		"strings"

		"golang.org/x/tools/go/packages"
	)

	// pkgData represents one loaded Go package with AST and metadata.
	type pkgData struct {
		PkgName    string
		ImportPath string
		Fset       *token.FileSet
		ASTPkg     *ast.Package
		TypeFilter string
	}

	// loadFromInput loads and parses Go package(s) for a given input string.
	// It supports:
	// - a single .go file: parses its directory and builds the package from all non-test files in the same package.
	// - a local directory path: parses as a package.
	// - an import path: uses go/packages to locate, then parses files with go/parser.
	// - an import path with ".Type" at the end: loads package and sets TypeFilter to that type name.
	// Only one primary package is loaded and returned as a single-element slice.
	func loadFromInput(input string) ([]pkgData, error) {
		// Detect file input
		if strings.HasSuffix(strings.ToLower(input), ".go") {
			return loadFromFile(input)
		}

		// Detect probable local path (absolute or relative)
		if isLikelyPath(input) {
			return loadFromDir(input)
		}

		// Otherwise treat as import path or pkg.Type
		return loadFromImportOrType(input)
	}

	func isLikelyPath(p string) bool {
		switch {
		case strings.HasPrefix(p, "./"):
			return true
		case strings.HasPrefix(p, "../"):
			return true
		case strings.HasPrefix(p, "/"):
			return true
		case strings.Contains(p, string(os.PathSeparator)):
			return true
		default:
			return false
		}
	}

	func loadFromFile(path string) ([]pkgData, error) {
		abs, err := filepath.Abs(path)
		if err != nil {
			return nil, errors.New("failed to resolve file path")
		}
		info, statErr := os.Stat(abs)
		if statErr != nil || info.IsDir() {
			return nil, errors.New("file path is invalid or points to a directory")
		}
		dir := filepath.Dir(abs)
		return loadFromDir(dir)
	}

	func loadFromDir(dir string) ([]pkgData, error) {
		abs, err := filepath.Abs(dir)
		if err != nil {
			return nil, errors.New("failed to resolve directory path")
		}
		info, statErr := os.Stat(abs)
		if statErr != nil || !info.IsDir() {
			return nil, errors.New("directory path is invalid")
		}

		// Parse all files in the directory with go/parser
		fset := token.NewFileSet()
		pkgs, err := parser.ParseDir(fset, abs, func(fi os.FileInfo) bool {
			name := fi.Name()
			// skip test files
			if strings.HasSuffix(name, "_test.go") {
				return false
			}
			return strings.HasSuffix(name, ".go")
		}, parser.ParseComments)
		if err != nil {
			return nil, errors.New("failed to parse directory")
		}
		if len(pkgs) == 0 {
			return nil, errors.New("no go packages found in the directory")
		}

		// Select primary package (ignore *_test package names)
		var chosen *ast.Package
		var pkgName string
		for name, p := range pkgs {
			if strings.HasSuffix(name, "_test") {
				continue
			}
			chosen = p
			pkgName = name
			break
		}
		if chosen == nil {
			for name, p := range pkgs {
				chosen = p
				pkgName = name
				break
			}
		}

		return []pkgData{{
			PkgName:    pkgName,
			ImportPath: abs, // best-effort for local; not module path
			Fset:       fset,
			ASTPkg:     chosen,
			TypeFilter: "",
		}}, nil
	}

	func loadFromImportOrType(path string) ([]pkgData, error) {
		// Support pkg.Type filter
		imp, typeName := splitImportType(path)

		cfg := &packages.Config{
			Mode: packages.NeedName |
				packages.NeedFiles |
				packages.NeedCompiledGoFiles |
				packages.NeedModule |
				packages.NeedSyntax,
		}
		pkgs, err := packages.Load(cfg, imp)
		if err != nil {
			return nil, errors.New("failed to load package by import path")
		}
		if packages.PrintErrors(pkgs) > 0 {
			return nil, errors.New("package contains build or load errors")
		}
		if len(pkgs) == 0 {
			return nil, errors.New("no packages matched the import path")
		}

		p := pkgs[0]
		// Re-parse to ensure we have comments collected via parser.ParseComments (packages.Syntax is already parsed but comments mode typically enabled; we keep usage consistent)
		fset := token.NewFileSet()
		astFiles := make(map[string]*ast.File)
		for _, fname := range p.GoFiles {
			file, perr := parser.ParseFile(fset, fname, nil, parser.ParseComments)
			if perr != nil {
				return nil, errors.New("failed to parse package files")
			}
			astFiles[fname] = file
		}
		astPkg := &ast.Package{
			Name:  p.Name,
			Files: astFiles,
		}
		return []pkgData{{
			PkgName:    p.Name,
			ImportPath: p.PkgPath,
			Fset:       fset,
			ASTPkg:     astPkg,
			TypeFilter: typeName,
		}}, nil
	}

	func splitImportType(s string) (importPath string, typeName string) {
		// If there is a dot after the last slash, interpret as pkg.Type
		lastSlash := strings.LastIndex(s, "/")
		lastDot := strings.LastIndex(s, ".")
		switch {
		case lastDot > lastSlash && lastDot != -1:
			return s[:lastDot], s[lastDot+1:]
		default:
			return s, ""
		}
	}