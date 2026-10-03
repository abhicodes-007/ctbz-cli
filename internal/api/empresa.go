package api

import "context"

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
