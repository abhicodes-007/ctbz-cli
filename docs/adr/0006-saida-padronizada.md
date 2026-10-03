# ADR-0006: Saída padronizada com List/Record e tipos de valor

- **Status:** aceita
- **Data:** 2026-10-03

## Contexto

Cada comando imprime dados que precisam servir a pessoas (terminal), scripts (`jq`) e
planilhas. Formatar à mão em cada comando gera inconsistência (datas e valores em formatos
diferentes) e multiplica o trabalho por três.

## Decisão

- Comandos descrevem os dados como `output.List` (várias linhas) ou `output.Record` (um registro)
  e chamam `output.Write`. Formatos: `table` (padrão), `json`, `csv`.
- Valores semânticos usam tipos próprios — `output.Money`, `output.Date`, `output.DateTime`,
  `output.CNPJ` — e cada formato decide a apresentação (`R$ 1.234,56` × `1234.56`).
- Chaves (`Key`) são `snake_case` em português sem acento e formam um **contrato público** com
  scripts: renomear uma chave é mudança incompatível (ver [ADR-0003](0003-versionamento-e-changelog.md)).
- Precedência do formato: `-o` do comando > `CTBZ_OUTPUT` > `table`. `ctbz api` usa JSON por padrão
  e ignora `CTBZ_OUTPUT`.
- stdout só recebe dados; mensagens, progresso e avisos vão para stderr.
- CSV de valores aninhados: JSON dentro da célula; listas de texto separadas por `; `.

## Consequências

- Comando novo ganha os três formatos sem código extra.
- Golden files em `internal/output/testdata/` protegem o formato.

## Alternativas consideradas

- **`text/template` por comando** — flexível, mas sem formatos de máquina consistentes.
- **YAML** — pouco usado com `jq`/planilhas; pode ser acrescentado depois sem quebrar nada.
