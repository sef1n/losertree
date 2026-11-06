package losertree

// leaf is a participant of a LoserTree.
//
// After reading the winner, the caller must update the
// winner leaf via Upd or Kill then call Pop to re-establish the winner.
// This package imposes no constraints on participant values
//
// Leaf keeps current int value and bool alive flag.
type leaf struct {
	v     int
	alive bool
}

// leafInd is an alias for int introduced for readability.
type leafInd = int

// Winner describes the current minimum among alive leaves.
// It holds the index of the winning leaf for value updation
// and a validity flag.
//
// I is only meaningful when valid is true.
// A Winner with valid == false means that all leaves have been
// killed (the merge is exhausted) and must not be treated
// as a real minimum. Callers should check valid before using I.
type Winner struct {
	I     leafInd
	valid bool
}

// LoserTree is a tournament tree in which each internal node
// stores the index of the leaf that lost the comparison.
// The overall winner is kept in the winner field.
type LoserTree struct {
	// nodes stores the indices of leaves.
	nodes []leafInd

	// leafs are the comparison participants.
	leafs []leaf

	// winner is the current minimum, exposed to the caller.
	winner Winner
}
