//go:build !safe

package losertree

import "unsafe"

type (
	nodeSlice struct{ p *node }
	leafSlice struct{ p *leaf }
)

func sliceNodes(s []node) nodeSlice {
	return nodeSlice{unsafe.SliceData(s)}
}

func sliceLeaves(s []leaf) leafSlice {
	return leafSlice{unsafe.SliceData(s)}
}

func (s nodeSlice) ind(i int) node {
	return *(*node)(unsafe.Add(unsafe.Pointer(s.p), uintptr(i)*unsafe.Sizeof(*s.p)))
}

func (s leafSlice) ind(i node) leaf {
	return *(*leaf)(unsafe.Add(unsafe.Pointer(s.p), uintptr(i)*unsafe.Sizeof(*s.p)))
}

func (s nodeSlice) indP(i int) *node {
	return (*node)(unsafe.Add(unsafe.Pointer(s.p), uintptr(i)*unsafe.Sizeof(*s.p)))
}

func (s leafSlice) indP(i node) *leaf {
	return (*leaf)(unsafe.Add(unsafe.Pointer(s.p), uintptr(i)*unsafe.Sizeof(*s.p)))
}
