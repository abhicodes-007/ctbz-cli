# ADR-0003: SemVer, Keep a Changelog e uma versão por épico

- **Status:** aceita
- **Data:** 2026-10-03

## Contexto

O roadmap divide o trabalho em versões por contexto (v0.1 login e empresa, v0.2 impostos…),
cada uma com um épico e várias sub-issues entregues em PRs separados.

## Decisão

- Versões seguem [SemVer 2.0](https://semver.org/lang/pt-BR/). Antes da 1.0, uma mudança
  incompatível incrementa o *minor*.
- Cada épico do roadmap é uma versão *minor* (`0.2.0` = Impostos). Correções depois de
  publicada a versão viram *patch*.
- O [CHANGELOG](https://github.com/edusouza/ctbz-cli/blob/main/CHANGELOG.md) segue
  [Keep a Changelog 1.1](https://keepachangelog.com/pt-BR/1.1.0/): cada PR acrescenta sua linha
  em `[Unreleased]`, na seção certa (`Added`, `Changed`, `Fixed`, `Removed`…), em português.
- O último PR de um épico "corta" a versão: move `[Unreleased]` para `[X.Y.0] - data` e
  atualiza os links de comparação.
- A tag `vX.Y.0` é criada no commit de merge do PR de release em `main`, não antes:
  PRs empilhados ainda não estão em `main`.
- O binário informa a versão por `-ldflags "-X main.version=…"`; sem isso, mostra a versão do
  módulo (`debug.ReadBuildInfo`, cobre `go install …@vX.Y.Z`) ou `dev`.

## Consequências

- O CHANGELOG é escrito para usuários: descreve comportamento, não commits.
- Nenhum PR de funcionalidade mexe em número de versão, só o PR de release do épico.

## Alternativas consideradas

- **Versão por PR** — dezenas de versões sem significado para o usuário.
- **Changelog gerado de commits** — mistura detalhes internos; Keep a Changelog pede curadoria.
