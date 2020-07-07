/**
 * Copyright 2020 Appvia Ltd <info@appvia.io>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package jsonutils_test

import (
	"github.com/appvia/kore/pkg/utils/jsonutils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CompactWithOrderedMapKeys", func() {

	It("should compact null as is", func() {
		res, err := jsonutils.Compact([]byte(` null `))
		Expect(err).ToNot(HaveOccurred())
		Expect(string(res)).To(Equal(`null`))
	})

	It("should compact a literal value as is", func() {
		res, err := jsonutils.Compact([]byte(` "foo" `))
		Expect(err).ToNot(HaveOccurred())
		Expect(string(res)).To(Equal(`"foo"`))
	})

	It("should compact an array and preserve the order", func() {
		res, err := jsonutils.Compact([]byte(` [ 3, 2, 1 ] `))
		Expect(err).ToNot(HaveOccurred())
		Expect(string(res)).To(Equal(`[3,2,1]`))
	})

	It("should compact an empty map", func() {
		res, err := jsonutils.Compact([]byte(` { } `))
		Expect(err).ToNot(HaveOccurred())
		Expect(string(res)).To(Equal(`{}`))
	})

	It("should compact a map with sorted keys", func() {
		res, err := jsonutils.Compact([]byte(` { "c": 3, "b": 2, "a": 1 } `))
		Expect(err).ToNot(HaveOccurred())
		Expect(string(res)).To(Equal(`{"a":1,"b":2,"c":3}`))
	})

	It("should compact a complex value", func() {
		res, err := jsonutils.Compact([]byte(`{
			"b": {
				"b2": "b2v",
				"b1": "b1v"
			},
			"a": {
				"a1": {
					"a12": "a12v",
					"a11": "a11v"
				}
			}
		}`))
		Expect(err).ToNot(HaveOccurred())
		Expect(string(res)).To(Equal(`{"a":{"a1":{"a11":"a11v","a12":"a12v"}},"b":{"b1":"b1v","b2":"b2v"}}`))
	})

})
