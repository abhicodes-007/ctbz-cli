package cli

import "testing"

func TestPlano(t *testing.T) {
	withSession(t, fakeAPI(t, map[string]string{
		"/api/plataforma/contrato/buscarContratoServico/":    `{"html":"<html><body><p>CONTRATO</p><p>Cláusula&nbsp;1</p></body></html>","versao":"10.0","idHistoricoContrato":1}`,
		"/api/legado/contrato/buscarContratoPlanoPagamento/": `<html><body><p>PROPOSTA</p></body></html>`,
	}))

	out, stderr, code := execCLI(t, "", "plano", "-o", "csv")
	if code != ExitOK {
		t.Fatalf("código %d: %s", code, stderr)
	}
	if want := "plano,categoria,valor,ramos_atividade\n2025 - Simples - Serviço - Básico [139],BASICO,1000.00,SERVICO\n"; out != want {
		t.Errorf("CSV:\n%s\nesperado:\n%s", out, want)
	}

	if out, _, code := execCLI(t, "", "plano", "contrato", "--texto"); code != ExitOK || out != "CONTRATO\nCláusula 1\n" {
		t.Errorf("contrato --texto (código %d): %q", code, out)
	}
	if out, _, code := execCLI(t, "", "plano", "contrato"); code != ExitOK || out != "<html><body><p>CONTRATO</p><p>Cláusula&nbsp;1</p></body></html>" {
		t.Errorf("contrato em HTML (código %d): %q", code, out)
	}
	if out, _, code := execCLI(t, "", "plano", "proposta", "--texto"); code != ExitOK || out != "PROPOSTA\n" {
		t.Errorf("proposta --texto (código %d): %q", code, out)
	}
}

func TestPlanoContratoVazio(t *testing.T) {
	withSession(t, fakeAPI(t, map[string]string{"/api/plataforma/contrato/buscarContratoServico/": `{"html":"","versao":"","idHistoricoContrato":0}`}))
	if _, stderr, code := execCLI(t, "", "plano", "contrato"); code != ExitError {
		t.Errorf("contrato vazio deveria falhar, código %d: %s", code, stderr)
	}
}
