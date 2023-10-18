package parser

type Result struct {
	Nodes []*Node
	Error error
}

func newResult() Result {
	return Result{
		Nodes: make([]*Node, 0),
	}
}

func (r *Result) AddNodes(nodes ...*Node) {
	r.Nodes = append(r.Nodes, nodes...)
}

func (r *Result) SetError(err error) {
	r.Error = err
}
