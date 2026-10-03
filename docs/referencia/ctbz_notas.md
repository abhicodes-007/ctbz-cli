# ctbz notas

Lista as notas fiscais de serviço (NFS-e) emitidas

Lista as NFS-e emitidas pelo emissor da Contabilizei: número, data de emissão, tomador,
documento, valor, status e situação. Na tabela, o total do período vai para o stderr.

O período vai de --de até --ate (meses AAAA-MM; padrão: o mês atual), no máximo 24 meses.
--tomador filtra pelo nome do tomador ou, se for um CPF/CNPJ, pelo documento; --numero
filtra pelo número da nota.

O PDF e o XML da nota não estão disponíveis pela API: a Contabilizei os envia por e-mail
quando a prefeitura autoriza a nota.

## Uso

```
ctbz notas [flags]
```

## Exemplos

```sh
  ctbz notas
  ctbz notas --de 2026-01 --ate 2026-09 -o csv > notas.csv
  ctbz notas --tomador 11.222.333/0001-81
  ctbz notas --tomador "ACME" --de 2026-07
```

## Flags

```
      --ate string       último mês, AAAA-MM (padrão: --de ou o mês atual)
      --de string        primeiro mês, AAAA-MM (padrão: o mês atual)
      --numero string    número da nota
      --tomador string   nome do tomador ou CPF/CNPJ
```

## Flags globais

```
  -h, --help            mostra a ajuda do comando
  -o, --output string   formato de saída: table, json ou csv (padrão: CTBZ_OUTPUT ou table)
```

Veja também: [`ctbz`](ctbz.md).
