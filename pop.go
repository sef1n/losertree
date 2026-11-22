package losertree

import (
	"math"
)

// Upd sets the value of the current winning leaf to next.
// It does not modify alive.
// It does not update the entire tree; call Pop afterwards to recompute the winner.
//
// Do not call Upd on an exhausted tree (when Pop returned ok == false).
// Upd panics if t is nil.
func (t *LoserTree) Upd(next int) {
	indLeaf(t.leavesPtr, t.winner).v = next
}

// Kill replaces the current winning leaf with a sentinel value.
// A killed leaf always loses to any alive leaf.
// After killing, the leaf's value is set to [math.MaxInt] and the leaf is marked as dead.
// A killed leaf cannot be revived.
//
// Kill panics if t is nil.
func (t *LoserTree) Kill() {
	*indLeaf(t.leavesPtr, t.winner) = leaf{v: math.MaxInt}
}

// Pop updates the tree's winner and the path from the winning leaf to the root.
//
// Before calling Pop, the caller must supply the next value via t.Upd(next) or t.Kill().
//
// If no alive leaf remains, ok is false;
// otherwise ok is true and ind is the index of the winning leaf.
// ind is meaningless when ok is false.
// The caller must check the ok flag before using ind.
//
// Calling Pop twice without Upd/Kill recomputes the same path and returns the same ind.
//
// Pop panics if t is nil.
func (t *LoserTree) Pop() (ind Index, ok bool) {
	leaves := t.leavesPtr
	nodes := t.nodesPtr
	winner := t.winner

	for i := (winner + len(t.leavesSlice)) >> 1; i > 0; i >>= 1 {
		opp := indNode(nodes, i)
		oppI := *opp

		winner, *opp = popBattle(
			*indLeaf(leaves, winner), *indLeaf(leaves, oppI),
			winner, oppI)
	}

	t.winner = winner
	return winner, indLeaf(leaves, winner).alive
}

// popBattle compares v of both leaves and checks alive of the first operand only.
// It returns the winner and loser indices.
//
// Used by Pop; checks alive of the first operand so killed leaves do not win against alive leaves.
//
// WARN: A sentinel j with v < i.v can win even though j.alive is false;
// this is safe only while the invariant dead => v == math.MaxInt holds.
//
// If both operands are sentinels, the winner is arbitrary.
func popBattle(i, j leaf, ii, ji Index) (winner, loser Index) {
	if j.v < i.v || !i.alive {
		return ji, ii
	}

	return ii, ji
}
