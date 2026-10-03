# Autenticação

O login da Contabilizei é **server-side, baseado em formulários HTML e cookies**: não
existe endpoint JSON de login nem token Bearer. A CLI reproduz exatamente o que o
navegador faz, em quatro etapas:

| # | Etapa | Documento |
|---|---|---|
| 1 | Formulário de usuário e senha | [01-credenciais.md](01-credenciais.md) |
| 2 | Código de verificação (OTP) enviado por e-mail | [02-otp.md](02-otp.md) |
| 3 | Seleção da empresa (quando o usuário tem mais de uma) | [03-selecao-de-empresa.md](03-selecao-de-empresa.md) |
| 4 | Cookies, dados de sessão (`localStorage`) e expiração | [04-sessao-e-cookies.md](04-sessao-e-cookies.md) |

## Sequência

```mermaid
sequenceDiagram
    participant C as CLI / navegador
    participant S as app.contabilizei.com.br
    participant M as E-mail do usuário

    C->>S: GET /login
    S-->>C: 200 HTML (form-login, input hidden "token")
    C->>S: POST /login (user, password, token)
    S-->>C: 200 HTML "Insira o código de verificação" + Set-Cookie __C
    S->>M: e-mail de seguranca@contabilizei.com.br com 6 dígitos
    C->>S: POST /login  body "v=2&otp=NNNNNN"  (Cookie __C)
    alt código errado ou vencido
        S-->>C: 404
    else mais de uma empresa
        S-->>C: 206 fragmento "Selecione a empresa"
        C->>S: POST /selecionarempresa (cnpj, token)
        S-->>C: 200 <script> localStorage + Set-Cookie oauth-token
    else uma empresa (comportamento inferido do JS)
        S-->>C: 200 script final
    end
    C->>S: GET /api/plataforma/... (Cookie __C + oauth-token)
    S-->>C: 200 JSON
```

## Pontos de atenção

- **Um código por login.** Cada `POST /login` com credenciais gera e envia um código
  novo; códigos anteriores deixam de valer. Dois logins seguidos geram dois e-mails
  quase simultâneos e é fácil usar o errado (aconteceu durante a investigação).
- **Limite de sessões simultâneas.** O servidor pode recusar com `#limite-sessoes`.
  Não há endpoint de logout no servidor (ver [04](04-sessao-e-cookies.md)), então
  logins repetidos em sequência devem ser evitados.
- **Bloqueio temporário.** Muitas tentativas erradas levam a `#bloqueado`. A CLI
  nunca repete automaticamente um OTP recusado.
- **Login com Google** existe (`/loginGoogle?token=…`), mas não foi explorado: a CLI
  usa apenas usuário/senha.
