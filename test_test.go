package losertree

import "testing"

func TestFull(t *testing.T) {
	vals := []int{2, 3, 5, 4, 8, 6, 9, 10}
	tr, ind, ok := Init(vals)
	for ok {
		t.Log(vals[ind])
		tr.Kill()
		ind, ok = tr.Pop()
	}
}
