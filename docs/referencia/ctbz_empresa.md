# ctbz empresa

Mostra os dados da empresa selecionada

Mostra razão social, CNPJ, situação, regime tributário, plano, certificado digital
e as outras empresas do usuário.

A resposta crua da API está em "ctbz api dadosempresa/get".

## Uso

```
ctbz empresa [flags]
```

## Exemplos

```sh
  ctbz empresa
  ctbz empresa -o json | jq -r .certificado_validade
```

## Flags

```
      --json   atalho para -o json
```

## Flags globais

```
  -h, --help            mostra a ajuda do comando
  -o, --output string   formato de saída: table, json ou csv (padrão: CTBZ_OUTPUT ou table)
```

Veja também: [`ctbz`](ctbz.md).
