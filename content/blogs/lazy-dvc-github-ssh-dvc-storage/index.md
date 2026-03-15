+++
title = "Acessar storage DVC com as mesmas chaves SSH do GitHub"
date = 2026-03-15T10:00:00-03:00
draft = false
description = "Se seu time já usa GitHub, por que precisaria de credenciais separadas para armazenamento de dados? Este repositório resolve essa fricção."
author = "Guilherme de Castro"
tags = ["go", "dvc", "github", "ssh", "devops", "docker"]
aliases = ["/posts/lazy-dvc-github-ssh-dvc-storage/"]
+++

## O Problema  

Times de ML e Data Science que usam DVC para versionar datasets e modelos grandes enfrentam um problema recorrente: **autenticação duplicada**.

Você já tem chaves SSH no GitHub para acessar repositórios. Mas para usar DVC com armazenamento remoto, precisa:

1. Gerar um novo par de chaves SSH
2. Distribuir para cada membro do time
3. Configurar no servidor de armazenamento
4. Gerenciar revogação manual quando alguém sai do time

Quando um novo engenheiro entra no time, o processo de "onboarding de storage" se torna um ritual `separado do onboarding de código. Quando alguém sai, você lembra de remover o acesso do GitHub, mas esquece de revogar a chave do servidor S3.

## Alternativas Consideradas

**Git LFS** funciona, mas introduz problemas próprios:

- Requer credenciais HTTPS separadas para o storage LFS
- Limites rígidos no plano gratuito (10 GB storage + 10 GB bandwidth/mês)
- Rebase e filter-branch corrompem ponteiros LFS
- Quota excedida =Sem escrita, só leitura

**DVC padrão** resolve o problema de arquivos grandes, mas não resolve autenticação:

- Cada pessoa precisa configurar credenciais AWS/GCS/etc.
- Alternativamente, distribuir chaves SSH manualmente para cada pessoa
- Sem integração com controle de acesso existente

## A Solução

**lazy-dvc** usa o que você já tem: chaves SSH registradas no GitHub.

```bash
# Configure uma vez
export LDVC_GH_TOKEN=ghp_seu_token
export LDVC_GH_ORG_NAME=sua-organizacao

docker compose up -d

# Use DVC normalmente
dvc remote add -d storage ssh://dvc-storage@seu-servidor:2222/data
dvc push
```

Se você está na organização GitHub, você tem acesso ao storage. Se sai, perde acesso automaticamente.

## Como Funciona 

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

1. Você conecta via SSH/SFTP (como faria com qualquer servidor)
2. `authpubk` busca suas chaves públicas no GitHub, filtrado por org/team
3. Se a chave corresponde, acesso liberado ao mount S3
4. DVC funciona normalmente — `dvc push`, `dvc pull`

O container Docker inclui OpenSSH, rclone (para montar S3), e os binários Go para autenticação.

## Comparação Rápida

| Aspecto | Git LFS | DVC padrão | lazy-dvc |
|---------|---------|-------------|-----------|
| Auth para storage | HTTPS + PAT | Chaves SSH separadas | GitHub SSH |
| Onboarding | Adicionar ao repo | + distribuir chave SSH | + adicionar ao GitHub team |
| Offboarding | Remover do repo | + revogar chave do servidor | + remover do GitHub team |
| CI/CD | Configurar LFS + credenciais | Configurar credenciais storage | Usar deploy keys SSH existentes |
| Storage limits | Limitado pelo plano | Seu backend | Seu backend |

## Próximos Passos

- **Testar:** `docker compose up -d --build` com suas variáveis de ambiente
- **Repositório:** [github.com/guilycst/lazy-dvc](https://github.com/guilycst/lazy-dvc)
- **Contribuir:** Issues e PRs são bem-vindos
- **Docker:** `ghcr.io/guilycst/lazy-dvc:latest`

 funciona com qualquer backend S3-compatible — MinIO, Ceph, AWS S3, VersityGW. Se rclone suporta, funciona.
