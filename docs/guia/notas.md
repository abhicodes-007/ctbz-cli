# Notas fiscais

## Notas emitidas (NFS-e)

```sh
ctbz notas                                   # mês atual
ctbz notas --de 2026-01 --ate 2026-09        # um período (até 24 meses)
ctbz notas --tomador 11.222.333/0001-81      # por CPF/CNPJ do tomador
ctbz notas --tomador "ACME" --de 2026-07     # por nome do tomador
ctbz notas --numero 101
```

```text
Número  Emissão     Tomador    Documento                Valor  Status      Situação
101     04/09/2026  ACME LTDA  11.222.333/0001-81  R$ 1.000,00  AUTORIZADA  EMITIDA
Total: R$ 1.000,00 em 1 nota(s)
```

- A listagem é a do emissor de notas do painel (`novo-emissor/v2/listagem/notas/filtro`), uma
  consulta por mês e por página de 10 notas.
- `--tomador` decide sozinho: um CPF ou CNPJ (com ou sem pontuação) filtra pelo documento;
  qualquer outro texto filtra pelo nome. A API aceita um filtro por vez, por isso `--tomador`
  e `--numero` não podem ser usados juntos.
- Na tabela, o total do período vai para o stderr.
- Os campos das notas vêm do código do painel: a conta usada no desenvolvimento não tinha
  notas emitidas. Se algo vier diferente, abra uma issue com a saída de
  `ctbz api "novo-emissor/v2/listagem/notas/filtro?pagina=1&limite=10&ano=AAAA&mes=M"`.
