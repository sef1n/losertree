package losertree

import "math"

// Upd sets the value of the current winning leaf to next.
// It does not modify alive.
// It does not update the tree; call Pop afterwards to recompute the winner.
//
// Upd panics if t is nil.
//
// WARN: calling Upd on a killed leaf leaves alive set to false,
// but the new value may win comparisons where the leaf appears as the opponent.
// This can make Pop return inconsistent results (including ok=false while alive leaves remain).
// Do not call Upd on a leaf after Kill.
func (t *LoserTree) Upd(next int) {
	t.leaves[t.winner.I].v = next
}

// TODO: revive func?

// Kill replaces the current winning leaf with a sentinel value.
// A killed leaf always loses to any alive leaf.
// After killing, the leaf's value is set to [math.MaxInt] and the leaf is marked as dead.
// Killed leaf can't be revived.
//
// Kill panics if t is nil.
func (t *LoserTree) Kill() {
	t.leaves[t.winner.I] = leaf{v: math.MaxInt}
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
func (t *LoserTree) Pop() (ind leafInd, ok bool) {
	// NOTE: without checking the gain on small k
	// if !t.winner.valid {
	// 	return 0, false
	// }

	leaves := t.leaves
	winnerS := t.winner
	nodes := t.nodes

	winner := winnerS.I

	for i := (winner + len(leaves)) >> 1; i > 0; i >>= 1 {
		opp := nodes[i]
		winner, nodes[i] = popBattle(leaves[winner], leaves[opp], winner, opp)
	}

	t.winner.I, t.winner.valid = winner, leaves[winner].alive
	return winner, t.winner.valid
}

// initBattle compares v of both leaves and checks alive of the first operand only.
// It returns the winner and loser indices.
//
// Used by Pop; checks alive of the first operand so killed leaves do not win alive leaves.
//
// WARN: A sentinel j with v < i.v can win even though j.alive is false;
// this is safe only while the invariant dead >> v == math.MaxInt holds.
// (check Upd doc)
//
// If both operands are sentinels, the winner is arbitrary.
func popBattle(i, j leaf, ii, ji leafInd) (winner, loser leafInd) {
	if j.v < i.v || !i.alive {
		return ji, ii
	}

	return ii, ji
}
