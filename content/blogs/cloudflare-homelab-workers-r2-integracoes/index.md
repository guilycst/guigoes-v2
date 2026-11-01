+++
title = "Cloudflare para homelabs, parte 5: Workers, R2 e pequenos projetos"
description = "Sites estáticos, pequenas funções, arquivos e integrações na Cloudflare. O que Workers, R2, Turnstile e Email Routing oferecem gratuitamente e quando o uso gera cobrança."
language = "pt-BR"
authors = ["Guilherme de Castro"]
publishedAt = "2026-11-01T01:33:22-03:00"
tags = ["cloudflare", "homelab", "workers", "r2", "self-hosting"]
draft = false
series = "Cloudflare para homelabs"
seriesPart = 5
cover = "./assets/cover-cartoon.png"
cover_alt = "A nuvem laranja inspirada na Cloudflare organiza caixas numa bancada ao lado de uma página web ciano, representando pequenos projetos e seus arquivos."
margo = { imageViewer = true }
+++

O site de um grupo não precisa sair do ar quando o servidor de casa está desligado. Um formulário de contato ou uma pequena integração também pode ter uma vida separada do homelab.

Nas partes anteriores, a Cloudflare ajudou a encontrar e alcançar aplicações, com DNS, Tunnel e [controle de acesso](/blogs/cloudflare-homelab-access-rede-privada-gateway/). Nesta parte, vamos olhar o que pode rodar ou ficar armazenado na própria plataforma. O gratuito já atende a muitos projetos pequenos, mas cada produto mede o uso de uma forma.

## Workers executa funções pequenas

O [Cloudflare Workers](https://developers.cloudflare.com/workers/) executa código para atender requisições, sem você manter uma máquina para esse processo. Uma pequena API, um receptor de webhook ou uma integração entre serviços são exemplos que combinam com esse modelo.

Imagine uma ferramenta que recebe um evento externo, confere seus dados e o encaminha a outro serviço. Se o trabalho é curto, ela pode rodar num Worker e ficar independente da máquina que você usa para experimentar em casa. Se o destino ainda está no homelab, essa dependência continua existindo.

Os [principais limites do Workers Free](https://developers.cloudflare.com/workers/platform/limits/) incluem:

| Recurso | Limite gratuito |
|---|---|
| Requisições à função | 100 mil por dia, por conta. |
| Tempo de CPU | 10 ms por execução. |

O tempo de CPU mede o trabalho executado pelo código. Não é todo o tempo que a função passa esperando uma resposta de rede. Ainda assim, uma transformação pesada de arquivos ou uma análise demorada pode ultrapassar essa capacidade mesmo recebendo poucas visitas.

O limite diário reinicia à meia-noite UTC. Depois de esgotar a quota no Free, o Worker não executa o código para novas requisições. Em rotas ligadas a uma origem, o modo configurado pode deixar o acesso seguir sem a função ou devolver um erro. Migrar para o Workers Paid é uma contratação separada; ter o domínio num plano pago não faz essa mudança automaticamente.

Workers não é uma VPS para instalar qualquer programa que você já roda em casa. O projeto precisa funcionar no ambiente de execução e caber nas quotas. Isso ajuda a escolher funções pequenas para a plataforma, mantendo outros serviços onde fazem mais sentido.

## Um site estático não precisa executar uma função a cada visita

Um blog, uma documentação ou um site de apresentação pode ser composto por HTML, imagens, estilos e scripts já prontos. O [Workers Static Assets](https://developers.cloudflare.com/workers/static-assets/) hospeda e entrega esses arquivos, sem precisar consultar seu servidor de casa.

As [requisições servidas diretamente como arquivos estáticos são gratuitas e ilimitadas](https://developers.cloudflare.com/workers/static-assets/billing-and-limitations/). Quando uma requisição executa o código do Worker, entram as quotas da função. Assim, um site estático não fica limitado às 100 mil execuções diárias só por estar no mesmo projeto.

Existem limites de armazenamento desses arquivos na publicação. No Free, são até 20 mil arquivos por versão do Worker, com no máximo 25 MiB por arquivo. Isso cabe bem em muitos sites pequenos, mas não transforma Static Assets num repositório de grandes downloads.

Uma página que depende de uma API em casa ainda depende daquela API para mostrar os dados. Já os textos e arquivos publicados na Cloudflare continuam disponíveis sem a conexão do homelab. Essa separação costuma ser útil para um site público simples.

## R2 guarda os arquivos da aplicação

O [R2](https://developers.cloudflare.com/r2/) armazena objetos. Pense em arquivos identificados por nomes, acessados por uma API compatível com S3 ou por um Worker. Ele pode guardar anexos de um pequeno projeto ou uma cópia de dados fora do homelab.

É diferente de um disco montado no servidor ou de uma pasta compartilhada da rede. A aplicação precisa saber trabalhar com esse tipo de armazenamento, e a [compatibilidade com a API S3](https://developers.cloudflare.com/r2/api/s3/api/) tem seu próprio escopo.

A [franquia gratuita mensal do R2 Standard](https://developers.cloudflare.com/r2/pricing/) inclui:

| Uso | Franquia |
|---|---|
| Armazenamento | 10 GB-mês. |
| Operações de classe A, como gravações e listagens | 1 milhão por mês. |
| Operações de classe B, como leituras | 10 milhões por mês. |

Manter 10 GB durante o mês inteiro é um exemplo de 10 GB-mês. A conta considera o armazenamento ao longo do período, não uma franquia de 10 GB novos para enviar a cada mês. Essa faixa gratuita não se aplica à classe Infrequent Access.

O tráfego de saída do R2 não tem tarifa de transferência. Armazenamento e operações acima das franquias, porém, geram cobrança. Um arquivo pequeno consultado muitas vezes pode consumir a franquia de operações sem ocupar muito espaço.

O [R2 tem uma assinatura própria](https://developers.cloudflare.com/r2/get-started/), mesmo para aproveitar o uso gratuito incluído. As franquias não são um teto automático de gastos. Para uma cópia pequena de arquivos de projeto, elas podem bastar; uma biblioteca inteira de fotos ou backups do NAS pede outra conta de capacidade.

Os arquivos podem continuar privados. Se forem públicos, o [endereço de desenvolvimento `r2.dev`](https://developers.cloudflare.com/r2/buckets/public-buckets/) tem limite de taxa e não é destinado à produção. Um domínio próprio é o caminho previsto para a entrega pública de produção.

## Turnstile para formulários e cadastros

O [Turnstile](https://developers.cloudflare.com/turnstile/plans/) ajuda a distinguir interações legítimas de automação abusiva. Pode ser útil num formulário de contato ou num cadastro, inclusive se o site não usa o proxy da Cloudflare.

O Free inclui até 20 widgets, as configurações de integração usadas pelos sites, com até 10 domínios ou subdomínios autorizados por widget. Os desafios e as verificações não têm uma franquia de volume nesse plano.

O componente da página produz um token, e o servidor precisa [validá-lo antes de aceitar a operação](https://developers.cloudflare.com/turnstile/get-started/server-side-validation/). Só exibir o componente não protege o formulário. Isso é uma integração com a aplicação; não é necessário publicar todo o homelab para usá-la.

## Email Routing para um endereço do projeto

O [Email Routing](https://developers.cloudflare.com/email-service/) pode receber mensagens de um endereço como `contato@example.com` e encaminhá-las para uma caixa de e-mail que você já usa. É útil para dar ao projeto um contato no próprio domínio sem hospedar um servidor de e-mail em casa.

O [encaminhamento é gratuito](https://developers.cloudflare.com/email-service/platform/pricing/), mas não cria uma caixa postal própria. O recebimento depende dos registros de e-mail do domínio; esse uso precisa ser compatível com qualquer provedor de e-mail que já atende o mesmo nome.

Enviar mensagens a destinatários arbitrários é outro recurso, Email Sending, ainda em beta, que exige Workers Paid. Há uma modalidade gratuita de envio a endereços de destino já verificados na conta. Se você usar um Worker para processar as mensagens recebidas, também valem as condições de uso do Workers.

## Cada peça pode entrar sozinha

Um grupo pode hospedar só seu site estático e continuar usando ferramentas que já conhece. Um projeto pessoal pode aproveitar apenas o R2 ou acrescentar Turnstile ao formulário existente. Quem quer um endereço de contato pode usar o encaminhamento sem migrar sua aplicação inteira.

Essa independência é útil para experimentar sem transformar o homelab inteiro em uma migração. O plano Free do domínio, Workers Free, as franquias do R2 e o Zero Trust Free são coisas distintas. O que manda é o conjunto de produtos que você ativa e o uso de cada um.

Para indivíduos, homelabs e equipes pequenas, os recursos gratuitos podem sustentar serviços de uso real. Necessidades de suporte, disponibilidade e crescimento continuam entrando na escolha. Quando uma quota deixa de servir, dá para avaliar aquele componente sem presumir que todo o resto precisa mudar junto.

*Referências consultadas em 1º de novembro de 2026. Limites e condições podem mudar.*
