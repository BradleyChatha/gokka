# Overview

Gokka is a small mocking library for Go. It requires no external binary, and uses the expressiveness of
[CUE](https://cuelang.org/) to perform its validation and mocking.

- [Overview](#overview)
  - [Features](#features)
  - [Getting started](#getting-started)
  - [Examples](#examples)
  - [LICENSE](#license)

## Features

- Use the power of CUElang to perform validation; data templating, and more!
- Doesn't require an external binary.
- Easy interface to import Golang types into CUE, to keep everything statically typed.
- Each mocked function can:
  - Have multiple overloads, to customise the output based on any given input.
  - Specify the maximum amount of times it can be called before it returns an error (this is specified per overload currently).
- Uses generics to provide a sane interface to the developer.

## Getting started

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

    _, err = mock.HasName(Input{Name: "This should error out because empty list"}, []string{})
    if err == nil {
        t.Error("Expected error")
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

## Examples

All of the test cases also serve as examples, so please view the `exampleX_test.go` files.

## LICENSE

Most of the files in this repo are covered by the Mozilla Public License v 2.0, which you can find a copy of here: https://mozilla.org/MPL/2.0/

Please read this FAQ on how to handle source code with this license: https://www.mozilla.org/en-US/MPL/2.0/FAQ/

The major, meaty points:

* You must make available any user that has access to the distribution of your project any MPLv2 files, either modified or unmodified.
    * For unmodified code, making the `NOTICE.md` file available should suffice.
    * For modified code, the entire modified file, regardless of whether it has "proprietary" code (which cannot exist under an MPLv2 file by definition), must be made freely available to the user(s) with access to the distributed project.
    * This does not apply for projects that are not distributed directly to users, such as backend servers, even if those servers are publically accessible they are not "distributed" in the sense that MPLv2 applies itself to.
* You are free to statically and dynamically link this code into your project.
* You do **not** have to disclose anything else about your project. MPLv2 does not prevent you from keeping your proprietary code private, and it does not prevent you using a different license for your other files.

There's a bunch more points the FAQ can help you with, but the gist of this license is that you cannot modify the code without "giving back to the community". You **must** release in full any modified files that contain the MPLv2 license notification (which you cannot remove), but other than that it's a very permissive license, despite being copy-left.

This license encourages people to open MRs to improve the project, rather than keeping things to themselves, but does so in a way that it doesn't make it unviable for commercial use.