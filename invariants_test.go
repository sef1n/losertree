package losertree

import (
	"math"
	"testing"
	"unsafe"
)

// checkInvariants verifies every property that unsafe pointer arithmetic
// in Pop, Upd, and Kill relies on. O(k). Test/fuzz only.
//
// If any of these break, the tree may read or write foreign memory:
// indLeaf / indNode do no bounds checking.
func checkInvariants(t *testing.T, tree *LoserTree) {
	t.Helper()

	ll := len(tree.leavesSlice)
	nl := len(tree.nodesSlice)

	// (1) Slice && pointer aliasing.
	// Pop reads/writes through *Ptr but
	// GC keeps the array alive only via the Slice header.
	// Divergence here means a silent use-after-free.
	if tree.nodesPtr != unsafe.SliceData(tree.nodesSlice) {
		t.Fatalf("nodesPtr != unsafe.SliceData(nodesSlice)")
	}
	if tree.leavesPtr != unsafe.SliceData(tree.leavesSlice) {
		t.Fatalf("leavesPtr != unsafe.SliceData(leavesSlice)")
	}

	// (2) Layout.
	// Pop computes the walk start as (winner + nl) >> 1;
	// this requires len(nodes) == 2*len(leaves).
	if nl != 2*ll {
		t.Fatalf("len(nodes)=%d, want 2*len(leaves)=%d",
			nl, 2*ll)
	}

	// (3) Every node value is a valid leaf index.
	// popBattle returns one of its inputs,
	// so this holds iff winner stays in range (see (4)).
	for i, n := range tree.nodesSlice {
		if n < 0 || int(n) >= ll {
			t.Fatalf("nodes[%d]=%d out of [0,%d)", i, n, ll)
		}
	}

	// (4) Winner is a valid leaf index.
	if tree.winner < 0 || int(tree.winner) >= ll {
		t.Fatalf("winner=%d out of [0,%d)", tree.winner, ll)
	}

	// (5) Sentinel invariant used by popBattle: dead implies MaxInt.
	//    popBattle lets j win on j.v < i.v without checking j.alive;
	//    that is only sound while no dead leaf carries a finite value.
	//    Violated only by Upd after Kill (contract violation).
	for i, l := range tree.leavesSlice {
		if !l.alive && l.v != math.MaxInt {
			t.Fatalf("leaves[%d]: dead but v=%d, want MaxInt", i, l.v)
		}
	}
}
