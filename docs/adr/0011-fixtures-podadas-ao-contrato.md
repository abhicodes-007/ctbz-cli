# ADR-0011: Fixtures podadas ao contrato

- **Status:** aceita
- **Data:** 2026-10-03
- **Complementa:** [ADR-0010](0010-testes-de-contrato.md)

## Contexto

Algumas respostas da API trazem muito mais dados pessoais do que a CLI usa. A lista de sócios
(`socio/list`), por exemplo, tem ~100 campos: filiação, RG, CNH, título de eleitor, dependentes.
A anonimização por regras ([ADR-0010](0010-testes-de-contrato.md)) não garante cobrir campos
com nomes inesperados.

## Decisão

`tools/capture` passa a **podar** a resposta antes de anonimizar: a fixture guarda só os campos
que o tipo de contrato declara (`contract.Prune`). A flag `-full` mantém a resposta inteira
quando for útil investigar campos novos.

## Consequências

- Fixtures menores e com uma fração dos dados para revisar.
- O teste offline verifica exatamente o contrato. Campos novos da API só aparecem no modo ao
  vivo (`CTBZ_CONTRACT_LIVE=1`) e no relatório da captura.
- Para usar um campo novo, acrescente-o ao tipo e capture a fixture de novo.

## Alternativas consideradas

- **Só anonimizar** — depende de prever o nome de todo campo pessoal.
- **Fixtures escritas à mão** — perdem a fidelidade à resposta real.
