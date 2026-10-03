# Notas fiscais: o que a API oferece

Resultado da investigação das issues #29, #30 e #33 nos fronts do painel
(`NovoEmissorNF.*.js`, `notas*`, `listagem-notas`, `registr*`, `clients`, `app`) e de
notas de entrada (`/nota-entrada/static/js/app.*.js`).

## NFS-e emitidas: sem PDF nem XML

- A listagem (`novo-emissor/v2/listagem/notas/filtro`) não tem campo de link, e nenhum
  código do front baixa ou abre o documento da nota. As buscas por `pdf`, `xml`, `danfe`,
  `download`, `baixar`, `visualizar`, `window.open`, `linkPdf` e `url*` não encontraram nada
  ligado às NFS-e.
- Os textos da tela explicam o fluxo: *"Assim que a nota for autorizada pela prefeitura você
  receberá o arquivo por e-mail"* e *"Esse e-mail é o comprovante da emissão… estamos enviando
  o link para a visualização e impressão"*.
- Os PDFs gerados no navegador (pdfMake) são só dos relatórios contábeis e do informe de
  rendimentos.

**Conclusão:** não há como baixar o PDF ou o XML de uma NFS-e emitida pela API do painel. O
documento chega por e-mail, com o link da prefeitura. Por isso o `ctbz notas` só lista as
notas (#30 fechada como inviável).

## NF-e de entrada: XML e DANFE existem na base `/api/emissor/`

O front de notas de entrada (`/nota-entrada/#/`) usa uma base própria, `/api/emissor/`, que
tem download de arquivos **das notas recebidas** (NF-e de compra, não NFS-e):

| Método | Caminho | O que faz |
|---|---|---|
| GET | `/api/emissor/notasentrada/download/{id}` | XML da NF-e de entrada |
| GET | `/api/emissor/notasentrada/gerar/{id}` | DANFE (PDF) da NF-e de entrada |
| GET | `/api/emissor/nfeEmissao/getXml/{id}` | XML (mesmo fluxo, outro módulo) |
| GET | `/api/emissor/danfe/gerar/{id}` | DANFE (outro módulo) |
| GET | `/api/emissor/classificacaonotas/agendarExportarArquivos/{mes}/{ano}?tipo=&email=` | agenda o envio dos arquivos do mês por e-mail (efeito colateral: não usar) |

O formato dessas respostas não pôde ser verificado: a conta usada no desenvolvimento não tem
notas de entrada, e o front só repassa a resposta (`resolve(a)`) sem mostrar como a trata.
Os downloads ficam documentados aqui e não viram comando até haver um exemplo real.

## Notas de entrada: listagens

O front `/nota-entrada/` (Vue 1 com vue-resource, bundle `static/js/app.*.js`) tem duas
telas, cada uma com duas abas. Todas as listagens respondem `{list, total, cursor,
serializedList}` e foram chamadas de verdade (#33):

| Tela / aba | Caminho | Paginação |
|---|---|---|
| Manifestação / a manifestar | `GET /api/emissor/notasentrada/listar/0?mes&ano&empresa&qtdPagina&cursor` | `cursor` da resposta anterior |
| Manifestação / manifestadas | `GET /api/emissor/notasentrada/listar/1?…` | idem |
| Classificação / a classificar | `GET /api/emissor/classificacaonotas/listar/?mes&ano&empresa&tipo=0&limite&cursor=&offset` | `offset` |
| Classificação / classificadas | `GET /api/emissor/classificacaonotas/listar/?…&tipo=1…` | idem |

- `mes` vai de 1 a 12; `empresa` é um filtro **pela razão social do emitente** (vazio, sem
  filtro).
- Itens (do exemplo do tutorial do front): `{id, chave, cnpjEmitente, razaoSocial,
  inscricaoEstadual, dataEmissao (epoch ms), valor, situacao{id, descricao}}`.
- Escrita (fora do escopo até a 1.0): `POST notasentrada/manifestar/`,
  `POST classificacaonotas/salvarloteclassificacao/{tipo}`, `salvarclassificacao/`,
  `reclassificar?idNfe=`.
- Outros `GET`: `classificacaonotas/listarprodutos/{id}` (itens da nota),
  `/api/plataforma/parametrosEmpresa/list`, `/api/legado/empresa/getEmpresaLogada`.

## Notas tomadas

O front `/nota-tomada/` não existe mais: responde **404** mesmo com sessão válida (verificado
em 2026-10-03). As notas de serviço tomadas não têm tela nem API conhecida no painel.
