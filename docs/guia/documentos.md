# Documentos e certificado

## Central de documentos

```sh
ctbz documentos                                           # tipos aceitos e quantos foram enviados
ctbz documentos --tipo EXTRATO_BANCARIO_MOVIMENTACOES     # documentos enviados de um tipo
```

```text
Tipo                            Nome              Enviados  Última modificação
EXTRATO_BANCARIO_MOVIMENTACOES  Extrato bancário        13  10/09/2026
```

- Sem `--tipo`, a lista mostra os tipos da central de documentos (extratos, contratos de
  empréstimo e financiamento, estoque, informes de investimentos…) e quantos já foram
  enviados.
- Com `--tipo`, aparecem os arquivos enviados: competência, nome do arquivo, data de envio,
  descrição, valor (quando o tipo tem) e o link do arquivo.
- Enviar documentos continua sendo feito pelo painel (a CLI só lê, ver
  [ADR-0002](../adr/0002-somente-leitura-ate-1-0.md)).
