package main

type Node interface {
	Reset()
	Execute() Result
}

type Result uint

const (
	NotDone Result = iota
	Success
	Failure
)

type sequence struct {
	cur int
	c   []Node
}

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

type NodeFunc func() Result

func (nf NodeFunc) Reset() {
}

func (nf NodeFunc) Execute() Result {
	return nf()
}

func RunBT(root Node) Result {
	r := root.Execute()
	if r != NotDone {
		root.Reset()
	}

	return r
}
