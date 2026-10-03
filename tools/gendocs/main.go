// Comando gendocs regenera docs/referencia a partir da árvore de comandos do ctbz.
//
//	go run ./tools/gendocs
package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/edusouza/ctbz-cli/internal/cli"
)

func main() {
	dir := filepath.Join("docs", "referencia")
	if err := os.RemoveAll(dir); err != nil {
		fail(err)
	}
	if err := cli.GenMarkdown(cli.NewRootCmd(""), dir); err != nil {
		fail(err)
	}
	fmt.Println("referência gerada em", dir)
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, "gendocs:", err)
	os.Exit(1)
}
