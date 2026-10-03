package cli

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

// execCLI roda a CLI como o binário faria e devolve stdout, stderr e o código de saída.
func execCLI(t *testing.T, stdin string, args ...string) (stdout, stderr string, code int) {
	t.Helper()
	var out, errb bytes.Buffer
	root := NewRootCmd("test")
	root.SetArgs(args)
	root.SetIn(strings.NewReader(stdin))
	root.SetOut(&out)
	root.SetErr(&errb)
	code = run(context.Background(), root)
	return out.String(), errb.String(), code
}

func TestUsageErrors(t *testing.T) {
	t.Setenv("CTBZ_HOME", t.TempDir())
	for _, args := range [][]string{
		{"comando-que-nao-existe"},
		{"status", "--flag-inexistente"},
		{"api"},
		{"status", "-o", "yaml"},
	} {
		_, stderr, code := execCLI(t, "", args...)
		if code != ExitUsage {
			t.Errorf("%v: código %d, quero %d (stderr: %s)", args, code, ExitUsage, stderr)
		}
		if !strings.HasPrefix(stderr, "erro:") {
			t.Errorf("%v: stderr sem prefixo de erro: %q", args, stderr)
		}
	}
}

func TestHelpInPortuguese(t *testing.T) {
	out, _, code := execCLI(t, "", "--help")
	if code != ExitOK {
		t.Fatalf("código %d", code)
	}
	for _, want := range []string{"Uso:", "Comandos:", "login", "empresa", "Flags:", "--output"} {
		if !strings.Contains(out, want) {
			t.Errorf("ajuda sem %q:\n%s", want, out)
		}
	}
}

func TestVersion(t *testing.T) {
	out, _, code := execCLI(t, "", "version", "-o", "json")
	if code != ExitOK || !strings.Contains(out, `"versao": "test"`) {
		t.Fatalf("código %d, saída:\n%s", code, out)
	}
	out, _, _ = execCLI(t, "", "--version")
	if strings.TrimSpace(out) != "ctbz test" {
		t.Errorf("--version = %q", out)
	}
}

func TestNoSessionError(t *testing.T) {
	t.Setenv("CTBZ_HOME", t.TempDir())
	_, stderr, code := execCLI(t, "", "empresa")
	if code != ExitError || !strings.Contains(stderr, "ctbz login") {
		t.Errorf("código %d, stderr %q", code, stderr)
	}
}
