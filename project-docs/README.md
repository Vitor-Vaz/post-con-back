# Contexto do projeto

Este diretório existe para manter contexto portátil entre sessões de trabalho e entre diferentes assistentes de IA.

## Repositório de referência (mesa-mestre)

O **post-con-back** quase sempre fica no **mesmo diretório pai** que o **mesa-mestre** (projeto Go de referência para Makefile, migrações, padrões de pastas).

- Caminho típico do mesa-mestre: `../mesa-mestre` em relação à raiz deste repo.
- Exemplo absoluto (máquina de desenvolvimento): `/home/vitorandrade/Documentos/projetos-pessoais/mesa-mestre` ao lado de `/home/vitorandrade/Documentos/projetos-pessoais/post-con-back`.

Antes de espelhar comportamento (Make, Docker, sqlc), comparar com o código do mesa-mestre nesse caminho.

## Como usar

Antes de começar:

- ler `ai-context.md`
- ler `current-state.md`
- verificar `next-tasks.md`

Ao terminar uma sessão:

- atualizar `current-state.md`
- atualizar `next-tasks.md`
- registrar decisões relevantes em `decisions.md`

## Objetivo

Permitir continuidade de desenvolvimento sem depender da memória de ferramenta, editor ou sessão.
