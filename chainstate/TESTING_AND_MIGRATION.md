# ChainState Indexer: Testing & Migration Plan

This document outlines the plan for validating the chainstate indexer and safely
migrating consumers off the operator-state subgraph. See the
[README](./README.md) for what the indexer does and its current status.

> **Status**: draft plan, AI-drafted and pending human review (per CLAUDE.md §4,
> humans own docs — treat this as a skeleton to edit, not a final plan).

## 1. Current Consumer Landscape

Services that consume indexed operator state today, all via
`thegraph.MakeIndexedChainState` (the `core.IndexedChainState` interface):

| Consumer | Wiring | Role |
|----------|--------|------|
| `disperser/controller` | `disperser/cmd/controller/main.go` | Attestation path (highest stakes) |
| `operators/churner` | `operators/churner/cmd/main.go` | Churn approval decisions |
| `relay` | `relay/cmd/lib/relay.go` | Serving path |
| `disperser/dataapi` | `disperser/cmd/dataapi/main.go` | Dashboards / API |
| `disperser/batcher` (v1) | `disperser/cmd/batcher/main.go` | Legacy attestation path |
| `tools/quorumscan`, `tools/semverscan` | tool mains | Operator tooling |

**Additionally**: `disperser/dataapi` queries the subgraph *directly* via
GraphQL (`disperser/dataapi/subgraph_client.go`) for richer history queries.
This is a separate, larger migration than the interface swap and needs an
endpoint-by-endpoint mapping to the REST API.

## 2. Phase 0 — Merge Readiness (this PR)

1. **Human review of AI-drafted tests** (CLAUDE.md §3). This applies to the
   e2e tests (`inabox/tests/chainstate_indexer_test.go`,
   `inabox/tests/chainstate_subgraph_parity_test.go`) and to the unit tests
   added in commit `961952eb` (`chainstate/indexer_test.go`,
   `indexer_startblock_test.go`, `indexed_chain_state_test.go`,
   `api/server_test.go`, `store/memory_store_test.go`,
   `store/json_persister_test.go`, `core/attestation_test.go`). A human must
   validate that the assertions encode *intended* behavior. Highest-scrutiny
   items:
   - The parity test's "latest APK snapshot <= block" lookup semantics.
   - The `deregisteredAsOf` boundary (`<=`, mirroring the subgraph's
     `deregistrationBlockNumber_gt` filter) in
     `indexed_chain_state_test.go`.
   - The intra-transaction event emission order
     (`NewPubkeyRegistration` / `OperatorSocketUpdate` before
     `OperatorRegistered`) encoded in `TestFirstRegistrationIntraTxOrder`,
     derived from reading `RegistryCoordinator._registerOperator` rather
     than from a spec.
2. **Unit test status**: regression tests for the review-round fixes exist
   and pass under `-race` — chronological event replay (dereg→rereg,
   quorum remove→re-add), skeleton-record merging, multi-quorum ejection
   dedup, start-block inclusivity and catch-up batching, snapshot/restore
   round-trip, API status codes, and golden JSON field names. Not unit
   tested (deliberately): `collectEvents` filterer plumbing,
   `snapshotQuorumAPKs` historical eth_calls, and service/main wiring —
   these are covered by the inabox e2e tests.
3. **Run the inabox suite** (`cd inabox && make run-e2e-tests`, requires
   Docker). The parity test compares at every membership-change block and is
   the strongest end-to-end signal in this PR. It has not been executed since
   the event-ordering rewrite — this is the highest-value single action.
4. Standard gates: lint clean, CI green, human PR review.

## 3. Phase 1 — Extended Validation (preprod / testnet)

1. **Full historical backfill**: run against Holesky/preprod with
   `StartBlockNumber` = contract deployment block and an archive RPC.
   Exercises the catch-up loop and archive-node dependency at real scale
   (the inabox devnet only covers a few hundred blocks).
2. **Offline parity sweep**: script that queries both the subgraph and the
   indexer REST API for operator state, APKs, ejections, and socket history
   at a large sample of historical blocks (ideally every membership-change
   block) and diffs the results. Real history contains re-registrations,
   ejections, and churn that the devnet cannot produce.
3. **Soak test (1–2 weeks)**: steady head-tracking, memory growth, persist
   file size. Include crash-recovery drills: `kill -9` mid-save (the fsync'd
   `AtomicWrite` should make this safe — prove it), delete/corrupt the state
   file, restart and re-backfill.
4. **Failure-mode drills**: kill the primary RPC (MultiHomingClient
   failover), point at a non-archive node (should stall loudly, not corrupt
   data), induce RPC flapping.
5. **API load test** on the REST endpoints at realistic consumer QPS.

## 4. Phase 2 — Prerequisites Before Production Traffic

Known gaps (documented in the README's Future Enhancements) that become
blocking once other components depend on this service:

1. **Reorg handling** — currently absent. The indexer trusts the unfinalized
   head; a reorged-out event is permanent state corruption with no
   self-healing. Minimum viable fix: index only up to
   `head - confirmationDepth` (new config value, e.g. 64 blocks). The
   subgraph handles reorgs via graph-node, so this is a regression until
   fixed.
2. **Prometheus metrics** — at minimum `last_indexed_block`, indexing error
   counts, and API request metrics. The shadow phase (below) cannot run
   without lag/divergence alerting.
3. **HA story** — one process + one JSON file is a single point of failure.
   Simplest viable model: 2+ independent replicas (indexing is
   deterministic, so they converge) behind a load balancer, each with its
   own persistence path. Decide before consumers depend on it.

## 5. Phase 3 — Shadow Deployment in Production

1. Deploy alongside the subgraph; no consumers pointed at it.
2. Run a **continuous parity checker** comparing indexer vs. subgraph at
   `head - confirmationDepth`; alert on any divergence and on indexed-block
   lag.
3. Bake for an agreed window (2–4 weeks) with zero unexplained divergences
   before migrating anything.

## 6. Phase 4 — Consumer Migration (lowest blast radius first)

Both implementations satisfy `core.IndexedChainState`, so each swap is
per-service wiring in `cmd/main.go`. Add a config flag per service selecting
the backend so rollback is a config flip, not a deploy.

| Order | Consumer | Blast radius | Notes |
|-------|----------|--------------|-------|
| 1 | `tools/quorumscan`, `semverscan` | None (operator tooling) | Cheap first validation |
| 2 | `disperser/dataapi` (IndexedChainState half) | Dashboards, no consensus role | Direct-GraphQL queries are a separate work item |
| 3 | `relay` | Serving path; degraded, not corrupted, on failure | |
| 4 | `operators/churner` | Wrong state = wrong churn decisions | |
| 5 | `disperser/controller` (and v1 `batcher` if still live) | Attestation path — highest stakes | Migrate last, after everything else has baked |

For each consumer: canary one instance → compare outputs/metrics against a
subgraph-backed peer → full rollout → next consumer. One consumer per change
window.

## 7. Phase 5 — Decommission

1. Keep the subgraph and parity checker running through a full bake period
   after the last consumer migrates (e.g. 30 days).
2. Archive final subgraph data; tear down graph-node infrastructure.
3. Delete the `core/thegraph` client and `disperser/dataapi/subgraph` code
   paths in a separate cleanup PR.

## 8. Open Decisions

1. **Confirmation depth vs. freshness** — the subgraph serves near-head
   data; consumers like the controller request state at specific reference
   blocks. Verify `head - confirmationDepth` indexing satisfies every
   consumer's reference-block choice, or the controller will see "no APK
   snapshot" errors for very recent blocks.
2. **Archive node cost** — permanent archive-RPC dependency vs. implementing
   incremental APK computation from the event stream (which removes the
   dependency entirely). Decide before infra contracts are signed.
3. **dataapi's rich queries** — confirm the REST API's filters cover what
   `subgraph_client.go` uses today (ejection history pagination,
   batch-related lookups); anything missing needs endpoints added before
   step 2 of Phase 4.
