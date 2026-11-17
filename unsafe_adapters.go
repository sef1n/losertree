//go:build !safe

package losertree

import "unsafe"

func sliceLeafInd(s []leafInd) *leafInd {
	return unsafe.SliceData(s)
}

func sliceLeaf(s []leaf) *leaf {
	return unsafe.SliceData(s)
}

func (li *leafInd) ind(i int) *leafInd {
	return (*leafInd)(unsafe.Add(unsafe.Pointer(li), uintptr(i)*unsafe.Sizeof(*li)))
}

func (l *leaf) ind(i leafInd) *leaf {
	return (*leaf)(unsafe.Add(unsafe.Pointer(l), uintptr(i)*unsafe.Sizeof(*l)))
}
