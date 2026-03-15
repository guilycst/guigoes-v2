+++
title = "Anunciando lazy-dvc: Chega de gerenciar múltiplas identidades para seus dados"
date = 2026-03-15T10:00:00-03:00
draft = false
description = "Se seu time já tem acesso ao código via GitHub, por que o acesso aos dados precisa ser um processo separado? Conheça o lazy-dvc."
author = "Guilherme de Castro"
tags = ["go", "dvc", "github", "ssh", "devops", "platform-engineering"]
+++

## A "Lacuna de Identidade" no Onboarding de Times

Todo mundo que trabalha em times de engenharia já passou por isso: um novo colega entra no time, você dá acesso ao repositório no GitHub e, cinco minutos depois, ele te manda uma mensagem no Slack:

> *"Consegui clonar o projeto, mas não consigo baixar as dependências/assets/datasets. Onde pego a chave de acesso do storage?"*

Nesse momento, começa o **ritual do segundo onboarding**: você precisa criar um usuário no IAM do provedor de cloud, gerar um par de chaves Access/Secret, enviar de forma segura para o colega e torcer para que ele saiba onde configurar isso localmente.

Quando esse colega sai do time, o processo se repete ao contrário — e é aqui que o risco de segurança mora: muitas vezes removemos o acesso ao código, mas esquecemos de revogar a chave do bucket S3.

**Por que o acesso aos assets não é automático, se o acesso ao código já foi concedido?**

## Apresentando o lazy-dvc

O **lazy-dvc** nasceu para fechar essa lacuna. A filosofia é simples: **Sua identidade no GitHub deve ser a única que você precisa.**

Ele atua como uma ponte de autenticação entre DVC e o Git provider. Em vez de você gerenciar chaves de storage para cada membro do time, o `lazy-dvc` utiliza as **chaves SSH públicas que seus desenvolvedores já cadastraram no GitHub**.

### Como a mágica acontece (sem esforço)

O `lazy-dvc` unifica o fluxo de trabalho. Se o seu desenvolvedor já usa chaves SSH para dar `git push`, ele agora usa **essas mesmas chaves** para dar `dvc pull`.

1. **Identidade como Código:** O `lazy-dvc` consulta a sua Organização ou Time no GitHub em tempo real.
2. **Autenticação Dinâmica:** Quando alguém tenta acessar o storage via DVC, o servidor verifica se a chave SSH daquela pessoa pertence a um membro autorizado no GitHub.
3. **Storage Agregado:** Uma vez autenticado, o usuário tem acesso a um mount (via `rclone`) que aponta diretamente para o seu backend (S3, MinIO, Ceph, etc).

```
┌─────────────┐                        ┌─────────────────┐
│   Developer │     SSH keys           │    GitHub       │
│             │ ─────────────────────► │   (org/team)    │
└─────────────┘                        └─────────────────┘
       │                                      │
       │                                      │ same keys
       │                                      │
       ▼                                      ▼
┌─────────────┐                        ┌─────────────────┐
│   dvc push  │ ──── SSH/SFTP ───────► │   lazy-dvc      │ 
│   dvc pull  │                        │                 │
└─────────────┘                        └─────────────────┘
                                              │
                                              │ S3 Backend
                                              │
                                              ▼
                                       ┌─────────────────┐
                                       │   AWS,B2,MinIO  |
                                       │                 │
                                       └─────────────────┘
```

## Como isso ajuda?

### 1. Onboarding em 0 Segundos

Adicionou o desenvolvedor ao time no GitHub? Pronto. Ele já pode rodar um `dvc pull`. Não há necessidade de distribuir credenciais adicionais.

### 2. Offboarding Garantido

Removeu o colaborador da Organização? O acesso ao storage é cortado instantaneamente. A segurança do seu storage de ativos pesados passa a herdar a mesma governança que você já aplica ao seu código.

### 3. Infraestrutura Transparente

Para o desenvolvedor, é apenas SSH. Ele não precisa saber se os dados estão no S3 da Amazon, em um servidor MinIO local ou em um storage distribuído. O `lazy-dvc` abstrai a complexidade do backend.

### 4. CI/CD Simplificado

Em pipelines, você pode usar as mesmas *Deploy Keys* que já usa para o Git.

## Testando o lazy-dvc

O setup é de teste é via docker compose e pode ser verificado em poucos minutos.

1. Clone https://github.com/guilycst/lazy-dvc 

2. Execute no seu shell

```bash
# Configure as variáveis de ambiente
export LDVC_GH_TOKEN=ghp_seu_token_com_leitura_de_org
export LDVC_GH_ORG_NAME=sua-organizacao-no-github

# Suba o serviço
docker compose up -d

# No lado do cliente, o DVC agora fala SSH puro com as chaves que já estão na máquina
dvc remote add -d meu-storage ssh://dvc-storage@meu-servidor:2222/data

```
---

## Conclusão

O objetivo do `lazy-dvc` é ser uma ferramenta para reduzir a carga cognitiva de quem opera infraestrutura e de quem interage com ela. Menos chaves para gerenciar significa menos chances de erro humano e mais tempo focando no que realmente importa: o produto.

O projeto é **Open Source** e você pode conferir o código, contribuir ou abrir issues aqui:

👉 **[github.com/guilycst/lazy-dvc](https://github.com/guilycst/lazy-dvc)**

