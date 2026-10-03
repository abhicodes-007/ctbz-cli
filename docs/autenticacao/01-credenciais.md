# 1. Credenciais (usuário e senha)

## `GET /` e `GET /login`

`GET /` devolve só um redirecionamento em JavaScript:

```html
<script>location.replace("/login?redirect="+encodeURIComponent(location.pathname+location.hash))</script>
```

`GET /login` devolve a página de login (≈36 KB, imagens inline em base64). O que importa:

```html
<script>
  localStorage.clear();
  sessionStorage.clear();
  function logarComGoogle() { window.location = '/loginGoogle?token=<token>'; }
</script>
<form class="form" id="form-login" action="/login" method="post">
  <input type="text"     name="user"     id="user" required>   <!-- e-mail ou CPF -->
  <input type="password" name="password" id="password" required>
  <input type="hidden"   name="token"    value="<token>"/>   <!-- anti-CSRF -->
</form>
```

- O `token` é um UUID em base64, diferente a cada carregamento. Ele é reutilizado
  depois na seleção de empresa.
- Nenhum cookie é emitido neste `GET`.
- O campo `user` aceita **e-mail ou CPF** (só dígitos).

## `POST /login` (credenciais)

```http
POST /login HTTP/2
Content-Type: application/x-www-form-urlencoded
Origin: https://app.contabilizei.com.br
Referer: https://app.contabilizei.com.br/login

user=00000000000&password=********&token=<token>
```

Resposta de sucesso: **200** com a tela de OTP (ver [02-otp.md](02-otp.md)) e o cookie
de pré-sessão:

```http
set-cookie: __C=<96 caracteres>; Path=/; HttpOnly; Secure; SameSite=Lax; Expires=<agora + 2h>
```

## Erros

A página de login mostra mensagens a partir do **fragmento da URL**. O servidor sinaliza
falhas redirecionando para `/login#<código>`. Tabela extraída do JavaScript da página:

| Fragmento | Mensagem exibida |
|---|---|
| `#incorreto` | E-mail ou senha incorretos. |
| `#bloqueado` | Acesso bloqueado temporariamente por medida de segurança. |
| `#limite-sessoes` | Você atingiu o limite de acessos simultâneos. |
| `#sem-permissao` | Sua empresa ainda não está pronta para usar a plataforma. |
| `#expirou` | Para sua segurança, sua sessão foi encerrada. |
| `#google` | Conta Google não permitida. |
| `#problema` | Encontramos um problema no seu acesso. |
| `#error-session` | (usado pelo painel ao receber 401 de uma API) |

A CLI desliga o seguimento automático de redirecionamentos para conseguir ler o
fragmento do `Location` e transformá-lo em erro legível
(`internal/ctbz/login.go`, `followRedirects`).

A página também tem um ponto para mensagem vinda do backend (`getBackendErrorMessage`),
vazio em todos os testes.

## Cabeçalhos

Os testes usaram `User-Agent` de Chrome, `Origin` e `Referer` iguais aos do navegador.
Não foi testado se o servidor rejeita requisições sem eles; a CLI os envia por precaução.
A borda é da Azion (`x-azion-request-id`); nenhum desafio anti-bot apareceu.
