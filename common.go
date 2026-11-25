package losertree

import (
	"math/bits"
)

// Index is an alias for int introduced for readability.
//
// It is shared by every generated tree variant because node indices are
// always plain ints, regardless of the key type.
type Index = int

// nextPowerOfTwo returns the smallest power of two not less than n.
func nextPowerOfTwo(n int) int {
	if n <= 1 {
		return 1
	}
	return 1 << bits.Len(uint(n-1))
}

// initNodes initializes the nodes slice of tree
// with indices of participants.
func initNodes(ext int) []Index {
	nodes := make([]Index, ext*2)
	// fill start level of tree with indexes
	for i := 0; i < ext; i++ {
		nodes[ext+i] = i
	}
	return nodes
}

// initWinners creates a temporary nodes slice by copying the nodes[ext:].
func initWinners(nodes []Index, ext int) []Index {
	w := make([]Index, ext*2)
	copy(w[ext:], nodes[ext:])
	return w
}
