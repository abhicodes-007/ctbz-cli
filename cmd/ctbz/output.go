package main

import (
	"errors"
	"flag"
	"os"
	"strings"
)

// globalOutput guarda o -o/--output informado antes do comando (ctbz -o json status).
var globalOutput string

// parseGlobalFlags consome as flags globais antes do nome do comando.
func parseGlobalFlags(args []string) ([]string, error) {
	for len(args) > 0 {
		a := args[0]
		switch {
		case a == "-o" || a == "--output" || a == "-output":
			if len(args) < 2 {
				return nil, errors.New(a + " exige um formato: table, json ou csv")
			}
			globalOutput, args = args[1], args[2:]
		case strings.HasPrefix(a, "-o=") || strings.HasPrefix(a, "--output=") || strings.HasPrefix(a, "-output="):
			globalOutput, args = a[strings.IndexByte(a, '=')+1:], args[1:]
		default:
			return args, nil
		}
	}
	return args, nil
}

// defaultFormat: -o global, depois CTBZ_OUTPUT, depois tabela.
func defaultFormat() string {
	if globalOutput != "" {
		return globalOutput
	}
	if v := os.Getenv("CTBZ_OUTPUT"); v != "" {
		return v
	}
	return "table"
}

// explicitFormatOr ignora CTBZ_OUTPUT: usado por comandos com padrão próprio.
func explicitFormatOr(def string) string {
	if globalOutput != "" {
		return globalOutput
	}
	return def
}

// addOutputFlag registra -o e --output no comando, ambos apontando para o mesmo valor.
func addOutputFlag(fs *flag.FlagSet, def string) *string {
	p := new(string)
	fs.StringVar(p, "output", def, "formato de saída: table, json ou csv (env CTBZ_OUTPUT)")
	fs.StringVar(p, "o", def, "atalho para --output")
	return p
}

func nilIfEmpty(s string) any {
	if s == "" {
		return nil
	}
	return s
}

func nonNil(s []string) []string {
	if s == nil {
		return []string{}
	}
	return s
}
