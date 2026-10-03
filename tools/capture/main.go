// Comando capture grava fixtures anonimizadas de respostas da API para os testes de
// contrato (internal/api/testdata).
//
//	go run ./tools/capture dadosempresa menu        # chama a API com a sessão do ctbz login
//	go run ./tools/capture -from resp.json appbar   # anonimiza uma resposta já salva
//
// Revise o arquivo gerado antes do commit: a anonimização cobre os campos conhecidos.
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/edusouza/ctbz-cli/internal/api"
	"github.com/edusouza/ctbz-cli/internal/contract"
	"github.com/edusouza/ctbz-cli/internal/ctbz"
)

func main() {
	from := flag.String("from", "", "lê a resposta deste arquivo em vez de chamar a API (um endpoint por vez)")
	path := flag.String("path", "", "caminho a chamar, para endpoints sem LivePath (ex.: com um ID)")
	dir := flag.String("dir", filepath.Join("internal", "api", "testdata"), "diretório das fixtures")
	full := flag.Bool("full", false, "mantém campos que o contrato não usa (por padrão a fixture é podada)")
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "uso: go run ./tools/capture [-from ARQUIVO | -path CAMINHO] ENDPOINT...")
		fmt.Fprintln(os.Stderr, "\nendpoints:")
		for _, e := range api.Endpoints() {
			fmt.Fprintf(os.Stderr, "  %-30s %s\n", e.Name, e.LivePath)
		}
		flag.PrintDefaults()
	}
	flag.Parse()
	if flag.NArg() == 0 || ((*from != "" || *path != "") && flag.NArg() != 1) {
		flag.Usage()
		os.Exit(2)
	}
	for _, name := range flag.Args() {
		if err := capture(name, *from, *path, *dir, *full); err != nil {
			fmt.Fprintf(os.Stderr, "%s: %v\n", name, err)
			os.Exit(1)
		}
	}
}

func capture(name, from, path, dir string, full bool) error {
	ep, ok := find(name)
	if !ok {
		return fmt.Errorf("endpoint desconhecido (veja a lista com -h)")
	}
	raw, err := fetch(ep, from, path)
	if err != nil {
		return err
	}
	report, err := contract.Check(raw, ep.Type)
	if err != nil {
		return err
	}
	if !full {
		if raw, err = contract.Prune(raw, ep.Type); err != nil {
			return err
		}
	}
	anon, err := contract.Anonymize(raw)
	if err != nil {
		return err
	}
	for _, f := range report.Broken() {
		fmt.Fprintf(os.Stderr, "%s: CONTRATO QUEBRADO: %s\n", name, f)
	}
	if n := len(report.Added()); n > 0 {
		fmt.Fprintf(os.Stderr, "%s: %d campos que a CLI não usa (informativo)\n", name, n)
	}
	file := filepath.Join(dir, name+".json")
	if err := os.WriteFile(file, anon, 0o644); err != nil {
		return err
	}
	fmt.Println("gravado:", file, "(revise antes do commit)")
	return nil
}

func find(name string) (api.Endpoint, bool) {
	for _, e := range api.Endpoints() {
		if e.Name == name {
			return e, true
		}
	}
	return api.Endpoint{}, false
}

func fetch(ep api.Endpoint, from, path string) ([]byte, error) {
	if from != "" {
		return os.ReadFile(from)
	}
	if path == "" {
		path = ep.LivePath
	}
	if path == "" {
		return nil, fmt.Errorf("endpoint sem LivePath: informe -path")
	}
	store, err := ctbz.DefaultStore()
	if err != nil {
		return nil, err
	}
	sess, err := store.LoadSession()
	if err != nil {
		return nil, err
	}
	resp, err := sess.Client().API(context.Background(), "GET", path, nil)
	if err != nil {
		return nil, err
	}
	if resp.Status != 200 {
		return nil, &ctbz.HTTPError{Step: path, Status: resp.Status, Body: resp.Body}
	}
	return resp.Body, nil
}
