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

type ExampleOne struct {
	Value string
}

const example1Definition = `
funcs: ReturnsAStruct: #MockFunction & {
	returns: [#ExampleOne & {Value: "abba"}]
}
`

var _ = Describe("Example1", func() {
	gokka.RegisterType[ExampleOne]("ExampleOne")

	When("we make a mock", func() {
		It("should return abba", func() {
			mock, err := gokka.NewMock(example1Definition)
			Expect(err).ToNot(HaveOccurred())
			Expect(mock).ToNot(BeNil())

			ret, err := gokka.Exec1[ExampleOne](mock, "ReturnsAStruct")
			Expect(err).ToNot(HaveOccurred())
			Expect(ret).To(Equal(ExampleOne{"abba"}))
		})
	})
})
