package heap_bench_test

import (
	"container/heap"
	"fmt"
	"testing"
)

type seq struct {
	cur  int
	rest []int
}

func (s *seq) Next() bool {
	if len(s.rest) == 0 {
		return false
	}
	s.cur = s.rest[0]
	s.rest = s.rest[1:]
	return true
}

type seqsHeap []seq

// sort.Interface
func (h seqsHeap) Len() int           { return len(h) }
func (h seqsHeap) Less(i, j int) bool { return h[i].cur < h[j].cur }
func (h seqsHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

// heap.Interface
func (h *seqsHeap) Push(x any) {
	*h = append(*h, x.(seq))
}

func (h *seqsHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

func BenchmarkMergePop(b *testing.B) {
	for _, k := range []int{4, 16, 64, 256} {
		for _, n := range []int{100, 1_000} {
			b.Run(fmt.Sprintf("k=%d,n=%d", k, n), func(b *testing.B) {
				template := makeLeafs(k, n)

				for i := 0; i < b.N; i++ {
					b.StopTimer()
					h := cloneSeqs(template)
					heap.Init(&h)
					b.StartTimer()

					for h.Len() > 0 {
						popped := heap.Pop(&h).(seq)
						if popped.Next() {
							heap.Push(&h, popped)
						}
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

func cloneSeqs(src [][]int) seqsHeap {
	dst := make(seqsHeap, len(src))
	for i, l := range src {
		dst[i] = seq{cur: l[0], rest: l[1:]}
	}
	return dst
}
