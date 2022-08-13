/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 * Author: Bradley Chatha
 */
package gokka_test

import (
	"github.com/BradleyChatha/gokka"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type ExampleZero struct {
	Value string
}

const validDefinition = `
funcs: Expect: #MockFunction & {
	args: [
		#ExampleZero,
		"a" | "b"
	]

	returns: [args[0].Value == args[1]]
}
`

var _ = Describe("Example0", func() {
	gokka.RegisterType[ExampleZero]("ExampleZero")

	When("we make a mock", func() {
		It("should not error out with no schema", func() {
			mock, err := gokka.NewMock("")
			Expect(err).ToNot(HaveOccurred())
			Expect(mock).ToNot(BeNil())
		})

		It("should error if we violate the builtin schema", func() {
			mock, err := gokka.NewMock("funcs: 200")
			Expect(err).To(HaveOccurred())
			Expect(mock).To(BeNil())
		})

		It("shouldn't error when accessing a builtin definition", func() {
			mock, err := gokka.NewMock("#Abc: #MockFunction")
			Expect(err).ToNot(HaveOccurred())
			Expect(mock).ToNot(BeNil())
		})

		It("shouldn't error when accessing a registered definition", func() {
			mock, err := gokka.NewMock("ExampleZero: #ExampleZero")
			Expect(err).ToNot(HaveOccurred())
			Expect(mock).ToNot(BeNil())
		})

		It("shouldn't error when passing in a fully valid definition", func() {
			mock, err := gokka.NewMock(validDefinition)
			Expect(err).ToNot(HaveOccurred())
			Expect(mock).ToNot(BeNil())
		})

		It("should return the correct values when calling our mock", func() {
			mock, err := gokka.NewMock(validDefinition)
			Expect(err).ToNot(HaveOccurred())
			Expect(mock).ToNot(BeNil())

			v, err := gokka.Exec1[bool](mock, "Expect", ExampleZero{Value: "a"}, "a")
			Expect(err).ToNot(HaveOccurred())
			Expect(v).To(BeTrue())

			v, err = gokka.Exec1[bool](mock, "Expect", ExampleZero{Value: "b"}, "b")
			Expect(err).ToNot(HaveOccurred())
			Expect(v).To(BeTrue())

			v, err = gokka.Exec1[bool](mock, "Expect", ExampleZero{Value: "a"}, "b")
			Expect(err).ToNot(HaveOccurred())
			Expect(v).To(BeFalse())
		})

		It("should error if we validate the user defined schema", func() {
			mock, err := gokka.NewMock(validDefinition)
			Expect(err).ToNot(HaveOccurred())
			Expect(mock).ToNot(BeNil())

			_, err = gokka.Exec1[bool](mock, "Expect", ExampleZero{Value: "invalid"}, "invalid")
			Expect(err).To(HaveOccurred())
		})
	})
})
