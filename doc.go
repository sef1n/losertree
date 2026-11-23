// Package losertree implements a loser tree algorithm for efficient k-way merging.
// It uses unsafe to avoid bounds checking during slice indexing.
// This package intentionally ships only the unsafe build.
//
// A bounds-checked "safe" build existed during development (behind
// `//go:build safe`) but was removed.
//
// Maintaining two parallel implementations of Pop/Init without a
// concrete consumer increases the number of possible bugs without any benefit.
//
// If you need a safe build (e.g. a "no unsafe" policy), use go-loser
// instead: it is pure Go and roughly 10–30% slower for small k, which
// is the honest cost of dropping unsafe in the hot loop.
//
// The unsafe implementation is guarded by
// `go test -gcflags=all=-d=checkptr=2` and a fuzz test against sort.Ints.
package losertree
