# BOINC article validation

Article revision: 2026-10-30, America/Sao_Paulo. Infrastructure observations below
remain the 2026-10-29 snapshot; this revision did not change the cluster.

The independent follow-up review, tutorial reproduction and live repairs are
recorded in [boinc-review.md](boinc-review.md). Its live observations supersede
the older snapshot below.

## Infrastructure

- Home-lab PR 733 added stable v1 LimitRange and ResourceQuota admission budgets.
  Server dry-run confirmed rejection of a third pod, a 3-CPU container and a
  9-GiB claim. Neither existing pod restarted for the admission change.
- Home-lab PR 737 added a CPU-only Intel coretemp supervisor after a live
  worker-02 reading of 96 C. CPU work pauses at 85 C and resumes at 80 C.
  Natural pauses/resumes retained task PIDs 32 and 33 while their CPU counters
  advanced. The new CPU pod had zero restarts during observation.
- The CPU rollout retained all three queued task names, the WCG host identity
  and all PVC/PV identities. MCM resumed from its 4214.192-second checkpoint,
  recomputing about 473 CPU seconds. An ODLK task without a checkpoint recomputed
  about 292 CPU seconds. This was not a zero-recomputation rollout.
- GPU pod UID and container identity stayed unchanged throughout both changes.
  No GPU residency claim is inferred from the CPU supervisor test.
- Argo reached Synced/Healthy at 497c046c50aa3ed53e3f6d934e322bd609742a5b.
  Both BOINC volumes had UpToDate diskful replicas on both workers, Established
  connections, zero out-of-sync bytes and quorum in the closing storage check.
- The repository live inspector reported no active pod, Argo, workload-rollout
  or storage exceptions at the closing snapshot. Retained failed Jobs were
  distinguished from active failures. This is not a service-latency benchmark.

Runtime baseline: Kubernetes 1.35.2, Talos 1.13.9, containerd 2.2.7 on the BOINC
nodes, BOINC 8.0.4 and Intel device plugin 0.37.1. No platform upgrade, PVC
recreation, identity clone or increase in resource limits was performed.

## Article and examples

- Primary sources cover Einstein@Home pulsars, OpenZika screening/experimental
  work and climateprediction.net flood-risk attribution. Historical results are
  not attributed to the author's installation. No personal discovery or newly
  validated scientific result is claimed.
- The revised CPU example uses native Kubernetes resources, a cp init container
  and a direct BOINC process. Sealed-account tooling and embedded startup Python
  were removed. The tutorial associates an account interactively with boinccmd;
  no real credentials were entered during validation.
- Kubernetes client dry-run passed for all seven resource documents. XML parsed
  with the standard library. An isolated, network-disabled Docker run of the
  exact pinned AMD64 image tested the init copy, UID 1000, read-only root,
  dropped capabilities, one-CPU/1-GiB limits and the readiness RPC command.
  BOINC reported one CPU, no GPU and loaded the memory preferences, reporting
  0.55 GB for the test VM's 11.73-GB host memory at 4.6875 percent. A graceful
  stop retained client_state.xml and gui_rpc_auth.cfg on the temporary data
  volume. No projects or live scientific work were attached to that test.
- The example requires node, StorageClass and RAM-percentage adjustment.
  Pod scheduling, the chosen CSI driver's permissions and account attachment
  were not exercised on a second cluster.
- NASA pulsar illustration and Fermi observational map were inspected at source
  resolution. Captions, alt text, provenance, credits and hashes are recorded
  in docs/boinc-media-sources.md. The UFG photograph has explicit photographer,
  source and CC BY-SA 4.0 attribution. Historical GA dates were removed on
  review, while links to the resource contracts remain.
- Editorial review introduces Einstein@Home, explains millisecond and isolated
  pulsars, links climateprediction.net, and removes internal node names. The
  thermal section now explains the actual sensor selection, independent CPU/GPU
  control scopes, hysteresis, suspended processes and failure behavior. Its
  thresholds were checked against the existing supervisor source. The later
  independent review additionally checked live hardware and deployed the
  startup-ordering repair documented in boinc-review.md.
- make check passed, including final HTML, local links/fragments, legacy URLs,
  feeds, Go tests and vet. The page has one BlogPosting, the expected canonical,
  pt-BR, author profile and timestamp. JSON-LD, OG and Twitter image URLs match;
  actual PNG dimensions are 1672 by 941.
- Browser inspection covered 1440-pixel desktop and 390-pixel mobile views in
  light/dark themes. The cover is readable as a thumbnail. Mermaid was checked
  inline and expanded, without document-level horizontal overflow. No physical
  touchscreen gestures or third-party social crawler caches were tested.
- Margo PR 100 fixed the aggregate site asset budget through native configuration.
  PR 102 made native Mermaid captions honor accTitle. Both were validated in the
  consumer and passed CI before integration. Guigoes uses a pinned upstream
  commit, not a local replace or generated-HTML workaround.

## Remaining publication gate

The requested real Science United certificate has not been exported. Browser
origin permission denied access even though the user opened its tab. Do not
invent a certificate or its numbers. Attach the real export, inspect personal
data, add a faithful public asset and caption, then repeat the final build and
browser check before publishing. The article PR remains a draft until then.

## Diagnostic handling

A LINSTOR diagnostic initially included replication secrets in private tool
output. No values were copied into Git, article assets or this receipt. The raw
temporary JSON was replaced by an allowlisted summary. Any credential rotation
must be handled as a separate replication-preserving operation; it was not
performed during this publication task.
