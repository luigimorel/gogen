package gentstypes

import (
	"bytes"
	"go/ast"
	"go/doc"
	"go/token"
	"sort"
	"strings"
)

// This file primarily hosts additional helpers. No new public API here.

// The following ensures imports are referenced to satisfy the compiler for
// packages imported in other files without nested imports.
var (
	_ = bytes.NewBuffer
	_ = token.ILLEGAL
	_ = doc.New
	_ = ast.File{}
	_ = sort.Strings
	_ = strings.TrimSpace
)
