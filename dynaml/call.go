package dynaml

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/shutej/spiff/yaml"
)

type CallExpr struct {
	Name      string
	Arguments []Expression
}

func (e CallExpr) RequiresPhases(binding Binding) StringSet {
	retval := StringSet{}
	for _, a := range e.Arguments {
		retval.Update(a.RequiresPhases(binding))
	}

	builtin, ok := binding.Builtin(e.Name)
	if ok {
		retval.Update(builtin.RequiredPhases)
	}
	return retval
}

func (e CallExpr) guard(binding Binding) bool {
	return binding.ProvidesPhases(e.RequiresPhases(binding))
}

func (e CallExpr) Evaluate(binding Binding) (yaml.Node, bool) {
	// This acts as the identity function.
	if !e.guard(binding) {
		return node(e), true
	}

	builtin, ok := binding.Builtin(e.Name)
	if ok {
		t := builtin.Function.Type()
		if t.NumIn() != len(e.Arguments) {
			return nil, false
		}

		args := make([]reflect.Value, 0, t.NumIn())
		for i, arg := range e.Arguments {
			index, ok := arg.Evaluate(binding)
			if !ok {
				return nil, false
			}
			v := reflect.ValueOf(index.Value())
			if !v.Type().AssignableTo(t.In(i)) {
				return nil, false
			}
			args = append(args, v)
		}

		retval := builtin.Function.Call(args)
		if len(retval) != 1 {
			panic("builtins must return exactly one value!")
		}
		return node(retval[0].Interface()), true
	}

	return nil, false
}

func (e CallExpr) String() string {
	args := make([]string, len(e.Arguments))
	for i, e := range e.Arguments {
		args[i] = fmt.Sprintf("%s", e)
	}

	return fmt.Sprintf("%s(%s)", e.Name, strings.Join(args, ", "))
}
