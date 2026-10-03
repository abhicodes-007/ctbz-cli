package cli

import "testing"

func TestDocumentos(t *testing.T) {
	tipo := "EXTRATO_BANCARIO_MOVIMENTACOES"
	withSession(t, fakeAPI(t, map[string]string{
		"/api/plataforma/documentos/tipos/init?area=DOCUMENTOS_CONTABEIS": `[{"tipo":"` + tipo + `","nome":"Extrato bancário","dataModificacao":"10/09/2026","quantidade":13}]`,
		"/api/plataforma/documentos/listar-enviados?tipoDocumento=" + tipo + "&limit=12&offset=0": `{"totalPages":2,"content":[
{"nomeArquivo":"extrato-ago.ofx","competencia":{"mes":8,"ano":2026},"dataEnvio":"2026-09-10","urlArquivo":"https://exemplo.invalid/a",
 "metadados":{"nome":"Conta PJ","valor":null}}]}`,
		"/api/plataforma/documentos/listar-enviados?tipoDocumento=" + tipo + "&limit=12&offset=1": `{"totalPages":2,"content":[
{"nomeArquivo":"extrato-jul.pdf","competencia":"07/2026","dataEnvio":1784635200000}]}`,
	}))
	out, stderr, code := execCLI(t, "", "documentos", "-o", "csv")
	want := "tipo,nome,quantidade,modificado\n" + tipo + ",Extrato bancário,13,2026-09-10\n"
	if code != ExitOK || out != want {
		t.Errorf("tipos (código %d, %s):\n%s\nesperado:\n%s", code, stderr, out, want)
	}
	out, stderr, code = execCLI(t, "", "documentos", "--tipo", tipo, "-o", "csv")
	want = "competencia,arquivo,enviado,descricao,valor,link\n" +
		"08/2026,extrato-ago.ofx,2026-09-10,Conta PJ,,https://exemplo.invalid/a\n" +
		"07/2026,extrato-jul.pdf,2026-07-21,,,\n"
	if code != ExitOK || out != want {
		t.Errorf("enviados (código %d, %s):\n%s\nesperado:\n%s", code, stderr, out, want)
	}
}

func TestDocumentosFixtures(t *testing.T) {
	withSession(t, fakeAPI(t, map[string]string{
		"/api/plataforma/documentos/tipos/init?area=DOCUMENTOS_CONTABEIS":                                         fixture(t, "documentos_tipos"),
		"/api/plataforma/documentos/listar-enviados?tipoDocumento=EXTRATO_APLICACAO_FINANCEIRA&limit=12&offset=0": fixture(t, "documentos_enviados"),
	}))
	if _, stderr, code := execCLI(t, "", "documentos", "-o", "json"); code != ExitOK {
		t.Errorf("tipos: código %d: %s", code, stderr)
	}
	if out, _, code := execCLI(t, "", "documentos", "--tipo", "EXTRATO_APLICACAO_FINANCEIRA", "-o", "json"); code != ExitOK || out != "[]\n" {
		t.Errorf("enviados vazios (código %d): %q", code, out)
	}
}
