// Package goloser_bench measures the go-loser merge loop as a baseline
// for benchmarks/ltree. See the fairness contract in
// benchmarks/ltree/bench_test.go for the full list of invariants.
package goloser_bench_test

import (
	"fmt"
	"math"
	"testing"

	"github.com/bboreham/go-loser"
)

func BenchmarkMergePop(b *testing.B) {
	for _, k := range []int{4, 16, 64, 256} {
		for _, n := range []int{100, 1000} {
			b.Run(fmt.Sprintf("k=%d,n=%d", k, n), func(b *testing.B) {
				template := makeSeqs(k, n)

				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					b.StopTimer()
					seqs := cloneSeqs(template)
					tree := loser.New[int, *seq](seqs, math.MaxInt)
					b.StartTimer()

					for tree.Next() {
						_ = tree.At()
					}
				}
			})
		}
	}
}

func BenchmarkMergeInit(b *testing.B) {
	for _, k := range []int{4, 16, 64, 256, 512, 1_000, 10_000, 100_000} {
		b.Run(fmt.Sprintf("k=%d", k), func(b *testing.B) {
			initial := genInitial(k)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = loser.New[int, *seq](initial, math.MaxInt)
			}
		})
	}
}

// seq adapts an int slice to loser.Sequence[int].
//
// cur — value returned by At(); rest — remaining values.
type seq struct {
	cur  int
	rest []int
}

func (s *seq) At() int { return s.cur }

func (s *seq) Next() bool {
	if len(s.rest) == 0 {
		return false
	}
	s.cur = s.rest[0]
	s.rest = s.rest[1:]
	return true
}

func makeSeqs(k, n int) []*seq {
	seqs := make([]*seq, k)
	for i := range seqs {
		vals := make([]int, n)
		for j := range vals {
			vals[j] = i*n + j
		}
		seqs[i] = &seq{cur: -1, rest: vals}
	}
	return seqs
}

func cloneSeqs(src []*seq) []*seq {
	dst := make([]*seq, len(src))
	for i, s := range src {
		rest := make([]int, len(s.rest))
		copy(rest, s.rest)
		dst[i] = &seq{cur: s.cur, rest: rest}
	}
	return dst
}

func genInitial(k int) []*seq {
	initial := make([]*seq, k)
	for i := range k {
		vals := []int{i}
		initial[i] = &seq{cur: -1, rest: vals}
	}
	return initial
}
