# Contabilidade

Relatórios contábeis, caixa, extratos e plano de contas. Os relatórios usam o mês como
`AAAA-MM`; o painel mostra por padrão dezembro do ano anterior.

## Balancete

```sh
ctbz balancete 2026-08
ctbz balancete 2026-08 -o csv > balancete-2026-08.csv
```

```text
Conta          Descrição                          Nível  Saldo anterior      Débitos     Créditos       Saldo
1              ATIVO                                  1         R$ 0,00  R$ 1.815,00  R$ 1.419,83    R$ 395,17
1.01           CIRCULANTE                             2         R$ 0,00  R$ 1.815,00  R$ 1.419,83    R$ 395,17
1.01.01.01.00          Caixa Geral                    5         R$ 0,00  R$ 1.139,00  R$ 1.280,83  -R$ 141,83
```

- Cada linha é uma conta da árvore contábil; na tabela, a descrição é recuada pelo nível.
- Em CSV e JSON a descrição vem sem recuo e o nível fica na coluna `nivel`, o que facilita
  filtrar (ex.: só o nível 1 para ATIVO, PASSIVO e PATRIMÔNIO LÍQUIDO).
- O PDF do painel é gerado no navegador a partir dos mesmos dados; a API não oferece PDF.

## Balanço patrimonial

```sh
ctbz balanco 2025          # fechamento do exercício (dezembro)
ctbz balanco 2026-09       # posição em um mês
```

```text
Conta  Descrição            Nível  Grupo            Saldo  Exercício anterior
1      ATIVO                    1  ATIVO        R$ 395,17             R$ 0,00
1.01     CIRCULANTE             2  ATIVO        R$ 395,17             R$ 0,00
2      PASSIVO                  1  PASSIVO     -R$ 395,17             R$ 0,00
```

- Ativo, passivo e patrimônio líquido, com o saldo do exercício e o do exercício anterior.
- As contas de resultado (receitas e despesas) ficam de fora, como no painel.

## Razão

```sh
ctbz razao --de 2026-01 --ate 2026-09               # todas as contas no período
ctbz razao --conta 1.01.01.01.00 --de 2026-07       # uma conta
ctbz razao --conta 1.01 -o csv > circulante.csv     # prefixo: todas as contas de 1.01
```

```text
Data        Conta          Descrição da conta  Histórico       Contrapartida                                    Débito  Crédito        Saldo
21/07/2026  1.01.01.01.00  Caixa Geral         Capital social  2.07.01.01.00 Capital Social Realizado no País  R$ 1.000,00         R$ 1.000,00
```

- Um lançamento por linha, com a conta de contrapartida e o saldo acumulado da conta no
  exercício depois do lançamento.
- A API devolve um mês por vez; a CLI consulta cada mês do período (até 24).

## Caixa do mês

```sh
ctbz caixa 2026-09
ctbz caixa 2026-09 -o csv > caixa-2026-09.csv
```

```text
Data        Descrição               Conta                         Classificação        Valor  Situação    Automático
18/09/2026  Pagamento Contabilizei  Mensalidade de contabilidade  DESPESA         -R$ 141,83  CONFIRMADO  sim
Total do mês: -R$ 141,83 em 1 lançamento(s)
```

- Entradas e saídas classificadas do mês, com a conta de classificação (de
  `ctbz contas`) e o tipo (receita, despesa…). Saídas têm valor negativo.
- `automatico` indica lançamentos feitos pelo sistema (ex.: a mensalidade da Contabilizei).
- Como no painel, a CLI pede até 1000 lançamentos por mês e avisa se houver mais.
