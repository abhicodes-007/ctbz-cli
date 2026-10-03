package api

import "testing"

func TestPathNotasEntrada(t *testing.T) {
	f := FiltroNotasEntrada{Ano: 2026, Mes: 9, Emitente: "ACME LTDA"}
	for lista, want := range map[string]string{
		ListaAManifestar:   "/api/emissor/notasentrada/listar/0?mes=9&ano=2026&empresa=ACME+LTDA&qtdPagina=10&cursor=c1",
		ListaManifestadas:  "/api/emissor/notasentrada/listar/1?mes=9&ano=2026&empresa=ACME+LTDA&qtdPagina=10&cursor=c1",
		ListaAClassificar:  "/api/emissor/classificacaonotas/listar/?mes=9&ano=2026&empresa=ACME+LTDA&tipo=0&limite=10&cursor=&offset=20",
		ListaClassificadas: "/api/emissor/classificacaonotas/listar/?mes=9&ano=2026&empresa=ACME+LTDA&tipo=1&limite=10&cursor=&offset=20",
	} {
		f.Lista = lista
		if got, err := pathNotasEntrada(f, "c1", 20); err != nil || got != want {
			t.Errorf("%s: %q (%v), quero %q", lista, got, err, want)
		}
	}
	f.Lista = "outra"
	if _, err := pathNotasEntrada(f, "", 0); err == nil {
		t.Error("lista desconhecida deveria dar erro")
	}
}
