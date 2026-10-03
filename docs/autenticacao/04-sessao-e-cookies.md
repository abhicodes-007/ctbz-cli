# 4. Sessão, cookies e dados do `localStorage`

## Cookies

| Cookie | Emitido em | Validade observada | Papel |
|---|---|---|---|
| `__C` | `POST /login` (credenciais) | 2 h a partir das credenciais | identifica a sessão desde a pré-autenticação |
| `oauth-token` | `POST /selecionarempresa` | 4 h a partir das credenciais | token de acesso (base64 opaco) |

Ambos com `Path=/; HttpOnly; Secure; SameSite=Lax`.

Testes com `GET /api/plataforma/appbar/get`:

| Cookies enviados | Resultado |
|---|---|
| `__C` + `oauth-token` | 200 |
| só `oauth-token` | 401 |
| só `__C` | 401 |
| nenhum | 401 |

**A validade não desliza:** depois de várias chamadas às APIs, o servidor não reemitiu
os cookies; as datas de expiração continuaram as do login. Na prática a sessão útil é
de **2 horas** (limite do `__C`). A CLI ainda assim absorve qualquer `Set-Cookie` que
vier nas respostas.

O valor de `oauth-token` é o mesmo que aparece em `localStorage.l.token`.

## `localStorage`

O script final grava três chaves. Cada valor é
`base64( encodeURIComponent( JSON ) )`. Para decodificar:

```python
json.loads(urllib.parse.unquote(base64.b64decode(valor).decode()))
```

```go
raw, _ := base64.StdEncoding.DecodeString(v)
s, _ := url.PathUnescape(string(raw))   // PathUnescape: não converte '+' em espaço
```

### `l` — dados do login

```jsonc
{
  "token": "<base64>",                 // = cookie oauth-token
  "userId": 1234567890123456,
  "email": "fulano@exemplo.com",       // e-mail do usuário da plataforma
  "redirectUrl": "/painel-de-controle",
  "acessoInterno": false,
  "empresa": { /* mesmo conteúdo de "e" */ },
  "responsavel": { /* mesmo conteúdo de "r" */ },
  "listaEmpresa": [ /* empresas do usuário */ ]
}
```

### `r` — responsável

```jsonc
{
  "id": 2345678901234567,
  "nome": "FULANO DE TAL",
  "ref": "<hash>",
  "cpf": "00000000000",
  "email": "fulano@gmail.com",          // e-mail que recebe o OTP
  "listaEmpresaOnboarding": [],
  "listaEmpresa": [ /* 2 itens no teste */ ]
}
```

### `e` — empresa selecionada (campos principais)

| Campo | Exemplo / tipo |
|---|---|
| `id` | número (id interno da empresa; usado como `userId` em `l`) |
| `cnpj`, `razaoSocial`, `nomeFantasia` | texto |
| `status.id` | `ATIVO` |
| `regimeTributario`, `optanteSimples` | `SIMPLES`, `true` |
| `ramoAtividadeServico`, `ramoAtividadeComercio`, `ramosAtividade[]` | — |
| `naturezaJuridica` | `{cod, codReceita:"206-2", descricao:"Sociedade Empresária Limitada"}` |
| `planoPagamentoEmpresa.pagtoPlano` | `{descricao, categoria:"BASICO", valor}` |
| `inscricaoMunicipal`, `inscricaoEstadual`, `codIbge`, `uf` | — |
| `endereco` | logradouro, número, bairro, CEP, município (com `codIbge`, `codTributario`) |
| `responsavel`, `socioResponsavel` | nome e CPF |
| `responsabilidadeInicial` | `{mes, ano, periodo:202607, text:"7/2026"}`: início da contabilidade |
| `dataAbertura`, `dataAtivacao` | epoch em ms |
| `configuracao` | certificado digital: `dataValidadeCertificado`, `possuiCertificadoA3`, `proxNumeroRps`… |
| `permiteEmitirNota`, `permiteImportarNota`, `importacaoNotasAutomatica` | flags de nota fiscal |
| `perfil` | `listaMunicipio[]` (26 municípios), regras de emissão por município |
| `jiraIssueId`, `jiraFinalizado`, `timelineFinalizada`, `novoOnboarding`… | controle interno de onboarding |

A CLI guarda as três chaves em `session.json → storage` e usa `l`/`e` para mostrar
empresa e usuário no `ctbz status`.

## Logout

Não há chamada de logout no servidor. A ação "sair" do painel
(`store.dispatch("auth/logout")`) só limpa o estado local e navega para
`/login?redirect=…`; a página de login executa `localStorage.clear()`. A sessão
continua válida no servidor até expirar. `ctbz logout` apenas apaga o
`session.json` local.

## Sessão expirada

Toda API devolve **401**. O painel tem um interceptor axios que, ao receber 401, manda
para `/login?redirect=…#error-session`. A CLI converte em `ErrUnauthorized` e, se
`CTBZ_OTP_CMD` estiver definido, refaz o login sozinha e repete a chamada.
