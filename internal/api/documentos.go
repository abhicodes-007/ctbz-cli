package api

import (
	"context"
	"fmt"
	"net/url"
)

// AreaDocumentosContabeis é a única área usada pelo painel na central de documentos.
const AreaDocumentosContabeis = "DOCUMENTOS_CONTABEIS"

// PathTiposDocumento lista os tipos de documento de uma área, com a quantidade enviada.
func PathTiposDocumento(area string) string {
	return "documentos/tipos/init?area=" + url.QueryEscape(area)
}

// PathDocumentosEnviados é uma página (índice começando em 0) dos documentos enviados de
// um tipo. A API exige tipoDocumento.
func PathDocumentosEnviados(tipo string, limite, pagina int) string {
	return fmt.Sprintf("documentos/listar-enviados?tipoDocumento=%s&limit=%d&offset=%d", url.QueryEscape(tipo), limite, pagina)
}

// TipoDocumento é um tipo aceito na central de documentos.
type TipoDocumento struct {
	Tipo            string `json:"tipo"` // ex.: EXTRATO_BANCARIO_MOVIMENTACOES
	Nome            string `json:"nome"`
	DataModificacao string `json:"dataModificacao"`
	Quantidade      int    `json:"quantidade"`
}

// BuscarTiposDocumento lê os tipos de documento de uma área.
func BuscarTiposDocumento(ctx context.Context, g Getter, area string) ([]TipoDocumento, error) {
	t, err := get[[]TipoDocumento](ctx, g, PathTiposDocumento(area))
	if err != nil {
		return nil, err
	}
	return *t, nil
}

// DocumentoEnviado é um arquivo enviado. Campos lidos pelo front; a conta verificada não
// tinha documentos, por isso todos são opcionais.
type DocumentoEnviado struct {
	NomeArquivo string `json:"nomeArquivo" contract:"optional"`
	Competencia any    `json:"competencia" contract:"optional"` // formato não verificado
	DataEnvio   any    `json:"dataEnvio" contract:"optional"`
	URLArquivo  string `json:"urlArquivo" contract:"optional"`
	Metadados   *struct {
		Nome  string `json:"nome" contract:"optional"`
		Valor any    `json:"valor" contract:"optional"`
	} `json:"metadados" contract:"optional"`
}

// PaginaDocumentos é uma página do Spring (content, totalPages).
type PaginaDocumentos struct {
	Content    []DocumentoEnviado `json:"content"`
	TotalPages int                `json:"totalPages"`
}

// tamanhoPaginaDocumentos é o tamanho de página do painel.
const tamanhoPaginaDocumentos = 12

// BuscarDocumentosEnviados lê todas as páginas de documentos enviados de um tipo.
func BuscarDocumentosEnviados(ctx context.Context, g Getter, tipo string) ([]DocumentoEnviado, error) {
	var out []DocumentoEnviado
	for pagina := 0; pagina < maxPaginas; pagina++ {
		p, err := get[PaginaDocumentos](ctx, g, PathDocumentosEnviados(tipo, tamanhoPaginaDocumentos, pagina))
		if err != nil {
			return nil, err
		}
		out = append(out, p.Content...)
		if pagina+1 >= p.TotalPages {
			break
		}
	}
	return out, nil
}
