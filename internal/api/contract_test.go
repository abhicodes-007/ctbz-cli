package api

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/edusouza/ctbz-cli/internal/contract"
	"github.com/edusouza/ctbz-cli/internal/ctbz"
)

// TestContracts verifica as fixtures anonimizadas contra os tipos de resposta.
// Para atualizar uma fixture: go run ./tools/capture NOME.
func TestContracts(t *testing.T) {
	for _, e := range Endpoints() {
		t.Run(e.Name, func(t *testing.T) {
			data, err := os.ReadFile(filepath.Join("testdata", e.Name+".json"))
			if err != nil {
				t.Fatalf("fixture ausente: rode `go run ./tools/capture %s` (%v)", e.Name, err)
			}
			checkContract(t, data, e)
		})
	}
}

// TestFixturesHaveEndpoints evita fixtures órfãs depois de renomear um endpoint.
func TestFixturesHaveEndpoints(t *testing.T) {
	names := map[string]bool{}
	for _, e := range Endpoints() {
		if names[e.Name] {
			t.Errorf("endpoint duplicado: %s", e.Name)
		}
		names[e.Name] = true
	}
	files, _ := filepath.Glob(filepath.Join("testdata", "*.json"))
	for _, f := range files {
		if name := strings.TrimSuffix(filepath.Base(f), ".json"); !names[name] {
			t.Errorf("fixture sem endpoint: %s", f)
		}
	}
}

// TestContractsLive verifica os contratos contra a API real, com a sessão do
// `ctbz login`. Só faz GET. Rode com:
//
//	CTBZ_CONTRACT_LIVE=1 go test ./internal/api -run Live -v
func TestContractsLive(t *testing.T) {
	if os.Getenv("CTBZ_CONTRACT_LIVE") == "" {
		t.Skip("defina CTBZ_CONTRACT_LIVE=1 para verificar contra a API real")
	}
	store, err := ctbz.DefaultStore()
	if err != nil {
		t.Fatal(err)
	}
	sess, err := store.LoadSession()
	if err != nil {
		t.Fatal(err)
	}
	client := sess.Client()
	for _, e := range Endpoints() {
		t.Run(e.Name, func(t *testing.T) {
			if e.LivePath == "" {
				t.Skip("endpoint depende de dados de outra resposta")
			}
			resp, err := client.API(context.Background(), "GET", e.LivePath, nil)
			if err != nil {
				t.Fatal(err)
			}
			if resp.Status != 200 {
				t.Fatalf("HTTP %d", resp.Status)
			}
			checkContract(t, resp.Body, e)
		})
	}
}

func checkContract(t *testing.T, data []byte, e Endpoint) {
	t.Helper()
	r, err := contract.Check(data, e.Type)
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range r.Broken() {
		t.Errorf("%s", f)
	}
	if added := r.Added(); len(added) > 0 && testing.Verbose() {
		t.Logf("%d campos que a CLI não usa", len(added))
	}
}
