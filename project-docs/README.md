# Contexto do projeto

Este diretório existe para manter **contexto portátil** entre sessões de trabalho e entre diferentes assistentes de IA, sem depender só do histórico do editor ou da ferramenta.

## Repositório de referência (mesa-mestre)

O **post-con-back** quase sempre fica no **mesmo diretório pai** que o **mesa-mestre** — projeto Go usado como referência para **Makefile**, **migrações**, **CI** e organização de pastas.

- Caminho relativo típico: **`../mesa-mestre`**
- Exemplo absoluto (ajuste ao seu usuário e pasta):  
  `/home/vitorandrade/Documentos/projetos-pessoais/mesa-mestre`  
  ao lado de  
  `/home/vitorandrade/Documentos/projetos-pessoais/post-con-back`

Antes de espelhar comportamento (comandos `make`, fluxo de testes, workflow de PR), vale abrir o mesa-mestre nesse caminho e comparar.

## Como usar

**Antes de começar uma sessão**

- Ler `ai-context.md` (projeto, stack, arquitetura, contratos, CI).
- Ler `current-state.md` (o que já existe, arquivos-chave, riscos).
- Verificar `next-tasks.md` (fila sugerida).

**Ao terminar uma sessão**

- Atualizar `current-state.md` (o que mudou, o que ainda falta).
- Atualizar `next-tasks.md` (itens fechados ou repriorizados).
- Registrar decisões relevantes em `decisions.md` (ADRs curtos).

## Objetivo

Permitir **continuidade de desenvolvimento** com menos retrabalho: qualquer assistente ou pessoa pode ler estes arquivos e retomar o trabalho no mesmo lugar, com o mesmo vocabulário e as mesmas restrições já acordadas.
