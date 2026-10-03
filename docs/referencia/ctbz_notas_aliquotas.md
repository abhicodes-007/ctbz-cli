# ctbz notas aliquotas

Lista as alíquotas e os códigos de serviço por atividade

Lista, para cada atividade (CNAE) da empresa, o item da lista de serviços usado na nota, a
alíquota do Simples (%), a parte que é ISS (%), o Fator R considerado (%) e se o anexo é
fixo. As alíquotas são separadas por mercado: tomadores no Brasil (interno) e no exterior
(externo, sem ISS).

## Uso

```
ctbz notas aliquotas
```

## Exemplos

```sh
  ctbz notas aliquotas
  ctbz notas aliquotas -o csv
```

## Flags globais

```
  -h, --help            mostra a ajuda do comando
  -o, --output string   formato de saída: table, json ou csv (padrão: CTBZ_OUTPUT ou table)
```

Veja também: [`ctbz notas`](ctbz_notas.md).
