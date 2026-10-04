package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/spf13/cobra"

	"github.com/edusouza/ctbz-cli/internal/api"
	"github.com/edusouza/ctbz-cli/internal/output"
)

// execWrite roda um comando de escrita de teste ("ctbz escrita-teste") com a entrada dada.
func execWrite(t *testing.T, tty bool, stdin string, op operacao, enviar func(api.Sender) error, args ...string) (stdout, stderr string, code int) {
	t.Helper()
	old := isTerminal
	isTerminal = func(io.Reader) bool { return tty }
	t.Cleanup(func() { isTerminal = old })
	cmd := &cobra.Command{
		Use: "escrita-teste",
		RunE: func(cmd *cobra.Command, _ []string) error {
			enviado, err := escrever(cmd, op, enviar)
			if err != nil || !enviado {
				return err
			}
			return output.Write(cmd.OutOrStdout(), output.FormatJSON, resultadoEscrita("remover", "removido", 123))
		},
	}
	addWriteFlags(cmd)
	var out, errb bytes.Buffer
	root := NewRootCmd("test")
	root.AddCommand(cmd)
	root.SetArgs(append([]string{"escrita-teste"}, args...))
	root.SetIn(strings.NewReader(stdin))
	root.SetOut(&out)
	root.SetErr(&errb)
	code = run(context.Background(), root)
	return out.String(), errb.String(), code
}

func removerLancamento(s api.Sender) error {
	return s.Send(context.Background(), "DELETE", "caixa/lancamentousuario/remover/2026/9/123", nil, nil)
}

var (
	opMedio = operacao{Risco: riscoMedio, Resumo: "Excluir o lançamento 123 de 15/09/2026"}
	opAlto  = operacao{Risco: riscoAlto, Resumo: "Aceitar o termo", Consequencia: "Não dá para desfazer."}
)

func TestEscreverConfirmacao(t *testing.T) {
	for _, tc := range []struct {
		name       string
		tty        bool
		stdin      string
		op         operacao
		args       []string
		wantCode   int
		wantWrites int32
		wantErr    string
	}{
		{"tty s", true, "s\n", opMedio, nil, ExitOK, 1, ""},
		{"tty sim", true, "Sim\n", opMedio, nil, ExitOK, 1, ""},
		{"tty n", true, "n\n", opMedio, nil, ExitError, 0, "operação cancelada"},
		{"tty vazio", true, "", opMedio, nil, ExitError, 0, "operação cancelada"},
		{"tty alto confirmo", true, "confirmo\n", opAlto, nil, ExitOK, 1, ""},
		{"tty alto s não basta", true, "s\n", opAlto, nil, ExitError, 0, "operação cancelada"},
		{"sem tty sem --yes", false, "s\n", opMedio, nil, ExitUsage, 0, "--yes"},
		{"sem tty com --yes", false, "", opMedio, []string{"--yes"}, ExitOK, 1, ""},
		{"sem tty alto com -y", false, "", opAlto, []string{"-y"}, ExitOK, 1, ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			writes := writeServer(t, http.StatusOK, func(http.ResponseWriter, *http.Request) {})
			out, stderr, code := execWrite(t, tc.tty, tc.stdin, tc.op, removerLancamento, tc.args...)
			if code != tc.wantCode || writes.Load() != tc.wantWrites {
				t.Fatalf("código %d, %d escrita(s); quero %d e %d\nstderr: %s", code, writes.Load(), tc.wantCode, tc.wantWrites, stderr)
			}
			if !strings.Contains(stderr, tc.op.Resumo+" (risco "+tc.op.Risco.String()+")") {
				t.Errorf("stderr sem o resumo: %q", stderr)
			}
			if tc.op.Risco == riscoAlto && !strings.Contains(stderr, tc.op.Consequencia) {
				t.Errorf("stderr sem a consequência: %q", stderr)
			}
			if tc.wantErr != "" && !strings.Contains(stderr, tc.wantErr) {
				t.Errorf("stderr sem %q: %q", tc.wantErr, stderr)
			}
			if code == ExitOK && !strings.Contains(out, `"acao": "remover"`) {
				t.Errorf("saída sem o resultado: %q", out)
			}
		})
	}
}

func TestEscreverDryRun(t *testing.T) {
	writes := writeServer(t, http.StatusOK, func(http.ResponseWriter, *http.Request) {})
	enviar := func(s api.Sender) error {
		if err := s.Send(context.Background(), "POST", "conta-usuario/alterar-senha",
			map[string]any{"senhaAtual": "x", "novaSenha": "y", "codigo": "123456", "dados": map[string]any{"valor": -150.25, "token": nil}}, nil); err != nil {
			return err
		}
		return s.SendMultipart(context.Background(), "POST", "/api/plataforma/upload", []api.Campo{
			{Nome: "competencia", Valor: "2026-09"},
			{Nome: "arquivo", Arquivo: "extrato.ofx", Conteudo: []byte("conteúdo secreto")},
		}, nil)
	}
	// --dry-run não pede confirmação nem exige --yes, mesmo sem terminal.
	out, stderr, code := execWrite(t, false, "", opAlto, enviar, "--dry-run")
	if code != ExitOK || writes.Load() != 0 {
		t.Fatalf("código %d, %d escrita(s): %s", code, writes.Load(), stderr)
	}
	for _, want := range []string{
		"POST /api/plataforma/conta-usuario/alterar-senha\nContent-Type: application/json\n",
		`"senhaAtual": "***"`, `"novaSenha": "***"`, `"codigo": "***"`, `"valor": -150.25`, `"token": null`,
		"POST /api/plataforma/upload\nContent-Type: multipart/form-data\n\ncompetencia: 2026-09\narquivo: arquivo extrato.ofx (17 bytes)\n",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("saída sem %q:\n%s", want, out)
		}
	}
	if strings.Contains(out, "conteúdo secreto") || strings.Contains(out, `"x"`) || strings.Contains(out, "acao") {
		t.Errorf("saída com segredo, conteúdo do arquivo ou resultado:\n%s", out)
	}
	if !strings.Contains(stderr, "nada foi enviado") {
		t.Errorf("stderr = %q", stderr)
	}

	out, _, code = execWrite(t, false, "", opMedio, removerLancamento, "--dry-run", "-o", "json")
	var reqs []map[string]any
	if code != ExitOK || json.Unmarshal([]byte(out), &reqs) != nil || len(reqs) != 1 ||
		reqs[0]["metodo"] != "DELETE" || reqs[0]["caminho"] != "/api/plataforma/caixa/lancamentousuario/remover/2026/9/123" || reqs[0]["corpo"] != nil {
		t.Errorf("código %d, JSON:\n%s", code, out)
	}
}

func TestEscreverErroDoServidor(t *testing.T) {
	writeServer(t, http.StatusOK, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		io.WriteString(w, `{"message":"Não é possível realizar a assinatura como admin!"}`)
	})
	out, stderr, code := execWrite(t, false, "", opMedio, removerLancamento, "--yes")
	if code != ExitError || out != "" || !strings.Contains(stderr, "sem permissão: Não é possível realizar a assinatura como admin! (HTTP 403)") {
		t.Errorf("código %d, stdout %q, stderr %q", code, out, stderr)
	}
}
