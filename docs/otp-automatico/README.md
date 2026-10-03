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

[`scripts/otp-gmail-gws.sh`](https://github.com/edusouza/ctbz-cli/blob/main/scripts/otp-gmail-gws.sh) usa o
[Google Workspace CLI](https://github.com/googleworkspace/cli) e `jq`:

1. `gws gmail users messages list` com a busca
   `from:seguranca@contabilizei.com.br after:$CTBZ_OTP_SINCE` e `maxResults: 1`;
2. `gws gmail users messages get` (`format: full`) da mensagem mais recente;
3. procura 6 dígitos no `snippet` e, se não achar, nos corpos `text/plain`/`text/html`
   (base64url decodificado com `jq @base64d`).

```sh
export CTBZ_OTP_CMD="$PWD/scripts/otp-gmail-gws.sh"
export CTBZ_OTP_QUERY='from:seguranca@contabilizei.com.br'   # padrão
export GWS=gws                                                # binário, se não estiver no PATH
```

O operador `after:` do Gmail aceita epoch em segundos, por isso o script repassa
`CTBZ_OTP_SINCE` direto.

> Testado com um `gws` simulado (três cenários: código no corpo, código no snippet e
> nenhum e-mail). Ainda não rodou contra um Gmail real; confirme se os subcomandos
> `users messages list/get` batem com a sua versão do `gws`.

## Encaminhando para outra caixa

Se o e-mail de segurança for reencaminhado por uma regra do Gmail para uma caixa que a
CLI acessa:

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
