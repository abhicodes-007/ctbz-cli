package main

import (
	"os"
	"strings"
	"testing"
)

const changelog = `# Changelog

## [Unreleased]

## [0.2.0] - 2026-10-03

### Added

- Comando b.

## [0.1.0] - 2026-10-02

### Added

- Comando a.

[Unreleased]: https://exemplo/compare/v0.2.0...HEAD
[0.2.0]: https://exemplo/compare/v0.1.0...v0.2.0
`

func TestSecao(t *testing.T) {
	for versao, want := range map[string]string{
		"v0.2.0": "### Added\n\n- Comando b.\n",
		"0.1.0":  "### Added\n\n- Comando a.\n",
	} {
		got, err := secao(strings.NewReader(changelog), versao)
		if err != nil || got != want {
			t.Errorf("secao(%s) = %q, %v; quero %q", versao, got, err, want)
		}
	}
	if _, err := secao(strings.NewReader(changelog), "v9.9.9"); err == nil {
		t.Error("versão inexistente deveria dar erro")
	}
}

// Toda versão do CHANGELOG.md real precisa ter notas, senão a release sairia vazia.
func TestChangelogDoRepositorio(t *testing.T) {
	data, err := os.ReadFile("../../CHANGELOG.md")
	if err != nil {
		t.Fatal(err)
	}
	for _, l := range strings.Split(string(data), "\n") {
		if !strings.HasPrefix(l, "## [") || strings.HasPrefix(l, "## [Unreleased]") {
			continue
		}
		versao := strings.TrimPrefix(strings.SplitN(l, "]", 2)[0], "## [")
		if notas, err := secao(strings.NewReader(string(data)), versao); err != nil || notas == "" {
			t.Errorf("versão %s sem notas (%v)", versao, err)
		}
	}
}
