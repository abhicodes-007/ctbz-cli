# 2. Código de verificação (OTP)

## A tela

Depois das credenciais, o `POST /login` devolve uma página com o container
`<div id="content">`, que é trocado via `fetch` nas etapas seguintes:

```html
<span class="title">Insira o código de verificação</span>
<span class="subtitle">
  Enviamos um código de 6 dígitos para o seu e-mail: <strong>fu***@ex***com</strong>
</span>
<input class="input-token" type="text" id="otp" name="otp" maxlength="6">
<div id="otp-err">O código está incorreto ou inválido!</div>
<button id="confirm"     onclick="verify(2, this.id)">Confirmar</button>
<button id="resend-code" onclick="send(2, this.id)" disabled>Reenviar código</button>
<button id="cancel"      onclick="send(0, this.id)">Voltar</button>
```

- O e-mail vem de **`seguranca@contabilizei.com.br`** e traz um código de **6 dígitos**.
- A CLI reconhece esta tela pela presença de `id="otp"` e extrai o e-mail mascarado
  de `<strong>` para mostrar ao usuário.
- O botão de reenvio fica desabilitado por 60 s (só no front).

## Verificação: `verify(m)`

JavaScript da página (resumido):

```js
const res = await fetch("/login", { method: "POST", body: "v=" + m + "&otp=" + $("otp").value });
if (res.status === 206)      { $("content").innerHTML = await res.text(); rd_check(); } // próxima tela
else if (res.status === 200) { eval-like: <script> com o texto da resposta }            // login concluído
else if (res.status === 404) { mostra "O código está incorreto ou inválido!" }
else                         { "Sistema temporariamente indisponível" }
```

Requisição equivalente:

```http
POST /login HTTP/2
Content-Type: text/plain;charset=UTF-8      ← fetch() com body string
Cookie: __C=<pré-sessão>

v=2&otp=123456
```

| Status | Significado | Ação da CLI |
|---|---|---|
| **206** | Novo fragmento de tela, observado: "Selecione a empresa" | interpreta o HTML ([03](03-selecao-de-empresa.md)) |
| **200** | Script final (grava `localStorage` e redireciona) | extrai a sessão ([04](04-sessao-e-cookies.md)) |
| **404** | Código incorreto ou expirado | `ErrInvalidOTP`; o login fica pendente para outra tentativa |
| outro | Falha do sistema | erro |

O valor `v=2` corresponde ao canal e-mail. A função `rd_check()` marca o primeiro
`<input class="mdl-radio__button">`, usada tanto na seleção de empresa quanto,
provavelmente, numa tela de escolha de canal (e-mail/SMS) que não apareceu nos testes.

## Reenvio e cancelamento: `send(m)`

```js
fetch("/login", { method: "POST", body: "s=" + m })
```

| `s=` | Efeito (pelo JS) |
|---|---|
| `0` | Volta: o front recarrega `/login` (descarta o login em andamento) |
| `2` | Reenvia o código por e-mail e reinicia o contador de 60 s |
| `3` | Reenvio por outro canal (há botões `resend-sms`/`resend-email` em variantes da tela) |

A CLI não usa `send`: para pedir outro código ela recomeça o login
(`ctbz login --restart`). O reenvio (`s=2`) não foi testado.

## Validade

- O código vale só para o login que o gerou. Observado: com dois logins em sequência
  (cerca de 1 minuto entre eles), o código recebido primeiro foi recusado com 404 e o do
  login seguinte funcionou. Não foi testado se o código antigo vale enquanto nenhum outro
  login acontece.
- O tempo exato de expiração não foi medido. Um código usado cerca de 1 minuto depois
  do envio funcionou. A pré-sessão `__C` dura 2 h.
- Não foi testado quantos erros levam a `#bloqueado`.
