# OTP automático

O código de verificação chega só por e-mail. Para a CLI rodar sem ninguém digitar,
ela executa um **comando externo** que vai buscar o código e o imprime.

## Contrato do `CTBZ_OTP_CMD`

```sh
export CTBZ_OTP_CMD='/caminho/para/meu-script'
ctbz login --cnpj 22222222000122
```

- A CLI roda `sh -c "$CTBZ_OTP_CMD"` a cada **5 s** até ele imprimir um número de
  **6 dígitos** no stdout (usa o primeiro que aparecer), ou até `CTBZ_OTP_TIMEOUT`
  (padrão `3m`).
- Saída sem código ou status de saída diferente de zero significam "ainda não
  chegou": a CLI tenta de novo. A primeira linha do stderr aparece como motivo
  (veja com `CTBZ_VERBOSE=1`).
- Variáveis recebidas pelo comando:

  | Variável | Conteúdo |
  |---|---|
  | `CTBZ_OTP_SINCE` | epoch (segundos) do pedido do código, **menos 60 s** de margem para diferença de relógio |
  | `CTBZ_OTP_SINCE_ISO` | o mesmo instante em RFC 3339 (UTC) |

  Use-as para ignorar e-mails de logins anteriores: só o código mais recente vale.
- Se o tempo esgotar, o login fica **pendente** e pode ser concluído com
  `ctbz login --otp NNNNNN`.
- Com `CTBZ_OTP_CMD` definido, `ctbz empresa`/`ctbz api` refazem o login sozinhos ao
  receber 401 (sessão expirada) e repetem a chamada.

## Gmail via Google Workspace CLI (`gws`)

[`scripts/otp-gmail-gws.sh`](https://github.com/edusouza/ctbz-cli/blob/main/scripts/otp-gmail-gws.sh)
usa o [Google Workspace CLI](https://github.com/googleworkspace/cli) e `jq`:

1. `gws gmail users messages list` com a busca
   `from:seguranca@contabilizei.com.br after:$CTBZ_OTP_SINCE` e `maxResults: 1`;
2. `gws gmail users messages get` da mensagem mais recente;
3. procura 6 dígitos no `snippet` e, se não achar, nos corpos `text/plain`/`text/html`
   (base64url decodificado com `jq @base64d`).

A sintaxe (`gws gmail users messages list --params '{…}'`, saída em JSON) segue a
documentação do `gws`. O script é testado automaticamente contra um `gws` simulado
(`internal/otp/gmail_script_test.go`): código no corpo, código no snippet, nenhum e-mail,
e-mail sem código e `gws` ausente.

### Configurar

```sh
gws auth setup     # primeira vez (precisa do gcloud)
gws auth login
export CTBZ_OTP_CMD="$PWD/scripts/otp-gmail-gws.sh"
export CTBZ_OTP_QUERY='from:seguranca@contabilizei.com.br'   # padrão
export GWS=gws                                                # binário, se não estiver no PATH
```

O operador `after:` do Gmail aceita epoch em segundos, por isso o script repassa
`CTBZ_OTP_SINCE` direto. Sem `gws` ou `jq` no `PATH`, o script sai com código 127 e a
mensagem `comando não encontrado`.

### Validar

1. Peça um código (`ctbz login --restart`, sem `CTBZ_OTP_CMD`) e rode o script sozinho:
   `scripts/otp-gmail-gws.sh` deve imprimir os 6 dígitos (ele procura nos últimos 10 minutos).
2. Login de ponta a ponta, mostrando as tentativas:
   `CTBZ_VERBOSE=1 ctbz login --restart`. As linhas `aguardando OTP (tentativa N)` mostram
   quanto o e-mail demorou (uma tentativa a cada 5 s).
3. Se o e-mail costuma demorar mais de 3 minutos, aumente `CTBZ_OTP_TIMEOUT` (ex.: `5m`).

## Encaminhando para outra caixa

Para que a CLI leia uma caixa dedicada (e não a sua caixa pessoal), crie um filtro no Gmail
pessoal: **Configurações → Filtros e endereços bloqueados → Criar filtro**, com
`De: seguranca@contabilizei.com.br`, ação **Encaminhar para** o endereço da caixa da CLI
(o Gmail pede para confirmar o endereço de encaminhamento antes). Depois:

- O remetente visto pela outra caixa pode deixar de ser `seguranca@contabilizei.com.br`
  (depende de como o encaminhamento é feito). Ajuste `CTBZ_OTP_QUERY`, por exemplo:
  `CTBZ_OTP_QUERY='subject:(código de verificação) to:bot@seudominio.com'`.
- O encaminhamento adiciona atraso. Aumente `CTBZ_OTP_TIMEOUT` se necessário.
- A margem de 60 s em `CTBZ_OTP_SINCE` cobre pequenas diferenças de relógio; atrasos
  maiores não precisam de margem, porque o filtro é pela chegada.

## Outros provedores

Qualquer comando que imprima o código serve:

```sh
# IMAP com Python (exemplo de esqueleto)
CTBZ_OTP_CMD='python3 ~/bin/ler-otp-imap.py --desde "$CTBZ_OTP_SINCE_ISO"'

# Um serviço próprio que recebe o e-mail por webhook e guarda o último código
CTBZ_OTP_CMD='curl -fsS https://meu-servidor/ultimo-otp?desde=$CTBZ_OTP_SINCE'

# Manual, pedindo numa janela gráfica
CTBZ_OTP_CMD='zenity --entry --text "Código da Contabilizei"'
```

## Sem automação

- Terminal interativo: a CLI pergunta o código.
- Sem terminal (scripts, agentes): `ctbz login` termina com status `3` e o login
  pendente fica em `$CTBZ_HOME/pending.json`; `ctbz login --otp NNNNNN` conclui.
