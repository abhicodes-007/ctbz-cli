# Front-end

## Aplicações

O site é um conjunto de fronts independentes, servidos em subcaminhos do mesmo host.
Todos dependem dos cookies de sessão; sem eles, o HTML redireciona para `/login`.

| Caminho | O que é |
|---|---|
| `/login`, `/selecionarempresa` | páginas server-side (JSP: há links como `view-recuperar-senha.jsp`) |
| `/painel-de-controle/` | **painel principal** (Vue 2 + Vuex + axios, rotas por hash `#/home`) |
| `/sistema/#/` | front antigo (ex.: "Como emitir notas", "Importar notas fiscais") |
| `/plataforma/` | plataforma nova |
| `/conciliacao/#/` | extratos e conciliação bancária |
| `/nota-entrada/#/`, `/nota-tomada/#/` | notas de entrada / tomadas |
| `/autopilot/#/` | autopilot |
| `/onboarding/…`, `/onboarding-plataforma/#/` | onboarding |
| `/checkout/#/` | checkout de planos |
| `/widget-window/` | widget de CRM |

O `menu/get` informa, para cada item do menu, o app de destino (`application`) e a rota
(`route`), o que dá um mapa de navegação pronto.

## Bundles do painel

`GET /painel-de-controle/` traz dois scripts de entrada e ≈165 chunks em
`<link rel="prefetch">`:

```text
js/chunk-vendors.<hash>.js   ≈2,3 MB  bibliotecas
js/app.<hash>.js             ≈680 KB  núcleo do app, store, axios
js/<Tela>.<hash>.js          chunks por tela (CentralDeRotinas, NovoEmissorNF, balancete…)
```

- Os chunks **também exigem cookie**: sem sessão, recebem o HTML de redirecionamento
  (106 bytes).
- Os nomes dos chunks seguem as telas (`balancete`, `balanco`, `certificadoDigital`,
  `central-de-documentos`, `NovoEmissorNF`, `Multiusuario`, …) e servem de índice
  para achar uma funcionalidade.

## Configuração exposta (`process.env` embutido)

Valores `VUE_APP_*` ficam em texto no `app.<hash>.js`. Os relevantes:

| Chave | Valor |
|---|---|
| `VUE_APP_CONTABILIZEI_BFF_URL` | `/api/plataforma/` |
| `VUE_APP_CONTABILIZEI_MONOLITO_URL` | `/api/legado/` |
| `VUE_APP_CONTABILIZEI_MONOLITO_PUBLIC_URL` / `VUE_APP_BACK_PUBLIC_BASE_URL` | `/api/public/` |
| `VUE_APP_CONTABILIZEI_MULTIUSUARIO_URL` | `/api/multiusuario/` |
| `VUE_APP_CONTABILIZEI_PAGAMENTO` | `/api/pagamentos/` |
| `VUE_APP_FINTECH_BASE_URL` | `/api/fintech/` |
| `VUE_APP_BASE_URL` | `https://app.contabilizei.com.br/painel-de-controle` |
| `VUE_APP_VERSION` | `2.0.88` (na data da investigação) |
| `VUE_APP_INTERCEPT_401` | liga o interceptor que manda para `/login#error-session` |

Também aparecem integrações de terceiros (Datadog RUM, Mixpanel via proxy, Appcues,
Iugu, Adyen, BS2, Pluggy, Zendesk), irrelevantes para a CLI.

## Cliente HTTP

Um módulo webpack cria as instâncias axios e as exporta por letra:

```js
// módulo de instâncias (exports a–f)
d = axios.create({ baseURL: "/api/legado/",       withCredentials: true })
c = axios.create({ baseURL: "/api/plataforma/",   withCredentials: true })
e = axios.create({ baseURL: "/api/multiusuario/", withCredentials: true })
f = axios.create({ baseURL: "/api/fintech/",      withCredentials: true })
b = axios.create({ baseURL: "/api/legado/",       withCredentials: true })
a = axios.create({ baseURL: "/api/leads/hubspot/", withCredentials: false })
```

Nos demais módulos as chamadas aparecem como `X["c"].get("dashboard/fatura")`, em que
`X` é o nome local do módulo importado e a letra indica a base. É nisso que se apoia o
[gerador de catálogo](https://github.com/edusouza/ctbz-cli/blob/main/scripts/extrair-endpoints.py). Alguns serviços guardam o
caminho numa variável antes (`var e="/novo-emissor/listagem/notas"; return u["c"].get(e)`),
e o gerador trata esse padrão também.

Interceptors (se `VUE_APP_INTERCEPT_401`):

```js
request:  config => ({ ...config, credentials: true })
response: erro 401 (ou falha de rede online) → location.replace("/login?redirect=…#error-session")
```

Não há injeção de cabeçalho de autenticação: o navegador manda os cookies `HttpOnly` sozinho.

## Logout

`store.dispatch("auth/logout")` limpa o estado local e vai para `/login?redirect=…`.
Nenhuma requisição ao servidor (ver [sessão](../autenticacao/04-sessao-e-cookies.md#logout)).
