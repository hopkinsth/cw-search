package main

import "github.com/codegangsta/cli"

func init() {
	cli.AppHelpTemplate = `NAME:
   {{.Name}} - {{.Usage}}
USAGE:
   {{.Name}} [options] [log-group] [log-stream]
VERSION:
   {{.Version}}
AUTHOR(S): 
   {{range .Authors}}{{ . }}
   {{end}}
   {{if .Flags}}
GLOBAL OPTIONS:
   {{range .Flags}}{{.}}
   {{end}}{{end}}
	`
}
