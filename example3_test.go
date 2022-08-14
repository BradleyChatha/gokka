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

const example3Definition = `
funcs: UseInjectedVar: #MockFunction & {
	returns: [vars.injected]
}
`

var _ = Describe("Example3", func() {
	When("we make a mock", func() {
		It("should use our injected variable", func() {
			mock, err := gokka.NewMockWithVars(example3Definition, map[string]any{"injected": "value"})
			Expect(err).ToNot(HaveOccurred())
			Expect(mock).ToNot(BeNil())

			ret, err := gokka.Exec1[string](mock, "UseInjectedVar")
			Expect(err).ToNot(HaveOccurred())
			Expect(ret).To(Equal("value"))
		})
	})
})
