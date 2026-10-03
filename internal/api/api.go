// Package api descreve as leituras da API da Contabilizei usadas pela CLI: o caminho
// de cada endpoint, o tipo Go da resposta e uma função que busca e decodifica.
//
// Os tipos são também o contrato verificado por internal/contract: um campo com tag
// json é obrigatório na resposta, a menos que tenha `contract:"optional"`.
// Declare só os campos que a CLI usa.
package api

import "context"

// Getter faz um GET autenticado e decodifica a resposta JSON em v. A implementação
// (sessão, re-login, cookies) fica na camada de cima.
type Getter interface {
	GetJSON(ctx context.Context, path string, v any) error
}

// Endpoint liga um caminho ao tipo da resposta, para os testes de contrato.
type Endpoint struct {
	// Name identifica a fixture: testdata/<Name>.json.
	Name string
	// LivePath é o caminho chamado no modo ao vivo; vazio quando depende de dados
	// de outra resposta (ex.: um ID) e não pode ser chamado isoladamente.
	LivePath string
	// Type é o valor zero do tipo que decodifica a resposta.
	Type any
}

// Endpoints lista todas as leituras tipadas, na ordem do roadmap.
func Endpoints() []Endpoint {
	return []Endpoint{
		{Name: "dadosempresa", LivePath: PathDadosEmpresa, Type: DadosEmpresa{}},
		{Name: "appbar", LivePath: PathAppBar, Type: AppBar{}},
		{Name: "menu", LivePath: PathMenu, Type: []MenuItem{}},
		{Name: "sessao_empresa", Type: EmpresaSessao{}}, // localStorage "e" do login
		{Name: "certificado_status", LivePath: PathCertificadoStatus, Type: CertificadoStatus{}},
		{Name: "socios", LivePath: PathSocios, Type: []Socio{}},
		{Name: "cnaes", LivePath: PathCNAEs, Type: []CNAEEmpresa{}},
		{Name: "guias_a_pagar", LivePath: PathGuiasAPagar, Type: GuiasAPagar{}},
	}
}

func get[T any](ctx context.Context, g Getter, path string) (*T, error) {
	var v T
	if err := g.GetJSON(ctx, path, &v); err != nil {
		return nil, err
	}
	return &v, nil
}
