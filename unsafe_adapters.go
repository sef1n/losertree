package losertree

import "unsafe"

func newNodes(l int) *leafInd {
	return unsafe.SliceData(make([]leafInd, l))
}

func newFilledNodes(l int, filler []leafInd) *leafInd {
	s := make([]leafInd, l)
	copy(s, filler)
	return unsafe.SliceData(s)
}

func (li *leafInd) ind(i int) *leafInd {
	return (*leafInd)(unsafe.Add(unsafe.Pointer(li), uintptr(i)*unsafe.Sizeof(*li)))
}

func (l *leaf) ind(i leafInd) *leaf {
	return (*leaf)(unsafe.Add(unsafe.Pointer(l), uintptr(i)*unsafe.Sizeof(*l)))
}
