# BOINC article review and verification

Reviewed on 2026-10-30, America/Sao_Paulo. Article/tutorial acceptance identity:
`1143dc3dbf2720908086089f711908e7d45bfdfd`.

## Independent review

Five fresh reviewers received Git targets without the producer conversation,
session memory, previous findings or author validation receipts. Two reviewed
reader clarity independently; the others checked scientific facts, technical
contracts and tutorial execution. The technical reviewer also audited the
actual thermal-supervisor source and its startup failure paths.

Both clarity reviewers accepted the revised article. The scientific reviewer
reported no remaining factual blocker. The tutorial reviewer accepted the
changed executable paths. The thermal implementation was accepted after its
startup ordering and a stale README statement were corrected.

Corrections from review:

- Introduce prerequisites, including cluster-scoped permissions, Linux AMD64,
  Kubernetes support floor, storage, network and optional Metrics Server.
- Distinguish BOINC RPC suspension from stopping/recreating a pod. The example
  now keeps suspended CPU applications in memory.
- Put the tutorial's lack of a thermal supervisor at the start of that section.
- Define checkpoint and the calculate/upload/report/validate lifecycle.
- Make the description promise a CPU tutorial and a separate GPU-pilot account.
- State OpenZika's negative antiviral endpoint, alongside its cytoprotective hits.
- Complete the OpenZika explanation before the requested contextual UFG image.
- Describe Science United's selection options without claiming a fresh audit of
  the author's server-side preferences.

## Claim ledger

| Claim group | Evidence and result |
| --- | --- |
| BOINC and Science United roles, selection and compatibility | Supported by [BOINC](https://boinc.berkeley.edu/), [Science United introduction](https://scienceunited.org/intro.php), [about page](https://scienceunited.org/su_about.php), and [BOINC project guidance](https://github.com/BOINC/boinc/wiki/Choosing-and-joining-projects). Seven top-level categories are also defined in BOINC's `html/inc/keywords.inc`. |
| Einstein@Home discovery | The [2018 primary paper](https://pmc.ncbi.nlm.nih.gov/articles/PMC5829974/) supports two isolated millisecond pulsars, approximately 348 and 213 rotations per second, one undetected in radio in those observations, and volunteer-computer participation. |
| OpenZika | The [primary paper](https://pubmed.ncbi.nlm.nih.gov/36239304/) and [WCG update](https://worldcommunitygrid.org/about_us/article.s?articleId=779) support the UFG participation, nearly 404 million results, 61 prioritized compounds and experimental assays. Eight compounds were cytoprotective; none showed direct antiviral activity. No clinical benefit is claimed. |
| Climate attribution | The [primary study record](https://eprints.lancs.ac.uk/id/eprint/78221/) and [Oxford experiment description](https://eng.ox.ac.uk/climateprediction/conduct-a-study/projects/2014-uk-winter-floods) support ensemble/counterfactual modeling and uncertain changes in rainfall/flood risk. The article makes no single-cause or local-weather forecast claim. |
| Images | NASA pulsar is an illustration; the [Fermi map](https://svs.gsfc.nasa.gov/11342/) is five years of observed LAT data above 1 GeV. UFG photograph is credited to Laczkowski under CC BY-SA 4.0. Source bytes, captions and attribution were independently checked; see `boinc-media-sources.md`. Cover generation is declared in `boinc-cover-prompt.md`; the external local character-reference library was outside the immutable review target. |
| Energy example | `0.020 kW × 8 h × 30 = 4.8 kWh`; explicitly hypothetical, not a meter reading or cost estimate. |
| CPU/memory budgets and scheduling | Supported by [Kubernetes resource management](https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/), [non-preempting priority](https://kubernetes.io/docs/concepts/scheduling-eviction/pod-priority-preemption/#non-preempting-priorityclass), and live pod/cgroup inspection. CPU request is not a utilization cap; memory may OOM. |
| Init container, quota and PVC retention | Supported by [init resource accounting](https://kubernetes.io/docs/concepts/workloads/pods/init-containers/#resource-sharing-within-containers), [quota](https://kubernetes.io/docs/concepts/policy/resource-quotas/) and [StatefulSet retention](https://kubernetes.io/docs/concepts/workloads/controllers/statefulset/#persistentvolumeclaim-retention). Tested through real Kubernetes controllers and negative admission checks. |
| BOINC preferences and RPC | Supported by the pinned 8.0.4 source, [preferences](https://github.com/BOINC/boinc/wiki/Preferences), [boinccmd](https://github.com/BOINC/boinc/wiki/Boinccmd-tool) and execution of the exact image. Host-RAM percentage arithmetic, local-only RPC, persistent account-manager configuration and suspension modes were checked. |
| GPU allocation and sharing | Supported by [Kubernetes GPU scheduling](https://kubernetes.io/docs/tasks/manage-gpus/scheduling-gpus/) and [Intel plugin](https://github.com/intel/intel-device-plugins-for-kubernetes/tree/v0.37.1/cmd/gpu_plugin). Live plugin uses version 0.37.1 with shared-device slots; one allocation is not a utilization, power or exclusive-memory fraction. |
| Thermal policy | Existing source and fault-injection tests confirm max-valid `coretemp`, 85/80 C hysteresis, conservative initial state and termination after a failed control attempt. Live logs confirm natural pause/resume cycles. The policy remains cooperative, not a physical temperature ceiling. |
| Scientific progress and validation | Existing clients remained associated with Science United and had actual project tasks. After the GPU rollout, both a GPU task and a CPU task advanced. This is execution evidence, not a newly validated scientific result. |

## Tutorial execution

The exact pinned AMD64 image was executed. An isolated Kubernetes v1.35.0
cluster exercised:

- Namespace, PriorityClass, LimitRange, ResourceQuota, ConfigMap, Service,
  StatefulSet and PVC creation under restricted Pod Security admission.
- Successful init copy, readiness RPC, UID/GID 1000, no capabilities, no service
  account token, RuntimeDefault seccomp, read-only root and writable data/tmp.
- Effective one-CPU/1-GiB cgroups and exact declared quota consumption.
- Local loopback RPC, with no working pod-IP RPC listener.
- BOINC config loading, one simulated CPU, no GPU, RAM-percentage interpretation
  and final `leave_apps_in_memory=1` preference.
- Exact RPC `never`/`auto` suspend/resume commands and optional Metrics API failure.
- Graceful SIGTERM/exit, scale down/up, retained PVC/state marker, StatefulSet
  deletion/recreation with the retained claim, and deliberate ConfigMap rollout.
- Rejection of unsafe pods, oversized containers/claims and extra replicas.

The isolated cluster and dedicated kubeconfig were removed after verification.
The reviewer-owned detached Git worktree and private test receipt were retained.

## Live installation and repairs

Live Talos nodes run Kubernetes v1.35.2. Actual clients match the article's
installation profile: two separate clients, 25m/2-GiB requests and two-CPU/2-GiB
limits per client, independent 8-GiB volumes, one GPU allocation on the GPU
client, CPU counts of one/two, low non-preempting priority, resident applications,
and a 1.75-GiB internal task-memory budget per client.

Two real discrepancies were repaired:

1. The BOINC namespace's committed restricted Pod Security labels had not been
   reconciled. Refreshing its owning Argo application applied the Git state
   without restarting either client. Server dry-run then rejected privileged
   pods, oversized CPU/storage requests and a third client.
2. The GPU supervisor allowed 60 seconds for initial RPC retries while BOINC
   could start applications after 30 seconds. Its application delay is now
   90 seconds. A regression first failed on the old behavior at 60/64 seconds,
   then passed after repair. CPU remains 30/20 seconds; GPU is 90/60 seconds,
   with four-second RPC timeouts. CI runs both failure modes for both guards.

The deployed GitOps revision is `61edaabb2c965e7ff27cf212a440facbbc807ed1`.
The GPU-only ConfigMap rollout reached Argo Synced/Healthy with zero container
restarts in the new pod. The CPU pod identity stayed unchanged. Image, resource
budgets and all PVC/PV identities were preserved. One CPU task in the GPU client
resumed from its previous checkpoint, about 598 CPU-seconds behind the pre-rollout
snapshot; uncheckpointed work was recomputed, not claimed preserved.
The article's two PromQL queries returned data from the live monitoring system.

## Coverage limits

- Native AMD64 kind bootstrap was unavailable under the local emulation/seccomp
  combination. Controller/admission/storage tests used a native ARM64 kind node
  with only its architecture selector adapted; the pinned AMD64 client ran via
  binfmt. Real AMD64 deployment/image/cgroups/supervisors were separately
  inspected in the live v1.35.2 cluster. This is not a clean-room AMD64 kernel test.
- Kubernetes v1.32 itself was not provisioned. It is the stated support floor
  for the stable explicit PVC-retention policy, not the oldest possible API.
- No new real Science United account was entered in the test cluster. The join
  command contract was checked, while existing live associations/project tasks
  were verified separately. New-account authentication and future task
  availability remain dependent on the reader's account and scientific servers.
- No new task was completed merely for this test, and no new scientific discovery
  or validation is inferred from task execution. Temperature sampling does not
  prove service-latency neutrality or a hard maximum temperature.

## Public-image scientific GPU follow-up

The article now incorporates the physical-N100 observation recorded in
[the public verification receipt](https://github.com/guilycst/boinc-k8s-intel-n100/blob/d02e0e3f35f747bc24141954409da621a97e85bc/VALIDATION.md#scientific-gpu-task-follow-up).
The scientific application's own i915 engine counter advanced by 16.159 seconds
over 113.774 seconds, with increasing task fraction and a later saved checkpoint.
This is stronger execution evidence than OpenCL enumeration or the separate
standalone-kernel smoke test. The task remained in progress.

A separate ten-minute thermal-policy window authorized GPU work for 80.210
seconds. The prose distinguishes permission to run from measured GPU execution
and explains the effect on elapsed completion time. The follow-up kept both pod
identities, volume bindings, resource limits and thermal protection unchanged.
The new paragraphs do not claim task completion or server-side validation.
Earlier independent reviews above retain their original acceptance identity.
