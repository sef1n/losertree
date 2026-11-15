//go:build !safe

package losertree

import (
	"math"
	"unsafe"
)

// Init creates a LoserTree from the given values.
// [math.MaxInt] is treated as an ordinary (live) value.
//
// If values is nil or empty, it returns (nil, 0, false).
// Otherwise it returns the [*LoserTree], the index of the winning leaf, and ok.
// Both the [*LoserTree] and the index are meaningful only when ok is true.
//
// Init copies the input values and pads the leaves slice with dead entries to
// the smallest power of two not less than len(values).
// The nodes slice is twice as long as the leaves slice
func Init(values []int) (t *LoserTree, ind leafInd, ok bool) {
	if values == nil || len(values) == 0 {
		return // nil, 0, false
	}
	t, ext := initiateStruct(values)

	leaves := t.leaves
	nodes := t.nodes

	winnersSlice := make([]leafInd, ext*2)
	winners := unsafe.SliceData(winnersSlice)
	copy(winnersSlice[ext:], t.nodesSlice[ext:])

	for p := ext - 1; p >= 1; p-- {
		i, j := *winners.unsafeI(p * 2), *winners.unsafeI(p*2 + 1)
		winner, loser := initBattle(*leaves.unsafeI(i), *leaves.unsafeI(j), i, j)

		*winners.unsafeI(p) = winner
		*nodes.unsafeI(p) = loser
	}
	t.winner.I = *winners.unsafeI(1)
	t.winner.valid = true

	return t, t.winner.I, t.winner.valid
}

func initiateStruct(values []int) (*LoserTree, int) {
	// len of lt.leaves with sentinels
	ext := 1
	for ext < len(values) {
		ext <<= 1
	}

	lt := LoserTree{
		leavesSlice: make([]leaf, ext),
		nodesSlice:  make([]leafInd, ext*2),
	}
	lt.leaves = unsafe.SliceData(lt.leavesSlice)
	lt.nodes = unsafe.SliceData(lt.nodesSlice)

	// vals of values to lt.leaves vals
	for i, val := range values {
		lt.leaves.unsafeI(leafInd(i)).v = val
		lt.leaves.unsafeI(leafInd(i)).alive = true
	}

	// create sentinels
	for i := len(values); i < ext; i++ {
		lt.leaves.unsafeI(leafInd(i)).v = math.MaxInt
	}

	// fill start level of tree with indexes
	for i := 0; i < ext; i++ {
		*lt.nodes.unsafeI(ext + i) = leafInd(i)
	}
	return &lt, ext
}

// initBattle compares only the v fields of leaf structs and returns the winner and loser indices.
//
// Used during Init; faster because it skips the alive check.
//
// It is safe to call only during Init, before the first Kill or Upd call, because
// of the layout invariant maintained during Init: a sentinel can't be on the left side of an alive leaf.
func initBattle(i, j leaf, ii, ji leafInd) (winner, loser leafInd) {
	if j.v < i.v {
		return ji, ii
	}
	return ii, ji
}
