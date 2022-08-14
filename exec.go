/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 * Author: Bradley Chatha
 */
package gokka

import (
	"errors"
	"fmt"

	"cuelang.org/go/cue"
)

func decode(value cue.Value, ret any) error {
	if ptr, ok := ret.(*error); ok {
		if value.Null() == nil {
			*ptr = nil
			return nil
		}

		msg := value.LookupPath(cue.ParsePath("Error"))
		if err := msg.Err(); err != nil {
			panic(err)
		}

		str, err := msg.String()
		if err != nil {
			panic(err)
		}

		*ptr = errors.New(str)
		return nil
	} else {
		return value.Decode(ret)
	}
}

func Exec1[T any](mock *Mock, name string, args ...any) (T, error) {
	var ret1 T

	f, ok := mock.funcs[name]
	if !ok {
		return ret1, errors.New("function not found")
	}

	ret, err := execFunction(mock.callCount[name], f, args...)
	if err != nil {
		return ret1, err
	}

	if len(ret) != 1 {
		return ret1, fmt.Errorf("expected 1 return value, got %d", len(ret))
	}

	err = decode(ret[0], &ret1)
	return ret1, err
}

func Exec2[T1 any, T2 any](mock *Mock, name string, args ...any) (T1, T2, error) {
	var ret1 T1
	var ret2 T2

	f, ok := mock.funcs[name]
	if !ok {
		return ret1, ret2, errors.New("function not found")
	}

	ret, err := execFunction(mock.callCount[name], f, args...)
	if err != nil {
		return ret1, ret2, err
	}

	if len(ret) != 2 {
		return ret1, ret2, fmt.Errorf("expected 2 return values, got %d", len(ret))
	}

	err = decode(ret[0], &ret1)
	if err != nil {
		return ret1, ret2, err
	}

	err = decode(ret[1], &ret2)
	return ret1, ret2, err
}

func Exec3[T1 any, T2 any, T3 any](mock *Mock, name string, args ...any) (T1, T2, T3, error) {
	var ret1 T1
	var ret2 T2
	var ret3 T3

	f, ok := mock.funcs[name]
	if !ok {
		return ret1, ret2, ret3, errors.New("function not found")
	}

	ret, err := execFunction(mock.callCount[name], f, args...)
	if err != nil {
		return ret1, ret2, ret3, err
	}

	if len(ret) != 3 {
		return ret1, ret2, ret3, fmt.Errorf("expected 3 return values, got %d", len(ret))
	}

	err = decode(ret[0], &ret1)
	if err != nil {
		return ret1, ret2, ret3, err
	}

	err = decode(ret[1], &ret2)
	if err != nil {
		return ret1, ret2, ret3, err
	}

	err = decode(ret[2], &ret3)
	return ret1, ret2, ret3, err
}

func MustExec1[T any](mock *Mock, name string, args ...any) T {
	ret, err := Exec1[T](mock, name, args...)
	if err != nil {
		panic(err)
	}
	return ret
}

func MustExec2[T1 any, T2 any](mock *Mock, name string, args ...any) (T1, T2) {
	r1, r2, err := Exec2[T1, T2](mock, name, args...)
	if err != nil {
		panic(err)
	}
	return r1, r2
}

func MustExec3[T1 any, T2 any, T3 any](mock *Mock, name string, args ...any) (T1, T2, T3) {
	r1, r2, r3, err := Exec3[T1, T2, T3](mock, name, args...)
	if err != nil {
		panic(err)
	}
	return r1, r2, r3
}
