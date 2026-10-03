# ADR-0007: Login reproduzindo o navegador, sessão em arquivo e OTP por comando externo

- **Status:** aceita
- **Data:** 2026-10-03

## Contexto

A Contabilizei não tem API de login: o navegador envia formulários, recebe um código de 6
dígitos por e-mail e escolhe a empresa. A sessão são dois cookies (`__C`, `oauth-token`)
válidos por cerca de 2 horas. Detalhes em [Autenticação](../autenticacao/README.md).

## Decisão

- O cliente HTTP reproduz o fluxo do navegador (mesmos campos, `Content-Type` e cabeçalhos),
  com cookies e redirecionamentos tratados manualmente para poder salvá-los e ler o fragmento
  de erro (`/login#incorreto`).
- A sessão fica em `$CTBZ_HOME/session.json` (`0600`); credenciais vêm só de variáveis de ambiente
  e nunca são gravadas.
- O OTP vem de um **comando externo** (`CTBZ_OTP_CMD`) executado com polling, do terminal ou de
  uma segunda chamada (`ctbz login --otp`), com o login intermediário salvo em `pending.json`.
  A CLI não fala com provedores de e-mail diretamente.
- Com `CTBZ_OTP_CMD` definido, uma resposta 401 dispara re-login automático e a chamada é repetida uma vez.

## Consequências

- Funciona com qualquer provedor de e-mail (Gmail via `gws`, IMAP, webhook…).
- Mudanças no HTML de login quebram o fluxo; os testes de contrato e o modo ao vivo detectam isso.

## Alternativas consideradas

- **Navegador headless** — pesado e frágil para uma CLI.
- **Integração nativa com Gmail API** — prende a um provedor e exige OAuth do Google na CLI.
