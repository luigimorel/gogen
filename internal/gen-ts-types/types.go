package gentstypes

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/doc"
	"go/token"
	"reflect"
	"strings"

	"gopkg.hlmpn.dev/pkg/go-logger"
)

// We centralize number and basic type mapping here.
var basicTypeMap = map[string]string{
	"string":     "string",
	"bool":       "boolean",
	"byte":       "number",
	"rune":       "number",
	"int":        "number",
	"int8":       "number",
	"int16":      "number",
	"int32":      "number",
	"int64":      "number",
	"uint":       "number",
	"uint8":      "number",
	"uint16":     "number",
	"uint32":     "number",
	"uint64":     "number",
	"uintptr":    "number",
	"float32":    "number",
	"float64":    "number",
	"complex64":  "number",
	"complex128": "number",
	"error":      "string", // often marshaled as string; safer than any
	"any":        "any",
}

// genContext carries state for generating a single package's TS output.
type genContext struct {
	fset        *token.FileSet
	pkgName     string
	importPath  string
	registry    map[string]*ast.TypeSpec
	typeDocMap  map[string][]string
	typeFilter  string
	seenInline  map[string]bool
	seenTypes   map[string]bool
	parentPkg   pkgData
	pkgDoc      *doc.Package
	pkgsByName  map[string]*ast.Package
	importedPkg map[string]string // alias -> full import path (best-effort)
}

// fieldContext captures field-specific options for TS type generation.
type fieldContext struct {
	jsonName            string
	isPtr               bool
	forceNullable       bool
	omitempty           bool
	omitzero            bool
	overrideTSType      string
	insideStruct        bool
	tsTagHasNullable    bool
	hasExplicitTSOrType bool
}

// generateInterface emits "export interface Name { ... }" for a struct.
func generateInterface(buf *bytes.Buffer, g *genContext, name string, st *ast.StructType) {
	fields := collectStructFields(g, st)
	buf.WriteString(fmt.Sprintf("export interface %s {\n", name))
	for _, f := range fields {
		buf.WriteString("  ")
		buf.WriteString(f)
		buf.WriteString("\n")
	}
	buf.WriteString("}\n")
}

// collectStructFields returns a slice of "prop?: T" or "prop: T | null" lines for a struct type.
func collectStructFields(g *genContext, st *ast.StructType) []string {
	out := []string{}
	if st.Fields == nil || len(st.Fields.List) == 0 {
		return out
	}
	for _, fld := range st.Fields.List {
		// Embedded or named?
		switch {
		case len(fld.Names) == 0:
			// Embedded field
			embeddedLines := expandEmbeddedField(g, fld)
			switch {
			case len(embeddedLines) == 0:
				continue
			default:
				out = append(out, embeddedLines...)
				continue
			}
		default:
			// Named fields can have multiple names sharing same type: handle each
			for _, nameIdent := range fld.Names {
				line := generateFieldLine(g, nameIdent.Name, fld)
				if line == "" {
					continue
				}
				out = append(out, line)
			}
		}
	}
	return out
}

// expandEmbeddedField tries to flatten embedded fields.
// We flatten when:
// - anonymous struct type: inline its fields
// - named type defined in the same package and is a struct: inline its fields
func expandEmbeddedField(g *genContext, fld *ast.Field) []string {
	t := fld.Type
	switch tt := t.(type) {
	case *ast.Ident:
		// Same-package type?
		spec := g.registry[tt.Name]
		switch {
		case spec == nil:
			return nil
		default:
			st, ok := spec.Type.(*ast.StructType)
			switch {
			case !ok:
				return nil
			default:
				return collectStructFields(g, st)
			}
		}
	case *ast.StarExpr:
		// Pointer to ident or struct
		switch ut := tt.X.(type) {
		case *ast.Ident:
			spec := g.registry[ut.Name]
			switch {
			case spec == nil:
				return nil
			default:
				st, ok := spec.Type.(*ast.StructType)
				switch {
				case !ok:
					return nil
				default:
					return collectStructFields(g, st)
				}
			}
		case *ast.StructType:
			return collectStructFields(g, ut)
		default:
			return nil
		}
	case *ast.StructType:
		return collectStructFields(g, tt)
	default:
		return nil
	}
}

// generateFieldLine generates one TS property line for a single named field.
func generateFieldLine(g *genContext, goFieldName string, fld *ast.Field) string {
	// Field-level comment "ts ignore //" detection with logger warning rules
	hasIgnoreComment := hasTSIgnoreOnField(fld)
	tags := parseStructTag(fld)
	hasTsTag := tags.ts != ""
	hasTsTypeTag := tags.tsType != ""

	switch {
	case hasIgnoreComment && (hasTsTag || hasTsTypeTag):
		logger.Warnf("field %s has a 'ts ignore //' comment and ts or ts_type tag; ignoring the comment", goFieldName)
	case hasIgnoreComment:
		return ""
	default:
	}

	// Reflect tags
	if tags.ignore {
		return ""
	}
	// JSON tag "-" is ignore
	if tags.jsonName == "-" {
		return ""
	}

	// Determine property name (from json tag or Go field)
	propName := tags.jsonName
	switch {
	case propName == "":
		propName = goFieldName
	default:
	}

	// Determine type characteristics
	isPtr := isPtrType(fld.Type)
	isSlice := isSliceType(fld.Type)
	isMap := isMapType(fld.Type)

	// Optional/nullable decision
	// Priority: ts tag overrides json. Follow the detailed rules mentioned.
	optional := false
	nullable := false

	// Start with base derived from kind
	if isPtr {
		// pointer default
		optional = true
		nullable = true
	}
	if isSlice || isMap {
		// slices and maps can be nil => nullable by default
		nullable = true
	}

	// Apply json omitempty/omitzero -> optional
	if tags.omitempty || tags.omitzero {
		optional = true
	}

	// Apply ts:"nullable" => required and nullable for non-pointer/non-omitempty as per rules
	if tags.tsNullable {
		nullable = true
	}

	// Compose type string
	ctx := fieldContext{
		jsonName:            propName,
		isPtr:               isPtr,
		forceNullable:       nullable,
		omitempty:           tags.omitempty,
		omitzero:            tags.omitzero,
		overrideTSType:      tags.tsType,
		insideStruct:        true,
		tsTagHasNullable:    tags.tsNullable,
		hasExplicitTSOrType: hasTsTag || hasTsTypeTag,
	}
	tsType := g.exprToTSType(fld.Type, ctx)

	// Property line
	qs := ""
	if optional {
		qs = "?"
	}
	// Append null explicitly if nullable
	if nullable && !strings.Contains(tsType, "| null") {
		tsType = tsType + " | null"
	}

	return fmt.Sprintf("%s%s: %s;", propName, qs, tsType)
}

type parsedTags struct {
	jsonName  string
	omitempty bool
	omitzero  bool

	// ts:"ignore" or ts:"nullable"
	ts          string
	ignore      bool
	tsNullable  bool
	tsType      string
	otherRawTag string
}

// parseStructTag extracts json, ts, ts_type directives from field tag.
func parseStructTag(fld *ast.Field) parsedTags {
	var pt parsedTags
	if fld.Tag == nil {
		return pt
	}
	// BasicLit.Value is quoted with backticks or double quotes. Strip them.
	raw := fld.Tag.Value
	trim := strings.Trim(raw, "`\"")
	st := reflect.StructTag(trim)

	jsonTag := st.Get("json")
	tsTag := st.Get("ts")
	tsType := st.Get("ts_type")

	// json tag parse
	if jsonTag != "" {
		parts := strings.Split(jsonTag, ",")
		switch {
		case len(parts) > 0:
			pt.jsonName = parts[0]
		default:
		}
		for _, p := range parts[1:] {
			switch p {
			case "omitempty":
				pt.omitempty = true
			case "omitzero": // Go 1.23 addition
				pt.omitzero = true
			}
		}
	}

	// ts tag parse
	if tsTag != "" {
		pt.ts = tsTag
		if tsTag == "ignore" {
			pt.ignore = true
		}
		if tsTag == "nullable" {
			pt.tsNullable = true
		}
	}

	// ts_type override
	if tsType != "" {
		pt.tsType = strings.TrimSpace(tsType)
	}

	pt.otherRawTag = trim
	return pt
}

// hasTSIgnoreOnField detects a field-level "ts ignore //" comment (exact) on either Doc or Comment.
func hasTSIgnoreOnField(fld *ast.Field) bool {
	lines := []string{}

	if fld.Doc != nil {
		for _, c := range fld.Doc.List {
			lines = append(lines, normalizeComment(c.Text)...)
		}
	}
	if fld.Comment != nil {
		for _, c := range fld.Comment.List {
			lines = append(lines, normalizeComment(c.Text)...)
		}
	}

	if len(lines) == 0 {
		return false
	}
	last := lines[len(lines)-1]
	switch {
	case last == "ts ignore //":
		return true
	default:
		return false
	}
}

func normalizeComment(text string) []string {
	out := []string{}
	switch {
	case strings.HasPrefix(text, "/*"):
		trim := strings.Trim(text, "/*")
		parts := strings.Split(trim, "\n")
		for _, p := range parts {
			p = strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(p), "*"))
			if p != "" {
				out = append(out, p)
			}
		}
	default:
		line := strings.TrimSpace(strings.TrimPrefix(text, "//"))
		out = append(out, line)
	}
	return out
}

func isPtrType(expr ast.Expr) bool {
	_, ok := expr.(*ast.StarExpr)
	return ok
}

func isSliceType(expr ast.Expr) bool {
	_, ok := expr.(*ast.ArrayType)
	return ok
}

func isMapType(expr ast.Expr) bool {
	_, ok := expr.(*ast.MapType)
	return ok
}

// exprToTSType converts an AST expression to a TypeScript type string.
// The ctx provides field-level hints (tags, optional/nullable overrides).
func (g *genContext) exprToTSType(expr ast.Expr, ctx fieldContext) string {
	switch t := expr.(type) {
	case *ast.Ident:
		// Built-in or named type
		name := t.Name
		// Special-cases
		switch name {
		case "any":
			return "any"
		default:
		}
		// Basic type mapping
		if ts, ok := basicTypeMap[name]; ok {
			return ts
		}
		// Named type defined in same package
		if g.registry[name] != nil {
			return name
		}
		// Known special aliases
		switch name {
		case "Time":
			// time.Time would be SelectorExpr usually. If ident used here, we don't know; leave as any.
			return "any"
		default:
			return "any"
		}
	case *ast.ArrayType:
		// Slice or array
		elt := t.Elt
		// Special case: []byte
		if isIdentByte(elt) {
			// []byte as base64 string by default; allow ts_type override
			switch {
			case ctx.overrideTSType != "":
				return tsOverrideToType(ctx.overrideTSType)
			default:
				return "string"
			}
		}
		eltTs := g.exprToTSType(elt, fieldContext{
			jsonName:            ctx.jsonName,
			isPtr:               false,
			forceNullable:       false,
			omitempty:           ctx.omitempty,
			omitzero:            ctx.omitzero,
			overrideTSType:      "",
			insideStruct:        ctx.insideStruct,
			tsTagHasNullable:    ctx.tsTagHasNullable,
			hasExplicitTSOrType: ctx.hasExplicitTSOrType,
		})
		return fmt.Sprintf("%s[]", wrapIfUnion(eltTs))
	case *ast.StarExpr:
		// Pointer: unwrap and map
		ts := g.exprToTSType(t.X, fieldContext{
			jsonName:            ctx.jsonName,
			isPtr:               true,
			forceNullable:       ctx.forceNullable,
			omitempty:           ctx.omitempty,
			omitzero:            ctx.omitzero,
			overrideTSType:      ctx.overrideTSType,
			insideStruct:        ctx.insideStruct,
			tsTagHasNullable:    ctx.tsTagHasNullable,
			hasExplicitTSOrType: ctx.hasExplicitTSOrType,
		})
		return ts
	case *ast.MapType:
		// map[K]V => Record<Kts, Vts>
		keyTs := g.mapKeyToTS(t.Key)
		valTs := g.exprToTSType(t.Value, fieldContext{
			jsonName:            ctx.jsonName,
			isPtr:               false,
			forceNullable:       false,
			omitempty:           ctx.omitempty,
			omitzero:            ctx.omitzero,
			overrideTSType:      "",
			insideStruct:        ctx.insideStruct,
			tsTagHasNullable:    ctx.tsTagHasNullable,
			hasExplicitTSOrType: ctx.hasExplicitTSOrType,
		})
		return fmt.Sprintf("Record<%s, %s>", keyTs, valTs)
	case *ast.StructType:
		// Anonymous inline struct
		fields := collectStructFields(g, t)
		if len(fields) == 0 {
			return "{}"
		}
		var b strings.Builder
		b.WriteString("{ ")
		for i, line := range fields {
			// line is "prop?: T;" - we need to inline without leading indentation
			b.WriteString(strings.TrimSpace(line))
			if i < len(fields)-1 {
				b.WriteString(" ")
			}
		}
		b.WriteString(" }")
		return b.String()
	case *ast.SelectorExpr:
		// Qualified identifier: pkg.Type
		pkgIdent, _ := t.X.(*ast.Ident)
		sel := t.Sel.Name
		pkgName := ""
		if pkgIdent != nil {
			pkgName = pkgIdent.Name
		}
		_ = pkgName
		qname := sel
		if pkgName != "" {
			qname = pkgName + "." + sel
		}

		// Special-cases
		switch qname {
		case "json.RawMessage", "encoding/json.RawMessage":
			if ctx.overrideTSType == "array" {
				return "any[]"
			}
			if ctx.overrideTSType == "object" {
				return "Record<string, any>"
			}
			return "any"
		case "time.Time":
			return "string"
		default:
			// External types: map to any
			return "any"
		}
	case *ast.InterfaceType:
		return "any"
	case *ast.FuncType:
		return "any"
	case *ast.ChanType:
		return "any"
	case *ast.Ellipsis:
		eltTs := g.exprToTSType(t.Elt, fieldContext{})
		return fmt.Sprintf("%s[]", wrapIfUnion(eltTs))
	case *ast.IndexExpr:
		// generic type usage T[U] => map to 'any' for simplicity
		return "any"
	case *ast.IndexListExpr:
		return "any"
	default:
		return "any"
	}
}

func wrapIfUnion(ts string) string {
	if strings.Contains(ts, " ") || strings.Contains(ts, "|") || strings.HasPrefix(ts, "{") {
		return "(" + ts + ")"
	}
	return ts
}

func isIdentByte(expr ast.Expr) bool {
	id, ok := expr.(*ast.Ident)
	if !ok {
		return false
	}
	return id.Name == "byte" || id.Name == "uint8"
}

func tsOverrideToType(v string) string {
	switch v {
	case "string":
		return "string"
	case "number":
		return "number"
	case "boolean":
		return "boolean"
	case "object":
		return "Record<string, any>"
	case "array":
		return "any[]"
	case "any":
		return "any"
	case "Uint8Array":
		return "Uint8Array"
	case "ArrayBuffer":
		return "ArrayBuffer"
	default:
		return "any"
	}
}

func (g *genContext) mapKeyToTS(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		name := t.Name
		switch name {
		case "string":
			return "string"
		case "bool":
			// TS index signatures do not allow boolean keys, but instruction requires to use T as keys.
			// We still return 'string' to ensure valid TS, as JSON keys are strings.
			return "string"
		default:
			if _, ok := basicTypeMap[name]; ok {
				// number-like keys become number
				switch name {
				case "string":
					return "string"
				default:
					return "number"
				}
			}
			// Unknown key type => string
			return "string"
		}
	case *ast.StarExpr:
		// pointers cannot be map keys in Go; fallback
		return "string"
	case *ast.SelectorExpr:
		// qualified names used as map keys are often handled by encoding as string
		return "string"
	default:
		return "string"
	}
}

// Additional helper: safe string representation for AST expr (debug only)
func _debugExpr(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.ArrayType:
		return "[]" + _debugExpr(t.Elt)
	case *ast.StarExpr:
		return "*" + _debugExpr(t.X)
	case *ast.MapType:
		return "map[" + _debugExpr(t.Key) + "]" + _debugExpr(t.Value)
	case *ast.SelectorExpr:
		return _debugExpr(t.X) + "." + t.Sel.Name
	case *ast.StructType:
		return "struct{...}"
	default:
		return fmt.Sprintf("%T", expr)
	}
}
