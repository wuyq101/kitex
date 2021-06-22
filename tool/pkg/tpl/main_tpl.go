package tpl

// MainTpl is the template for generating main.go.
var MainTpl = `
package main

import (
	{{- range $path, $alias := .Imports}}
	{{$alias }}"{{$path}}"
	{{- end}}
)

func main() {
    svr := {{.PkgRefName}}.NewServer(new({{.ServiceName}}Impl))

    err := svr.Run()

    if err != nil {
        log.Println(err.Error())
    }
}
`
