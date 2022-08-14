/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 * Author: Bradley Chatha
 */
package gokka_test

import (
	"testing"

	"github.com/BradleyChatha/gokka"
)

type Input struct {
	Name string
}

type MyInterface interface {
	HasName(input Input, oneOf []string) (bool, error)
}

const myInterfaceCue = `
import "list" // Cue has a built in standard library

funcs: HasName: [
    // Gokka will match overloads from top to bottom, 
    // and uses the first one that doesn't cause an evaluation error

    #MockFunction & {
        args: [
            #InterfaceInput & { Name: !="" },
            [] // Matches only an empty array
        ]
        returns: [
            false,
            #GoError & {
                Error: "The list provided is empty!"
            }
        ]
    },

    #MockFunction & {
        args: [
            #InterfaceInput & { Name: !="" },
            [...string] // Matches a string array of any length
        ]
        returns: [
            list.Contains(args[1], args[0].Name),
            null
        ]
    }
]
`

type MyInterfaceMock struct {
	mock *gokka.Mock
}

func NewMyInterfaceMock() *MyInterfaceMock {
	mock, err := gokka.NewMock(myInterfaceCue)
	if err != nil {
		panic(err)
	}
	return &MyInterfaceMock{
		mock: mock,
	}
}

func (m *MyInterfaceMock) HasName(input Input, oneOf []string) (bool, error) {
	return gokka.MustExec2[bool, error](m.mock, "HasName", input, oneOf)
}

func TestMyInterfaceMock(t *testing.T) {
	gokka.RegisterType[Input]("InterfaceInput")

	mock := NewMyInterfaceMock()

	res, err := mock.HasName(Input{Name: "foo"}, []string{"foo", "bar"})
	if err != nil {
		t.Error(err)
	}
	if !res {
		t.Error("Expected true, got false")
	}

	res, err = mock.HasName(Input{Name: "baz"}, []string{"foo", "bar"})
	if err != nil {
		t.Error(err)
	}
	if res {
		t.Error("Expected false, got true")
	}

	_, err = mock.HasName(Input{Name: "This should error out because empty list"}, []string{})
	if err == nil {
		t.Error("Expected error")
	}

	// Use something other than the testing module to catch this panic.
	// Note: The cue we define below forbids `Name` from being empty in both overloads.
	//       This means no overload will be called, since they all cause an evaluation error.
	//       Because our mock uses `MustExec2` the `Must` means it'll panic upon evaluation error.
	//       Not the most ideal, but the best I can come up with for now.
	// mock.HasName(Input{Name: ""}, []string{"This should panic because of an evaluation error"})
}
