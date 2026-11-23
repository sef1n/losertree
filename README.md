# losertree

Loser tree (tournament tree) for k-way merge in Go.

Zero dependencies. `int` keys.
No generics, no callbacks, no interface dispatch in the hot path.

## Why

`k`-way merge is the bottleneck of every external merge sort:
for each emitted record the tree does `log₂ k` comparisons.
`losertree` is a minimal,allocation-free loser tree tuned exactly for
that loop — `int` keys, two contiguous slices, a single inlined `battle`.
It is 16–90% faster than `go-loser` in the merge loop (48% geomean)
and 11–14× faster than `container/heap` (see [Benchmarks](#benchmarks)).

## Features

- **Heap-ordered layout.**
`nodes[1]` is the root, leaves occupy `nodes[ext..2·ext-1]`.
The path from any leaf to the root is just `i >> 1` -
no level arithmetic in the hot loop.

- **Two flat slices.**
`leaves []leaf` and `nodes []leafInd`.
No pointers, no per-node allocations, no GC scanning of the tree.

- **`int` keys, no abstractions.**
`battle` is a single inlined function of two compares.
No `Sequence` interface, no function-pointer callback, no type switch.

- **Streaming API.**
`Peek` → `Upd`/`Kill` → `Pop`.
The caller owns the sequences (files, readers, in-memory slices)
and advances them by index.

- **`math.MaxInt` is a valid key.**
Liveness is an `alive` flag on each leaf, not a sentinel value.

- **Deterministic tie-break.**
On equal keys the left leaf wins, so the merge preserves source order.

- **Not goroutine-safe.**
One tree per goroutine; no locks, no atomics.

## Install

```sh
go get github.com/sef1n/losertree
```

## Quick Start

```go
package main

import (
 "github.com/sef1n/losertree"
)

// streams is k sorted int slices.
func merge(streams [][]int) []int {
 k := len(streams)
 initial := make([]int, k)
 for i := range streams {
  initial[i] = streams[i][0]
  streams[i] = streams[i][1:]
 }

 tree, ind, ok := losertree.Init(initial)
 if tree == nil {
  return nil
 }

 out := make([]int, 0)
 for ok {
  out = append(out, initial[ind]) // or track per-source values

  if len(streams[ind]) == 0 {
   tree.Kill()
  } else {
   initial[ind] = streams[ind][0]
   tree.Upd(initial[ind])
   streams[ind] = streams[ind][1:]
  }
  ind, ok = tree.Pop()
 }
 return out
}
```

## Benchmarks

Merge of `k` sorted `int` streams of `n` values against
[`go-loser`](https://github.com/bboreham/go-loser) and the standard
library's [`container/heap`](https://pkg.go.dev/container/heap).
Only the merge loop is timed — exactly `k·n` steps per iteration,
construction is outside the measured region. All three benchmarks share
the same names and shapes; see [`benchmarks/`](benchmarks/) for the full
fairness contract.

Environment: `linux/amd64`, Intel Core i5-11400F,
`BCount=10`, `BTime=1s`, `benchstat` medians.

| `k` | `n` | `losertree` | `go-loser` | `container/heap` |
| ---: | ---: | ---: | ---: | ---: |
| 4 | 100 | **1.98 µs** | 3.60 µs (+82%) | 21.4 µs (+984%) |
| 4 | 1 000 | **17.8 µs** | 33.8 µs (+90%) | 208 µs (+1 073%) |
| 16 | 100 | **9.98 µs** | 15.6 µs (+56%) | 106 µs (+967%) |
| 16 | 1 000 | **101 µs** | 153 µs (+52%) | 1.08 ms (+970%) |
| 64 | 100 | **50.3 µs** | 69.4 µs (+38%) | 544 µs (+981%) |
| 64 | 1 000 | **592 µs** | 685 µs (+16%) | 6.49 ms (+996%) |
| 256 | 100 | **243 µs** | 317 µs (+30%) | 3.18 ms (+1 208%) |
| 256 | 1 000 | **2.34 ms** | 3.11 ms (+33%) | 32.9 ms (+1 304%) |
| **geomean** | | **70.3 µs** | **104 µs (+48%)** | **812 µs (+1 055%)** |

Values in parentheses are slowdown relative to `losertree`. All rows
are significant (`p=0.000`, `n=10`). `losertree` and `go-loser` do
`0 allocs/op` in the merge loop.

`container/heap` is measured with the straightforward `heap.Interface`
implementation from [`benchmarks/heap`](benchmarks/heap): values are
pushed and popped through `any`, so interface boxing dominates its
allocations. Treat it as a common baseline, not a tuned heap.

Reproduce (results go to `AllMergeBench.bench`):

```sh
make test-bench BName="Merge" BCount=10 BTime=1s BOut="AllMergeBench.bench" \
  BPkgs="./benchmarks/losertree/ ./benchmarks/goloser/ ./benchmarks/heap/"
```
