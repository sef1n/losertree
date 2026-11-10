package losertree

import (
	"math"
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
	winners := make([]leafInd, ext*2)
	copy(winners[ext:], t.nodes[ext:])

	for p := ext - 1; p >= 1; p-- {
		i, j := winners[p*2], winners[p*2+1]
		winner, loser := initBattle(leaves[i], leaves[j], i, j)

		winners[p] = winner
		nodes[p] = loser
	}
	t.winner.I = winners[1]
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
		leaves: make([]leaf, ext),
		nodes:  make([]leafInd, ext*2),
	}

	// vals of values to lt.leaves vals
	for i, val := range values {
		lt.leaves[i].v = val
		lt.leaves[i].alive = true
	}

	// create sentinels
	for i := len(values); i < ext; i++ {
		lt.leaves[i].v = math.MaxInt
	}

	// fill start level of tree with indexes
	for i := 0; i < ext; i++ {
		lt.nodes[ext+i] = i
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
