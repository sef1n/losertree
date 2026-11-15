//go:build !safe

package losertree

func (t *LoserTree) String() string {
	if t == nil {
		return "{<nil LoserTree>}"
	}
	return "stringer is unfinished for unsafe tree"
}
