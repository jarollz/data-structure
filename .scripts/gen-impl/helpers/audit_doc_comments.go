package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func main() {
	folder := flag.String("folder", "", "target data structure folder relative to repo root")
	flag.Parse()

	if *folder == "" {
		fmt.Fprintln(os.Stderr, "missing --folder")
		os.Exit(1)
	}

	apiMethods, err := loadAPIMethods(*folder)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	var failures []string
	paths, err := filepath.Glob(filepath.Join(*folder, "*.go"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "glob %s: %v\n", *folder, err)
		os.Exit(1)
	}
	sort.Strings(paths)

	for _, path := range paths {
		base := filepath.Base(path)
		switch {
		case base == "api_contract.go":
			continue
		case strings.HasSuffix(base, "_test.go"):
			continue
		case strings.HasSuffix(base, "_bench_test.go"):
			continue
		case base == "helpers_test.go":
			continue
		case base == "bench_policy_test.go":
			continue
		}

		fileSet := token.NewFileSet()
		file, err := parser.ParseFile(fileSet, path, nil, parser.ParseComments)
		if err != nil {
			failures = append(failures, fmt.Sprintf("%s: parse failed: %v", path, err))
			continue
		}

		for _, decl := range file.Decls {
			funcDecl, ok := decl.(*ast.FuncDecl)
			if !ok {
				continue
			}
			if !funcDecl.Name.IsExported() {
				continue
			}

			name := funcDecl.Name.Name
			docText := strings.TrimSpace(funcDecl.Doc.Text())
			if docText == "" {
				failures = append(failures, fmt.Sprintf("%s: %s is missing a doc comment", path, name))
				continue
			}

			firstLine := firstNonEmptyLine(docText)
			if firstLine == "" {
				failures = append(failures, fmt.Sprintf("%s: %s doc comment is empty", path, name))
				continue
			}

			if funcDecl.Recv != nil && apiMethods[name] {
				want := name + " implements the API interface."
				if firstLine != want {
					failures = append(failures, fmt.Sprintf("%s: %s first doc line = %q, want %q", path, name, firstLine, want))
				}
			} else {
				if firstLine == name+" implements the API interface." {
					failures = append(failures, fmt.Sprintf("%s: %s must not claim to implement the API interface", path, name))
				}
				if !strings.HasPrefix(firstLine, name+" ") {
					failures = append(failures, fmt.Sprintf("%s: %s first doc line must start with %q", path, name, name+" "))
				}
			}

			if !strings.Contains(docText, "Example:") {
				failures = append(failures, fmt.Sprintf("%s: %s doc comment must contain Example:", path, name))
			}
		}
	}

	if len(failures) > 0 {
		for _, failure := range failures {
			fmt.Fprintln(os.Stderr, failure)
		}
		os.Exit(1)
	}

	fmt.Println("doc comment audit passed")
}

func loadAPIMethods(folder string) (map[string]bool, error) {
	path := filepath.Join(folder, "api_contract.go")
	fileSet := token.NewFileSet()
	file, err := parser.ParseFile(fileSet, path, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}

	methods := make(map[string]bool)
	for _, decl := range file.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.TYPE {
			continue
		}
		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok || typeSpec.Name.Name != "API" {
				continue
			}
			iface, ok := typeSpec.Type.(*ast.InterfaceType)
			if !ok {
				continue
			}
			for _, field := range iface.Methods.List {
				for _, name := range field.Names {
					methods[name.Name] = true
				}
			}
		}
	}

	return methods, nil
}

func firstNonEmptyLine(text string) string {
	for _, line := range strings.Split(text, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			return trimmed
		}
	}
	return ""
}
