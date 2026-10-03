package api

import "context"

const (
	PathEmissorInit      = "novo-emissor/listagem/init"
	PathVersaoEmissor    = "novo-emissor/v2/versao-emissor"
	PathAliquotasEmissor = "notafiscal/listaliquotaatividade"
)

// EmissorInit é a configuração do emissor de notas (novo-emissor/listagem/init). A resposta
// também traz dados do responsável (CPF, e-mail), que a CLI não lê.
type EmissorInit struct {
	EmissorEnabled         bool `json:"emissorEnabled"`
	HasInstability         bool `json:"hasInstability"`
	PermiteEmissaoExterior bool `json:"permiteEmissaoExterior"`
	CertificadoDigital     *struct {
		Vencido        bool   `json:"vencido"`
		DataVencimento string `json:"dataVencimento"` // AAAA-MM-DD
		EmRenovacao    bool   `json:"emRenovacao"`
	} `json:"certificadoDigital"`
	Endereco *struct {
		Municipio *struct {
			Nome string `json:"nome"`
			UF   struct {
				ID string `json:"id"`
			} `json:"uf"`
		} `json:"municipio"`
	} `json:"endereco"`
}

// BuscarEmissorInit lê a configuração do emissor.
func BuscarEmissorInit(ctx context.Context, g Getter) (*EmissorInit, error) {
	return get[EmissorInit](ctx, g, PathEmissorInit)
}

// VersaoEmissor indica qual emissor a empresa usa (ex.: "V2").
type VersaoEmissor struct {
	VersaoNovoEmissor string `json:"versaoNovoEmissor"`
}

// BuscarVersaoEmissor lê a versão do emissor.
func BuscarVersaoEmissor(ctx context.Context, g Getter) (*VersaoEmissor, error) {
	return get[VersaoEmissor](ctx, g, PathVersaoEmissor)
}

// AliquotaAtividade é a alíquota de uma atividade (CNAE) e o item da lista de serviços.
type AliquotaAtividade struct {
	CodigoCnae           string   `json:"codigoCnae"`
	DescricaoCnae        string   `json:"descricaoCnae"`
	CodigoItemServico    string   `json:"codigoItemServico"`
	DescricaoItemServico string   `json:"descricaoItemServico"`
	AliquotaBase         *float64 `json:"aliquotaBase"` // % do Simples
	AliquotaISS          *float64 `json:"aliquotaISS"`  // % de ISS dentro da alíquota
	FatorR               *float64 `json:"fatorR"`       // % de Fator R usado no cálculo
	AnexoFixo            bool     `json:"anexoFixo"`
}

// AliquotasMercado separa as alíquotas de serviço e de comércio.
type AliquotasMercado struct {
	Servico  []AliquotaAtividade `json:"servico"`
	Comercio []AliquotaAtividade `json:"comercio"`
}

// AliquotasEmissor é a resposta de notafiscal/listaliquotaatividade: alíquotas para
// tomadores no Brasil (interno) e no exterior (externo).
type AliquotasEmissor struct {
	RegimeTributario string           `json:"regimeTributario"`
	Interno          AliquotasMercado `json:"interno"`
	Externo          AliquotasMercado `json:"externo"`
}

// BuscarAliquotasEmissor lê as alíquotas por atividade.
func BuscarAliquotasEmissor(ctx context.Context, g Getter) (*AliquotasEmissor, error) {
	return get[AliquotasEmissor](ctx, g, PathAliquotasEmissor)
}
