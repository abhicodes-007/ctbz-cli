package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

// TestReferenceDocsUpToDate falha quando a CLI muda sem regenerar docs/referencia.
func TestReferenceDocsUpToDate(t *testing.T) {
	want := t.TempDir()
	if err := GenMarkdown(NewRootCmd(""), want); err != nil {
		t.Fatal(err)
	}
	got := filepath.Join("..", "..", "docs", "referencia")
	wantFiles, _ := filepath.Glob(filepath.Join(want, "*.md"))
	gotFiles, _ := filepath.Glob(filepath.Join(got, "*.md"))
	if len(wantFiles) != len(gotFiles) {
		t.Fatalf("docs/referencia tem %d arquivos, esperado %d: rode `go run ./tools/gendocs`", len(gotFiles), len(wantFiles))
	}
	for _, w := range wantFiles {
		name := filepath.Base(w)
		a, _ := os.ReadFile(w)
		b, err := os.ReadFile(filepath.Join(got, name))
		if err != nil || !bytes.Equal(a, b) {
			t.Errorf("docs/referencia/%s desatualizado: rode `go run ./tools/gendocs`", name)
		}
	}
}
