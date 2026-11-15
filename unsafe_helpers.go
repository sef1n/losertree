//go:build !safe

package losertree

import "unsafe"

// inline pls
func (l *leaf) unsafeI(i leafInd) *leaf {
	return (*leaf)(unsafe.Add(unsafe.Pointer(l), uintptr(i)*unsafe.Sizeof(*l)))
}

// inline pls
func (l *leafInd) unsafeI(i int) *leafInd {
	return (*leafInd)(unsafe.Add(unsafe.Pointer(l), uintptr(i)*unsafe.Sizeof(*l)))
}
