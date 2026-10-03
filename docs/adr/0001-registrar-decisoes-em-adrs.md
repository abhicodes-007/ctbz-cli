# ADR-0001: Registrar decisões em ADRs

- **Status:** aceita
- **Data:** 2026-10-03

## Contexto

O projeto é desenvolvido em PRs pequenos e empilhados, muitas vezes por um agente que não
guarda memória entre sessões. Sem registro, a mesma escolha (biblioteca, formato, nome de
campo) seria refeita a cada funcionalidade, com risco de inconsistência.

## Decisão

Toda decisão que afete mais de um PR, ou que tenha alternativas razoáveis, vira um ADR em
`docs/adr/NNNN-titulo.md`, usando o [modelo](template.md). Convenções do dia a dia (nomes,
estilo, testes) ficam em [Contribuindo](../contribuindo.md).

- ADRs são imutáveis depois de aceitos: mudar de ideia é um ADR novo que substitui o antigo.
- O ADR entra no mesmo PR da primeira implementação que depende dele.
- Antes de decidir algo, consultar o [índice](README.md).

## Consequências

- Revisores entendem o porquê sem reconstruir a discussão.
- Custo pequeno: um arquivo curto por decisão.

## Alternativas consideradas

- **Decisões só em descrições de PR** — se perdem e não aparecem na documentação publicada.
