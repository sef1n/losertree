package losertree

import (
	"math"
	"math/bits"
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
func Init(values []int) (t *LoserTree, ind node, ok bool) {
	// ok for identity with the Pop
	if len(values) == 0 {
		return // nil, 0, false
	}

	t = &LoserTree{}
	// len of lt.leaves with sentinels
	ext := nextPowerOfTwo(len(values))

	t.leavesSlice = initLeaves(values, ext)
	t.nodesSlice = initNodes(ext)
	winnersSlice := initWinners(t.nodesSlice, ext)

	leaves := sliceLeaves(t.leavesSlice)
	nodes := sliceNodes(t.nodesSlice)
	winners := sliceNodes(winnersSlice)

	for p := ext - 1; p >= 1; p-- {
		i, j := *indNode(winners, p*2), *indNode(winners, p*2+1)
		winner, loser := initBattle(*indLeaf(leaves, i), *indLeaf(leaves, j), i, j)

		*indNode(winners, p) = winner
		*indNode(nodes, p) = loser
	}

	t.winner = *indNode(winners, 1)
	t.leavesPtr = leaves
	t.nodesPtr = nodes

	return t, t.winner, true
}

func nextPowerOfTwo(n int) int {
	if n <= 1 {
		return 1
	}
	return 1 << bits.Len(uint(n-1))
}

func initLeaves(values []int, ext int) []leaf {
	leaves := make([]leaf, ext)

	// vals of values to lt.leaves vals
	for i, val := range values {
		leaves[i].v = val
		leaves[i].alive = true
	}

	// create sentinels
	for i := len(values); i < ext; i++ {
		leaves[i].v = math.MaxInt
	}

	return leaves
}

func initNodes(ext int) []node {
	nodes := make([]node, ext*2)
	// fill start level of tree with indexes
	for i := 0; i < ext; i++ {
		nodes[ext+i] = i
	}
	return nodes
}

func initWinners(nodes []node, ext int) []node {
	w := make([]node, ext*2)
	copy(w[ext:], nodes[ext:])
	return w
}

// initBattle compares only the v fields of leaf structs and returns the winner and loser indices.
//
// Used during Init; faster because it skips the alive check.
//
// It is safe to call only during Init, before the first Kill or Upd call, because
// of the layout invariant maintained during Init: a sentinel can't be on the left side of an alive leaf.
func initBattle(i, j leaf, ii, ji node) (winner, loser node) {
	if j.v < i.v {
		return ji, ii
	}
	return ii, ji
}
