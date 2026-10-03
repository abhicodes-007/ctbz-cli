package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
)

const (
	PathDadosEmpresa = "dadosempresa/get"
	PathAppBar       = "appbar/get"
	PathMenu         = "menu/get"
)

// DadosEmpresa é a resposta de dadosempresa/get: a empresa da sessão e as demais
// empresas do usuário.
type DadosEmpresa struct {
	EmpresaAtual struct {
		CNPJ               string   `json:"cnpj"`
		RazaoSocial        string   `json:"razaoSocial"`
		InscricaoMunicipal string   `json:"inscricaoMunicipal"`
		RegimeTributario   string   `json:"regimeTributario"`
		StatusEmpresa      string   `json:"statusEmpresa"`
		RamosAtividade     []string `json:"ramosAtividade"`
		Plano              string   `json:"plano"`
		Certificado        *struct {
			Status struct {
				Descricao string `json:"descricao"`
			} `json:"status"`
			DataValidade string `json:"dataValidade"`
		} `json:"certificado"`
	} `json:"empresaAtual"`
	Empresas []struct {
		CNPJ          string `json:"cnpj"`
		RazaoSocial   string `json:"razaoSocial"`
		StatusEmpresa string `json:"statusEmpresa"`
	} `json:"empresas"`
}

// BuscarDadosEmpresa lê dadosempresa/get.
func BuscarDadosEmpresa(ctx context.Context, g Getter) (*DadosEmpresa, error) {
	return get[DadosEmpresa](ctx, g, PathDadosEmpresa)
}

// AppBar é a resposta de appbar/get (barra superior do painel). É a leitura mais
// leve da API e serve para testar se a sessão ainda vale.
type AppBar struct {
	ContaDigital struct {
		HasRollout bool `json:"hasRollout"`
	} `json:"contaDigital"`
}

// BuscarAppBar lê appbar/get.
func BuscarAppBar(ctx context.Context, g Getter) (*AppBar, error) {
	return get[AppBar](ctx, g, PathAppBar)
}

// MenuItem é um item do menu lateral do painel (menu/get).
type MenuItem struct {
	ID          string     `json:"id"`
	Type        string     `json:"type"`
	Label       string     `json:"label"`
	Route       string     `json:"route"`
	Application string     `json:"application"`
	Children    []MenuItem `json:"children" contract:"optional"`
}

// BuscarMenu lê menu/get.
func BuscarMenu(ctx context.Context, g Getter) ([]MenuItem, error) {
	items, err := get[[]MenuItem](ctx, g, PathMenu)
	if err != nil {
		return nil, err
	}
	return *items, nil
}

// EmpresaSessao é a empresa selecionada, como o login a grava no localStorage do
// navegador (chave "e"). Não é um endpoint: vem de ctbz.Session.Storage e reflete o
// momento do login. Tem dados cadastrais que nenhuma leitura da API devolve.
type EmpresaSessao struct {
	CNPJ             string `json:"cnpj"`
	RazaoSocial      string `json:"razaoSocial"`
	NomeFantasia     string `json:"nomeFantasia" contract:"optional"`
	RegimeTributario string `json:"regimeTributario"`
	OptanteSimples   bool   `json:"optanteSimples"`
	NaturezaJuridica struct {
		CodReceita string `json:"codReceita"`
		Descricao  string `json:"descricao"`
	} `json:"naturezaJuridica"`
	InscricaoMunicipal string   `json:"inscricaoMunicipal"`
	InscricaoEstadual  string   `json:"inscricaoEstadual"`
	DataAbertura       int64    `json:"dataAbertura"`
	RamosAtividade     []string `json:"ramosAtividade"`
	// Competência em que a Contabilizei assumiu a contabilidade.
	ResponsabilidadeInicial struct {
		Mes int `json:"mes"`
		Ano int `json:"ano"`
	} `json:"responsabilidadeInicial"`
	Endereco struct {
		Logradouro  string `json:"logradouro"`
		Numero      string `json:"numero"`
		Complemento string `json:"complemento"`
		Bairro      string `json:"bairro"`
		CEP         string `json:"cep"`
		Municipio   struct {
			Nome string `json:"nome"`
			UF   struct {
				ID string `json:"id"`
			} `json:"uf"`
		} `json:"municipio"`
	} `json:"endereco"`
}

// EmpresaDaSessao decodifica a empresa gravada pelo login.
func EmpresaDaSessao(storage map[string]json.RawMessage) (*EmpresaSessao, error) {
	raw, ok := storage["e"]
	if !ok {
		return nil, errors.New("a sessão não tem os dados da empresa; rode `ctbz login` de novo")
	}
	var e EmpresaSessao
	if err := json.Unmarshal(raw, &e); err != nil {
		return nil, fmt.Errorf("dados da empresa na sessão: %w", err)
	}
	return &e, nil
}

const PathCertificadoStatus = "certificado/status"

// CertificadoStatus é a resposta de certificado/status.
type CertificadoStatus struct {
	Situacao       string `json:"situacao"`
	Valido         bool   `json:"valido"`
	AptoRenovacao  bool   `json:"aptoRenovacao"`
	DataVencimento *int64 `json:"dataVencimento"` // epoch em ms; nulo sem certificado
	MensagemErro   string `json:"mensagemErro"`
}

// BuscarCertificadoStatus lê certificado/status.
func BuscarCertificadoStatus(ctx context.Context, g Getter) (*CertificadoStatus, error) {
	return get[CertificadoStatus](ctx, g, PathCertificadoStatus)
}
