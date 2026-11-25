package losertree

import "unsafe"

// sliceNodes returns a pointer to the first element of s.
// Similar to unsafe.SliceData(s).
func sliceNodes(s []Index) *Index {
	return unsafe.SliceData(s)
}

// indNode returns a pointer to the ith node, starting from s.
// It is equivalent to &s[i], but skips bounds checking.
func indNode(s *Index, i int) *Index {
	return (*Index)(unsafe.Add(unsafe.Pointer(s), uintptr(i)*unsafe.Sizeof(*s)))
}
