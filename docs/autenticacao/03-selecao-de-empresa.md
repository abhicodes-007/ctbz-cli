# 3. Seleção de empresa

Quando o usuário tem mais de uma empresa (inclusive inativas), a verificação do OTP
responde **206** com um fragmento HTML que o front injeta em `#content`:

```html
<div class="title">Selecione a empresa</div>
<form action="/selecionarempresa" method="post">
  <label class="mdl-radio" for="11111111000111-1">
    <input type="radio" id="11111111000111-1" class="mdl-radio__button" name="cnpj" value="11111111000111">
    <div><b><span class="mdl-radio__label">11111111000111</span></b></div>
    <div><span>FULANO DE TAL</span></div>                       <!-- razão social -->
  </label>
  <label class="mdl-radio" for="22222222000122-1">
    <input type="radio" … name="cnpj" value="22222222000122">
    <div><span>FULANO TECNOLOGIA LTDA</span></div>
  </label>
  <input type="hidden" name="token" value="<token>" />        <!-- mesmo token do /login -->
  <button type="submit">Selecionar</button>
</form>
```

A lista mostra todas as empresas vinculadas ao CPF, inclusive as **inativas** (no teste,
um MEI antigo aparecia ao lado da LTDA ativa). A situação de cada uma só aparece depois,
em `dadosempresa/get` (campo `empresas[].statusEmpresa`).

## `POST /selecionarempresa`

```http
POST /selecionarempresa HTTP/2
Content-Type: application/x-www-form-urlencoded
Cookie: __C=…

cnpj=22222222000122&token=<token>
```

Resposta **200**:

```http
set-cookie: oauth-token=<base64>; Path=/; HttpOnly; Secure; SameSite=Lax; Expires=<login + 4h>
```

```html
<script>
  localStorage.clear(); sessionStorage.clear();
  localStorage.setItem("l", "JTdCJTIydG9rZW4…");   // dados do login
  localStorage.setItem("r", "JTdCJTIyaWQlMjI…");   // responsável (usuário)
  localStorage.setItem("e", "JTdCJTIyZnJhbnF…");   // empresa selecionada
  location.replace("/painel-de-controle/#/home");
</script>
```

O conteúdo de `l`, `r` e `e` está descrito em [04-sessao-e-cookies.md](04-sessao-e-cookies.md).

## Na CLI

- A empresa vem de `--cnpj` ou `CTBZ_CNPJ` (aceita com ou sem pontuação).
- Com uma só empresa na lista, ela é escolhida automaticamente.
- Num terminal, a CLI mostra a lista numerada; sem terminal, o login fica pendente e
  termina com `ctbz login --cnpj <CNPJ>`.
- Um CNPJ fora da lista é recusado antes de ir ao servidor.

## Trocar de empresa

Não foi encontrado endpoint para trocar de empresa sem refazer o login. O painel exibe
`dadosempresa/get → empresas[]`, mas a troca não foi investigada. Hoje, trocar exige
novo login com outro `--cnpj` (e um novo OTP).
