# ctbz mensalidade historico

Lista os pagamentos anteriores da mensalidade e a situação do débito automático

Lista os pagamentos anteriores da mensalidade (data, competência, valor e status) e mostra a
situação do débito automático: se está habilitado e ativo, a competência e a próxima cobrança.

Só campos conhecidos são copiados para a saída: chaves do gateway de pagamento e dados de cartão
nunca são mostrados. O formato de cada pagamento não é documentado; se nenhum campo for
reconhecido, o comando avisa e indica "ctbz api payments/recorrencia/historico".

## Uso

```
ctbz mensalidade historico
```

## Exemplos

```sh
  ctbz mensalidade historico
  ctbz mensalidade historico -o json | jq '.pagamentos[].valor'
```

## Flags globais

```
  -h, --help            mostra a ajuda do comando
  -o, --output string   formato de saída: table, json ou csv (padrão: CTBZ_OUTPUT ou table)
```

Veja também: [`ctbz mensalidade`](ctbz_mensalidade.md).
