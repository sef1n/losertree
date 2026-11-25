package losertree

// Index is an alias for int introduced for readability.
//
// It is shared by every generated tree variant because node indices are
// always plain ints, regardless of the key type.
type Index = int

// leaf is a participant of a tournament tree.
//
// V is the key type of a generated tree variant, so only the leaf alias
// (for example, leafInt = leaf[int]) has to be generated.
//
// After reading the winner, the caller must update the winning leaf
// via Upd or Kill, then call Pop to re-establish the winner.
// This package imposes no constraints on participant values.
//
// A leaf keeps the current key value and an alive flag.
type leaf[V any] struct {
	v     V
	alive bool
}
