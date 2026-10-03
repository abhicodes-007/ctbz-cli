package api

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"
)

// fixtureGetter responde com as fixtures de testdata, por caminho.
type fixtureGetter map[string]string

func (g fixtureGetter) GetJSON(_ context.Context, path string, v any) error {
	name, ok := g[path]
	if !ok {
		return errors.New("caminho inesperado: " + path)
	}
	data, err := os.ReadFile(filepath.Join("testdata", name+".json"))
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

func TestBuscarDadosEmpresa(t *testing.T) {
	d, err := BuscarDadosEmpresa(context.Background(), fixtureGetter{PathDadosEmpresa: "dadosempresa"})
	if err != nil {
		t.Fatal(err)
	}
	if d.EmpresaAtual.StatusEmpresa != "ATIVO" || d.EmpresaAtual.Certificado == nil || len(d.Empresas) != 2 {
		t.Errorf("decodificação inesperada: %+v", d)
	}
}

func TestBuscarMenu(t *testing.T) {
	items, err := BuscarMenu(context.Background(), fixtureGetter{PathMenu: "menu"})
	if err != nil {
		t.Fatal(err)
	}
	if len(items) == 0 || items[0].ID != "home" {
		t.Errorf("menu inesperado: %+v", items)
	}
}

func TestGetPropagatesError(t *testing.T) {
	if _, err := BuscarAppBar(context.Background(), fixtureGetter{}); err == nil {
		t.Error("esperava erro")
	}
}
