//go:build safe

package losertree

type (
	nodeSlice []node
	leafSlice []leaf
)

func sliceNodes(s []node) nodeSlice {
	return s
}

func sliceLeaves(s []leaf) leafSlice {
	return s
}

func (s nodeSlice) ind(i int) node {
	return s[i]
}

func (s leafSlice) ind(i node) leaf {
	return s[i]
}

func (s nodeSlice) indP(i int) *node {
	return &s[i]
}

func (s leafSlice) indP(i node) *leaf {
	return &s[i]
}
