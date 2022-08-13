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
)

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

	err = ret[0].Decode(&ret1)
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

	err = ret[0].Decode(&ret1)
	if err != nil {
		return ret1, ret2, err
	}

	err = ret[1].Decode(&ret2)
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

	err = ret[0].Decode(&ret1)
	if err != nil {
		return ret1, ret2, ret3, err
	}

	err = ret[1].Decode(&ret2)
	if err != nil {
		return ret1, ret2, ret3, err
	}

	err = ret[2].Decode(&ret3)
	return ret1, ret2, ret3, err
}
