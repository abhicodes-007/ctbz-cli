# ADR-0002: Somente leitura até a 1.0

- **Status:** substituída por [ADR-0018](0018-escrita-com-confirmacao.md)
- **Data:** 2026-10-03

## Contexto

A CLI opera a contabilidade real de uma empresa. A API é privada e não documentada (obtida
por engenharia reversa), então o efeito exato de uma escrita (emitir nota, confirmar
pagamento, aceitar termo) não pode ser testado sem consequências reais.

## Decisão

Até a versão 1.0, os comandos de domínio só fazem `GET`. Ações que alteram dados ficam fora
do roadmap (seção "Depois da 1.0" do [ROADMAP](https://github.com/edusouza/ctbz-cli/blob/main/ROADMAP.md)).
A exceção é `ctbz api -X MÉTODO`, a ferramenta genérica de baixo nível, que continua
aceitando qualquer método por decisão explícita do usuário.

## Consequências

- Os clientes de domínio em `internal/ctbz` só expõem leituras.
- Investigações de API usam apenas `GET`; endpoints de `GET` com efeito colateral pelo nome
  (`alterarParametroEmpresa`, `cancelarAPI`) ou que expõem segredos (`certificado/senha`) não são chamados.

## Alternativas consideradas

- **Escrita com confirmação (`--yes`, `--dry-run`)** — adiada; reavaliar depois da 1.0.
