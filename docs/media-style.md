# Capas do Guigoes

Direção escolhida em 27/10/2026: cartoon fofo de objetos de tecnologia, inspirado no humor e nos personagens de [Charm](https://charm.land/), com personagens próprios. Paleta ciano, teal, lilás, coral e creme; contorno ameixa, sombras simples e textura discreta de papel. Capas horizontais, sem texto, legíveis em miniatura. O tema do artigo determina a metáfora.

Geradas com a ferramenta integrada `imagegen`, uma chamada por capa. Os originais anteriores permanecem nas mesmas URLs. Novas capas em `content/blogs/<slug>/assets/cover-cartoon.png`, consumidas pelo frontmatter `cover` e `cover_alt`, pela home, listagens e JSON-LD.

## Referências para as próximas capas

Use o [catálogo de personagens](media/characters/README.md) e anexe as fichas PNG correspondentes à geração. Ele registra anatomia, símbolos, poses e sugestões de elenco por assunto da bio e dos projetos. Os prompts históricos abaixo documentam a origem; as fichas atuais definem a identidade visual para novas cenas.

## Prompt comum da primeira geração

```text
Create one finished editorial blog cover, landscape 16:9, 1536x864 or equivalent. Original cute cartoon illustration for Guigoes, a Brazilian programming blog. Friendly anthropomorphic technology objects, rounded chunky silhouettes, tiny dark oval eyes, small smiles, blush cheeks, stubby limbs. Refined 2D editorial cartoon with dark plum outlines, flat pastel fills, a hint of soft cel shading and fine paper grain, no photorealism, no glossy 3D. Palette warm cream #f8f3e8 background, Go cyan #67d3ed and teal #007d9c main subjects, lavender #b9a2ef, coral #f593a3 and butter yellow accents, plum #292331 ink. One clear visual metaphor, generous negative space, central composition with all important details within inner 75 percent for responsive crops. Warm playful developer-tool personality. No words, letters, labels, logos, watermarks, interface screenshots, gradients or decorative border. Do not copy any existing brand mascot.
```

O prompt acima registra a geração original. Para novas capas, remova `stubby limbs` e descreva explicitamente os membros necessários, seguindo as regras de revisão abaixo.

## Assuntos usados

### lazy-dvc-github-ssh-dvc-storage

Arquivo atual: `content/blogs/lazy-dvc-github-ssh-dvc-storage/assets/cover-holding-hands.png`.

Edição aprovada: pasta, chave e banco de dados de mãos dadas. Removidas as duas chaves pequenas e as fechaduras. A chave central conserva dois braços e uma única haste dentada, sem pernas. A versão anterior permanece em `cover-cartoon.png`. Prompt da edição em `docs/lazy-dvc-holding-hands-prompt.md`.

Prompt original:

```text
Subject: A smiling cyan SSH key character in the middle opens two friendly adjacent boxes: one little lavender code-folder with simple branching line symbol, and one teal stacked data-cylinder with a tiny heart. Both have matching keyholes. A gentle curved path joins them. The single key grants access to both code and data. Small stars convey relief and simplicity. Three main objects only.
```

Texto alternativo: Uma pasta, uma chave e um banco de dados sorridentes de mãos dadas, representando acesso unificado a código e dados.

### hyperscript_for_the_rescue

Arquivo: `content/blogs/hyperscript_for_the_rescue/assets/cover-cartoon.png`.

```text
Subject: A friendly small cyan terminal-window character with a tiny lavender superhero cape gently untangles a coral spaghetti-like cable into one clean simple line. The tangled cable sits left, relaxed terminal hero center-right, tiny joyful sparks. Terminal face has just eyes and smile, no code or text. Visual metaphor: Hyperscript rescues developers from JavaScript complexity.
```

Texto alternativo: Um terminal de capa de super-herói desembaraça fios, simbolizando a simplicidade do Hyperscript.

### html_meta_tags_social_media_ctr

Arquivo: `content/blogs/html_meta_tags_social_media_ctr/assets/cover-cartoon.png`.

```text
Subject: A cheerful cyan little webpage-card character with a small landscape thumbnail and simple abstract content bars, being neatly dressed by two floating lavender and coral gift-tag characters. A large friendly magnifying glass leans alongside it. The tags improve the card's appearance and discoverability. No letters or readable text. Clear metaphor for metadata making rich search and social previews.
```

Texto alternativo: Tags coloridas organizam um cartão de página ao lado de uma lupa, representando metadados e descoberta.

### creating_a_cheap_blog_golang_htmx

Arquivo: `content/blogs/creating_a_cheap_blog_golang_htmx/assets/cover-cartoon.png`.

```text
Subject: A friendly cyan cloud character holds a tiny cream blog-window like a house, with a little lavender piggy bank sitting beside it and two small butter-yellow coins. Blog window has a single coral picture block and simple dark content lines, no text. Cozy simple composition communicating an inexpensive lightweight serverless blog. No brand mascot.
```

Texto alternativo: Uma nuvem sustenta um pequeno blog ao lado de um cofrinho, representando hospedagem serverless econômica.

## Revisão

Inspecione cada imagem em tamanho real antes de integrar, depois em miniatura no site. Conte membros, cabos e encaixes. Confira silhueta reconhecível, contornos, sobreposições, metáfora, paleta, texto acidental, recorte em desktop/celular e leitura nos temas claro/escuro. Atualize a skill local `guigoes` quando o usuário mudar a direção de arte. Não altere datas de publicação para trocar uma capa.

A primeira geração acrescentou um pé à haste da chave e uma perna sob a lente da lupa. Edições localizadas removeram essas peças, preservando as composições. Nas próximas gerações, substitua a instrução genérica de membros curtos por estrutura explícita: chave com uma cabeça e uma haste dentada, lupa com uma lente e um cabo, ambas sem pernas. Braços e pernas são opcionais e só entram quando necessários à cena.

Revisão das quatro capas: chave com uma haste contínua após correção; lupa com um único cabo diagonal após correção; terminal com dois braços e dois pés intencionais; nuvem com dois braços e dois pés, cofrinho com patas coerentes com a vista lateral. Não foram encontrados outros membros extras nessa revisão. Reinspecione depois de qualquer edição.


## Prompts das correções

### lazy-dvc-github-ssh-dvc-storage

```text
Make one precise correction to this existing Guigoes cartoon blog cover. The central cyan key has an erroneous extra leg/foot branching diagonally down-right from its toothed lower shaft. Remove ONLY that extra diagonal leg and rounded foot. Reconstruct the underlying cream background/shadow. Keep the central key as one circular bow with exactly one continuous toothed shaft ending in a simple rounded tip; the teeth stay short rectangular key teeth, never feet. Its two deliberate arms remain unchanged. Do not add legs or feet to the central key. Preserve all other characters, their poses, the smaller keys held in the hands, palette, faces, outlines, composition, aspect ratio, paper texture and every other detail. No new objects or text.
```

### html_meta_tags_social_media_ctr

```text
Make one precise correction to this existing Guigoes cartoon blog cover. The magnifying glass on the right has an erroneous extra vertical leg/foot descending directly under its circular lens. Remove ONLY that vertical leg and rounded foot, reconstruct cream background and shadow naturally. Magnifying glass must have exactly one round lens/rim and ONE single handle angled toward the lower-right. The genuine diagonal handle remains attached to the rim. No legs or feet on the magnifying glass. Preserve the page character, both tag characters, the face on the lens, their intentional arms and feet, the palette, layout, aspect ratio, outlines, paper grain and all unrelated details. No additional objects or text.
```
