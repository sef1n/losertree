//go:build !safe

package losertree

import (
	"fmt"
	"strings"
)

// String formats the tree for debugging.
// Levels are printed from the root (L#0) down; nodes[0] is not used.
func (t *LoserTree) String() string {
	if t == nil {
		return "{<nil LoserTree>}"
	}

	var sb strings.Builder
	sb.Grow(64 + 32*len(t.nodesSlice) + 32*len(t.leavesSlice))

	fmt.Fprint(&sb, "{\nLeaves: {")
	for i, l := range t.leavesSlice {
		if i > 0 {
			fmt.Fprint(&sb, " ")
		}
		fmt.Fprintf(&sb, "(i:%d,v:%d,alive:%v)", i, l.v, l.alive)
	}
	fmt.Fprint(&sb, "};\nNodes: {\n")

	level := 0
	levelStart := 1 // 0th level begin
	for i := 1; i < len(t.nodesSlice); i++ {
		if i == levelStart {
			if i > 1 {
				fmt.Fprint(&sb, "\n")
			}
			fmt.Fprintf(&sb, "L#%d: ", level)
			level++
			levelStart <<= 1
		} else {
			fmt.Fprint(&sb, ",")
		}
		fmt.Fprintf(&sb, "%d", t.nodesSlice[i])
	}
	if len(t.nodesSlice) > 1 {
		fmt.Fprint(&sb, "\n")
	}

	fmt.Fprintf(&sb, "}\nwinner: %v}", t.winner)
	return sb.String()
}
