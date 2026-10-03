// Comando ctbz: CLI para a plataforma web da Contabilizei.
package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/edusouza/ctbz-cli/internal/ctbz"
	"github.com/edusouza/ctbz-cli/internal/output"
	"golang.org/x/term"
)

const usage = `ctbz — CLI para a Contabilizei

Uso:
  ctbz [-o FORMATO] COMANDO [flags]

Comandos:
  login [--otp N] [--otp-cmd CMD] [--cnpj CNPJ]   autentica (usuário/senha + OTP)
  status                                          mostra a sessão atual e testa se ainda vale
  empresa                                         dados da empresa selecionada
  api [-X MÉTODO] [-d CORPO] CAMINHO              chama uma URL da plataforma com a sessão
  logout                                          apaga a sessão local

Saída:
  -o, --output FORMATO   table (padrão), json ou csv; vale antes ou depois do comando.
                         Dados vão para stdout; mensagens e progresso, para stderr.

Variáveis de ambiente:
  CTBZ_USER, CTBZ_PASSWORD   credenciais (e-mail ou CPF, e senha)
  CTBZ_OTP_CMD               comando que imprime o OTP (habilita login e re-login automáticos)
  CTBZ_CNPJ                  empresa a selecionar quando houver mais de uma
  CTBZ_HOME                  diretório da sessão (padrão: ~/.config/ctbz)
  CTBZ_OUTPUT                formato de saída padrão (table, json ou csv)
  CTBZ_VERBOSE=1             mostra progresso do --otp-cmd
`

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	rest, err := parseGlobalFlags(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, "erro:", err)
		os.Exit(2)
	}
	if len(rest) == 0 {
		fmt.Fprint(os.Stderr, usage)
		os.Exit(2)
	}
	args := rest[1:]
	switch rest[0] {
	case "login":
		err = cmdLogin(ctx, args)
	case "logout":
		err = cmdLogout()
	case "status":
		err = cmdStatus(ctx, args)
	case "empresa":
		err = cmdEmpresa(ctx, args)
	case "api":
		err = cmdAPI(ctx, args)
	case "help", "-h", "--help":
		fmt.Print(usage)
	default:
		fmt.Fprintf(os.Stderr, "comando desconhecido: %s\n\n%s", rest[0], usage)
		os.Exit(2)
	}
	switch {
	case err == nil:
	case errors.Is(err, errPending):
		os.Exit(3)
	case errors.Is(err, flag.ErrHelp):
		os.Exit(0)
	default:
		fmt.Fprintln(os.Stderr, "erro:", err)
		os.Exit(1)
	}
}

func cmdLogout() error {
	store, err := ctbz.DefaultStore()
	if err != nil {
		return err
	}
	if err := store.Clear(); err != nil {
		return err
	}
	fmt.Fprintln(os.Stderr, "Sessão local removida.")
	return nil
}

func cmdStatus(ctx context.Context, args []string) error {
	fs := flag.NewFlagSet("status", flag.ContinueOnError)
	format := addOutputFlag(fs, defaultFormat())
	if err := fs.Parse(args); err != nil {
		return err
	}
	f, err := output.ParseFormat(*format)
	if err != nil {
		return err
	}
	store, err := ctbz.DefaultStore()
	if err != nil {
		return err
	}
	sess, err := store.LoadSession()
	if err != nil {
		return err
	}
	situacao := "válida"
	c := sess.Client()
	if _, err := c.API(ctx, "GET", "appbar/get", nil); err != nil {
		if !errors.Is(err, ctbz.ErrUnauthorized) {
			return err
		}
		situacao = "expirada"
		fmt.Fprintln(os.Stderr, "Sessão expirada: rode `ctbz login`.")
	} else {
		saveCookies(store, sess, c)
	}
	info := sessionInfo(sess)
	rec := &output.Record{}
	rec.Add("empresa", "Empresa", info.RazaoSocial).
		Add("cnpj", "CNPJ", output.NewCNPJ(sess.CNPJ)).
		Add("usuario", "Usuário", info.Email).
		Add("login_em", "Login em", output.DateTime{Time: sess.CreatedAt}).
		Add("expira_em", "Expira em", output.DateTime{Time: sess.ExpiresAt()}).
		Add("situacao", "Situação", situacao)
	return output.Write(os.Stdout, f, rec)
}

func cmdEmpresa(ctx context.Context, args []string) error {
	fs := flag.NewFlagSet("empresa", flag.ContinueOnError)
	format := addOutputFlag(fs, defaultFormat())
	fs.BoolFunc("json", "atalho para --output json", func(string) error { *format = "json"; return nil })
	if err := fs.Parse(args); err != nil {
		return err
	}
	f, err := output.ParseFormat(*format)
	if err != nil {
		return err
	}
	resp, err := authedAPI(ctx, "GET", "dadosempresa/get", nil)
	if err != nil {
		return err
	}
	if resp.Status != 200 {
		return &ctbz.HTTPError{Step: "dadosempresa/get", Status: resp.Status, Body: resp.Body}
	}
	var data struct {
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
	if err := json.Unmarshal(resp.Body, &data); err != nil {
		return fmt.Errorf("resposta inesperada de dadosempresa/get: %w", err)
	}
	e := data.EmpresaAtual
	var certSituacao any
	var certValidade output.Date
	if e.Certificado != nil {
		certSituacao = e.Certificado.Status.Descricao
		certValidade, _ = output.ParseDate(e.Certificado.DataValidade)
	}
	outras := []output.Record{}
	for _, o := range data.Empresas {
		if ctbz.OnlyDigits(o.CNPJ) == ctbz.OnlyDigits(e.CNPJ) {
			continue
		}
		r := output.Record{}
		r.Add("cnpj", "CNPJ", output.NewCNPJ(o.CNPJ)).
			Add("razao_social", "Razão social", o.RazaoSocial).
			Add("situacao", "Situação", o.StatusEmpresa)
		outras = append(outras, r)
	}
	rec := &output.Record{}
	rec.Add("razao_social", "Razão social", e.RazaoSocial).
		Add("cnpj", "CNPJ", output.NewCNPJ(e.CNPJ)).
		Add("situacao", "Situação", e.StatusEmpresa).
		Add("regime_tributario", "Regime tributário", e.RegimeTributario).
		Add("inscricao_municipal", "Inscrição municipal", nilIfEmpty(e.InscricaoMunicipal)).
		Add("ramos_atividade", "Ramos de atividade", nonNil(e.RamosAtividade)).
		Add("plano", "Plano", nilIfEmpty(e.Plano)).
		Add("certificado_situacao", "Certificado digital", certSituacao).
		Add("certificado_validade", "Validade do certificado", certValidade).
		Add("outras_empresas", "Outras empresas", outras)
	return output.Write(os.Stdout, f, rec)
}

func cmdAPI(ctx context.Context, args []string) error {
	fs := flag.NewFlagSet("api", flag.ContinueOnError)
	method := fs.String("X", "GET", "método HTTP")
	data := fs.String("d", "", "corpo JSON da requisição (use @arquivo para ler de um arquivo ou @- para stdin)")
	raw := fs.Bool("raw", false, "imprime o corpo da resposta como veio, sem formatar")
	// Para explorar a API o padrão é JSON; só um -o explícito muda o formato.
	format := addOutputFlag(fs, explicitFormatOr("json"))
	fs.Usage = func() {
		fmt.Fprint(fs.Output(), `Uso: ctbz api [-X MÉTODO] [-d CORPO] CAMINHO

CAMINHO relativo é resolvido contra o BFF da plataforma (/api/plataforma/):
  ctbz api dadosempresa/get
  ctbz api menu/get
Caminhos absolutos alcançam as demais APIs do app:
  ctbz api /api/legado/...

A resposta sai em JSON formatado. Com -o table ou -o csv, listas de objetos
viram tabelas (uma coluna por campo). --raw imprime o corpo exatamente como veio.

Flags:
`)
		fs.PrintDefaults()
	}
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() != 1 {
		fs.Usage()
		return errors.New("informe exatamente um CAMINHO")
	}
	f, err := output.ParseFormat(*format)
	if err != nil {
		return err
	}
	body, err := readBody(*data)
	if err != nil {
		return err
	}
	resp, err := authedAPI(ctx, strings.ToUpper(*method), fs.Arg(0), body)
	if err != nil {
		return err
	}
	if parsed, perr := output.FromJSON(resp.Body); !*raw && perr == nil {
		if err := output.Write(os.Stdout, f, parsed); err != nil {
			return err
		}
	} else {
		os.Stdout.Write(resp.Body)
		if !*raw && len(resp.Body) > 0 && resp.Body[len(resp.Body)-1] != '\n' {
			os.Stdout.WriteString("\n")
		}
	}
	if resp.Status >= 400 {
		return fmt.Errorf("HTTP %d", resp.Status)
	}
	return nil
}

// authedAPI faz uma chamada com a sessão salva. Se a sessão tiver expirado e
// CTBZ_OTP_CMD estiver configurado, refaz o login sozinho e repete a chamada.
func authedAPI(ctx context.Context, method, path string, body []byte) (*ctbz.Response, error) {
	store, err := ctbz.DefaultStore()
	if err != nil {
		return nil, err
	}
	sess, err := store.LoadSession()
	if err != nil {
		return nil, err
	}
	c := sess.Client()
	resp, err := c.API(ctx, method, path, bodyReader(body))
	if errors.Is(err, ctbz.ErrUnauthorized) && os.Getenv("CTBZ_OTP_CMD") != "" {
		fmt.Fprintln(os.Stderr, "Sessão expirada; refazendo login com CTBZ_OTP_CMD…")
		sess, err = login(ctx, store, loginOpts{
			otpCmd:     os.Getenv("CTBZ_OTP_CMD"),
			otpTimeout: envDuration("CTBZ_OTP_TIMEOUT", 3*time.Minute),
			cnpj:       sess.CNPJ,
			restart:    true,
		})
		if err != nil {
			return nil, err
		}
		c = sess.Client()
		resp, err = c.API(ctx, method, path, bodyReader(body))
	}
	if err != nil {
		return nil, err
	}
	saveCookies(store, sess, c)
	return resp, nil
}

// saveCookies persiste cookies renovados pelo servidor durante as chamadas.
func saveCookies(store *ctbz.Store, sess *ctbz.Session, c *ctbz.Client) {
	sess.Cookies = c.Cookies
	if err := store.SaveSession(sess); err != nil {
		fmt.Fprintln(os.Stderr, "aviso: não foi possível atualizar a sessão:", err)
	}
}

func bodyReader(b []byte) io.Reader {
	if b == nil {
		return nil
	}
	return strings.NewReader(string(b))
}

func readBody(arg string) ([]byte, error) {
	switch {
	case arg == "":
		return nil, nil
	case arg == "@-":
		return io.ReadAll(os.Stdin)
	case strings.HasPrefix(arg, "@"):
		return os.ReadFile(arg[1:])
	default:
		return []byte(arg), nil
	}
}

type sessInfo struct {
	Email       string
	RazaoSocial string
	CNPJ        string
}

// sessionInfo lê usuário e empresa dos dados do localStorage salvos no login.
func sessionInfo(sess *ctbz.Session) sessInfo {
	var login struct {
		Email   string `json:"email"`
		Empresa struct {
			CNPJ        string `json:"cnpj"`
			RazaoSocial string `json:"razaoSocial"`
		} `json:"empresa"`
	}
	if raw, ok := sess.Storage["l"]; ok {
		_ = json.Unmarshal(raw, &login)
	}
	info := sessInfo{Email: login.Email, RazaoSocial: login.Empresa.RazaoSocial, CNPJ: login.Empresa.CNPJ}
	if info.CNPJ == "" {
		info.CNPJ = sess.CNPJ
	}
	return info
}

// describeSession resume a sessão numa linha (mensagens em stderr).
func describeSession(sess *ctbz.Session) string {
	info := sessionInfo(sess)
	parts := []string{}
	if info.RazaoSocial != "" {
		parts = append(parts, fmt.Sprintf("%s (CNPJ %s)", info.RazaoSocial, output.FormatCNPJ(output.NewCNPJ(info.CNPJ))))
	} else if info.CNPJ != "" {
		parts = append(parts, "CNPJ "+output.FormatCNPJ(output.NewCNPJ(info.CNPJ)))
	}
	if info.Email != "" {
		parts = append(parts, "usuário "+info.Email)
	}
	if len(parts) == 0 {
		return "sessão ativa"
	}
	return strings.Join(parts, " — ")
}

func currentCNPJ(sess *ctbz.Session) string {
	for _, key := range []string{"e", "l"} {
		raw, ok := sess.Storage[key]
		if !ok {
			continue
		}
		var v struct {
			CNPJ    string `json:"cnpj"`
			Empresa struct {
				CNPJ string `json:"cnpj"`
			} `json:"empresa"`
		}
		if json.Unmarshal(raw, &v) == nil {
			if v.CNPJ != "" {
				return ctbz.OnlyDigits(v.CNPJ)
			}
			if v.Empresa.CNPJ != "" {
				return ctbz.OnlyDigits(v.Empresa.CNPJ)
			}
		}
	}
	return ""
}

func baseURL() string {
	if v := os.Getenv("CTBZ_BASE_URL"); v != "" {
		return strings.TrimRight(v, "/")
	}
	return ctbz.DefaultBaseURL
}

func stdinIsTerminal() bool {
	return term.IsTerminal(int(os.Stdin.Fd()))
}

func verboseWriter() io.Writer {
	if os.Getenv("CTBZ_VERBOSE") != "" {
		return os.Stderr
	}
	return nil
}
