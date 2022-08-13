/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 * Author: Bradley Chatha
 */
package gokka

import (
	"fmt"
	"math"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/ast"
	"cuelang.org/go/cue/cuecontext"
)

type Mock struct {
	funcs     map[string]cue.Value
	callCount map[string]*uint64
}

type mockContext struct {
	cueContext *cue.Context
	scope      cue.Value
}

var globalContext *mockContext

func init() {
	globalContext = &mockContext{
		cueContext: cuecontext.New(),
	}

	builtin := globalContext.cueContext.CompileString(`
#GoError: {
	Error: string @go(s,string)
}

#MockFunction: {
	args: [..._]
	maxCalls: int | *18446744073709551615
	returns: [..._]
}

funcs: [string]: #MockFunction
	`)

	globalContext.scope = builtin
}

func RegisterType[T any](as string) {
	var x T
	value := globalContext.cueContext.EncodeType(x)

	// We don't want the type to be registered in the top-level,
	// so we hide it behind a definition, which requires a tiny amount of AST manip
	valueNode := value.Syntax()
	switch expr := valueNode.(type) {
	case ast.Expr:
		node := ast.NewStruct(
			&ast.Field{
				Label: ast.NewIdent("#" + as),
				Value: expr,
			},
		)
		value = globalContext.cueContext.BuildExpr(node)
	default:
		panic(fmt.Errorf("type %T is not an expression. This is 100%% a bug, please report", x))
	}

	globalContext.scope = globalContext.scope.Unify(value)
}

func NewMock(cueCode string) (*Mock, error) {
	mock := &Mock{
		funcs:     make(map[string]cue.Value),
		callCount: make(map[string]*uint64),
	}

	compiled := globalContext.cueContext.CompileString(cueCode, cue.Scope(globalContext.scope))
	if compiled.Err() != nil {
		return nil, compiled.Err()
	}

	if strct, err := compiled.LookupPath(cue.ParsePath("funcs")).Struct(); err == nil {
		iter := strct.Fields()
		for iter.Next() {
			if iter.Selector().IsDefinition() || iter.Selector().PkgPath() != "" {
				continue
			}

			mock.funcs[iter.Label()] = iter.Value()
			mock.callCount[iter.Label()] = new(uint64)
		}
		return mock, nil
	} else {
		return nil, err
	}
}

func execFunction(count *uint64, f cue.Value, args ...any) ([]cue.Value, error) {
	*count++
	maxCalls, err := f.LookupPath(cue.ParsePath("maxCalls")).Uint64()
	if err != nil {
		return nil, err
	}

	if *count > maxCalls && *count < math.MaxUint64 {
		return nil, fmt.Errorf("max calls exceeded")
	}

	argCount, _ := f.LookupPath(cue.ParsePath("args")).Len().Int64()
	if int(argCount) != len(args) {
		return nil, fmt.Errorf("expected %d arguments, got %d", argCount, len(args))
	}

	for i, arg := range args {
		argValue := globalContext.cueContext.Encode(arg)
		f = f.FillPath(cue.ParsePath(fmt.Sprintf("args[%d]", i)), argValue)
	}

	err = f.Validate(cue.Concrete(true))
	if err != nil {
		return nil, err
	}

	retValues, err := f.LookupPath(cue.ParsePath("returns")).List()
	if err != nil {
		return nil, err
	}

	var ret []cue.Value
	for retValues.Next() {
		ret = append(ret, retValues.Value())
	}

	return ret, nil
}
