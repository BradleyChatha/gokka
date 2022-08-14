# Overview

Gokka is a small mocking library for Go. It requires no external binary, and uses the expressiveness of
[CUE](https://cuelang.org/) to perform its validation and mocking.

# Getting started

First, install the package:

```bash
go get -u github.com/BradleyChatha/gokka
```

Second, create your interface to mock:

```go
// interface.go
package interface

type Input struct {
    Name string
}

type MyInterface interface {
    HasName(input Input, oneOf []string) (bool, error)
}
```

Third, create your mock, which simply calls into a gokka.Mock:

```go
// interface_test.go
package interface_test

import (
    "github.com/BradleyChatha/gokka"
    "me.local/interface"
)

const myInterfaceCue = `...[DEFINED FURTHER BELOW]`

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

func (m *MyInterfaceMock) HasName(input interface.Input, oneOf []string) (bool, error) {
    return gokka.MustExec2[bool, error](m.mock, "HasName", input, oneOf)
}
```

Fourth, create your test, and register any Go types you need to make use of:

```go
...

func TestMyInterfaceMock(t *testing.T) {
    gokka.RegisterType[interface.Input]("InterfaceInput")

    mock := NewMyInterfaceMock()

    res, err := mock.HasName(interface.Input{Name: "foo"}, []string{"foo", "bar"})
    if err != nil {
        t.Error(err)
    }
    if !res {
        t.Error("Expected true, got false")
    }

    res, err = mock.HasName(interface.Input{Name: "baz"}, []string{"foo", "bar"})
    if err != nil {
        t.Error(err)
    }
    if res {
        t.Error("Expected false, got true")
    }

    _, err = mock.HasName(interface.Input{Name: "This should error out because empty list"}, []string{})
    if err != nil {
        t.Error(err)
    }

    // Use something other than the testing module to catch this panic.
    // Note: The cue we define below forbids `Name` from being empty in both overloads.
    //       This means no overload will be called, since they all cause an evaluation error.
    //       Because our mock uses `MustExec2` the `Must` means it'll panic upon evaluation error.
    //       Not the most ideal, but the best I can come up with for now.
    mock.HasName(interface.Input{Name: ""}, []string{"This should panic because of an evaluation error"})
}
```

Finally, define your cue mock:

```go
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
```