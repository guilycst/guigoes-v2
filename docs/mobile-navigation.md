# Menu mobile

Até 700 px, o botão no cabeçalho abre um painel pela esquerda. Os mesmos cinco
links são movidos para o painel, sem duplicar links nem o controle de tema.
Acima desse limite, retornam ao cabeçalho. Sem JavaScript ou `showModal`, os links
continuam no cabeçalho e o botão fica oculto.

## Escolha do componente

Referência visual: menu lateral do catálogo em demo.manja.araihu.com.
Consultados `drawer.Drawer` e `sidebar` do Goshtoso v0.4.5, `MobileNavigation`
do App Shells v0.1.9-0.20260910224508-5b2222e54637 e a configuração de layouts
do Margo v0.0.25. O Drawer depende de Alpine e seus assets; MobileNavigation
oferece painel superior. O CLI do Margo expõe componentdocshell, mas não um
Drawer isolado para o layout article usado aqui.

Este blog gera HTML com o CLI e não carrega Alpine. Usamos `<dialog>` nativo
para modalidade, Escape e retorno de foco, com código local para abrir/fechar,
mover os links no breakpoint e circular Tab nos limites. Não copiamos o runtime
dos shells nem adicionamos dependências para a navegação.

## Verificação no navegador

- Em 390 px, abrir pelo botão, conferir cinco links e foco em Fechar menu.
- Shift+Tab no primeiro controle chega ao RSS; Tab no RSS volta ao primeiro.
- Fechar por botão, Escape e fundo libera a rolagem e restaura o foco.
- Seguir Blog fecha o painel e navega para /blogs/.
- Redimensionar com o painel aberto para 1440 px fecha o diálogo, devolve os
  links ao cabeçalho e foca Blog.
- Conferir temas claro/escuro, largura de 320 px, ausência de overflow e fallback
  com os scripts removidos de uma cópia temporária do HTML gerado.
