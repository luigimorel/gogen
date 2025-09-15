package gentstypes

import (
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"go/doc"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.hlmpn.dev/pkg/go-logger"
)

// packOut is a small struct carrying per-package output data.
	type packOut struct {
		pkgName string
		content []byte
	}

	// Generate is the entry point for generating TypeScript .d.ts types from Go types.
	// Both input and output are required. The caller (CLI) is responsible for setting default values.
	// - input: a go package import path (e.g., "github.com/user/pkg"), a local path ("./local/dir"),
	//          a single file ("/path/to/file.go"), or package.Type pattern similar to go doc ("fmt.Stringer").
	// - output: a .d.ts file path or a directory. If directory, files will be named "packagename.d.ts".
	// This function never exits or panics; it returns errors with clear, human-readable messages.
	func Generate(input string, output string) error {
		if input == "" {
			return errors.New("input is required, provide a go package path, local dir/file, or package.Type")
		}
		if output == "" {
			return errors.New("output is required, provide a .d.ts file or a directory")
		}

		// Determine if output is a single file or directory
		outLower := strings.ToLower(output)
		isSingleFile := strings.HasSuffix(outLower, ".d.ts")

		// Prepare output target
		// If dir: ensure it's empty; if not empty, remove it in the simplest way, then re-create.
		createdDir := false
		switch {
		case isSingleFile:
			dir := filepath.Dir(output)
			err := os.MkdirAll(dir, 0o755)
			if err != nil {
				return errors.New("failed to create output directory")
			}
		default:
			exists, notEmpty, statErr := checkDirStatus(output)
			if statErr != nil {
				return errors.New("failed to check output directory status")
			}
			switch {
			case exists && notEmpty:
				rmErr := os.RemoveAll(output)
				if rmErr != nil {
					return errors.New("failed to clean output directory")
				}
				mkErr := os.MkdirAll(output, 0o755)
				if mkErr != nil {
					return errors.New("failed to create output directory after cleaning")
				}
				createdDir = true
			case exists && !notEmpty:
				createdDir = true
			default:
				mkErr := os.MkdirAll(output, 0o755)
				if mkErr != nil {
					return errors.New("failed to create output directory")
				}
				createdDir = true
			}
		}

		// Load package(s) via loader
		loadRes, err := loadFromInput(input)
		if err != nil {
			if !isSingleFile && createdDir {
				_ = os.RemoveAll(output)
			}
			return err
		}
		if len(loadRes) == 0 {
			if !isSingleFile && createdDir {
				_ = os.RemoveAll(output)
			}
			return errors.New("no packages found for the given input")
		}

		// Generate for each package
		results := make([]packOut, 0, len(loadRes))
		for _, pkg := range loadRes {
			// Build doc.Package to get comments
			docPkg := doc.New(pkg.ASTPkg, pkg.ImportPath, doc.AllDecls)

			typeDocMap := buildTypeDocMap(docPkg)
			registry := buildTypeRegistry(pkg.ASTPkg)
			gCtx := &genContext{
				fset:        pkg.Fset,
				pkgName:     pkg.PkgName,
				importPath:  pkg.ImportPath,
				registry:    registry,
				typeDocMap:  typeDocMap,
				typeFilter:  pkg.TypeFilter, // optional
				seenInline:  map[string]bool{},
				seenTypes:   map[string]bool{},
				parentPkg:   pkg,
				pkgDoc:      docPkg,
				pkgsByName:  map[string]*ast.Package{pkg.PkgName: pkg.ASTPkg},
				importedPkg: map[string]string{},
			}

			// Build output buffer for the package
			buf := &bytes.Buffer{}
			// Header comment
			writeHeader(buf, pkg.PkgName, pkg.ImportPath)

			// Determine which types to output
			typeSpecs := collectTypeSpecs(pkg.ASTPkg)
			if len(typeSpecs) == 0 {
				logger.Warnf("no types found in package %s", pkg.PkgName)
			}

			// Filter by specific type if input includes package.Type
			filtered := typeSpecs
			if pkg.TypeFilter != "" {
				f := make(map[string]*ast.TypeSpec)
				for n, ts := range typeSpecs {
					if n == pkg.TypeFilter {
						f[n] = ts
					}
				}
				filtered = f
				if len(filtered) == 0 {
					logger.Warnf("no matching type named %s in package %s", pkg.TypeFilter, pkg.PkgName)
				}
			}

			// For deterministic order
			names := make([]string, 0, len(filtered))
			for name := range filtered {
				names = append(names, name)
			}
			sortStrings(names)

			// Emit types
			for _, name := range names {
				spec := filtered[name]

				// Type-level ignore via special comment "ts ignore //"
				ignore := hasTSIgnoreType(spec)
				switch {
				case ignore:
					continue
				default:
				}

				// Use GoDoc comments above the type if available
				docLines := gCtx.typeDocMap[name]

				// Determine struct vs alias
				switch t := spec.Type.(type) {
				case *ast.StructType:
					emitDocLines(buf, docLines)
					generateInterface(buf, gCtx, name, t)
					buf.WriteString("\n")
				default:
					emitDocLines(buf, docLines)
					tsType := gCtx.exprToTSType(t, fieldContext{
						jsonName:             "",
						isPtr:                false,
						forceNullable:        false,
						omitempty:            false,
						omitzero:             false,
						overrideTSType:       "",
						insideStruct:         false,
						tsTagHasNullable:     false,
						hasExplicitTSOrType:  false,
					})
					stmt := fmt.Sprintf("export type %s = %s;\n\n", name, tsType)
					buf.WriteString(stmt)
				}
			}

			results = append(results, packOut{pkgName: pkg.PkgName, content: buf.Bytes()})
		}

		// Write outputs
		switch {
		case isSingleFile:
			var all bytes.Buffer
			for i, r := range results {
				if i > 0 {
					all.WriteString("\n")
				}
				// Add a package section header
				all.WriteString("// =====================================================\n")
				all.WriteString("// Package: " + r.pkgName + "\n")
				all.WriteString("// =====================================================\n\n")
				all.Write(r.content)
			}

			err := os.WriteFile(output, all.Bytes(), 0o644)
			if err != nil {
				return errors.New("failed to write output file")
			}
			return nil
		default:
			writeErr := writePackagesToDir(output, results)
			if writeErr != nil {
				_ = os.RemoveAll(output)
				return writeErr
			}
			return nil
		}
	}

	// checkDirStatus checks whether a directory exists and whether it is empty.
	func checkDirStatus(path string) (exists bool, notEmpty bool, err error) {
		info, statErr := os.Stat(path)
		if statErr != nil {
			switch {
			case os.IsNotExist(statErr):
				return false, false, nil
			default:
				return false, false, statErr
			}
		}
		switch {
		case !info.IsDir():
			return false, false, errors.New("output path exists but is not a directory")
		default:
		}
		entries, readErr := os.ReadDir(path)
		if readErr != nil {
			return true, false, readErr
		}
		return true, len(entries) > 0, nil
	}

	// writePackagesToDir writes each package content into a file named "packagename.d.ts".
	func writePackagesToDir(dir string, results []packOut) error {
		for _, r := range results {
			filename := filepath.Join(dir, r.pkgName+".d.ts")
			err := os.WriteFile(filename, r.content, 0o644)
			if err != nil {
				return errors.New("failed to write package file")
			}
		}
		return nil
	}

	// writeHeader emits a standard header for generated content.
	func writeHeader(buf *bytes.Buffer, pkgName, importPath string) {
		buf.WriteString("// Code generated by go-ts-types. DO NOT EDIT.\n")
		buf.WriteString("// Package: " + pkgName + "\n")
		if importPath != "" {
			buf.WriteString("// Import Path: " + importPath + "\n")
		}
		buf.WriteString("\n")
	}

	func emitDocLines(buf *bytes.Buffer, docLines []string) {
		if len(docLines) == 0 {
			return
		}
		for _, l := range docLines {
			if l == "" {
				buf.WriteString("//\n")
				continue
			}
			buf.WriteString("// " + l + "\n")
		}
	}

	// hasTSIgnoreType detects the exact "ts ignore //" directive for type declarations.
	// It must be the last line before the "type Name T" spec and exactly equal to "ts ignore //".
	func hasTSIgnoreType(spec *ast.TypeSpec) bool {
		// In Go AST, doc comments can be attached to GenDecl.Doc or TypeSpec.Doc.
		texts := make([]string, 0, 4)
		if spec.Doc != nil {
			for _, c := range spec.Doc.List {
				line := strings.TrimSpace(strings.TrimPrefix(c.Text, "//"))
				// For block comments the parser includes /* */; normalize to lines.
				if strings.HasPrefix(c.Text, "/*") {
					lines := strings.Split(strings.Trim(c.Text, "/*"), "\n")
					for _, ln := range lines {
						t := strings.TrimSpace(strings.TrimPrefix(strings.TrimPrefix(ln, "//"), "*"))
						if t != "" {
							texts = append(texts, t)
						}
					}
					continue
				}
				texts = append(texts, line)
			}
		}
		if len(texts) == 0 {
			return false
		}
		last := texts[len(texts)-1]
		switch {
		case last == "ts ignore //":
			return true
		default:
			return false
		}
	}

	// buildTypeDocMap extracts doc comments for each type via go/doc.
	func buildTypeDocMap(pkg *doc.Package) map[string][]string {
		m := make(map[string][]string)
		if pkg == nil {
			return m
		}
		for _, t := range pkg.Types {
			name := t.Name
			docText := strings.TrimSpace(t.Doc)
			if docText == "" {
				continue
			}
			lines := strings.Split(docText, "\n")
			m[name] = lines
		}
		return m
	}

	// buildTypeRegistry constructs a registry of type specs in the package.
	func buildTypeRegistry(astPkg *ast.Package) map[string]*ast.TypeSpec {
		reg := make(map[string]*ast.TypeSpec)
		if astPkg == nil {
			return reg
		}
		for _, f := range astPkg.Files {
			for _, decl := range f.Decls {
				gd, ok := decl.(*ast.GenDecl)
				if !ok {
					continue
				}
				if gd.Tok != token.TYPE {
					continue
				}
				for _, spec := range gd.Specs {
					ts, ok := spec.(*ast.TypeSpec)
					if !ok {
						continue
					}
					reg[ts.Name.Name] = ts
				}
			}
		}
		return reg
	}

	// collectTypeSpecs collects all type specs in the package by name.
	func collectTypeSpecs(astPkg *ast.Package) map[string]*ast.TypeSpec {
		return buildTypeRegistry(astPkg)
	}

	// sortStrings sorts a list of strings.
	func sortStrings(ss []string) {
		if len(ss) < 2 {
			return
		}
		sort.Strings(ss)
	}