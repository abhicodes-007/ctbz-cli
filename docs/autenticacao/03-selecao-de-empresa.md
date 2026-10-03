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

O seletor de empresas do painel não usa uma chamada de API: ele navega para

```http
GET /api/public/requestselecaoempresa?redirect=/&cnpj=22222222000122
```

que responde `302` para um **host de SSO separado**:

```http
Location: https://sso.contabilizei.com/selecionarempresa?token=<token>&cnpj=22222222000122
```

O SSO tem sessão própria (cookie `JSESSIONID` em `sso.contabilizei.com`), que a CLI não tem:
o login da CLI acontece em `app.contabilizei.com.br`. Sem essa sessão, o SSO devolve a tela de
login (verificado em 03/10/2026, pedindo a troca para a própria empresa atual; a sessão antiga
continuou válida).

Por isso `ctbz empresa usar CNPJ` troca de empresa **refazendo o login** com `--cnpj`, o que
envia um novo código OTP. Com `CTBZ_OTP_CMD` definido, a troca é automática.

A lista de empresas para escolher vem de `dadosempresa/get → empresas[]` (`ctbz empresas`).
