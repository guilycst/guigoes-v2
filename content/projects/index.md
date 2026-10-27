+++
title = "Projetos"
description = "Projetos de Guilherme de Castro e da organização Araihu: interfaces em Go, publicação, infraestrutura e ferramentas para desenvolvedores."
language = "pt-BR"
+++
Uma seleção dos projetos que construo e mantenho. O fio condutor é Go, software livre e ferramentas que simplificam o trabalho de quem desenvolve e opera sistemas.

## Na Araihu

Sou mantenedor da [Araihu](https://github.com/araihu), organização onde desenvolvemos bibliotecas e aplicações em Go. Estes são alguns dos destaques da organização.

<ul class="guigoes-projects" aria-label="Projetos da Araihu">
<li>
<div class="guigoes-project-heading">
<img class="guigoes-project-icon" src="./assets/goshtoso.svg" alt="" width="64" height="64" loading="lazy">
<div>
<h3><a href="https://github.com/araihu/goshtoso">Goshtoso</a></h3>
<a class="guigoes-project-version" href="https://github.com/araihu/goshtoso/releases/tag/v0.4.5" aria-label="Release v0.4.5 · beta de Goshtoso">v0.4.5 · beta</a>
</div>
</div>
<p class="guigoes-project-kind">Interfaces web · Go, templ, HTMX</p>
<p>Biblioteca de componentes para aplicações que renderizam HTML no servidor. Integra templ, HTMX e Alpine.js, com estilos e assets distribuídos junto da biblioteca. A API ainda está em evolução.</p>
<p class="guigoes-project-links"><a href="https://goshtoso.araihu.com/" aria-label="Site oficial de Goshtoso">Site oficial</a> · <a href="https://github.com/araihu/goshtoso" aria-label="Repositório de Goshtoso no GitHub">GitHub</a></p>
</li>
<li>
<div class="guigoes-project-heading">
<img class="guigoes-project-icon" src="./assets/margo.svg" alt="" width="64" height="64" loading="lazy">
<div>
<h3><a href="https://github.com/araihu/margo">Margo</a></h3>
<a class="guigoes-project-version" href="https://github.com/araihu/margo/releases/tag/v0.0.25" aria-label="Release v0.0.25 · alpha de Margo">v0.0.25 · alpha</a>
</div>
</div>
<p class="guigoes-project-kind">Publicação · Markdown, HTML, PDF</p>
<p>Transforma Markdown em páginas HTML, sites estáticos, PDFs e apresentações. Pode ser usado pela linha de comando ou como biblioteca Go. É o motor de publicação deste blog.</p>
<p class="guigoes-project-links"><a href="https://margo.araihu.com/" aria-label="Site oficial de Margo">Site oficial</a> · <a href="https://github.com/araihu/margo" aria-label="Repositório de Margo no GitHub">GitHub</a></p>
</li>
<li>
<div class="guigoes-project-heading">
<img class="guigoes-project-icon" src="./assets/manja.svg" alt="" width="64" height="64" loading="lazy">
<div>
<h3><a href="https://github.com/araihu/manja">Manja</a></h3>
<a class="guigoes-project-version" href="https://github.com/araihu/manja/releases/tag/v0.1.17" aria-label="Release v0.1.17 · beta de Manja">v0.1.17 · beta</a>
</div>
</div>
<p class="guigoes-project-kind">Documentação de APIs · OpenAPI, Goshtoso</p>
<p>Renderiza e publica documentação OpenAPI com busca por operações e schemas. Exporta um site estático que funciona sem servidor Manja e também oferece execução como servidor de documentação.</p>
<p class="guigoes-project-links"><a href="https://manja.araihu.com/" aria-label="Site oficial de Manja">Site oficial</a> · <a href="https://github.com/araihu/manja" aria-label="Repositório de Manja no GitHub">GitHub</a></p>
</li>
<li>
<div class="guigoes-project-heading">
<img class="guigoes-project-icon" src="./assets/balemoh.svg" alt="" width="64" height="64" loading="lazy">
<div>
<h3><a href="https://github.com/araihu/balemoh">Balemoh</a></h3>
<a class="guigoes-project-version" href="https://github.com/araihu/balemoh/releases/tag/v0.0.8" aria-label="Release v0.0.8 · alpha de Balemoh">v0.0.8 · alpha</a>
</div>
</div>
<p class="guigoes-project-kind">Infraestrutura · Kubernetes, Docker, Podman</p>
<p>Descobre serviços e organiza uma página inicial com os que interessam. Combina catálogo persistente, revisão das descobertas e uma interface em Go, com plugins para integrar outras ferramentas.</p>
<p class="guigoes-project-links"><a href="https://balemoh.araihu.com/" aria-label="Site oficial de Balemoh">Site oficial</a> · <a href="https://github.com/araihu/balemoh" aria-label="Repositório de Balemoh no GitHub">GitHub</a></p>
</li>
<li>
<div class="guigoes-project-heading">
<div>
<h3><a href="https://github.com/araihu/talaria-mem">Talaria-Mem</a></h3>
<span class="guigoes-project-version">Sem release</span>
</div>
</div>
<p class="guigoes-project-kind">Memória para agentes · Go, SQLite, MCP</p>
<p>Serviço de memória local para o Codex, separado por workspace. Guarda contexto persistente e oferece consulta por CLI e MCP. Está em desenvolvimento; a implementação atual ainda não é uma versão lançada.</p>
</li>
</ul>

## Na minha conta

Ferramentas e experimentos publicados em [guilycst](https://github.com/guilycst).

<ul class="guigoes-projects" aria-label="Projetos de Guilherme de Castro">
<li>
<div class="guigoes-project-heading">
<div>
<h3><a href="https://github.com/guilycst/lazy-dvc">lazy-dvc</a></h3>
<a class="guigoes-project-version" href="https://github.com/guilycst/lazy-dvc/releases/tag/2026.03.16-cd6269d" aria-label="Release 2026.03.16-cd6269d de lazy-dvc">2026.03.16-cd6269d</a>
</div>
</div>
<p class="guigoes-project-kind">Dados e identidade · DVC, SSH, GitHub</p>
<p>Usa as chaves SSH e a participação em organizações e times do GitHub para autenticar o acesso a um remote DVC. Aproxima o acesso aos dados do fluxo que o time já usa para código.</p>
<p><a href="/blogs/lazy-dvc-github-ssh-dvc-storage/">Leia o artigo sobre lazy-dvc</a></p>
</li>
<li>
<div class="guigoes-project-heading">
<div>
<h3><a href="https://github.com/guilycst/gitmeout">gitmeout</a></h3>
<a class="guigoes-project-version" href="https://github.com/guilycst/gitmeout/releases/tag/v0.0.1-beta.1" aria-label="Release v0.0.1-beta.1 · alpha de gitmeout">v0.0.1-beta.1 · alpha</a>
</div>
</div>
<p class="guigoes-project-kind">Portabilidade · GitHub, Forgejo, Codeberg</p>
<p>CLI em Go para espelhar repositórios do GitHub em instâncias Forgejo ou Codeberg. Permite selecionar contas, organizações ou repositórios e trabalhar com espelhos por push ou pull.</p>
</li>
</ul>

Os repositórios têm código, instruções de uso e detalhes sobre o estágio de cada projeto. Veja outros projetos nos perfis de [guilycst](https://github.com/guilycst?tab=repositories) e da [Araihu](https://github.com/orgs/araihu/repositories).

<p>Versões consultadas em 27/10/2026. Neste blog, 0.0.x indica alpha; demais versões 0.x.x indicam beta. Releases por data não recebem essa classificação.</p>
