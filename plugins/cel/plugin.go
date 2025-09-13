package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/overloads"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/ext"
	"sync"
	"sync/atomic"
)

const varNameSelf = "self"

type Plugin struct {
	mtx    sync.Mutex
	isInit atomic.Bool
	prog   cel.Program
}

var Export = Plugin{}

type Config struct {
	Expr string `yaml:"expr"`
}

func (p *Plugin) Init(cfg []byte) error {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	celCfg := &Config{}
	if err := json.Unmarshal(cfg, celCfg); err != nil {
		return err
	}

	env, err := cel.NewEnv(
		cel.Variable(varNameSelf, cel.MapType(cel.StringType, cel.DynType)),
		cel.Function(overloads.TypeConvertString, cel.Overload(
			"map_to_string", []*cel.Type{cel.MapType(cel.StringType, cel.DynType)}, cel.StringType,
			cel.UnaryBinding(func(value ref.Val) ref.Val {
				b, _ := json.Marshal(value.Value())
				return types.String(b)
			}),
		)),
		cel.OptionalTypes(),
		ext.Regex(),
		ext.Strings(),
		ext.Encoders(),
		ext.Math(),
		ext.Sets(),
		ext.Lists(),
		ext.TwoVarComprehensions(),
	)
	if err != nil {
		return fmt.Errorf("cel: init: %v", err)
	}

	ast, iss := env.Compile(celCfg.Expr)
	if iss != nil && iss.Err() != nil {
		return fmt.Errorf("cel: init: %v", iss.Err())
	}

	prog, err := env.Program(ast)
	if err != nil {
		return fmt.Errorf("cel: init: %v", err)
	}

	p.prog = prog

	p.isInit.Store(true)
	return nil
}

func (p *Plugin) IsInit() bool {
	return p.isInit.Load()
}

func (p *Plugin) Eval(obj []byte) (any, error) {
	if !p.IsInit() {
		return nil, fmt.Errorf("cel: eval: not init")
	}

	var data any
	data = make(map[string]any)
	if err := json.Unmarshal(obj, &data); err != nil {
		data = string(obj)
	}

	eval, _, err := p.prog.Eval(map[string]any{varNameSelf: data})
	if err != nil {
		return nil, fmt.Errorf("cel: eval: %v", err)
	}

	switch eV := eval.Value().(type) {
	case bool,
		string,
		int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64, uintptr,
		float32, float64,
		complex64, complex128:
		return eV, nil
	case map[ref.Val]ref.Val:
		return json.Marshal(convert(eV))
	default:
		return json.Marshal(eV)
	}
}

func convert(src map[ref.Val]ref.Val) map[string]any {
	dst := make(map[string]any)
	for k, v := range src {
		switch t := v.Value().(type) {
		case map[ref.Val]ref.Val:
			dst[k.Value().(string)] = convert(t)
		default:
			dst[k.Value().(string)] = t
		}
	}
	return dst
}
