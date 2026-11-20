package losertree

import "unsafe"

func sliceNodes(s []node) *node {
	return unsafe.SliceData(s)
}

func sliceLeaves(s []leaf) *leaf {
	return unsafe.SliceData(s)
}

func indNode(s *node, i int) *node {
	return (*node)(unsafe.Add(unsafe.Pointer(s), uintptr(i)*unsafe.Sizeof(*s)))
}

func indLeaf(s *leaf, i int) *leaf {
	return (*leaf)(unsafe.Add(unsafe.Pointer(s), uintptr(i)*unsafe.Sizeof(*s)))
}
