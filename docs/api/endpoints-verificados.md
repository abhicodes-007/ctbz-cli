# Endpoints verificados

Chamados com uma sessão real em 03/10/2026, empresa do Simples Nacional (serviços),
plano Básico. Todos são `GET` em `/api/plataforma/` e foram chamados com
`ctbz api <caminho>`. As colunas listam só os **campos de primeiro nível** da resposta.

Legenda: ✅ 200 · ⚠️ responde, mas precisa de parâmetro ou não está liberado · ❌ erro

## Empresa, conta e navegação

| Caminho | | Resposta |
|---|---|---|
| `dadosempresa/get` | ✅ | `empresaAtual{cnpj, razaoSocial, inscricaoMunicipal, regimeTributario, statusEmpresa, ramosAtividade, plano, certificado{status, dataValidade, valido}}`, `empresas[]`, `mgm`, `chamados`, `pendencias` |
| `socio/dados-logado` | ✅ | `cidade, uf, idade, email` |
| `conta-usuario/init` | ✅ | `email, telefone, metodo` |
| `menu/get` | ✅ | lista de itens `{id, type, icon, label, route, application, children[]}`: o menu lateral inteiro, com rotas de cada app |
| `appbar/get` | ✅ | `zendeskUrl, contaDigital{ctbzBank{banco, statusAbertura, conta, agencia…}, bs2}, uploadDocumentoAtivo` |
| `appshell/get` | ✅ | `alerta` |
| `contrato/buscarContratoServico` | ✅ | `html, versao, idHistoricoContrato`: texto do contrato de serviço |
| `empresa/dadosacesso/init` | ⚠️ 403 | `{"message": …}` |
| `/api/legado/socio/list` | ✅ | lista de sócios: `id, nome, cpf, administrador, responsavelReceita, possuiProLabore, salarioBase, dataAdmissao, categoria{descricao}, situacaoColaborador{descricao}` e ~90 outros campos pessoais — usado por `ctbz empresa socios` |
| `/api/legado/notafiscal/cnaeanexosmultiplos/list` | ✅ | CNAEs: `cnae{codigo, descricao, tipoRamoAtividade}, anexos[{codTabelaSimples, ativo, principal}]` — usado por `ctbz empresa atividades` |
| `/api/public/requestselecaoempresa?redirect=/&cnpj=` | ⚠️ 302 | redireciona para o SSO (`sso.contabilizei.com/selecionarempresa`), que exige sessão própria (ver [seleção de empresa](../autenticacao/03-selecao-de-empresa.md#trocar-de-empresa)) |

## Dashboard (tela inicial)

| Caminho | | Resposta |
|---|---|---|
| `dashboard/v2/init` | ✅ | `contaDigital, impostos, notasFiscais, mensalidade, prolabore, prolaboreV2, modalInitDashboard, checklistPrimeirosPassos…` (resumo de tudo) |
| `dashboard/init` | ✅ | `versao` |
| `dashboard/situacao-app` | ✅ | `elegivel, cadastroFinalizado, companyInfo{cnpj, cpf, cadastroStatus}, tipoPerfil, variacaoBannerCbank` |
| `dashboard/fatura` | ✅ | `total, vencimento, competencia, status, label, possuiCartaoPrincipal, botaoAcao…`: mensalidade da Contabilizei |
| `dashboard/v1/mensalidade` | ✅ | `status, plano, valor, dataVencimento, competenciaAnteriorAtrasada…` |
| `dashboard/rotinas-mensais` | ✅ | `cardInfo` |
| `dashboard/v2/central-rotinas` | ✅ | `tipoExibicao, tipoPendencia, pendencias, rotinas, rotinasContabilizei` |
| `dashboard/prolabore` | ✅ | `valorProlabore, competenciaAtual, competenciaAnterior, calculando, sociedade…` |
| `dashboard/card-certificado` | ✅ | `tipoCard, prazo, certificadoVencido, prazoFinalizado, temEmissor` |
| `dashboard/informerendimento/debitos-federais` | ✅ | `possuiDebitosFederais` |
| `home/pendencia/pendenciasEmpresa` | ✅ | lista `{id, dataCriacao, tipoPendencia, detalhe, dataLimite, situacaoPendencia…}` |

## Impostos

| Caminho | | Resposta |
|---|---|---|
| `impostos/v5/impostos-a-pagar/guias` | ✅ | `emAtraso[], esteMes[], proximoMes[]`; cada guia: `id, nome{label}, origem, tipo, identificadorImposto, vencimento, vencimentoOriginal, valor{label}, competencia, status{label}, acaoBotao, pendencias[]` |
| `impostos/v5/impostos-a-pagar/guia/{id}` | ✅ | `id, origem, tipo, identificadorImposto, vencimento, valorTotal, valorEstimado, valorOriginal, valorJurosEMulta, valorEmAtraso, oraculo, banner…` |
| `impostos/v3/impostos-a-pagar/guia/{id}/baixar-guia` | ✅ | `url`: link para o PDF da guia (DAS etc.) |
| `impostos/v5/impostos-a-pagar/init` | ✅ | `podeGerenciarDebitoAutomatico, exibirComoMeuImpostoFoiCalculado, exibirMemoriaDeCalculo, dadosMemoriaDeCalculo` |
| `impostos/v5/impostos-a-pagar/banners` | ✅ | `banners` |
| `impostos/v5/historico-impostos/guias` | ✅ | `impostos` |
| `impostos/v5/historico-impostos/anos-vigentes` | ✅ | lista de anos |
| `impostos/v5/historico-impostos/faturamento-mensal` | ✅ | `faturamentoMensal` |
| `impostos/v2/historico-impostos/init` | ✅ | `emDia, quantidadeGuiasVencidas` |
| `impostos/v2/historico-impostos/guias?pagina=1` | ✅ | `paginaAtual, totalPaginas, competencias`; aceita também `status`, `mes`, `ano` |
| `impostos/v2/historico-impostos/guias` (sem `pagina`) | ❌ 560 | `{identificador, detalhe, dataHora}` |
| `impostos/rollout` | ✅ | `versao` |
| `impostos/` · `impostos/impostos-a-pagar/` · `impostos/parcelamentos` | ❌ 404 | não são rotas de API |

Fluxo para baixar as guias do mês:

```sh
ctbz api impostos/v5/impostos-a-pagar/guias | jq '.esteMes[] | {id, nome: .nome.label, vencimento, valor: .valor.label}'
ctbz api impostos/v3/impostos-a-pagar/guia/<id>/baixar-guia | jq -r .url
```

## Contabilidade e finanças

| Caminho | | Resposta |
|---|---|---|
| `relatorios-ms/gerar-balancete/{ano}/{mes}` | ✅ | lista de contas `{id, idContaPai, descricao, natureza, tipo, classificacaoConta, totalCredito, totalDebito, saldoExercicioAnterior, saldoAnterior, saldoExercicio…}` |
| `caixa/listpaginada/{ano}/{mes}/{porPagina}/{pagina}` | ✅ | `cursor, total, list, serializedList, responseCode`: lançamentos do caixa |
| `caixa/listpaginada/` (sem parâmetros) | ❌ 404 | — |
| `movimentacao-financeira/v2/extratos` | ✅ | lista `{ano, mes, idContaBancaria, banco, numeroConta, agencia, situacao, statusIntegracao, exibirBotaoImportar…}` |
| `movimentacao-financeira/contasUsuario` | ✅ | lista (≈236) de contas de classificação `{id, descricao, situacao, descricaoContaContabil, classificacao…}` |
| `contabancaria/list` | ✅ | `bancos, contasBancarias, processandoCtbzBank` |
| `informerendimento/recuperardadosdistribuicaocliente` | ✅ | `ano, saldo, totalDistribuido, totalAdiantamentos, exercicioFechado, lucrosSocios, dataLimite…`: distribuição de lucros |
| `conciliacao-fiscal/v2/init` | ✅ | `qtdNotasFiscaisPendentes, qtdRecebimentosPendentes, qtdConciliacoesAutomaticasMesAnterior…` |
| `simulador-impostos-avancado/init` | ✅ | `disponibilidade, atividades, primeiroCiclo, motorFatorR…` |

## Pró-labore

| Caminho | | Resposta |
|---|---|---|
| `prolabore/init` | ✅ | `salarioMinimo, valorMaximoContribuicaoInss, valorMaximoProlabore, porcentagemInss, isMotorFatorR…` |
| `prolabore/central/init` | ✅ | `tipoGerenciamento, socios, totalProLabore, valorMaximoInss, baseCalculoIrrf, zerarProlabore…` |

## Notas fiscais

| Caminho | | Resposta |
|---|---|---|
| `novo-emissor/listagem/init` | ✅ | `certificadoDigital, hasInstability, user, endereco, permiteEmissaoExterior, emissorEnabled…` |
| `novo-emissor/listagem/notas` | ✅ | lista de notas emitidas |
| `novo-emissor/tomadores/init` | ✅ | `emissaoSemTomador, tomadores, permiteEmissaoExterior` |
| `novo-emissor/v2/versao-emissor` | ✅ | `versaoNovoEmissor` |
| `notafiscal/listaliquotaatividade` | ✅ | `regimeTributario, temCodigoServicoItemServico, interno, externo` |

## Certificado digital, documentos, pagamentos e atendimento

| Caminho | | Resposta |
|---|---|---|
| `certificado/status` | ✅ | `situacao, dataVencimento` (epoch ms), `mensagemErro, valido, aptoRenovacao` — usado por `ctbz empresa certificado` |
| `documentos/envio-documento/init` | ✅ | `tiposPermitidos, documentos` |
| `documentos/listar-enviados?…` | ⚠️ 400 | exige `tipoDocumento=` (um ou mais) + `limit` + `offset` |
| `payments/recorrencia/init` | ✅ | `status, competencia, dataProximoPagamento, habilitado, ativado…` |
| `payments/recorrencia/historico` | ✅ | `pagamentos` |
| `inadimplencia/consultasituacaomensalidadeempresa` | ✅ | texto `OK` |
| `atendimento/chamados?em-andamento=true` | ✅ | lista de chamados abertos |
| `atendimento/chamados?finalizados=true` | ⚠️ 400 | texto `timeout` (lento no servidor) |
| `escritorio-virtual/recuperar-mensagens` | ❌ 560 | serviço não contratado |

## Não chamados de propósito

Alteram estado, fazem pagamentos ou expõem segredos, e por isso ficaram de fora dos testes:

- tudo que é `POST`/`PUT`/`PATCH`/`DELETE` (emitir/cancelar nota, aceitar termos,
  confirmar pagamento de guia, alterar pró-labore, remover certificado…);
- `GET` com efeito colateral pelo nome: `certificado/fluxo-one-click/alterarParametroEmpresa`,
  `notafiscal/cancelarAPI`;
- `GET` que expõem segredos: `certificado/senha`, `certificado/download`,
  `billing/gestao-pagamentos/get-client-key`, `conta-usuario/qr-code-aplicativo-autenticacao`.
