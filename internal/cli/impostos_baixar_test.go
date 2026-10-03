package cli

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestImpostosBaixar(t *testing.T) {
	downloads := 0
	pdf := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Cookie") != "" {
			t.Error("o download não deve enviar cookies da sessão")
		}
		downloads++
		w.Write([]byte("%PDF-1.4 guia"))
	}))
	t.Cleanup(pdf.Close)
	withSession(t, fakeAPI(t, map[string]string{
		"/api/plataforma/impostos/v5/impostos-a-pagar/guias":                             fixture(t, "guias_a_pagar"),
		"/api/plataforma/impostos/v5/impostos-a-pagar/guia/1000000000000001":             fixture(t, "guia_detalhe"),
		"/api/plataforma/impostos/v3/impostos-a-pagar/guia/1000000000000001/baixar-guia": `{"url":"` + pdf.URL + `/guia.pdf"}`,
	}))
	dir := t.TempDir()

	out, stderr, code := execCLI(t, "", "impostos", "baixar", "--pendentes", "-d", dir, "-o", "csv")
	if code != ExitOK {
		t.Fatalf("código %d: %s", code, stderr)
	}
	// A fixture anonimizada troca o nome do imposto por "FULANO DE TAL".
	file := filepath.Join(dir, "2026-07-fulano-de-tal-2026-10-06.pdf")
	if out != "id,arquivo,situacao\n1000000000000001,"+file+",baixado\n" {
		t.Errorf("saída:\n%s", out)
	}
	if b, _ := os.ReadFile(file); !strings.HasPrefix(string(b), "%PDF") {
		t.Errorf("arquivo não é o PDF: %q", b)
	}

	out, _, _ = execCLI(t, "", "impostos", "baixar", "1000000000000001", "-d", dir, "-o", "csv")
	if !strings.Contains(out, ",já existe\n") || downloads != 1 {
		t.Errorf("deveria pular arquivo existente (downloads=%d):\n%s", downloads, out)
	}
	execCLI(t, "", "impostos", "baixar", "1000000000000001", "-d", dir, "--force")
	if downloads != 2 {
		t.Errorf("--force deveria baixar de novo (downloads=%d)", downloads)
	}

	for _, args := range [][]string{{"impostos", "baixar"}, {"impostos", "baixar", "1", "--pendentes"}} {
		if _, _, code := execCLI(t, "", args...); code != ExitUsage {
			t.Errorf("%v: código %d, quero %d", args, code, ExitUsage)
		}
	}
}

func TestSlug(t *testing.T) {
	if got := slug("DARF Unificado — Ação/Previdência"); got != "darf-unificado-acao-previdencia" {
		t.Errorf("slug = %q", got)
	}
}
