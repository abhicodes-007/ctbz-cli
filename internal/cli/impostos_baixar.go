package cli

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/edusouza/ctbz-cli/internal/api"
	"github.com/edusouza/ctbz-cli/internal/output"
)

func newImpostosBaixarCmd() *cobra.Command {
	var dir string
	var pendentes, force bool
	cmd := &cobra.Command{
		Use:   "baixar [ID...]",
		Short: "Baixa o PDF de guias de imposto",
		Long: `Baixa o PDF das guias informadas (IDs de "ctbz impostos") ou, com --pendentes, de todas
as guias a pagar (em atraso, do mês e do próximo mês).

Os arquivos são nomeados COMPETÊNCIA-IMPOSTO-VENCIMENTO.pdf (ex.:
2026-07-darf-unificado-2026-10-06.pdf). Arquivos existentes não são sobrescritos sem --force.
A lista de arquivos gravados sai em stdout (no formato de -o).`,
		Example: `  ctbz impostos baixar 1000000000000001
  ctbz impostos baixar --pendentes -d ~/guias`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if pendentes == (len(args) > 0) {
				return usageError{errors.New("informe IDs de guias ou --pendentes (um dos dois)")}
			}
			f, err := outputFormat(cmd, "")
			if err != nil {
				return err
			}
			s := streamsOf(cmd)
			g := sessionGetter{s}
			ids, err := idsParaBaixar(cmd.Context(), g, args, pendentes)
			if err != nil {
				return err
			}
			if err := os.MkdirAll(dir, 0o755); err != nil {
				return err
			}
			l := &output.List{Columns: []output.Column{
				{Key: "id", Header: "ID"},
				{Key: "arquivo", Header: "Arquivo"},
				{Key: "situacao", Header: "Situação"},
			}}
			for _, id := range ids {
				arquivo, situacao, err := baixarGuia(cmd.Context(), g, id, dir, force)
				if err != nil {
					return fmt.Errorf("guia %d: %w", id, err)
				}
				l.Append(id, arquivo, situacao)
			}
			return output.Write(s.out, f, l)
		},
	}
	cmd.Flags().StringVarP(&dir, "dir", "d", ".", "diretório onde salvar os PDFs")
	cmd.Flags().BoolVar(&pendentes, "pendentes", false, "baixa todas as guias a pagar")
	cmd.Flags().BoolVar(&force, "force", false, "sobrescreve arquivos existentes")
	return cmd
}

func idsParaBaixar(ctx context.Context, g api.Getter, args []string, pendentes bool) ([]int64, error) {
	if !pendentes {
		ids := make([]int64, len(args))
		for i, a := range args {
			id, err := parseID(a)
			if err != nil {
				return nil, err
			}
			ids[i] = id
		}
		return ids, nil
	}
	guias, err := api.BuscarGuiasAPagar(ctx, g)
	if err != nil {
		return nil, err
	}
	var ids []int64
	for _, grupo := range [][]api.GuiaResumo{guias.EmAtraso, guias.EsteMes, guias.ProximoMes} {
		for _, guia := range grupo {
			ids = append(ids, guia.ID)
		}
	}
	return ids, nil
}

// baixarGuia grava o PDF de uma guia em dir e devolve o caminho e a situação
// ("baixado" ou "já existe").
func baixarGuia(ctx context.Context, g api.Getter, id int64, dir string, force bool) (string, string, error) {
	detalhe, err := api.BuscarGuia(ctx, g, id)
	if err != nil {
		return "", "", err
	}
	path := filepath.Join(dir, nomeArquivoGuia(detalhe))
	if _, err := os.Stat(path); err == nil && !force {
		return path, "já existe", nil
	}
	link, err := api.BuscarLinkGuia(ctx, g, id, detalhe.Tipo)
	if err != nil {
		return "", "", err
	}
	if err := download(ctx, link.URL, path); err != nil {
		return "", "", err
	}
	return path, "baixado", nil
}

// nomeArquivoGuia monta "AAAA-MM-imposto-AAAA-MM-DD.pdf" a partir do detalhe da guia.
func nomeArquivoGuia(g *api.GuiaDetalhe) string {
	parts := []string{}
	if c, ok := competenciaDeTexto(g.Competencia).(string); ok && len(c) == 7 && c[2] == '/' {
		parts = append(parts, c[3:]+"-"+c[:2]) // MM/AAAA → AAAA-MM
	}
	parts = append(parts, slug(firstNonEmpty(g.Oraculo.Nome, g.IdentificadorImposto, fmt.Sprint(g.ID))))
	if v, ok := output.ParseDate(g.Vencimento); ok {
		parts = append(parts, v.Format("2006-01-02"))
	}
	return strings.Join(parts, "-") + ".pdf"
}

var reNaoSlug = regexp.MustCompile(`[^a-z0-9]+`)

// slug simplifica um texto para nome de arquivo: minúsculas, sem acentos, com hífens.
func slug(s string) string {
	s = strings.NewReplacer("á", "a", "à", "a", "â", "a", "ã", "a", "é", "e", "ê", "e", "í", "i",
		"ó", "o", "ô", "o", "õ", "o", "ú", "u", "ç", "c").Replace(strings.ToLower(s))
	return strings.Trim(reNaoSlug.ReplaceAllString(s, "-"), "-")
}

func firstNonEmpty(ss ...string) string {
	for _, s := range ss {
		if strings.TrimSpace(s) != "" {
			return s
		}
	}
	return ""
}

// downloadClient baixa arquivos de links assinados; não envia os cookies da sessão.
var downloadClient = &http.Client{Timeout: 2 * time.Minute}

// download grava o conteúdo de url em path, via arquivo temporário.
func download(ctx context.Context, url, path string) error {
	if !strings.HasPrefix(url, "https://") && !strings.HasPrefix(url, "http://") {
		return fmt.Errorf("link de download inválido")
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	resp, err := downloadClient.Do(req)
	if err != nil {
		return fmt.Errorf("baixando arquivo: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("baixando arquivo: HTTP %d", resp.StatusCode)
	}
	tmp, err := os.CreateTemp(filepath.Dir(path), ".ctbz-*.part")
	if err != nil {
		return err
	}
	defer os.Remove(tmp.Name())
	if _, err := io.Copy(tmp, resp.Body); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	return os.Rename(tmp.Name(), path)
}
