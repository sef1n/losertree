package losertree

import (
	"math"
	"reflect"
	"testing"
	"unsafe"
)

func TestInit(t *testing.T) {
	inf := func() leaf { return leaf{v: math.MaxInt, alive: false} }
	rn := func(v int) leaf { return leaf{v: v, alive: true} }

	type wantTree struct {
		leavesSlice []leaf
		nodesSlice  []Index
		winner      Index
	}

	tests := []struct {
		name string
		in   []int
		want *wantTree
	}{
		{name: "nil input", in: nil, want: nil},
		{name: "empty slice", in: []int{}, want: nil},

		{
			name: "k=1 single leaf",
			in:   []int{42},
			want: &wantTree{
				leavesSlice: []leaf{rn(42)},
				nodesSlice:  []Index{0, 0},
				winner:      0,
			},
		},
		{
			name: "k=2 sorted",
			in:   []int{1, 2},
			want: &wantTree{
				leavesSlice: []leaf{rn(1), rn(2)},
				nodesSlice: []Index{
					0, 1,
					0, 1,
				},
				winner: 0,
			},
		},
		{
			name: "k=2 reversed",
			in:   []int{2, 1},
			want: &wantTree{
				leavesSlice: []leaf{rn(2), rn(1)},
				nodesSlice: []Index{
					0, 0,
					0, 1,
				},
				winner: 1,
			},
		},
		{
			name: "k=3 padded to 4",
			in:   []int{3, 2, 1},
			want: &wantTree{
				leavesSlice: []leaf{rn(3), rn(2), rn(1), inf()},
				nodesSlice: []Index{
					0, 1,
					0, 3,
					0, 1, 2, 3,
				},
				winner: 2,
			},
		},
		{
			name: "k=3 equal vals",
			in:   []int{1, 1, 1},
			want: &wantTree{
				leavesSlice: []leaf{rn(1), rn(1), rn(1), inf()},
				nodesSlice: []Index{
					0, 2,
					1, 3,
					0, 1, 2, 3,
				},
				winner: 0,
			},
		},
		{
			name: "k=4 descending",
			in:   []int{4, 3, 2, 1},
			want: &wantTree{
				leavesSlice: []leaf{rn(4), rn(3), rn(2), rn(1)},
				nodesSlice: []Index{
					0, 1,
					0, 2,
					0, 1, 2, 3,
				},
				winner: 3,
			},
		},
		{
			name: "k=4 all equal, stable tie-break",
			in:   []int{5, 5, 5, 5},
			want: &wantTree{
				leavesSlice: []leaf{rn(5), rn(5), rn(5), rn(5)},
				nodesSlice:  []Index{0, 2, 1, 3, 0, 1, 2, 3},
				winner:      0,
			},
		},
		{
			name: "k=5 padded to 8",
			in:   []int{1, 2, 3, 4, 5},
			want: &wantTree{
				leavesSlice: []leaf{
					rn(1), rn(2), rn(3), rn(4), rn(5), inf(), inf(), inf(),
				},
				nodesSlice: []Index{
					0, 4,
					2, 6,
					1, 3, 5, 7,
					0, 1, 2, 3, 4, 5, 6, 7,
				},
				winner: 0,
			},
		},
		{
			// Regression: winner = 3 while nodes[3] = 6 (padding, dead).
			// Older code derived valid from leafs[nodes[winner]].alive and
			// returned false even though the winner was the live leaf 3.
			name: "k=5 winner at leaf 3, nodes[3] is padding loser",
			in:   []int{5, 3, 2, 1, 4},
			want: &wantTree{
				leavesSlice: []leaf{
					rn(5), rn(3), rn(2), rn(1), rn(4), inf(), inf(), inf(),
				},
				nodesSlice: []Index{
					0, 4,
					1, 6,
					0, 2, 5, 7,
					0, 1, 2, 3, 4, 5, 6, 7,
				},
				winner: 3,
			},
		},
		{
			name: "MaxInt is a valid value",
			in:   []int{math.MaxInt},
			want: &wantTree{
				leavesSlice: []leaf{rn(math.MaxInt)},
				nodesSlice:  []Index{0, 0},
				winner:      0,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, wI, wV := Init(test.in)

			if test.want == nil {
				if got != nil {
					t.Errorf("Init() = %+v, want nil", got)
				}
				if wI != 0 || wV != false {
					t.Errorf("Init() = (%d, %v), want (0, false)", wI, wV)
				}
				return
			}

			if got == nil {
				t.Fatal("Init() = nil, want non-nil")
			}

			checkInvariants(t, got)

			if !reflect.DeepEqual(got.leavesSlice, test.want.leavesSlice) {
				t.Errorf("leaves mismatch\n got: %+v\nwant: %+v",
					got.leavesSlice, test.want.leavesSlice)
			}
			if !reflect.DeepEqual(got.nodesSlice, test.want.nodesSlice) {
				t.Errorf("nodes mismatch\n got: %+v\nwant: %+v",
					got.nodesSlice, test.want.nodesSlice)
			}
			if got.winner != test.want.winner {
				t.Errorf("winner = %d, want %d", got.winner, test.want.winner)
			}

			// Pointer invariants: *Ptr must alias the slice header's data.
			if got.nodesPtr != unsafe.SliceData(got.nodesSlice) {
				t.Errorf("nodesPtr = %p, want %p (== unsafe.SliceData(nodesSlice))",
					got.nodesPtr, unsafe.SliceData(got.nodesSlice))
			}
			if got.leavesPtr != unsafe.SliceData(got.leavesSlice) {
				t.Errorf("leavesPtr = %p, want %p (== unsafe.SliceData(leavesSlice))",
					got.leavesPtr, unsafe.SliceData(got.leavesSlice))
			}

			// Return values must match internal state.
			if wI != got.winner {
				t.Errorf("Init returned ind=%d, tree.winner=%d", wI, got.winner)
			}
			if !wV {
				t.Error("Init returned ok=false for non-nil tree")
			}
		})
	}
}
