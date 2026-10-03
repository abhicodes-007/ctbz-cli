package otp

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestExtract(t *testing.T) {
	for in, want := range map[string]string{
		"123456\n": "123456",
		"Seu código é 654321. Não compartilhe": "654321",
		"CEP 01001000 código 111222":           "111222",
		"nada aqui":                            "",
	} {
		got, _ := Extract(in)
		if got != want {
			t.Errorf("Extract(%q) = %q, quero %q", in, got, want)
		}
	}
}

func TestCommandRetriesUntilCode(t *testing.T) {
	dir := t.TempDir()
	counter := filepath.Join(dir, "n")
	// Falha nas duas primeiras execuções, depois imprime o código e o "since" recebido.
	script := `n=$(cat "` + counter + `" 2>/dev/null || echo 0); n=$((n+1)); echo $n > "` + counter + `"
[ $n -lt 3 ] && { echo "ainda não chegou" >&2; exit 1; }
echo "codigo 987654 since=$CTBZ_OTP_SINCE"`
	since := time.Unix(1700000000, 0)
	code, err := (&Command{Cmd: script, Since: since, Interval: 10 * time.Millisecond, Timeout: 5 * time.Second}).Fetch(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if code != "987654" {
		t.Fatalf("code = %q", code)
	}
	n, _ := os.ReadFile(counter)
	if strings.TrimSpace(string(n)) != "3" {
		t.Fatalf("execuções = %s", n)
	}
}

func TestCommandTimeout(t *testing.T) {
	_, err := (&Command{Cmd: "echo sem codigo", Interval: 10 * time.Millisecond, Timeout: 100 * time.Millisecond}).Fetch(context.Background())
	if err == nil || !strings.Contains(err.Error(), "6 dígitos") {
		t.Fatalf("esperava timeout com último erro, veio %v", err)
	}
}

func TestPrompt(t *testing.T) {
	var out strings.Builder
	code, err := Prompt(strings.NewReader("abc\n 123 456\n111111\n"), &out, "OTP: ")
	if err != nil || code != "111111" {
		t.Fatalf("code=%q err=%v", code, err)
	}
}
