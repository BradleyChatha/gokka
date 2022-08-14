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

const example4Definition = `
funcs: Overloads: [#MockFunction & {
	args: ["a"]
	returns: ["a"]
}, #MockFunction & {
	args: ["b"]
	maxCalls: 1
	returns: ["b"]
}]
`

var _ = Describe("Example4", func() {
	When("we make a mock", func() {
		It("should succeed when calling an overload", func() {
			mock, err := gokka.NewMock(example4Definition)
			Expect(err).ToNot(HaveOccurred())
			Expect(mock).ToNot(BeNil())

			ret, err := gokka.Exec1[string](mock, "Overloads", "a")
			Expect(err).ToNot(HaveOccurred())
			Expect(ret).To(Equal("a"))

			ret, err = gokka.Exec1[string](mock, "Overloads", "b")
			Expect(err).ToNot(HaveOccurred())
			Expect(ret).To(Equal("b"))
		})

		It("should error when calling an overload that doesn't match", func() {
			mock, err := gokka.NewMock(example4Definition)
			Expect(err).ToNot(HaveOccurred())
			Expect(mock).ToNot(BeNil())

			ret, err := gokka.Exec1[string](mock, "Overloads", "c")
			Expect(err).To(HaveOccurred())
			Expect(ret).To(BeEmpty())
		})

		It("should respect an overload's max calls", func() {
			mock, err := gokka.NewMock(example4Definition)
			Expect(err).ToNot(HaveOccurred())
			Expect(mock).ToNot(BeNil())

			ret, err := gokka.Exec1[string](mock, "Overloads", "b")
			Expect(err).ToNot(HaveOccurred())
			Expect(ret).To(Equal("b"))
			ret, err = gokka.Exec1[string](mock, "Overloads", "b")
			Expect(err).To(HaveOccurred())
			Expect(ret).To(BeEmpty())

			ret, err = gokka.Exec1[string](mock, "Overloads", "a")
			Expect(err).ToNot(HaveOccurred())
			Expect(ret).To(Equal("a"))
		})
	})
})
