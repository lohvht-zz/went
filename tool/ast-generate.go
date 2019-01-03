package main

import (
	"bufio"
	"bytes"
	"flag"
	"go/format"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
)

// nodeType models the string data needed for the Node interface
type nodeType struct {
	DirName  string // parent directory name
	BaseName string // name of the node type
	Decls    []nodeImpl
}

// nodeImpl models the string data needed for Node implementations
type nodeImpl struct {
	Name   string            // name of the node implementation
	Fields map[string]string // a mapping of fieldnames to types
}

type visitorData struct {
	DirName string
	types   []nodeType
}

func main() {
	var outdir string
	flag.StringVar(&outdir, "outdir", "", "The output directory where the generated files will be saved to. (Required)")
	flag.Parse()

	if outdir == "" {
		flag.Usage()
		os.Exit(1)
	}

	if _, err := os.Stat(outdir); os.IsNotExist(err) {
		os.Mkdir(outdir, 0755)
	}
	expr := nodeType{
		DirName:  outdir,
		BaseName: "Expr",
		Decls: []nodeImpl{
			nodeImpl{Name: "BinExpr", Fields: map[string]string{
				"Left": "Expr", "Right": "Expr", "Op": "token.Token",
			}},
			nodeImpl{Name: "UnExpr", Fields: map[string]string{
				"Operand": "Expr", "Op": "token.Token",
			}},
			nodeImpl{Name: "BasicLit", Fields: map[string]string{
				"Text": "string",
			}},
		},
	}
	types := []nodeType{expr}
	for _, typ := range types {
		generateNodeFile(typ, nodeTemplate)
	}
}

// func generateVisitor(vd visitorData) {
// 	f, err := os.Create(filepath.Join(vd.DirName, "visitor.go"))
// 	if err != nil {
// 		panic(err) // TODO: HANDLE ERROR properly
// 	}
// 	defer f.Close()
// 	t := generateTemplate("visitor", visitorTemplate)

// 	var src bytes.Buffer
// 	t.Execute(&src, vd)
// }

// generateNodeFile generates a file that represents an AST node based on the
// noteType struct passed in.
func generateNodeFile(nt nodeType, templateText string) {
	f, err := os.Create(filepath.Join(nt.DirName, strings.ToLower(nt.BaseName)+".go"))
	if err != nil {
		panic(err) // TODO: HANDLE ERROR properly
	}
	defer f.Close()
	t := generateTemplate(nt.BaseName, templateText)

	var src bytes.Buffer
	t.Execute(&src, nt) // output template writer to pipe
	_, err = format.Source(src.Bytes())
	if err != nil {
		panic(err) // TODO: HANDLE ERROR properly
	}
	f.Sync() // NOTE: we may not need to include this
	fw := bufio.NewWriter(f)
	goimports := exec.Command("goimports")
	goimports.Stdin = &src
	goimports.Stdout = fw
	err = goimports.Run()
	if err != nil {
		panic(err) // TODO: HANDLE ERROR properly
	}
	fw.Flush()
}

// generateTemplate generates a *Template object with some common string
// manipulating functions baked into its FuncMap
func generateTemplate(name string, templateText string) *template.Template {
	funcMap := template.FuncMap{
		"ToUpper":      strings.ToUpper,
		"ToLower":      strings.ToLower,
		"JoinString":   strings.Join,
		"FilePathBase": filepath.Base,
		"MapKvJoin": func(m map[string]string, kvSep string, sep string) string {
			var buf bytes.Buffer
			first := true
			for k, v := range m {
				if !first {
					buf.WriteString(sep)
				}
				first = false
				buf.WriteString(k)
				buf.WriteString(kvSep)
				buf.WriteString(v)
			}
			return buf.String()
		},
	}
	return template.Must(template.New(name).Funcs(funcMap).Parse(templateText))
}

var javatemp = `
package com.craftinginterpreters.lox;

import java.util.List;

abstract class {{.BaseName}} {

	interface Visitor<R> {
		{{- range $i, $nodeImpl := .Decls}}
		R visit{{$nodeImpl.Name}}{{$.BaseName}}({{$nodeImpl.Name}} {{$.BaseName | ToLower}});
		{{- end}}
	}
	{{range $i, $nodeImpl := .Decls}}
	static class {{$nodeImpl.Name}} extends {{$.BaseName}} {
		{{- range $fn, $ft := $nodeImpl.Fields}}
		final {{$fn}};
		{{- end}}

		{{$nodeImpl.Name}}({{MapKvJoin $nodeImpl.Fields " " ", "}}) {
			{{- range $fn, $ft := $nodeImpl.Fields}}
			this.{{$fn}} = {{$fn}};
			{{- end}}
		}

		// Visitor pattern
		<R> R accept(Visitor<R> visitor) {
			return visitor.visit{{$nodeImpl.Name}}{{$.BaseName}}(this);
		}
	}
	{{end}}
	abstract <R> R accept(Visitor<R> visitor);
}
`

var nodeTemplate = `
package {{.DirName | FilePathBase}}

import "github.com/lohvht/went/lang/token"

type {{.BaseName}} interface {
	{{.BaseName | ToLower}}()
}

// {{.BaseName}} nodes
type (
	{{- range $i, $nodeImpl := .Decls}}
	// {{$nodeImpl.Name}} node
	{{$nodeImpl.Name}} struct {
		{{MapKvJoin $nodeImpl.Fields " " "\n"}}
	}
	{{end}}
)

{{- range $i, $nodeImpl := .Decls}}
func (n *{{$nodeImpl.Name}}) {{$.BaseName | ToLower}}() {}
{{- end}}

{{- range $i, $nodeImpl := .Decls}}
// func (n *{{$nodeImpl.Name}}) accept(v ) {}
{{- end}}
`

var visitorTemplate = `
package {{.DirName | FilePathBase}}

import "github.com/lohvht/went/lang/token"

// Visitor is the interface used to implement visitor pattern for the AST
type Visitor interface {
	{{- range $i, $type := $types}}
	// visit {{$type.BaseName}} node functions
	{{- range $j, $nodeImpl := $type.Decls}}
	visit{{$nodeImpl.Name}}(*{{$nodeImpl.Name}})
	{{- end}}
	{{- end}}
}
`
