package bt_test

import (
	"testing"

	"github.com/DeedleFake/bt"
)

func NodeTester(out int, r bt.Result) bt.Node {
	return bt.NodeFunc(func() bt.Result {
		// TODO: Do something with out.

		return r
	})
}

func TestSequence(t *testing.T) {
	tests := []struct {
		name string
		tree bt.Node
		ex   bt.Result
	}{
		{
			name: "Simple Failure",
			tree: bt.Sequence(NodeTester(1, bt.Failure)),
			ex:   bt.Failure,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var r bt.Result
			for r == bt.NotDone {
				r = bt.Run(test.tree)
			}

			if r != test.ex {
				t.Fatalf("Expected %q. Got %q.", test.ex, r)
			}
		})
	}
}
