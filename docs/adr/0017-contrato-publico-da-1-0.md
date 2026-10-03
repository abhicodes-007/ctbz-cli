# ADR-0017: O que a 1.0 garante (contrato público da CLI)

- **Status:** aceita
- **Data:** 2026-10-03

## Contexto

Até a 0.8, mudanças incompatíveis subiam o *minor* (ADR-0003). Com a 1.0, scripts e alertas
em cron passam a depender da CLI, e o SemVer exige dizer o que é "API pública": uma
mudança incompatível nela só pode sair numa versão *major*.

## Decisão

A partir da 1.0.0, são API pública e só mudam de forma incompatível numa 2.0:

- nomes de comandos e subcomandos, e nomes e significado das flags;
- as **chaves** (`Key`) de JSON e CSV, os tipos dos valores (número, texto, data ISO,
  `null`) e os valores fixos documentados (ex.: `alerta` = `vencida`/`próxima`/`crítica`);
- os códigos de saída (0, 1, 2, 3, 4) e o que cada `--fail-on-*` sinaliza;
- variáveis de ambiente (`CTBZ_*`) e o contrato do `CTBZ_OTP_CMD`.

Não são API pública: a formatação da tabela (rótulos, alinhamento, recuo), mensagens em
stderr, o texto de ajuda, a ordem das linhas quando a API não garante ordem, o formato de
`session.json`/`pending.json` e os pacotes Go em `internal/`.

Acrescentar comandos, flags, colunas ou campos é *minor*; correções são *patch*.

## Consequências

- Quando a Contabilizei mudar a API (o monitoramento abre issue, ADR-0016), a CLI se adapta
  mantendo as chaves de saída; se não der, a mudança espera uma 2.0 ou vira uma chave nova.
- Ações de escrita (emitir nota, enviar documento…) entram como comandos novos, em *minors*
  depois da 1.0 (ver "Depois da 1.0" no ROADMAP).

## Alternativas consideradas

- **Tratar a tabela como contrato** — travaria melhorias de legibilidade; quem automatiza usa
  `-o json` ou `-o csv`.
- **Ficar em 0.x indefinidamente** — não dá garantia nenhuma a quem automatiza.
