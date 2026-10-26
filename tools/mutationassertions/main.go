// Command mutationassertions inventories assertion callsites without changing tests.
package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type site struct {
	File     string   `json:"file"`
	Line     int      `json:"line"`
	Kind     string   `json:"kind"`
	Function string   `json:"function"`
	Roots    []string `json:"roots"`
	SHA256   string   `json:"sha256"`
}
type function struct {
	sites []site
	calls []string
}

func testingNames(fields *ast.FieldList) map[string]bool {
	names := map[string]bool{}
	if fields == nil {
		return names
	}
	for _, f := range fields.List {
		typ := f.Type
		if p, ok := typ.(*ast.StarExpr); ok {
			typ = p.X
		}
		s, ok := typ.(*ast.SelectorExpr)
		if !ok {
			continue
		}
		pkg, ok := s.X.(*ast.Ident)
		if !ok || pkg.Name != "testing" || (s.Sel.Name != "T" && s.Sel.Name != "TB") {
			continue
		}
		for _, n := range f.Names {
			names[n.Name] = true
		}
	}
	return names
}
func main() {
	args := os.Args[1:]
	if len(args) > 0 && args[0] == "--" {
		args = args[1:]
	}
	if len(args) != 1 {
		panic("one package directory required")
	}
	directory := args[0]
	paths, err := filepath.Glob(filepath.Join(directory, "*_test.go"))
	if err != nil || len(paths) == 0 {
		panic("test files required")
	}
	functions := map[string]*function{}
	for _, path := range paths {
		raw, err := os.ReadFile(path)
		if err != nil {
			panic(err)
		}
		hash := sha256.Sum256(raw)
		digest := hex.EncodeToString(hash[:])
		fs := token.NewFileSet()
		file, err := parser.ParseFile(fs, path, raw, 0)
		if err != nil {
			panic(err)
		}
		for _, decl := range file.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Body == nil {
				continue
			}
			names := testingNames(fn.Type.Params)
			if len(names) == 0 {
				continue
			}
			f := &function{}
			functions[fn.Name.Name] = f
			// Test closures inherit the outer T and can declare a subtest T. Names are
			// used only to inventory source callsites, never to execute or rewrite AST.
			ast.Inspect(fn.Body, func(n ast.Node) bool {
				if lit, ok := n.(*ast.FuncLit); ok {
					for name := range testingNames(lit.Type.Params) {
						names[name] = true
					}
				}
				return true
			})
			ast.Inspect(fn.Body, func(n ast.Node) bool {
				call, ok := n.(*ast.CallExpr)
				if !ok {
					return true
				}
				pos := fs.Position(call.Pos())
				kind, target := "", ""
				if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
					if id, ok := sel.X.(*ast.Ident); ok && names[id.Name] {
						switch sel.Sel.Name {
						case "Fatal", "Fatalf", "Error", "Errorf":
							kind = sel.Sel.Name
						}
					}
					target = sel.Sel.Name
				} else if id, ok := call.Fun.(*ast.Ident); ok {
					target = id.Name
				}
				if kind == "" {
					for _, a := range call.Args {
						if id, ok := a.(*ast.Ident); ok && names[id.Name] {
							kind = "helper-call"
							f.calls = append(f.calls, target)
							break
						}
					}
				}
				if kind != "" {
					f.sites = append(f.sites, site{File: filepath.Base(path), Line: pos.Line, Kind: kind, Function: fn.Name.Name, SHA256: digest})
				}
				return true
			})
		}
	}
	roots := []string{}
	for name := range functions {
		if strings.HasPrefix(name, "Test") {
			roots = append(roots, name)
		}
	}
	sort.Strings(roots)
	all := map[string]*site{}
	for _, root := range roots {
		seen := map[string]bool{}
		var walk func(string)
		walk = func(name string) {
			if seen[name] {
				return
			}
			seen[name] = true
			f := functions[name]
			if f == nil {
				return
			}
			for _, s := range f.sites {
				key := s.File + ":" + strconv.Itoa(s.Line) + ":" + s.Kind
				if old := all[key]; old != nil {
					old.Roots = append(old.Roots, root)
				} else {
					copy := s
					copy.Roots = []string{root}
					all[key] = &copy
				}
			}
			for _, next := range f.calls {
				walk(next)
			}
		}
		walk(root)
	}
	out := []site{}
	for _, s := range all {
		out = append(out, *s)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].File != out[j].File {
			return out[i].File < out[j].File
		}
		if out[i].Line != out[j].Line {
			return out[i].Line < out[j].Line
		}
		return out[i].Kind < out[j].Kind
	})
	if err := json.NewEncoder(os.Stdout).Encode(out); err != nil {
		panic(err)
	}
}
