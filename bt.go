package bt

// A Node is a node of a behavior tree.
type Node interface {
	// Reset resets the node. It is the responsibility of a Node's
	// implementation to reset its children, if it has any.
	Reset()

	// Execute executes the node.
	Execute() Result
}

// Result represents the result of a node's execution.
type Result uint

const (
	// NotDone signals that a node needs more calls to its Execute
	// method in order to finish.
	NotDone Result = iota

	// Success and Failure indicate exactly what they sound like.
	Success
	Failure
)

type sequence struct {
	cur int
	c   []Node
}

// Sequence returns a node which iterates over each of its children
// before returning Success. This works according to the following
// rules:
//
//    * If a child returns NotDone, it is called again next time.
//    * If a child returns Success, the next child will be called next
//      time.
//    * If a child returns Failure, the Sequence returns Failure.
//
// Upon being reset, the Sequence returns to the first child.
func Sequence(children ...Node) Node {
	return &sequence{c: children}
}

func (s *sequence) Reset() {
	s.cur = 0
	for _, c := range s.c {
		c.Reset()
	}
}

func (s *sequence) Execute() Result {
	if s.cur >= len(s.c) {
		return Success
	}

	switch r := s.c[s.cur].Execute(); r {
	case Success:
		s.cur++
		if s.cur >= len(s.c) {
			return Success
		}
		return NotDone

	default:
		return r
	}
}

type selector struct {
	cur int
	c   []Node
}

// Selector returns a node which iterates over each of its children
// before returning Failure. This works according to the following
// rules:
//
//    * If a child returns NotDone, it is called again next time.
//    * If a child returns Success, the Selector returns Success.
//    * If a child returns Failure, the next child will be called next
//      time.
//
// Upon being reset, the Selector returns to the first child.
func Selector(children ...Node) Node {
	return &selector{c: children}
}

func (s *selector) Reset() {
	s.cur = 0
	for _, c := range s.c {
		c.Reset()
	}
}

func (s *selector) Execute() Result {
	if s.cur >= len(s.c) {
		return Failure
	}

	switch r := s.c[s.cur].Execute(); r {
	case Failure:
		s.cur++
		if s.cur >= len(s.c) {
			return Failure
		}
		return NotDone

	default:
		return r
	}
}

// A NodeFunc is a wrapper for functions that allows them to more
// easily be used as the Execute method of a Node. The resulting
// Node's Reset method is a no-op.
type NodeFunc func() Result

func (nf NodeFunc) Reset() {
}

func (nf NodeFunc) Execute() Result {
	return nf()
}

// Run runs a tree for one tick, returning the result of the root
// node's Execute method. If the root node signals either success or
// failure, the node is reset, potentially resetting the entire tree.
func Run(root Node) Result {
	r := root.Execute()
	if r != NotDone {
		root.Reset()
	}

	return r
}
