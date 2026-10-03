# ctbz notas tomadores

Lista os tomadores (clientes) cadastrados no emissor de notas

Lista os tomadores cadastrados no emissor de notas: nome, documento, e-mail, telefone,
inscrição municipal, município, UF e se é do exterior.

Para os dados cadastrais de um CNPJ qualquer na Receita (como o emissor faz ao cadastrar um
tomador), use "ctbz notas tomadores consulta CNPJ".

## Uso

```
ctbz notas tomadores
```

## Exemplos

```sh
  ctbz notas tomadores
  ctbz notas tomadores consulta 00.000.000/0001-91
```

## Subcomandos

- [`ctbz notas tomadores consulta`](ctbz_notas_tomadores_consulta.md): Consulta os dados cadastrais de um CNPJ na Receita

## Flags globais

```
  -h, --help            mostra a ajuda do comando
  -o, --output string   formato de saída: table, json ou csv (padrão: CTBZ_OUTPUT ou table)
```

Veja também: [`ctbz notas`](ctbz_notas.md).
