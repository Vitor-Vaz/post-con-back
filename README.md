# post-con-back

API em **Go** do **PostoConfiável** (avaliações de postos de combustível no Brasil). Visão de produto e roadmap: `technical-plan/SAAS-POSTOS-COMBUSTIVEL.md`.

## Contexto para humanos e IA

Este repositório prioriza documentação na **raiz** para continuidade entre sessões e entre ferramentas.

- **`AGENTS.md`** — leitura principal para assistentes de IA e para alinhar stack, camadas, comandos, CI e convenções do código.
- Este **`README.md`** — referência rápida ao fluxo de trabalho e ao repositório vizinho **mesa-mestre**.

## Repositório de referência (mesa-mestre)

O **post-con-back** costuma ficar no **mesmo diretório pai** que o **mesa-mestre** (Makefile, migrações, CI, organização de pastas).

- Caminho relativo típico: **`../mesa-mestre`**

Antes de espelhar comportamento (`make`, testes, workflow de PR), vale abrir o mesa-mestre nesse caminho e comparar.

## Como usar na prática

**Antes de começar**

- Ler **`AGENTS.md`**.

**Ao encerrar uma sessão relevante**

- Atualizar **`AGENTS.md`** (por exemplo fila resumida ou decisões estáveis) e/ou `technical-plan/SAAS-POSTOS-COMBUSTIVEL.md` quando o escopo de produto mudar.
- Preferir commits e PRs descritivos para histórico fino (o que mudou e por quê).

## Objetivo

Reduzir retrabalho: qualquer pessoa ou agente retoma o trabalho com o mesmo vocabulário e as mesmas restrições já acordadas, sem depender de uma pasta extra de documentação.
