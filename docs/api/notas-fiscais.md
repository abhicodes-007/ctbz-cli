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

As listagens de notas de entrada estão em [endpoints verificados](endpoints-verificados.md)
(#33).
