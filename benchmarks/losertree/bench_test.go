// Package ltree_bench measures the losertree merge loop.
//
// Fairness contract with benchmarks/goloser:
//   - both benchmarks are named BenchmarkMergePop
//   - both have identical sub-names k=%d,n=%d
//   - both do exactly k*(n-1) merge steps per b.N iteration
//   - both exclude tree construction from the timed section
package ltree_bench_test

import (
	"fmt"
	"testing"

	"losertree"
)

func BenchmarkMergePop(b *testing.B) {
	for _, k := range []int{4, 16, 64, 256} {
		for _, n := range []int{100, 1000} {
			b.Run(fmt.Sprintf("k=%d,n=%d", k, n), func(b *testing.B) {
				template := makeLeafs(k, n)
				iters := k * (n - 1)

				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					b.StopTimer()
					leafs := cloneLeafs(template)
					initial := make([]int, k)
					for j := range leafs {
						initial[j] = leafs[j][0]
						leafs[j] = leafs[j][1:]
					}
					tree := losertree.Init(initial)
					b.StartTimer()

					ind, ok := tree.Peek()
					for j := 0; j < iters && ok; j++ {
						if len(leafs[ind]) > 0 {
							tree.Upd(leafs[ind][0])
							leafs[ind] = leafs[ind][1:]
						} else {
							tree.Kill()
						}
						ind, ok = tree.Pop()
					}
				}
			})
		}
	}
}

func makeLeafs(k, n int) [][]int {
	leafs := make([][]int, k)
	for i := range leafs {
		l := make([]int, n)
		for j := range l {
			l[j] = i*n + j
		}
		leafs[i] = l
	}
	return leafs
}

func cloneLeafs(src [][]int) [][]int {
	dst := make([][]int, len(src))
	for i, l := range src {
		c := make([]int, len(l))
		copy(c, l)
		dst[i] = c
	}
	return dst
}
