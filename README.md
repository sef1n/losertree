# losertree

Loser tree (tournament tree) for k-way merge in Go.

Zero dependencies. `int` keys. No generics, no callbacks, no interface dispatch in the hot path.

## Why

`k`-way merge is the bottleneck of every external merge sort: for each emitted record the tree does `log₂ k` comparisons. `losertree` is a minimal, allocation-free loser tree tuned exactly for that loop — `int` keys, two contiguous slices, a single inlined `battle`. It is 10–30% faster than `go-loser` for small and medium `k` (4–16) and on par for large `k`.

## Features

- **Heap-ordered layout.** `nodes[1]` is the root, leaves occupy `nodes[ext..2·ext-1]`. The path from any leaf to the root is just `i >>= 1` — no level arithmetic in the hot loop.
- **Two flat slices.** `leafs []leaf` and `nodes []leafInd`. No pointers, no per-node allocations, no GC scanning of the tree.
- **`int` keys, no abstractions.** `battle` is a single inlined function of two compares. No `Sequence` interface, no function-pointer callback, no type switch.
- **Streaming API.** `Peek` → `Upd`/`Kill` → `Pop`. The caller owns the sequences (files, readers, in-memory slices) and advances them by index.
- **`math.MaxInt` is a valid key.** Liveness is an `alive` flag on each leaf, not a sentinel value.
- **Deterministic tie-break.** On equal keys the left leaf wins, so the merge preserves source order.
- **Not goroutine-safe.** One tree per goroutine; no locks, no atomics.

## Install

```sh
go get github.com/sef1n/losertree
```
