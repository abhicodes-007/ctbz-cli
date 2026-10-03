package otp

import (
	"encoding/base64"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestGmailScript roda scripts/otp-gmail-gws.sh contra um gws falso, cobrindo os
// formatos de resposta da Gmail API que o script precisa entender.
func TestGmailScript(t *testing.T) {
	if _, err := exec.LookPath("jq"); err != nil {
		t.Skip("jq não instalado")
	}
	script, _ := filepath.Abs(filepath.Join("..", "..", "scripts", "otp-gmail-gws.sh"))
	body := base64.RawURLEncoding.EncodeToString([]byte("Seu código de verificação é: 482913\nCEP 01001000"))

	for _, tc := range []struct {
		name, list, get, want string
		wantExit              int
	}{
		{"código no corpo", `{"messages":[{"id":"m1"}]}`,
			`{"snippet":"","payload":{"mimeType":"multipart/alternative","parts":[{"mimeType":"text/plain","body":{"data":"` + body + `"}}]}}`,
			"482913", 0},
		{"código no snippet", `{"messages":[{"id":"m1"}]}`, `{"snippet":"Use o código 111222 para entrar","payload":{}}`, "111222", 0},
		{"nenhum e-mail ainda", `{"resultSizeEstimate":0}`, ``, "", 1},
		{"e-mail sem código", `{"messages":[{"id":"m1"}]}`, `{"snippet":"olá","payload":{}}`, "", 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			fake := filepath.Join(dir, "gws")
			// Registra os argumentos para conferir a busca enviada ao Gmail.
			os.WriteFile(filepath.Join(dir, "list.json"), []byte(tc.list), 0o644)
			os.WriteFile(filepath.Join(dir, "get.json"), []byte(tc.get), 0o644)
			os.WriteFile(fake, []byte(`#!/bin/sh
echo "$*" >> "`+dir+`/args"
case "$4" in list) cat "`+dir+`/list.json";; get) cat "`+dir+`/get.json";; esac
`), 0o755)

			cmd := exec.Command("sh", script)
			cmd.Env = append(os.Environ(), "GWS="+fake, "CTBZ_OTP_SINCE=1700000000")
			out, err := cmd.Output()
			exit := 0
			if ee, ok := err.(*exec.ExitError); ok {
				exit = ee.ExitCode()
			} else if err != nil {
				t.Fatal(err)
			}
			if got := strings.TrimSpace(string(out)); got != tc.want || exit != tc.wantExit {
				t.Errorf("saída %q (exit %d), quero %q (exit %d)", got, exit, tc.want, tc.wantExit)
			}
			args, _ := os.ReadFile(filepath.Join(dir, "args"))
			if !strings.Contains(string(args), `from:seguranca@contabilizei.com.br after:1700000000`) {
				t.Errorf("busca enviada ao Gmail: %s", args)
			}
		})
	}
}

func TestGmailScriptMissingGws(t *testing.T) {
	script, _ := filepath.Abs(filepath.Join("..", "..", "scripts", "otp-gmail-gws.sh"))
	cmd := exec.Command("sh", script)
	cmd.Env = append(os.Environ(), "GWS=/nao/existe/gws")
	out, err := cmd.CombinedOutput()
	if ee, ok := err.(*exec.ExitError); !ok || ee.ExitCode() != 127 || !strings.Contains(string(out), "não encontrado") {
		t.Errorf("esperava exit 127 com mensagem, veio %v: %s", err, out)
	}
}
