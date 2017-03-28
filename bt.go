package main

type Node interface {
	Reset()
	Execute(v interface{}) Result
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

func (s *sequence) Execute(v interface{}) Result {
	if s.cur >= len(s.c) {
		return Success
	}

	switch r := s.c[s.cur].Execute(v); r {
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

type NodeFunc func(interface{}) Result

func (nf NodeFunc) Reset() {
}

func (nf NodeFunc) Execute(v interface{}) Result {
	return nf(v)
}

func RunBT(root Node, v interface{}) Result {
	r := root.Execute(v)
	if r != NotDone {
		root.Reset()
	}

	return r
}
