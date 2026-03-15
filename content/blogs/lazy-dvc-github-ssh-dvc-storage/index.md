+++
title = "lazy-dvc: Gerenciando dados grandes no DVC com autenticação GitHub SSH"
date = 2025-03-15T10:00:00-03:00
draft = false
description = "Um bridge de autenticação que permite usar chaves SSH do GitHub para acessar armazenamento DVC, filtrado por organização e time."
author = "Guilherme de Castro"
tags = ["go", "dvc", "github", "ssh", "devops", "docker"]
aliases = ["/posts/lazy-dvc-github-ssh-dvc-storage/"]
+++

Gerenciar arquivos grandes em projetos de Machine Learning e Data Science sempre foi um desafio. Git LFS? Limitado e caro. DVC padrão? Exige credenciais separadas para storage. O **lazy-dvc** resolve isso de forma elegante: use as mesmas chaves SSH que você já tem no GitHub.

## O Problema

Quando trabalhamos com datasets grandes, modelos treinados, ou outros artefatos de ML, o Git não consegue lidar bem com arquivos grandes. Duas soluções populares:

### Git LFS

Problemas:
- **Credenciais extras**: Requer HTTPS + PAT separado para storage
- **Limites rígidos**: 10 GB storage + 10 GB bandwidth/mês no plano gratuito
- **Problemas de histórico**: Rebase e filter-branch corrompem ponteiros LFS
- **CI/CD complicado**: Cada job precisa de `git lfs install` + credenciais

### DVC Padrão

Problemas:
- **Autenticação duplicada**: Chaves SSH para Git + credenciais separadas para storage
- **Onboarding trabalhoso**: Novo membro precisa gerar chave SSH, distribuir para servidor
- **Gestão manual de acesso**: Revogar acesso requer remover chave do servidor manualmente

## A Solução: lazy-dvc

**Uma única credencial para tudo.**

```
┌─────────────┐                        ┌─────────────────┐
│   Developer │     SSH keys           │    GitHub       │
│             │ ─────────────────────► │   (org/team)    │
└─────────────┘                        └─────────────────┘
       │                                      │
       │                                      │ mesmas chaves
       ▼                                      ▼
┌─────────────┐                        ┌─────────────────┐
│   dvc push  │ ──── SSH/SFTP ───────► │   lazy-dvc      │
│   dvc pull   │                        │   → S3 Backend   │
└─────────────┘                        └─────────────────┘
```

### Como Funciona

1. Membro da organização GitHub faz `dvc push` ou `dvc pull`
2. lazy-dvc recebe conexão SSH
3. `authpubk` busca chaves públicas do GitHub (org/team)
4. Se a chave corresponde → acesso concedido ao storage via S3

### Binários

| Binário | Propósito |
|---------|-----------|
| `lazypubk` | Busca chaves públicas SSH de membros da org/team GitHub |
| `authpubk` | Wrapper para SSH AuthorizedKeysCommand |
| `noshell` | Shell mínimo para sessões SSH/SFTP |

## Quick Start

```bash
# Clone e configure
git clone https://github.com/guilycst/lazy-dvc.git
cd lazy-dvc

# Configure suas variáveis de ambiente
export LDVC_GH_TOKEN=ghp_seu_token
export LDVC_GH_ORG_NAME=sua-organizacao

# Execute
docker compose up -d --build

# Configure DVC remote
dvc remote add -d storage ssh://dvc-storage@localhost:2222/data
dvc push
```

## Comparação

| Recurso | Git LFS | DVC Padrão | lazy-dvc |
|---------|---------|------------|----------|
| Autenticação | Git SSH + HTTPS PAT para LFS | Git SSH + credenciais storage | Apenas GitHub SSH |
| Onboarding | Adicionar ao repo | Adicionar ao repo + distribuir chave SSH | Adicionar ao GitHub team |
| Offboarding | Remover do repo | Remover do repo + revogar chave do servidor | Remover do GitHub team |
| Custo storage | Limitado (10-250 GB) | Seu backend S3 | Seu backend S3 |
| CI/CD | Configurar LFS + credenciais | Configurar credenciais storage | Usar deploy keys SSH existentes |

## Arquitetura

O projeto é distribuído como uma imagem Docker que inclui:

- **OpenSSH Server**: Para conexões SSH/SFTP
- **rclone**: Para montar backend S3
- **lazypubk/authpubk**: Binários Go para autenticação GitHub
- **Entrypoint**: Supervisão de processos com logging unificado

### Cache Inteligente

Para evitar rate limits da API do GitHub, as chaves são cacheadas localmente:

```
/var/cache/lazy-dvc/keys.json
```

Cache com TTL de 5 minutos e file locking para segurança em conexões concorrentes.

## Onde Usar

- Times de ML/Data Science usando DVC
- Organizações que já usam GitHub para código
- Projetos que precisam de controle de acesso baseado em times
- Qualquer projeto que prefira infra própria ao invés de quotas de vendor

## Links

- **Repositório**: [github.com/guilycst/lazy-dvc](https://github.com/guilycst/lazy-dvc)
- **Docker Image**: `ghcr.io/guilycst/lazy-dvc`
- **Licença**: MIT

## Conclusão

O lazy-dvc resolve um problema real de forma pragmática: se você já tem chaves SSH no GitHub por ser membro de uma organização, por que precisaria de outra credencial para storage? 

A implementação é relativamente simples, Dockerizada, e funciona com qualquer backend S3-compatible. Se seu time já usa GitHub e DVC, vale a pena dar uma olhada.