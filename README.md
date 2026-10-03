# ctbz — CLI para a Contabilizei

CLI que substitui a interface web da [Contabilizei](https://app.contabilizei.com.br):
faz o mesmo login que o navegador (usuário/senha → código por e-mail → escolha
da empresa) e chama as mesmas URLs internas que o painel usa, como se fossem uma API.

## Instalação

```sh
go install github.com/edusouza/ctbz-cli/cmd/ctbz@latest
# ou, a partir do clone:
go build -o ctbz ./cmd/ctbz
```

Requer Go 1.24+.

## Configuração

| Variável           | Uso                                                                       |
|--------------------|---------------------------------------------------------------------------|
| `CTBZ_USER`        | e-mail ou CPF do login                                                    |
| `CTBZ_PASSWORD`    | senha                                                                     |
| `CTBZ_OTP_CMD`     | comando que imprime o código OTP (habilita login e re-login automáticos)  |
| `CTBZ_OTP_TIMEOUT` | quanto esperar pelo `CTBZ_OTP_CMD` (padrão `3m`)                          |
| `CTBZ_CNPJ`        | empresa a selecionar quando o usuário tem mais de uma                     |
| `CTBZ_HOME`        | onde guardar a sessão (padrão `~/.config/ctbz`)                           |
| `CTBZ_VERBOSE=1`   | mostra as tentativas do `CTBZ_OTP_CMD`                                     |

A sessão (cookies `__C` e `oauth-token`) fica em `$CTBZ_HOME/session.json` com
permissão `0600`. As credenciais nunca são gravadas em disco.

## Login

A Contabilizei sempre envia um código de 6 dígitos por e-mail. Há três formas de entregá-lo à CLI:

### 1. Automático, via comando (`--otp-cmd` / `CTBZ_OTP_CMD`)

A CLI executa o comando a cada 5 s até ele imprimir um código de 6 dígitos (ou até
`CTBZ_OTP_TIMEOUT`). O comando recebe `CTBZ_OTP_SINCE` (epoch, em segundos) e
`CTBZ_OTP_SINCE_ISO` com o instante do pedido, para ignorar e-mails antigos.

O repositório traz [`scripts/otp-gmail-gws.sh`](scripts/otp-gmail-gws.sh), que
lê o Gmail com o [Google Workspace CLI](https://github.com/googleworkspace/cli) (`gws`) e `jq`:

```sh
export CTBZ_OTP_CMD="$PWD/scripts/otp-gmail-gws.sh"
# Se você encaminha o e-mail para outra caixa, ajuste a busca:
# export CTBZ_OTP_QUERY='subject:"código" to:bot@seudominio.com'
ctbz login --cnpj 00.000.000/0001-00
```

Qualquer comando serve: um script IMAP, uma chamada a um webhook etc. Basta imprimir o código.

Com `CTBZ_OTP_CMD` definido, `ctbz empresa` e `ctbz api` também **refazem o login
sozinhos** quando a sessão expira (HTTP 401) e repetem a chamada.

### 2. Interativo

Sem `CTBZ_OTP_CMD` e rodando num terminal, a CLI pergunta o código (e a empresa, se houver mais de uma).

### 3. Em duas etapas (scripts, agentes, CI)

Sem terminal, o login fica salvo como pendente e o processo sai com status `3`:

```sh
$ ctbz login
Código de verificação enviado para fu***@gm***com.
Conclua com: ctbz login --otp NNNNNN
$ ctbz login --otp 123456 --cnpj 00000000000100
Login concluído: FULANO LTDA (CNPJ 00000000000100) — usuário fulano@exemplo.com
```

`ctbz login --restart` descarta um login pendente e pede um código novo. Atenção:
cada `ctbz login` gera um e-mail novo, e só o código mais recente vale.

## Comandos

```sh
ctbz status              # dados da sessão e se ela ainda é válida
ctbz empresa             # resumo da empresa selecionada
ctbz empresa --json      # resposta completa de /api/plataforma/dadosempresa/get
ctbz api menu/get        # qualquer endpoint; caminho relativo → /api/plataforma/
ctbz api -X POST -d @corpo.json /api/plataforma/algum/endpoint
ctbz logout              # apaga a sessão local
```

Códigos de saída: `0` sucesso, `1` erro, `3` login pendente (aguardando OTP ou CNPJ).

## Como funciona (engenharia reversa)

1. `GET /login` → formulário com um `token` anti-CSRF.
2. `POST /login` (`user`, `password`, `token`) → tela de OTP e cookie de pré-sessão `__C`.
   Erros voltam como redirecionamento para `/login#incorreto`, `#bloqueado`, `#limite-sessoes`…
3. `POST /login` com corpo `v=2&otp=NNNNNN` (text/plain, como o `fetch` do navegador):
   `206` = próxima tela (seleção de empresa), `200` = login concluído, `404` = código inválido.
4. `POST /selecionarempresa` (`cnpj`, `token`) → cookie `oauth-token` e um `<script>` que grava
   no `localStorage` as chaves `l` (login), `r` (usuário) e `e` (empresa), em base64 de JSON
   passado por `encodeURIComponent`. A CLI decodifica e guarda esses dados na sessão.
5. O painel (`/painel-de-controle/`) chama as APIs com os cookies (`withCredentials`).
   Bases encontradas no bundle do front:

   | Base                    | Uso                                   |
   |-------------------------|---------------------------------------|
   | `/api/plataforma/`      | BFF do painel (padrão do `ctbz api`)  |
   | `/api/legado/`          | monólito                              |
   | `/api/public/`          | endpoints públicos                    |
   | `/api/multiusuario/`    | usuários e convites                   |
   | `/api/pagamentos/`      | pagamentos                            |
   | `/api/fintech/`         | conta digital                         |

   Endpoints `GET` do BFF já testados: `dadosempresa/get`, `menu/get`, `appbar/get`,
   `dashboard/situacao-app`, `atendimento/chamados?em-andamento=true`.
   Outros vistos no front (não testados): `impostos/v5/impostos-a-pagar/guia/{id}`,
   `relatorios-ms/gerar-balancete/`, `prolabore/central/init`, `conciliacao-fiscal/pendencias`,
   `certificado/status`, `notafiscal/buscarNotaEmpresa`…

As sessões duram algumas horas (o cookie `__C` expira cerca de 2 h após o login).

## Desenvolvimento

```sh
go test ./...
```

Os testes simulam o fluxo de login com um servidor local, sem tocar na Contabilizei.
