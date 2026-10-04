# ADR-0018: Escrita com confirmação, simulação e sem retentativa

- **Status:** aceita
- **Data:** 2026-10-03

## Contexto

A 1.0 estabilizou a leitura. A ADR-0002 adiou as ações que alteram dados (adicionar,
alterar, remover, marcar como concluído) porque a API é privada e o efeito de uma escrita
não pode ser testado sem consequências reais.

A análise estática do JavaScript do painel (ver [Escrita](../api/escrita/README.md))
levantou, para cada escrita, o método, o caminho, o corpo, de onde vêm os valores, as
validações do front, as mensagens e se existe um endpoint para desfazer. Com isso dá para
reproduzir as escritas com fidelidade sem chamá-las durante o desenvolvimento.

Os riscos são de natureza diferente da leitura:

- uma escrita repetida (retentativa, re-login no meio) pode duplicar um lançamento ou
  aceitar um termo duas vezes;
- várias escritas são irreversíveis (aceites, manifestações à SEFAZ, contratação de
  parcelamento, exclusões) e algumas geram cobrança;
- scripts e cron (o público da 1.0) não têm ninguém para responder a um prompt.

## Decisão

A partir da 1.1, a CLI escreve. Cada contexto entra numa versão *minor*, como a leitura
entrou de v0.1 a v0.8. Todo comando de escrita segue estas regras:

1. **Lê antes, valida e mostra.** O comando faz os `GET` de preparação, aplica as mesmas
   validações do front (ex.: lançamento `confirmadoViaSistema` não é editável) e mostra
   um resumo legível da ação.
2. **Confirmação.** No terminal, pergunta antes de enviar. Sem terminal, só envia com
   `--yes`; sem ele, é erro de uso (código 2) e nada é enviado. Recusar a confirmação
   termina com código 1 ("operação cancelada").
3. **Risco declarado.** Cada comando tem risco `baixo`, `médio` ou `alto`. No risco alto
   (irreversível, fiscal, legal ou com cobrança), a confirmação mostra a consequência
   (texto do termo, custo, "não dá para desfazer") e exige digitar `confirmo`.
4. **`--dry-run`.** Mostra método, caminho e corpo e termina com 0, sem chamar a escrita.
   Segredos (senhas, códigos) aparecem como `***`.
5. **Uma única tentativa.** Escrita não tem retentativa nem re-login automático depois de
   enviada (estende a ADR-0013). A sessão é conferida antes do envio.
6. **Lê depois.** Quando existe um `GET` para o estado, o comando relê e mostra o resultado.
   Em `-o json`, as chaves de resultado entram no contrato público da ADR-0017.
7. **Testes de requisição.** Cada escrita tem um teste que compara método, caminho e corpo
   com um golden. Escritas nunca são chamadas ao vivo em teste nem em investigação.
8. **Registro local.** Cada escrita enviada vira uma linha em `$CTBZ_HOME/acoes.jsonl`
   (data, empresa, comando, método, caminho, status), sem corpo.

Ficam fora do roadmap de escrita, por decisão do mantenedor: pagamentos e cartões
(o cartão é cifrado no navegador pelo Adyen ou pela Iugu), certificado digital e emissão,
replicação e cancelamento de NFS-e. Os endpoints continuam documentados.

## Consequências

- `ctbz api -X MÉTODO` continua sendo a ferramenta de baixo nível, sem essas proteções.
- A camada `internal/api` passa a ter tipos de requisição e funções de escrita, ao lado
  das leituras (ADR-0009), e uma lista `Escritas()` para os testes.
- Os comandos de escrita ficam mais longos que os de leitura: a preparação e a releitura
  fazem parte deles.
- O formato das escritas não pode ser verificado ao vivo; mudanças no front são detectadas
  pelo monitoramento do catálogo (ADR-0016), que passa a cobrir os caminhos de escrita.
- Substitui a [ADR-0002](0002-somente-leitura-ate-1-0.md).

## Alternativas consideradas

- **Continuar só leitura** — deixa de fora as tarefas mais frequentes do painel (confirmar
  pagamento de guia, lançar no caixa, enviar extrato).
- **Escrever sem confirmação, como `ctbz api -X`** — perigoso em ações irreversíveis e fácil
  de disparar por engano num script.
- **Testar escritas contra a conta real** — cada teste seria uma ação de verdade na
  contabilidade da empresa.
- **Retentativa com chave de idempotência** — a API não aceita chave de idempotência.
