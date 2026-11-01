# cloudflared article preparation

Prepared on October 31, 2026. The source remains `draft = true`; no publication, push, domain change, tunnel creation or Access configuration was performed.

## Research

The article links primary Cloudflare documentation beside factual claims. All 27 distinct linked URLs returned HTTP 200 during preparation. Key distinctions checked against the owning sources:

- Public hostnames and private network routes have different client requirements and CDN-content restrictions.
- gRPC is supported through private subnet routing, not public Tunnel hostnames.
- Registrar, authoritative DNS and Tunnel routing are separate functions. Full setup works without transferring registration; partial CNAME setup requires Business or Enterprise.
- The current 413 documentation lists 100 MB for Free/Pro, 200 MB for Business and up to 5 GB self-service for Enterprise.
- The Quick Tunnel limit is 200 in-flight requests, not 200 users or a limit for all tunnels; SSE is unsupported in Quick Tunnels.
- Workload examples are editorial applications of the linked terms, not named-product bans or measured deployment results.

The Free-plan emphasis was added at the author's request. Cloudflare's current plan comparison confirms DNS, Universal SSL and DDoS protection in the website Free tier; its Zero Trust offering includes a 50-user Free tier. The Tunnels for Platforms documentation confirms no-cost tunnel capacity in every account. The article distinguishes Zero Trust seats from public website visitors and identifies small-scale suitability as an architectural assessment, not a benchmark or SLA. Paid exceptions, domain/origin costs and existing content restrictions remain explicit.

The author requested a lighter treatment of security and a clearer choice between publication and private access. The introduction now recommends private access for an individual or small group when installing a client is acceptable, using Cloudflare One private routes or Tailscale. This is a scoped architectural recommendation based on avoiding a published application endpoint, not a claim that any VPN configuration is universally secure. Repeated warnings were removed; Access is explained as optional browser-based login. Protocol, billing and workload restrictions remain covered. Sources are Cloudflare's private-network documentation and Tailscale's own architecture explanation.

## Artwork

Integrated image generation produced the Cloudflare-inspired sheet at `docs/media/characters/v1/cloudflare.png`, documented in [the character record](media/characters/cloudflare.md). It is an editorial interpretation, not an official mascot. The cover is `content/blogs/cloudflared-tunnel-servicos-internet/assets/cover-cartoon.png`.

The initial generic cyan-cloud cover remains in the image tool's output directory outside published sources. The final cover uses the new orange character. Both selected characters retain two arms and two feet; the continuous tube is a separate object. Full-size images and rendered thumbnails were visually inspected.

Initial composition prompt, integrated image tool, with the generic cloud and mini-server sheets attached:

```text
Use case: illustration-story. Create ONE finished Guigoes editorial blog cover, landscape 16:9, 1536x864 or equivalent. Image 1 is the cloud character identity reference, image 2 is the mini-server character identity reference, NOT layouts to copy. Exactly ONE instance of each character. A small cyan mini-server is inside a simple open-cutaway cream house on the left; a friendly cyan cloud stands outside on the right. One smooth teal enclosed tube connects a socket on the house side to the cloud, an editorial metaphor for a secure outbound tunnel without opening the home's front door. The house front door stays closed, with a small plain keyhole and no face. Tube is one continuous path, not a limb. Mini-server: low cyan chassis, teal front panel with small face, two yellow LEDs, three side vents, exactly two arms and two feet. Cloud: high central dome cyan silhouette, soft eyebrows, exactly two arms and two feet. Both characters' arms rest at their own sides, no hands holding anything. No extra characters, no anatomical fusions, no dangling cables. Original cute rounded 2D technology-object cartoon, tiny eyes, small smiles, dark plum #292331 outlines, flat cyan #67d3ed and teal #007d9c fills, warm cream #f8f3e8 background, restrained lilac #b9a2ef and coral #f593a3 accents, simple shadows and fine paper grain. Generous negative space and crop-safe margins, main scene inside center 75 percent. No text, letters, logos, watermark, official brand mascots, glossy 3D or photorealism.
```

Final edit prompt, integrated image tool, with the initial cover and new Cloudflare sheet attached:

```text
Use case: precise-object-edit. Image 1 is the edit target, an existing Guigoes cover. Image 2 is the NEW Cloudflare editorial character identity sheet. Replace ONLY the cyan cloud character on the right with ONE instance of the orange Cloudflare-inspired character from image 2. Preserve the cloud's standing pose, approximate scale, placement, simple friendly face, exactly two arms relaxed at sides and two feet. Use orange main body, golden right puff and the two small cream airflow markings shown in its reference; do not reproduce three poses. Keep the house, cyan mini-server, closed front door, continuous teal tunnel tube, background, shadows, aspect ratio, texture and all other composition details unchanged. The tube should attach cleanly to one socket at the orange cloud's left side. No extra limbs, clouds, detached hands, text, logos, wordmarks or watermarks. Maintain original 2D Guigoes style, dark plum outlines and cream background. The server remains cyan. The cloud is an editorial interpretation, not an official mascot.
```

## Validation

- `make check` passed, including publication checks, uncached Go tests and `go vet`.
- A separate preview copy changes only its own draft flag, then builds with the installed Gazeta/Margo dependencies. It generated 55 HTML pages and 10 BlogPosting blocks, including the new article. The repository's production build keeps the draft and its assets excluded.
- Browser inspection covered the article cover in desktop and 390px layouts, light/dark themes, the expanded Mermaid diagram in light/dark desktop and dark mobile layouts, and home/archive thumbnails. No real-device gesture testing or live Cloudflare deployment was performed.
- The example commands are illustrative and were not used to expose a real application.

## Authenticated dashboard screenshots

Five viewport screenshots were captured from the authenticated Codex in-app browser on October 31, 2026 and added beside the relevant explanations:

- `cloudflare-tunnels.jpg`: the existing tunnel list, with connector status and the creation entry point.
- `cloudflare-route-types.jpg`: the unsaved route-type selector, distinguishing published applications from private routes.
- `cloudflare-published-application.jpg`: an unsaved example with subdomain `painel` and `http://127.0.0.1:8080`. The real domain and full hostname are blurred; the article uses the reserved example domain.
- `cloudflare-nameservers.jpg`: the zone's current `DNS > Settings` screen. The caption tells readers to use their own assigned nameservers.
- `cloudflare-access-policy.jpg`: an unsaved Allow policy using the fictitious address `voce@example.com`. The application and policy forms were discarded afterward.

The author requested crops and blur to hide account identity. Deterministic ImageMagick edits remove account headers, navigation, and surrounding page content. Blur hides tunnel names, the route's real domain/full hostname, and assigned nameservers. The fictitious `voce@example.com` and local example service remain readable. No generative edits were used. No token-installation screen was captured. The native Margo image viewer remains enabled.

Original captures remain outside the repository in the task's local screenshot evidence directory. That directory contains `redact-screenshots.py` and `redaction-manifest.json`, with crop/blur coordinates, output dimensions and SHA-256 hashes. Only the redacted derivatives belong in article sources and preview assets.

No route, tunnel, policy or application was saved. DNS and security settings were not changed. The post remains a local draft.

Validation after initial integration: `make check` passed. The isolated preview built 55 pages and 10 BlogPosting blocks with all five screenshot figures. Desktop expansion, Actual size/Fit controls, and the mobile route screenshot were inspected through the native viewer. Real-device touch gestures were not tested. After redaction, all five edited files were opened at full size and checked for readable account identifiers and preserved tutorial fields.

## Independent content review

Three isolated reviewers inspected the article and its six assets directly from the Git worktree, without producer transcript, notes or receipts. Lenses were technical accuracy, plans/terms/workloads, and editorial clarity/tutorial completeness.

Round 1 identified three omissions: Access token validation and negative-path checks, a Service Auth policy for service tokens, and an explicit DNSSEC/DS migration action. The article now includes those instructions, keeps the GUI screenshot caption accurate about its collapsed settings, and links the owning documentation.

All three reviewers accepted round 2 with no remaining actionable findings or deferred hardening. Each verified an unchanged scoped snapshot before and after review:

- Worktree: `/Users/guilhermecastro/repos/guigoes-v2`.
- Scope: `content/blogs/cloudflared-tunnel-servicos-internet`.
- Base/HEAD: `0662de3e4302d76d73ee78d279c0541ef18313d5`.
- Accepted snapshot: `a9e3aa0bb66c633d45b859e6be125bdbf769d8a1d9cd20bdf5cfa778cff0e9d6`.

The revised isolated preview built and validated 55 HTML pages and 10 BlogPosting blocks. This is acceptance of content and examples, not evidence of executing a new tunnel, DNS migration or Access deployment. The source remains a draft; nothing was published.

After that accepted snapshot, the author requested a short mention of dynamic DNS as another publication option. The article now explains automatic A/AAAA updates, Cloudflare API/ddclient support, the need for reachable inbound connectivity, and why DDNS does not bypass ISP CGNAT. The introduction and closing comparison were updated. Cloudflare's dynamic-IP documentation was checked directly; this later addition is not part of the independent-review snapshot above.

## Series introduction rewrite

The author found the long article too complex and requested a lighter introduction to free Cloudflare services for homelab users, individuals, small businesses and groups. The current draft supersedes the reviewed tutorial. It retains the existing slug and cover, adds the native `Cloudflare para homelabs` series metadata as part 1, and introduces DNS/CDN/HTTPS, DDNS, Tunnel, Access/private networking, Gateway, Workers/Static Assets, R2, Turnstile and Email Routing.

Installation commands, dashboard procedures, protocol matrices, DNS migration instructions, service-token details and the Mermaid flow were removed. One redacted route-selector screenshot remains to illustrate public versus private access. The other redacted images remain available for later parts without appearing in this introduction. The pre-rewrite Markdown is preserved outside published sources at `/Users/guilhermecastro/.codex/visualizations/2026/10/31/01a193fb-a9e0-7b02-ba9c-711fbf446a3f/cloudflared-before-series-rewrite.md`.

Planned follow-up subjects are DNS/domains/DDNS, Tunnel/public versus private routing, Access/Gateway, and Workers/R2/integrations. No placeholder posts or tutorial follow-ups were created. The previous independent-review acceptance does not apply to this rewritten introduction.

Official Cloudflare documentation and product pages were checked again for free availability, Workers execution/request quotas, R2 Standard storage allowances and overage billing, Turnstile availability, Email Routing versus sending, public Tunnel upload/media constraints, and private networking. The Zero Trust pricing page currently omits its plan table in the text response; the 50-user allowance was cross-checked against Cloudflare's current Mesh announcement and its published Free-plan explanations.

Validation of this rewrite: 1,369 body words, down from the previous roughly 3,600-word draft; no fenced examples or installation instructions remain. The isolated preview built and validated 53 HTML pages and 10 BlogPosting blocks, including the native series index. A normal source build validated 51 pages and 9 BlogPosting blocks and excluded this draft from its article URL and RSS. The browser showed the updated title, series block, introduction, retained route figure and surrounding prose at its default desktop size. No new independent review, real-device check or Cloudflare deployment was performed for this rewrite.

## Independent review of the series introduction

At the author's request, two new isolated reviewer tasks assessed editorial clarity/lightness and Free-tier accuracy. Both used the required reviewer profile, `gpt-5.6-sol` with high reasoning, without inherited producer history. An initial editorial task reported that its broad search had exposed this notes file; that task's review was discarded and replaced by a fresh task with source searches explicitly constrained to the article scope. The replacement did not inherit that task's context.

The valid editorial review identified two material ambiguities: small-business suitability needed an availability/support qualification, and the 50-user limit needed to describe shared active Zero Trust seats rather than an Access-only allowance. The article now links the seat-counting FAQ, explains application authentication or Cloudflare One Client enrollment, excludes ordinary public visitors, and states that the Free site plan has no uptime SLA. A repeated phrase about caching was shortened.

Both valid reviewers accepted round 2 with no unresolved blockers. Their independent read-only passes bracketed the same scoped Git snapshot:

- Target: `/Users/guilhermecastro/repos/guigoes-v2`.
- Scope: `content/blogs/cloudflared-tunnel-servicos-internet`.
- Base/HEAD: `0662de3e4302d76d73ee78d279c0541ef18313d5`.
- Accepted snapshot: `29c944ce73e82ab8ce5204f19d36c27b3a78ba8af06eb1c2fb0bfd8ba2d1e52e`.
- Boundary receipts: transcript not read, no producer polling, Git-native inspection, stable snapshots.

Non-blocking editorial observations remain: the title uses "Free" as shorthand for multiple plans/allowances, which the body distinguishes, and four redacted dashboard screenshots are retained for future parts without article references. Those assets should be reassessed when preparing publication. No reviewer found visible sensitive data in them.

The article explicitly enumerates planned parts 2 through 5: domains/DNS/HTTPS/DDNS; Tunnel behavior/protocols/limits; Access/private networking/Gateway; and Workers/static assets/R2/integrations. No placeholder articles were created. The final isolated preview validated 53 HTML pages and 10 BlogPosting blocks. The roadmap rendered as an ordered list starting at 2. This remains a local draft; no deployment or live Cloudflare configuration was performed.
