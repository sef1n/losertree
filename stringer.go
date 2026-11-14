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
	sb.Grow(64 + 32*len(t.nodes) + 32*len(t.leaves))

	fmt.Fprint(&sb, "{\nLeaves: {")
	for i, l := range t.leaves {
		if i > 0 {
			fmt.Fprint(&sb, " ")
		}
		fmt.Fprintf(&sb, "(i:%d,v:%d,alive:%v)", i, l.v, l.alive)
	}
	fmt.Fprint(&sb, "};\nNodes: {\n")

	level := 0
	levelStart := 1 // 0th level begin
	for i := 1; i < len(t.nodes); i++ {
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
		fmt.Fprintf(&sb, "%d", t.nodes[i])
	}
	if len(t.nodes) > 1 {
		fmt.Fprint(&sb, "\n")
	}

	fmt.Fprintf(&sb, "}\nwinner: {I:%d, valid:%v}}", t.winner.I, t.winner.valid)
	return sb.String()
}
