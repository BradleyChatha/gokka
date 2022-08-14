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
	"github.com/hashicorp/go-multierror"
)

type Mock struct {
	funcs     map[string]cue.Value
	callCount map[string][]uint64
	scope     cue.Value
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

vars: [string]: _

funcs: [string]: #MockFunction | [...#MockFunction]
	`)

	globalContext.scope = builtin
}

func (m *Mock) injectVar(as string, value any) {
	cueValue := globalContext.cueContext.Encode(value)

	valueNode := cueValue.Syntax()
	switch expr := valueNode.(type) {
	case ast.Expr:
		node := ast.NewStruct(
			&ast.Field{
				Label: ast.NewIdent(as),
				Value: expr,
			},
		)

		parent := ast.NewStruct(
			&ast.Field{
				Label: ast.NewIdent("vars"),
				Value: node,
			},
		)

		value := globalContext.cueContext.BuildExpr(parent)
		m.scope = m.scope.Unify(value)
	default:
		panic(fmt.Errorf("type %T is not an expression. This is 100%% a bug, please report", valueNode))
	}
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
	return NewMockWithVars(cueCode, nil)
}

func NewMockWithVars(cueCode string, vars map[string]any) (*Mock, error) {
	mock := &Mock{
		funcs:     make(map[string]cue.Value),
		callCount: make(map[string][]uint64),
		scope:     globalContext.scope,
	}

	for as, value := range vars {
		mock.injectVar(as, value)
	}

	compiled := globalContext.cueContext.CompileString(cueCode, cue.Scope(mock.scope))
	if compiled.Err() != nil {
		return nil, compiled.Err()
	}

	if strct, err := compiled.LookupPath(cue.ParsePath("funcs")).Struct(); err == nil {
		iter := strct.Fields()
		for iter.Next() {
			if iter.Selector().IsDefinition() || iter.Selector().PkgPath() != "" {
				continue
			}

			overloadCount, _ := iter.Value().Len().Int64()
			mock.funcs[iter.Label()] = iter.Value()
			mock.callCount[iter.Label()] = make([]uint64, overloadCount+1)
		}
		return mock, nil
	} else {
		return nil, err
	}
}

func execFunction(count []uint64, f cue.Value, args ...any) ([]cue.Value, error) {
	var mainErr error
	var overloads []cue.Value

	if l, err := f.List(); err == nil {
		for l.Next() {
			overloads = append(overloads, l.Value())
		}
	} else {
		overloads = append(overloads, f)
	}

	for i, overload := range overloads {
		count[i]++
		maxCalls, err := overload.LookupPath(cue.ParsePath("maxCalls")).Uint64()
		if err != nil {
			return nil, err
		}

		if count[i] > maxCalls && count[i] < math.MaxUint64 {
			return nil, fmt.Errorf("max calls exceeded")
		}

		argCount, _ := overload.LookupPath(cue.ParsePath("args")).Len().Int64()
		if int(argCount) != len(args) {
			mainErr = multierror.Append(mainErr, fmt.Errorf("expected %d arguments, got %d", argCount, len(args)))
			continue
		}

		for i, arg := range args {
			argValue := globalContext.cueContext.Encode(arg)
			overload = overload.FillPath(cue.ParsePath(fmt.Sprintf("args[%d]", i)), argValue)
		}

		err = overload.Validate(cue.Concrete(true))
		if err != nil {
			mainErr = multierror.Append(mainErr, err)
			continue
		}

		retValues, err := overload.LookupPath(cue.ParsePath("returns")).List()
		if err != nil {
			return nil, err
		}

		var ret []cue.Value
		for retValues.Next() {
			ret = append(ret, retValues.Value())
		}

		return ret, nil
	}

	return nil, multierror.Append(mainErr, fmt.Errorf("no overloads matched"))
}
