package losertree

import (
	"math"
	"testing"
)

type popStep struct {
	cur int
	ind int
	ok  bool
}

func TestPop(t *testing.T) {
	tests := []struct {
		name    string
		streams [][]int
		want    []popStep
	}{
		{name: "k=1 single element", streams: [][]int{{42}}, want: []popStep{
			{42, 0, true}, {0, 0, false},
		}},
		{name: "k=1 two elements", streams: [][]int{{1, 2}}, want: []popStep{
			{1, 0, true}, {2, 0, true}, {0, 0, false},
		}},
		{name: "k=2 merge", streams: [][]int{{1, 2}, {3, 4}}, want: []popStep{
			{1, 0, true}, {2, 0, true}, {3, 1, true}, {4, 1, true}, {0, 0, false},
		}},
		{name: "k=2 unbalanced", streams: [][]int{{1}, {2, 3, 4}}, want: []popStep{
			{1, 0, true}, {2, 1, true}, {3, 1, true}, {4, 1, true}, {0, 0, false},
		}},
		{name: "k=3 padded", streams: [][]int{{1, 2}, {3, 4}, {5, 6}}, want: []popStep{
			{1, 0, true},
			{2, 0, true},
			{3, 1, true},
			{4, 1, true},
			{5, 2, true},
			{6, 2, true},
			{0, 0, false},
		}},
		{name: "k=5 with three levels", streams: [][]int{
			{1, 10}, {2, 20}, {3, 30}, {4, 40}, {5, 50},
		}, want: []popStep{
			{1, 0, true},
			{2, 1, true},
			{3, 2, true},
			{4, 3, true},
			{5, 4, true},
			{10, 0, true},
			{20, 1, true},
			{30, 2, true},
			{40, 3, true},
			{50, 4, true},
			{0, 0, false},
		}},
		{name: "k=4 all equal", streams: [][]int{{5}, {5}, {5}, {5}}, want: []popStep{
			{5, 0, true}, {5, 1, true}, {5, 2, true}, {5, 3, true}, {0, 0, false},
		}},
		{name: "MaxInt is a valid value", streams: [][]int{{math.MaxInt, math.MaxInt}}, want: []popStep{
			{math.MaxInt, 0, true}, {math.MaxInt, 0, true}, {0, 0, false},
		}},
		{name: "MaxInt among regular values", streams: [][]int{{math.MaxInt}, {1}, {2}}, want: []popStep{
			{1, 1, true}, {2, 2, true}, {math.MaxInt, 0, true}, {0, 0, false},
		}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			initial := make([]int, len(test.streams))
			rest := make([][]int, len(test.streams))
			for i, s := range test.streams {
				initial[i] = s[0]
				rest[i] = s[1:]
			}
			tree, ind, ok := Init(initial)

			for i, want := range test.want {
				if ok != want.ok {
					t.Fatalf("step %d: ok=%v want %v", i, ok, want.ok)
				}
				if !ok {
					continue
				}
				if ind != want.ind {
					t.Errorf("step %d: ind=%d want %d", i, ind, want.ind)
				}
				if got := tree.leaves[ind].v; got != want.cur {
					t.Errorf("step %d: val=%d want %d", i, got, want.cur)
				}

				if len(rest[ind]) > 0 {
					tree.Upd(rest[ind][0])
					rest[ind] = rest[ind][1:]
				} else {
					tree.Kill()
				}
				ind, ok = tree.Pop()
			}

			// Calling Pop on an exhausted tree must be a no-op (0, false),
			// not a panic and not a stale winner.
			if _, ok := tree.Pop(); ok {
				t.Error("Pop on exhausted tree should return ok=false")
			}
		})
	}
}

// TestKill is a regression test for a bug where Kill used a leaf copy
// (l := t.leaves[i]; l.alive = false) instead of mutating the slice,
// leaving the tree unchanged. Verifies that Kill marks the winner dead,
// sets its value to MaxInt, leaves other leaves untouched, and that Pop
// advances to the next live winner.
func TestKill(t *testing.T) {
	tree, _, _ := Init([]int{1, 2, 3})

	tree.Kill() // winner is leaf 0 with value 1

	if tree.leaves[0].alive {
		t.Error("leaf 0 still alive after Kill")
	}
	if tree.leaves[0].v != math.MaxInt {
		t.Errorf("leaf 0 v = %d, want MaxInt", tree.leaves[0].v)
	}
	if !tree.leaves[1].alive || !tree.leaves[2].alive {
		t.Error("Kill affected leaves other than the current winner")
	}

	// After Kill, Pop must re-sort the tree and return the next live winner.
	ind, ok := tree.Pop()
	if !ok {
		t.Fatal("Pop after Kill returned ok=false, want a new winner")
	}
	if ind != 1 {
		t.Errorf("winner ind = %d, want 1", ind)
	}
}

// TestUpd verifies that Upd replaces the winner's value without touching
// its alive flag, and that the change propagates to the tree on the next Pop.
func TestUpd(t *testing.T) {
	tree, _, _ := Init([]int{1, 2, 3})

	tree.Upd(100) // winner is leaf 0; raise its value to 100
	if tree.leaves[0].v != 100 {
		t.Errorf("leaf 0 v = %d, want 100", tree.leaves[0].v)
	}
	if !tree.leaves[0].alive {
		t.Error("Upd marked the leaf dead")
	}

	// After raising leaf 0 to 100, the new winner should be leaf 1 (value 2).
	ind, ok := tree.Pop()
	if !ok || ind != 1 {
		t.Errorf("Pop = (%d, %v), want (1, true)", ind, ok)
	}
}
