+++
title = "Quando não dá para fugir do JavaScript: Hyperscript te salva"
date = 2024-04-30T01:37:09-03:00
draft = false
description = "As vezes não tem jeito, você vai ter que usar JavaScript. Mas não se preocupe, o Hyperscript vai te ajudar a escrever menos código e deixar seu código mais legível."
author = "Guilherme de Castro"
tags = ["web", "frontend", "JS", "hyperscript"]
aliases = ["/posts/hyperscript_for_the_rescue/"]
cover = "./assets/js_headache.webp"
cover_alt = "Dev de frente para o computador sofrendo com o JavaScript"
+++
Recentemente escrevi um post sobre [HTMX e como ele pode ajudar a escrever menos JS](/blogs/creating_a_cheap_blog_golang_htmx), mas a realidade cruel do desenvolvimento web é que SSR (Server Side Rendering) só te leva até certo ponto, algumas experiências desejáveis na Web só vão ser possibilitadas executando código do lado cliente.

É ai onde JavaScript normalmente entraria, mas felizmente hoje existem outras opções, nesse artigo vou falar sobre uma delas - O Hyperscript

## O que é Hyperscript

É uma linguagem de script feita para interceptar e responder à eventos do lado do cliente em uma web app, assim como realizar manipulações simples na DOM, veja esse exemplo:

<div class="border-[2px] border-zinc-50 p-2 mb-2">


```html
<button _="on click increment :x then put the result into the next <output/>">
Me clica!
</button>
<output>--</output>
```


<button class="bg-blue-599 p-2 text-white rounded-md" _="on click increment :x then put the result into the next <output/>">
Me clica!
</button>
<output>--</output>

</div>


Sintaxe peculiar, não? A primeira a coisa a se notar é que o script é definido diretamente no elemento usando o atributo _ (underscore).

Incorporar comportamento diretamente no elemento pode parecer estranho, mas é uma tendência de tecnologias mais recentes de favorecer [Locality of Behavior](https://htmx.org/essays/locality-of-behaviour/) ao invés de [Separation of Concerns](https://en.wikipedia.org/wiki/Separation_of_concerns). Exemplos de outras bibliotecas indo nesse mesmo caminho são Tailwind CSS, AlpineJS e HTMX.

As vantagens dessa sintaxe excêntrica é que fica claro que hyperscript está sendo usado na página. 

A linguagem também é declarativa e de fácil leitura para humanos, me lembra SQL, você não precisa definir de forma procedural como realizar comportamentos comuns em uma web page, Hyperscript abstrai a manipulação da DOM para você, assim como SQL abstrai o comportamento interno do banco de dados e a interação com o sistema de arquivos.

Para efeitos de comparação, segue o mesmo código usando Vanilla JS

<div class="border-[2px] border-zinc-50 p-2 mb-2">

```html
<button id="buttonjs">
Me clica!
</button>
<output>--</output>
<script>
    let x = 1;
    document.querySelector("#buttonjs").addEventListener("click", function() {
        const output = this.nextElementSibling;
        if (output.nodeName == "OUTPUT"){
            output.innerHTML = ++x;
        } 
    });
</script>
```

<button id="buttonjs" class="bg-blue-599 p-2 text-white rounded-md">
Me clica!
</button>
<output>--</output>
<script>
    let x = 1;
    document.querySelector("#buttonjs").addEventListener("click", function() {
        const output = this.nextElementSibling;
        if (output.nodeName == "OUTPUT"){
            output.innerHTML = ++x;
        } 
    });
</script>

</div>

Com Hyperscript fica bem mais clean 😊

# Emitindo eventos customizados

As vezes queremos lidar com eventos que não necessariamente são oriundos dos componentes nativos da DOM, mas sim eventos customizados introduzidos por código proprietário, Hyperscript permite trabalhar com esses tipos de eventos também, essa modal por exemplo:

<div class="border-[2px] border-zinc-50 p-2 mb-2">

```html
<button 
    _="on click toggle .hidden .modal on #modal" 
    class="bg-blue-599 p-2 text-white rounded-md">
    Abrir modal
</button>
<div id="modal" class="hidden" 
_="on closeModal add .closing then wait for animationend then toggle .hidden .modal .closing">
    <div class="modal-underlay" 
        _="on click trigger closeModal">
    </div>
    <div id="modal-content" class="modal-content">
        Hyperscript rocks 💙
        <button 
            _="on click trigger closeModal" 
            class="bg-blue-599 p-2 text-white rounded-md">
            Fechar
        </button>
    </div>
</div>
<style>
.modal {
    /**
    Exibe modal com posição fixa
    Anima a abertura
    **/
}

.modal > .modal-underlay {
	/**
    "Painel de vidro"
    cobre o conteúdo por trás da modal
    **/
}

.modal > .modal-content {
    /**
    Estilos do conteúdo da modal
    **/
}

.modal.closing {
	/* Animação de fechamento da modal */
}

.modal.closing > .modal-content {
	/* Animação de fechamento da modal */
}

</style>
```
<button _="on click toggle .hidden .modal on #modal" class="bg-blue-600 p-2 text-white rounded-md">Abrir modal</button>
<div id="modal" class="hidden" _="on closeModal add .closing then wait for animationend then toggle .hidden .modal .closing">
    <div class="modal-underlay" _="on click trigger closeModal"></div>
    <div id="modal-content" class="modal-content">
        Hyperscript rocks 💙
        <button _="on click trigger closeModal" class="text-zinc-300 bg-zinc-700 hover:text-white rounded-md text-sm p-2 font-medium self-center">Fechar</button>
    </div>
</div>

</div>

Segue o passo a passo:

1. Para exibir a modal, o botão remove a classe ```.hidden``` do elemento principal e adiciona a classe ```.modal```.
2. A ```div``` principal da modal lida com o evento customizado ```closeModal``` que adiciona a classe ```.closing``` que tem as animações de fechamento, remove as classe ```.modal .closing``` quando a animação termina e adiciona novamente a classe ```.hidden```.
3. A ```div``` underlay, que é o "painel de vidro" que cobre todo o view port, e o botão de fechar, disparam o evento ```closeModal``` que vai executar o handler do passo 2.

## Estado atual

Hyperscript ainda está em beta, a sintaxe e funcionalidades estão em maior parte completas, agora o desenvolvedores estão trabalhando em melhorias na documentação e testes para o release da versão 1.0.

Existe também uma dependência com promises o que impossibilita a compatibilidade com IE11.

## Conclusão

Aqui só arranhei a superfície do que a linguagem é capaz, não deixe de visitar o [site oficial](https://hyperscript.org/) e verificar as outras funcionalidades que podem ser úteis para você. No desenvolvimento desse blog usarei hyperscript para adicionar interatividade onde é necessário, até então tem sido um prazer. 

Pode ser que em seu projeto também seja 😄

<script src="https://unpkg.com/hyperscript.org@0.9.12"> </script>
