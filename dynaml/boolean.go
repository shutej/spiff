package dynaml

import (
	"fmt"

	"github.com/shutej/spiff/yaml"
)

type BooleanExpr struct {
	Value bool
}

func (e BooleanExpr) RequiresPhases(_ Binding) StringSet {
	return StringSet(nil)
}

func (e BooleanExpr) Evaluate(Binding) (yaml.Node, bool) {
	return node(e.Value), true
}

func (e BooleanExpr) String() string {
	return fmt.Sprintf("%v", e.Value)
}
