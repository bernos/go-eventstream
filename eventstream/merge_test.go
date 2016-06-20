package eventstream

import "testing"

func TestMerge(t *testing.T) {
	a := FromValues(1, 3, 5, 7, 9)
	b := FromValues(2, 4, 6, 8, 10)

	out := Merge(a, b)

	for event := range out.Events() {
		if event.Value() == 5 {
			out.Cancel()
		}
	}
}
