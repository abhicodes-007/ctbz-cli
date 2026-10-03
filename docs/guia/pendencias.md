# Pendências, rotinas e atendimento

## Pendências da empresa

```sh
ctbz pendencias
```

```text
ID                Tipo                                                Detalhe                                                       Criada em   Prazo       Situação  Alerta
1000000000000001  Cadastre o PIS para ativar o pró-labore automático  O pró-labore influencia no valor do imposto, por isso assi…  21/07/2026  21/07/2026  Pendente  vencida
```

- São as pendências do card da tela inicial do painel: tipo, detalhe, criação, prazo e situação.
- Por padrão aparecem só as **abertas**; `--todas` inclui as finalizadas.
- O detalhe vem da API em HTML; a CLI tira as tags. Na tabela ele é cortado em 60
  caracteres; use `-o json` para ler o texto completo.

### Alertas de prazo

A coluna `alerta` só é preenchida para pendências abertas
([ADR-0012](../adr/0012-prazos-e-alertas.md)):

| `alerta` | Quando |
|---|---|
| `vencida` | o prazo já passou |
| `próxima` | o prazo é hoje ou nos próximos 7 dias |
| vazio (`null` no JSON) | prazo distante, sem prazo ou pendência finalizada |

`--fail-on-vencidas` faz o comando terminar com **código 4** quando há pendência vencida:

```sh
ctbz pendencias --fail-on-vencidas -o csv > /dev/null || notify-send "Há pendências vencidas"
ctbz pendencias -o json | jq -r '.[] | select(.alerta != null) | "\(.prazo) \(.tipo)"'
```

## Conciliação fiscal

```sh
ctbz pendencias conciliacao
```

```text
Competência de referência:        09/2026
Notas fiscais sem recebimento:    0
Recebimentos sem nota fiscal:     0
Conciliações automáticas no mês:  0
```

- **Notas sem recebimento**: notas fiscais emitidas que não foram ligadas a um crédito no extrato.
- **Recebimentos sem nota**: créditos no extrato sem nota fiscal correspondente.
- A competência de referência é o mês anterior, o mesmo que o painel mostra.
- `--listar notas` ou `--listar recebimentos` lista os itens pendentes dos últimos 12 meses.
  O formato desses itens ainda não foi verificado (a conta de desenvolvimento não tinha
  pendências), por isso os campos saem como a API os devolve.
- `--fail-on-pendencias` termina com código 4 quando há algo a conciliar.

## Rotinas e obrigações do mês

```sh
ctbz rotinas
ctbz rotinas --mes 2026-11
```

```text
Responsável   Rotina                                 Prazo       Status     Valor     Alerta
empresa       Importar extrato bancário de setembro  05/10/2026  EM_ABERTO            próxima
empresa       DARF Unificada (Ativação do fator R)   06/10/2026  EM_ABERTO            próxima
empresa       Mensalidade da Contabilizei            15/10/2026  EM_ABERTO  R$ 15,90
contabilizei  eSocial                                15/10/2026  EM_ABERTO
contabilizei  DCTFWeb                                15/10/2026  EM_ABERTO
```

- `responsavel` separa o que a **empresa** precisa fazer do que a **Contabilizei** entrega
  (eSocial, DCTFWeb, EFD-Reinf…).
- O filtro é pelo mês do prazo. A API não aceita competência nem mês: devolve sempre o mês
  anterior, o atual e o próximo, e a CLI filtra. Fora dessa janela, a lista sai vazia.
- `alerta` segue a mesma regra das pendências: `vencida` ou `próxima` para rotinas não
  realizadas ([ADR-0012](../adr/0012-prazos-e-alertas.md)).
- `--fail-on-vencidas` (código 4) considera só as rotinas da **empresa**: atraso numa
  obrigação da Contabilizei não é algo que você resolve sozinho.

## Chamados de atendimento

```sh
ctbz chamados                  # em andamento
ctbz chamados --finalizados    # já resolvidos
```

- Colunas: `id`, `assunto`, `status`, `canal`, `criado`, `atualizado`, `previsao_retorno` e
  `link` (a página do chamado na central de ajuda).
- Em contas com muitos chamados, a Contabilizei não consegue listar os finalizados (o
  servidor responde com erro sempre, não adianta repetir). A CLI então mostra os finalizados
  entre os **100 chamados mais recentes** da empresa e avisa no stderr
  ([ADR-0013](../adr/0013-sem-retentativa-com-fonte-alternativa.md)). Nessa fonte não há
  canal, data de atualização nem previsão de retorno.
