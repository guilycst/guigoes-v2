# Ícones de contato

A página Contato e o rodapé compartilham SVGs de Bootstrap Icons v1.13.1, MIT, adquiridos pelo
Muamba v0.0.6. O commit de origem está fixado em `.muamba.yaml`; URLs, tamanhos
e hashes SHA-384 estão em `.muamba.lock.yaml`. A licença acompanha os SVGs em
`content/contact/assets/icons/LICENSE` e é publicada junto deles.

Decisão de integração: Heroicons não inclui as marcas sociais. O iconpack do
Goshtoso pode gerar um sprite e um pacote Go, mas esta página Margo só precisa
de imagens estáticas. Reutilizamos Bootstrap Icons via Muamba, sem acrescentar
Goshtoso às dependências do gerador do blog. SVGs revisados, sem scripts,
recursos externos ou alterações no desenho.

Os ícones são decorativos, com `alt=""`, junto de rótulos visíveis. O CSS mantém
contraste nos temas claro, escuro e do sistema. Discord preserva o nome de
usuário, sem inventar um endereço de perfil a partir dele.

Verificar arquivos locais, sem baixar assets:

```sh
go run github.com/araihu/muamba/cmd/muamba@v0.0.6 verify --strict
```

Restaurar arquivos apenas quando os bytes correspondem ao lock:

```sh
go run github.com/araihu/muamba/cmd/muamba@v0.0.6 sync --strict
```

`go run` pode baixar o executável e módulos para o cache Go. O build normal usa
os SVGs versionados e não executa Muamba. Para mudar a versão, revise a origem,
gere o novo lock explicitamente e confira ícones em desktop/celular e ambos
os temas antes de publicar. Não rode `lock` ou `update` automaticamente no CI.
