# ctbz impostos

Lista as guias de impostos a pagar (em atraso, do mês e do próximo mês)

Lista as guias de impostos a pagar, agrupadas em atraso, deste mês e do próximo mês,
com competência, vencimento, valor e situação. Na tabela, os totais por grupo vão para stderr.

Com --fail-on-atraso, o comando termina com código 4 quando há guias em atraso (útil em
scripts e alertas).

## Uso

```
ctbz impostos [flags]
```

## Exemplos

```sh
  ctbz impostos
  ctbz impostos --atrasadas
  ctbz impostos -o csv > guias.csv
  ctbz impostos --fail-on-atraso || notify-send "Há impostos em atraso"
```

## Flags

```
      --atrasadas        mostra só as guias em atraso
      --fail-on-atraso   termina com código 4 se houver guias em atraso
```

## Flags globais

```
  -h, --help            mostra a ajuda do comando
  -o, --output string   formato de saída: table, json ou csv (padrão: CTBZ_OUTPUT ou table)
```

Veja também: [`ctbz`](ctbz.md).
