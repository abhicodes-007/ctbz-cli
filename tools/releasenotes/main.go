// Comando releasenotes imprime a seção de uma versão do CHANGELOG.md (Keep a Changelog),
// usada como notas da release no GitHub:
//
//	go run ./tools/releasenotes v0.3.0
package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "uso: go run ./tools/releasenotes vX.Y.Z")
		os.Exit(2)
	}
	f, err := os.Open("CHANGELOG.md")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer f.Close()
	notas, err := secao(f, os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Print(notas)
}

// secao devolve o corpo da seção "## [X.Y.Z]" (sem o título), sem linhas em branco nas
// pontas. A versão pode vir com ou sem o "v" da tag.
func secao(r io.Reader, versao string) (string, error) {
	titulo := "## [" + strings.TrimPrefix(versao, "v") + "]"
	var linhas []string
	dentro := false
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		l := sc.Text()
		switch {
		case strings.HasPrefix(l, "## ["):
			if dentro {
				return juntar(linhas), nil
			}
			dentro = strings.HasPrefix(l, titulo)
		case dentro && strings.HasPrefix(l, "[") && strings.Contains(l, "]: "):
			return juntar(linhas), nil // links de referência no fim do arquivo
		case dentro:
			linhas = append(linhas, l)
		}
	}
	if err := sc.Err(); err != nil {
		return "", err
	}
	if !dentro {
		return "", fmt.Errorf("versão %s não encontrada no CHANGELOG.md", versao)
	}
	return juntar(linhas), nil
}

func juntar(linhas []string) string {
	s := strings.TrimSpace(strings.Join(linhas, "\n"))
	if s == "" {
		return ""
	}
	return s + "\n"
}
