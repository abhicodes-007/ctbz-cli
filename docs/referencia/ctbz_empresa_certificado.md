# ctbz empresa certificado

Mostra a situação do certificado digital da empresa

Mostra a situação do certificado digital (e-CNPJ) usado pela Contabilizei, a data de
vencimento, os dias que faltam e se já é possível renovar.

## Uso

```
ctbz empresa certificado
```

## Exemplos

```sh
  ctbz empresa certificado
  ctbz empresa certificado -o json | jq .dias_para_vencer
```

## Flags globais

```
  -h, --help            mostra a ajuda do comando
  -o, --output string   formato de saída: table, json ou csv (padrão: CTBZ_OUTPUT ou table)
```

Veja também: [`ctbz empresa`](ctbz_empresa.md).
