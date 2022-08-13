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

const example2Definition = `
funcs: CanOnlyCallOnce: #MockFunction & {
	maxCalls: 1
	returns: [true]
}
`

var _ = Describe("Example2", func() {
	When("we make a mock", func() {
		It("should only succeed once", func() {
			mock, err := gokka.NewMock(example2Definition)
			Expect(err).ToNot(HaveOccurred())
			Expect(mock).ToNot(BeNil())

			ret, err := gokka.Exec1[bool](mock, "CanOnlyCallOnce")
			Expect(err).ToNot(HaveOccurred())
			Expect(ret).To(BeTrue())

			ret, err = gokka.Exec1[bool](mock, "CanOnlyCallOnce")
			Expect(err).To(HaveOccurred())
			Expect(ret).To(BeFalse())
		})
	})
})
