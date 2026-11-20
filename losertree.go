// Package losertree implements a loser tree algorithm for efficient k-way merging.
// It uses unsafe to avoid borders checking during slice indexing.
package losertree

// leaf is a participant of a [LoserTree].
//
// After reading the winner, the caller must update the
// winning leaf via Upd or Kill, then call Pop to re-establish the winner.
// This package imposes no constraints on participant values.
//
// A leaf keeps current int value and an alive flag.
type leaf struct {
	v     int
	alive bool
}

// node is an alias for int introduced for readability.
type node = int

// [LoserTree] is a tournament tree in which each internal node
// stores the index of the leaf that lost the comparison.
// The overall winner is kept in the winner field.
type LoserTree struct {
	// nodesPtr stores the indices of leaves.
	nodesPtr *node
	// leavesPtr are the comparison participants.
	leavesPtr *leaf

	// headers of slices to prevent garbage collection.
	nodesSlice  []node
	leavesSlice []leaf

	// winner is the current minimum, exposed to the caller.
	winner node
}
