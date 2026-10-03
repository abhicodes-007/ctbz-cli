# Catálogo de endpoints (gerado)

> Gerado por `scripts/extrair-endpoints.py` em 2026-10-03 a partir de 165 arquivos JavaScript do painel. 335 chamadas encontradas.
>
> A base de cada chamada é inferida pela instância axios usada no código; caminhos terminados em `/` ou `=` recebem parâmetros concatenados pelo front.
> Endpoints já testados estão em [endpoints-verificados.md](endpoints-verificados.md).

Bases declaradas no front: `/api/fintech/`, `/api/leads/hubspot/`, `/api/legado/`, `/api/multiusuario/`, `/api/pagamentos/`, `/api/plataforma/`, `/api/public/`

## `/api/plataforma/` (318)

### api

| Método | Caminho |
|---|---|
| GET | `/api/pagamentos-fintech/card-cobranca` |
| GET | `api/pagamentos-fintech/card-impostos` |

### appbar

| Método | Caminho |
|---|---|
| GET | `appbar/get` |

### atendimento

| Método | Caminho |
|---|---|
| GET | `atendimento/chamados?em-andamento=true` |
| GET | `atendimento/chamados?finalizados=true` |

### ativacao

| Método | Caminho |
|---|---|
| POST | `/ativacao/timeline/aceite/mfa-gov` |
| GET | `ativacao/timeline` |
| POST | `ativacao/timeline/finalizar` |
| POST | `ativacao/timeline/finish-step-plano-saude` |
| POST | `ativacao/timeline/salvar-interacao` |

### autopilot

| Método | Caminho |
|---|---|
| DELETE | `/autopilot/clientes/` |
| GET | `/autopilot/clientes/` |
| GET | `/autopilot/clientes/consulta/` |

### billing

| Método | Caminho |
|---|---|
| GET | `billing/gestao-pagamentos/get-client-key` |

### caixa

| Método | Caminho |
|---|---|
| POST | `caixa/lancamentousuario/novo/` |
| DELETE | `caixa/lancamentousuario/remover/` |
| GET | `caixa/listpaginada/` |

### cartoes-de-credito

| Método | Caminho |
|---|---|
| GET | `/cartoes-de-credito/init` |
| GET | `cartoes-de-credito/gestao-pagamentos` |
| POST | `cartoes-de-credito/gestao-pagamentos` |
| DELETE | `cartoes-de-credito/gestao-pagamentos/` |
| PUT | `cartoes-de-credito/gestao-pagamentos/` |

### central-rotinas

| Método | Caminho |
|---|---|
| POST | `central-rotinas/aceitar-carta-responsabilidade` |
| POST | `central-rotinas/aceitar-termo-aceite-tp` |
| POST | `central-rotinas/aceitar-termo-debitos` |
| GET | `central-rotinas/init` |
| POST | `central-rotinas/resolver-pendencia-procuracao-ecac` |

### certificado

| Método | Caminho |
|---|---|
| GET | `/certificado/agendamento-contabilizei/horario` |
| POST | `certificado/agendamento-contabilizei` |
| GET | `certificado/agendamento-contabilizei/init` |
| POST | `certificado/agendamento-contabilizei/senha` |
| GET | `certificado/agendamento/dados-cliente` |
| POST | `certificado/agendamento/dados-cliente` |
| GET | `certificado/bird-id/init` |
| GET | `certificado/checkout` |
| GET | `certificado/download` |
| POST | `certificado/emitir-certificado` |
| GET | `certificado/emitir-certificado/aguarde/init` |
| GET | `certificado/emitir-certificado/init` |
| GET | `certificado/emitir-certificado/status` |
| GET | `certificado/fluxo-configurar-certificado/init` |
| GET | `certificado/fluxo-one-click/alterarParametroEmpresa` |
| POST | `certificado/fluxo-one-click/bird-id/emitir-certificado` |
| GET | `certificado/fluxo-one-click/bird-id/init` |
| POST | `certificado/fluxo-one-click/bird-id/valida-senha` |
| GET | `certificado/fluxo-one-click/sucesso/init` |
| GET | `certificado/pre-checkout` |
| GET | `certificado/processo-aquisicao/agendamento/init` |
| PUT | `certificado/processo-aquisicao/etapa` |
| GET | `certificado/processo-aquisicao/status` |
| GET | `certificado/senha` |
| GET | `certificado/solicitacao-certificado/init` |
| GET | `certificado/status` |

### certificado-digital

| Método | Caminho |
|---|---|
| POST | `certificado-digital/v2/agendar-atendimento` |
| GET | `certificado-digital/v2/buscar-horarios-agendamento` |
| PUT | `certificado-digital/v2/cancelar-solicitacao` |
| GET | `certificado-digital/v2/consultar-status-solicitacao-criada` |
| POST | `certificado-digital/v2/emitir-certificado` |
| PUT | `certificado-digital/v2/emitir-certificado-com-senha` |
| POST | `certificado-digital/v2/informar-cnh` |
| GET | `certificado-digital/v2/iniciar-atendimento-expresso` |
| GET | `certificado-digital/v2/init` |
| GET | `certificado-digital/v2/pontos-atendimento/buscar` |
| GET | `certificado-digital/v2/pontos-atendimento/buscar-bairros` |
| GET | `certificado-digital/v2/pontos-atendimento/buscar-cidades` |
| GET | `certificado-digital/v2/reiniciar-atendimento` |
| POST | `certificado-digital/v2/salvar-aceite` |
| POST | `certificado-digital/v2/solicitar-atendimento-pendencia` |
| PUT | `certificado-digital/v2/tipo-atendimento` |
| GET | `certificado-digital/v2/verificar-rollout` |

### checklist-onboarding

| Método | Caminho |
|---|---|
| POST | `/checklist-onboarding/aside/concluir-etapa` |
| GET | `/checklist-onboarding/aside/init` |
| GET | `/checklist-onboarding/aside/trilha/` |
| PATCH | `/checklist-onboarding/exibir-primeiros-passos/dispensar` |
| POST | `/checklist-onboarding/procuracoes/compartilhamentos` |
| POST | `/checklist-onboarding/pular-tarefas` |
| PATCH | `/checklist-onboarding/reativar-tarefas` |
| POST | `/checklist-onboarding/tooltip/guia-inicio/concluir` |

### conciliacao-fiscal

| Método | Caminho |
|---|---|
| PATCH | `conciliacao-fiscal/conciliar/alerta-receitas/fechar` |
| GET | `conciliacao-fiscal/conciliar/recebimentos` |
| GET | `conciliacao-fiscal/conciliar/receitas` |
| POST | `conciliacao-fiscal/conciliar/resolver-pendencia` |
| GET | `conciliacao-fiscal/init` |
| POST | `conciliacao-fiscal/pendencia/detalhes` |
| GET | `conciliacao-fiscal/pendencias` |
| GET | `conciliacao-fiscal/reclassificacao/consulta-periodo` |
| GET | `conciliacao-fiscal/v2/init` |
| GET | `conciliacao-fiscal/v2/pendencias` |
| GET | `conciliacao-fiscal/versao-experiencia` |

### conta-bancaria

| Método | Caminho |
|---|---|
| GET | `conta-bancaria/bs2` |

### conta-usuario

| Método | Caminho |
|---|---|
| POST | `conta-usuario/alterar-dados` |
| POST | `conta-usuario/alterar-senha` |
| POST | `conta-usuario/enviar-token-otp` |
| GET | `conta-usuario/init` |
| GET | `conta-usuario/qr-code-aplicativo-autenticacao` |
| POST | `conta-usuario/validar-codigo-app` |

### contabancaria

| Método | Caminho |
|---|---|
| GET | `/contabancaria/detalhes-da-conta/saldo/` |
| POST | `/contabancaria/divergencia/atualizar-conta` |
| POST | `/contabancaria/divergencia/cadastrar-conta` |
| GET | `contabancaria/detalhes-da-conta/init/` |
| GET | `contabancaria/divergencia/init/` |
| DELETE | `contabancaria/excluir/` |
| GET | `contabancaria/list` |
| POST | `contabancaria/salvar` |

### contrato

| Método | Caminho |
|---|---|
| GET | `/contrato/buscarContratoServico` |
| GET | `/contrato/empresaSemContratosComFaturaGeradaPeloAdmin` |
| GET | `contrato/buscarContratoServico/` |
| GET | `contrato/empresaSemContratosComFaturaGeradaPeloAdmin/` |

### cross-sell

| Método | Caminho |
|---|---|
| GET | `cross-sell/multibeneficios/cms/v1/landing-upsell` |
| GET | `cross-sell/multibeneficios/coachmark/v1/beneficios` |
| GET | `cross-sell/multibeneficios/modal/v2/beneficios` |
| POST | `cross-sell/multibeneficios/v1/beneficio/jornada/termo-adesao/assinar` |
| POST | `cross-sell/multibeneficios/v1/beneficios/escolher` |
| GET | `cross-sell/multibeneficios/v1/beneficios/guard?urlDestino=` |
| GET | `cross-sell/multibeneficios/v2/beneficio/detalhes?idBeneficio=` |
| GET | `cross-sell/multibeneficios/v2/beneficio/jornada/termo-adesao?idBeneficio=` |
| GET | `cross-sell/multibeneficios/v2/beneficio/jornada?idBeneficio=` |
| GET | `cross-sell/multibeneficios/v2/beneficios` |
| GET | `cross-sell/multibeneficios/v2/beneficios/escolhidos` |
| POST | `cross-sell/simulador/enviar-contato` |
| POST | `cross-sell/simulador/plano-saude/hubspot/enviar-evento` |
| GET | `cross-sell/simulador/planos-saude/init` |
| POST | `cross-sell/simulador/planos-saude/prestadores` |
| GET | `cross-sell/simulador/planos-saude/resultados/init` |
| POST | `cross-sell/simular` |
| POST | `cross-sell/simular/busca-cidade` |
| POST | `cross-sell/v2/simulador/planos-saude/prestadores` |
| POST | `cross-sell/v2/simular` |

### dadosempresa

| Método | Caminho |
|---|---|
| GET | `/dadosempresa/existeContaCadastrada` |
| GET | `dadosempresa/get` |

### dashboard

| Método | Caminho |
|---|---|
| GET | `/dashboard/cross-upsell/vitrine` |
| GET | `/dashboard/fatura` |
| GET | `/dashboard/rotinas-mensais` |
| GET | `/dashboard/situacao-app` |
| GET | `dashboard/card-certificado` |
| GET | `dashboard/informerendimento/debitos-federais` |
| GET | `dashboard/init` |
| POST | `dashboard/interacao-notificacao` |
| GET | `dashboard/primeiros-passos` |
| GET | `dashboard/produto/plano-de-saude/init` |
| GET | `dashboard/prolabore` |
| GET | `dashboard/v1/central-rotinas` |
| GET | `dashboard/v1/mensalidade` |
| GET | `dashboard/v2/central-rotinas` |
| GET | `dashboard/v2/init` |
| POST | `dashboard/v2/notas-fiscais/salvar-clique-emitir` |
| POST | `dashboard/v2/prolabore/salvar-decisao` |
| POST | `dashboard/v2/termo-exclusao` |

### documentos

| Método | Caminho |
|---|---|
| GET | `/documentos/listar-enviados?` |
| GET | `/documentos/tipos/init?area=` |
| POST | `documentos/envio-documento/enviar/sem-arquivo` |
| GET | `documentos/envio-documento/init` |
| POST | `documentos/reclassificar` |
| GET | `documentos/reclassificar/init` |

### empresa

| Método | Caminho |
|---|---|
| DELETE | `empresa/certificado/removercert` |
| POST | `empresa/dadosacesso/atualizarDadosAcesso` |
| PUT | `empresa/dadosacesso/confirmar-credencial-prefeitura-valida` |
| GET | `empresa/dadosacesso/init` |

### empresa-inativa

| Método | Caminho |
|---|---|
| GET | `/empresa-inativa/permite-registrar-solicitacao` |
| POST | `/empresa-inativa/registrar-solicitacao` |

### escritorio-virtual

| Método | Caminho |
|---|---|
| PUT | `escritorio-virtual/autorizar-recebimento/` |
| GET | `escritorio-virtual/endereco-entrega` |
| GET | `escritorio-virtual/gerar-link-download-comprovante/` |
| GET | `escritorio-virtual/gerar-link-download-correspondencia/` |
| GET | `escritorio-virtual/pesquisarEndereco/` |
| GET | `escritorio-virtual/recuperar-mensagens` |
| POST | `escritorio-virtual/salvar-endereco-correspondencia` |
| GET | `escritorio-virtual/verifica-autorizacao-recebimento` |

### experience-enhancement

| Método | Caminho |
|---|---|
| POST | `experience-enhancement/pode-exibir-flow-cnpj` |
| PATCH | `experience-enhancement/salvar-exibicao-cnpj` |

### experts

| Método | Caminho |
|---|---|
| GET | `experts/cashflow/dashboard` |
| PUT | `experts/cashflow/notas-fiscais` |
| GET | `experts/cashflow/notas-fiscais?mes=` |
| GET | `experts/cashflow/pagamentos?mes=` |
| GET | `experts/cashflow/status` |

### fintech

| Método | Caminho |
|---|---|
| POST | `fintech/conta-digital/aceite?origem=` |
| POST | `fintech/conta-digital/aceite?origem=pagamento-recorrente` |
| POST | `fintech/conta-digital/aceite?origem=timeline` |
| GET | `fintech/conta-digital/dashboard` |
| PUT | `fintech/conta-digital/desistir` |
| POST | `fintech/conta-digital/lead-atendimento` |
| POST | `fintech/conta-digital/notificacao-conclusao` |
| GET | `fintech/metodos-pagamento-impostos/init` |

### home

| Método | Caminho |
|---|---|
| PUT | `/home/pagtoplano/atualizar-plano` |
| POST | `home/nps` |
| GET | `home/pendencia/pendenciasEmpresa` |

### hubspot

| Método | Caminho |
|---|---|
| POST | `hubspot/enviar` |

### impostos

| Método | Caminho |
|---|---|
| GET | `impostos/` |
| GET | `impostos/como-imposto-foi-calculado/init` |
| GET | `impostos/como-imposto-foi-calculado/tabela-irrf` |
| GET | `impostos/impostos-a-pagar/` |
| POST | `impostos/impostos-a-pagar/memoria-de-calculo-visualizada` |
| GET | `impostos/parcelamento/banner-oferta-parcelamento` |
| GET | `impostos/parcelamento/detalhes/` |
| POST | `impostos/parcelamento/negociacao-automatica/pgfn-nao-previdenciario/contratar` |
| GET | `impostos/parcelamento/negociacao-automatica/pgfn-nao-previdenciario/detalhes/init/` |
| GET | `impostos/parcelamento/negociacao-automatica/pgfn-nao-previdenciario/init` |
| POST | `impostos/parcelamento/negociacao-automatica/pgfn-previdenciario/contratar` |
| GET | `impostos/parcelamento/negociacao-automatica/pgfn-previdenciario/detalhes/init/` |
| GET | `impostos/parcelamento/negociacao-automatica/pgfn-previdenciario/init` |
| POST | `impostos/parcelamento/negociacao-automatica/pgfn-simples-nacional/contratar` |
| GET | `impostos/parcelamento/negociacao-automatica/pgfn-simples-nacional/detalhes/init/` |
| GET | `impostos/parcelamento/negociacao-automatica/pgfn-simples-nacional/init` |
| POST | `impostos/parcelamento/negociacao-automatica/simples-nacional/contratar` |
| GET | `impostos/parcelamento/negociacao-automatica/simples-nacional/detalhes/init/` |
| GET | `impostos/parcelamento/negociacao-automatica/simples-nacional/init` |
| POST | `impostos/parcelamento/negociacao-especializada/contratar` |
| GET | `impostos/parcelamento/negociacao-especializada/detalhes/init/` |
| GET | `impostos/parcelamento/negociacao-especializada/init?tipoNegociacao=` |
| GET | `impostos/parcelamento/regularize-seus-impostos/init?tipoNegociacao=` |
| GET | `impostos/rollout` |
| GET | `impostos/v2/historico-impostos/guia/` |
| PUT | `impostos/v2/historico-impostos/guia/` |
| GET | `impostos/v2/historico-impostos/guias` |
| GET | `impostos/v2/historico-impostos/init` |
| GET | `impostos/v2/historico-impostos/parcela/` |
| PUT | `impostos/v2/impostos-a-pagar/guia/` |
| GET | `impostos/v2/impostos-a-pagar/recalculo/init?idGuia=` |
| GET | `impostos/v3/impostos-a-pagar/` |
| GET | `impostos/v3/impostos-a-pagar/guia/` |
| PUT | `impostos/v3/impostos-a-pagar/guia/` |
| GET | `impostos/v3/impostos-a-pagar/parcela/` |
| GET | `impostos/v5/impostos-a-pagar/banners` |
| GET | `impostos/v5/impostos-a-pagar/guia/` |
| PUT | `impostos/v5/impostos-a-pagar/guia/` |
| GET | `impostos/v5/impostos-a-pagar/guias` |
| GET | `impostos/v5/impostos-a-pagar/init` |

### inadimplencia

| Método | Caminho |
|---|---|
| GET | `inadimplencia/consultasituacaomensalidadeempresa` |

### informerendimento

| Método | Caminho |
|---|---|
| POST | `informerendimento/aceitar-carta-responsabilidade/` |
| POST | `informerendimento/aceitar-termo-debitos/` |
| POST | `informerendimento/aceite/` |
| POST | `informerendimento/historico` |
| GET | `informerendimento/listSocioInformeRendimentos/` |
| POST | `informerendimento/reabrir-balanco/` |
| GET | `informerendimento/recuperardadosdistribuicaocliente` |
| POST | `informerendimento/salvarconfiguracaocliente` |
| GET | `informerendimento/v2/` |
| GET | `informerendimento/visualizar-carta-responsabilidade/` |

### loans

| Método | Caminho |
|---|---|
| POST | `loans/aceite-oferta-dash/OFERTA_LOANS_FINTECH` |

### menu

| Método | Caminho |
|---|---|
| GET | `menu/get` |

### movimentacao-financeira

| Método | Caminho |
|---|---|
| PUT | `/movimentacao-financeira/classificar` |
| PUT | `/movimentacao-financeira/desmembrar` |
| GET | `/movimentacao-financeira/redirecionarNovoCaixaExtrato` |
| GET | `/movimentacao-financeira/v2/extratos` |
| GET | `movimentacao-financeira/contas-usuario/` |
| GET | `movimentacao-financeira/contasUsuario` |
| POST | `movimentacao-financeira/eventoUploadExtrato` |
| GET | `movimentacao-financeira/getContaBancaria` |
| GET | `movimentacao-financeira/info-extrato/` |
| POST | `movimentacao-financeira/interagiu-modal-integracao` |
| GET | `movimentacao-financeira/permiteImportacao/` |

### notafiscal

| Método | Caminho |
|---|---|
| GET | `/notafiscal/emitir/tomador/list/` |
| GET | `notafiscal/listaliquotaatividade` |

### novo-emissor

| Método | Caminho |
|---|---|
| GET | `/novo-emissor/assistente-nbs/session` |
| POST | `/novo-emissor/assistente-nbs/session` |
| GET | `/novo-emissor/emissao/atividades` |
| POST | `/novo-emissor/emissao/cancelar/` |
| GET | `/novo-emissor/listagem/init` |
| GET | `/novo-emissor/listagem/notas` |
| GET | `/novo-emissor/v2/emissao/atividades/cindop?` |
| GET | `/novo-emissor/v2/emissao/atividades/codigo-municipal?codigoIbge=` |
| GET | `/novo-emissor/v2/feature-flag` |
| GET | `/novo-emissor/v2/registro/atividades` |
| GET | `/novo-emissor/v2/tipo-emissor-disponivel` |
| GET | `/novo-emissor/v2/trilhas-empresa` |

### oferta-conta-integracao-bs2

| Método | Caminho |
|---|---|
| GET | `/oferta-conta-integracao-bs2` |
| POST | `/oferta-conta-integracao-bs2` |

### onboarding-plataforma

| Método | Caminho |
|---|---|
| GET | `onboarding-plataforma/rollout` |

### pagamento-cartao-credito

| Método | Caminho |
|---|---|
| GET | `pagamento-cartao-credito/faturas` |
| POST | `pagamento-cartao-credito/pagar` |

### parametrosEmpresa

| Método | Caminho |
|---|---|
| GET | `/parametrosEmpresa/list/` |

### payments

| Método | Caminho |
|---|---|
| POST | `payments/aceite-termo-recorrencia` |
| DELETE | `payments/cartao/` |
| POST | `payments/cartao/adicionar` |
| POST | `payments/recorrencia/ativar` |
| POST | `payments/recorrencia/ativar-conta-com-termo` |
| POST | `payments/recorrencia/desativar` |
| GET | `payments/recorrencia/historico` |
| GET | `payments/recorrencia/init` |

### planos

| Método | Caminho |
|---|---|
| POST | `planos/aceite-plano-manutencao` |

### prolabore

| Método | Caminho |
|---|---|
| POST | `prolabore/assessor/abandonar-fluxo` |
| PUT | `prolabore/assessor/escolher-socio/` |
| GET | `prolabore/assessor/init` |
| PUT | `prolabore/assessor/passo` |
| GET | `prolabore/assessor/passo/` |
| POST | `prolabore/assessor/rollout` |
| DELETE | `prolabore/assessor/socio/` |
| GET | `prolabore/assessor/socio/` |
| POST | `prolabore/assessor/socio/` |
| PUT | `prolabore/central/empresa/gestao-inteligente` |
| PATCH | `prolabore/central/empresa/zerar-prolabore` |
| PATCH | `prolabore/central/empresa/zerar-prolabore-gt` |
| GET | `prolabore/central/gestao/` |
| PUT | `prolabore/central/gestao/` |
| GET | `prolabore/central/historico/` |
| GET | `prolabore/central/init` |
| GET | `prolabore/central/rollout` |
| PUT | `prolabore/excluir-empresa-motor-fator-r` |
| GET | `prolabore/init` |
| GET | `prolabore/rollout` |
| POST | `prolabore/viu-dialog-informacao-prolabore` |

### questionario-preferencia-prolabore

| Método | Caminho |
|---|---|
| GET | `questionario-preferencia-prolabore/decisao-preferencia` |
| POST | `questionario-preferencia-prolabore/decisao-preferencia` |
| GET | `questionario-preferencia-prolabore/init` |
| POST | `questionario-preferencia-prolabore/salvar-resposta` |

### relatorios-ms

| Método | Caminho |
|---|---|
| GET | `relatorios-ms/gerar-balancete/` |
| GET | `relatorios-ms/gerarbalanco/` |
| GET | `relatorios-ms/gerarrazaotipoa/` |

### retencao

| Método | Caminho |
|---|---|
| POST | `/retencao/paywall/salvar-interacao` |

### simulador-impostos-avancado

| Método | Caminho |
|---|---|
| GET | `simulador-impostos-avancado/init` |
| POST | `simulador-impostos-avancado/simular` |

### socio

| Método | Caminho |
|---|---|
| GET | `socio/dados-logado` |
| GET | `socio/tipo-dependentes-esocial` |

### upload-documentos

| Método | Caminho |
|---|---|
| POST | `upload-documentos/extrato-aplicacao-financeira/enviar/sem-aplicacao-financeira` |
| POST | `upload-documentos/extrato/enviar/bucket` |
| GET | `upload-documentos/extrato/v2/init` |

### upsell

| Método | Caminho |
|---|---|
| POST | `upsell/simulador/checkpoint` |
| GET | `upsell/spot` |

## `/api/legado/` (11)

### contrato

| Método | Caminho |
|---|---|
| GET | `contrato/buscarContratoPlanoPagamento/` |

### empresa

| Método | Caminho |
|---|---|
| GET | `empresa/cnpj?cnpj=` |

### evento-tour

| Método | Caminho |
|---|---|
| POST | `evento-tour` |
| GET | `evento-tour/` |

### nota

| Método | Caminho |
|---|---|
| GET | `nota` |

### notafiscal

| Método | Caminho |
|---|---|
| GET | `/notafiscal/cnaeanexosmultiplos/list` |
| PUT | `/notafiscal/salvaranexosprincipais` |
| GET | `notafiscal/buscarNotaEmpresa` |
| GET | `notafiscal/cancelarAPI` |

### socio

| Método | Caminho |
|---|---|
| GET | `socio/list` |

### zendesk

| Método | Caminho |
|---|---|
| GET | `zendesk/artigos?filtro=` |

## `/api/multiusuario/` (4)

### invites

| Método | Caminho |
|---|---|
| POST | `/invites/enviar` |

### status

| Método | Caminho |
|---|---|
| GET | `/status/servico` |

### usuarios

| Método | Caminho |
|---|---|
| PUT | `/usuarios/ativar` |

### usuarios-e-invites

| Método | Caminho |
|---|---|
| GET | `/usuarios-e-invites/consultar` |

## `/api/leads/hubspot/` (2)

### send

| Método | Caminho |
|---|---|
| POST | `send` |

### track-event

| Método | Caminho |
|---|---|
| POST | `track-event` |

