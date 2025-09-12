package internal

import (
	"encoding/json"
	"fmt"
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/overloads"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/ext"
)

const varNameSelf = "self"

type Filter struct {
	prog cel.Program
}

func NewFilter(filter string) (*Filter, error) {
	env, err := cel.NewEnv(
		cel.Variable(varNameSelf, cel.MapType(cel.StringType, cel.AnyType)),
		cel.Function(overloads.TypeConvertString, cel.Overload(
			"map_to_string", []*cel.Type{cel.MapType(cel.StringType, cel.AnyType)}, cel.StringType,
			cel.UnaryBinding(func(value ref.Val) ref.Val {
				b, _ := json.Marshal(value.Value())
				return types.String(b)
			}),
		)),
		ext.Strings(),
		ext.Encoders(),
		ext.Math(),
		ext.Sets(),
		ext.Lists(),
		ext.TwoVarComprehensions(),
	)
	if err != nil {
		return nil, fmt.Errorf("filter: new: %v", err)
	}

	ast, iss := env.Compile(filter)
	if iss != nil && iss.Err() != nil {
		return nil, fmt.Errorf("filter: new: %v", iss.Err())
	}

	prog, err := env.Program(ast)
	if err != nil {
		return nil, fmt.Errorf("filter: new: %v", err)
	}

	return &Filter{prog: prog}, nil
}

func (f *Filter) Eval(obj any) (bool, error) {
	bytes, err := json.Marshal(obj)
	if err != nil {
		return false, err
	}

	data := make(map[string]any)
	if err := json.Unmarshal(bytes, &data); err != nil {
		return false, err
	}

	eval, _, err := f.prog.Eval(map[string]any{varNameSelf: data})
	if err != nil {
		return false, fmt.Errorf("filter: eval: %v", err)
	}
	return eval.Value().(bool), nil
}
