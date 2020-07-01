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

package kubernetes_test

import (
	"github.com/appvia/kore/pkg/utils/kubernetes"
	"github.com/appvia/kore/pkg/utils/kubernetes/kubernetesfakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func fakeObject(name string) *kubernetesfakes.FakeObject {
	o := &kubernetesfakes.FakeObject{}
	o.GetObjectKindReturns(schema.EmptyObjectKind)
	o.GetNamespaceReturns("testns")
	o.GetNameReturns(name)
	return o
}

func objectIndex(o kubernetes.Object, res []kubernetes.Object) int {
	for i, e := range res {
		if e == o {
			return i
		}
	}
	return -1
}

var _ = Describe("DepencencyResolver", func() {

	It("should return a single node as is", func() {
		n1 := fakeObject("n1")
		resolver := kubernetes.NewDependencyResolver()
		resolver.AddNode(n1)

		res, err := resolver.Resolve()
		Expect(err).ToNot(HaveOccurred())
		Expect(res).To(Equal([]kubernetes.Object{n1}))
	})

	It("should return nodes in strict dependency order", func() {
		n1 := fakeObject("n1")
		n2 := fakeObject("n2")
		n3 := fakeObject("n3")
		resolver := kubernetes.NewDependencyResolver()
		resolver.AddNode(n1, n2)
		resolver.AddNode(n2, n3)
		resolver.AddNode(n3)

		res, err := resolver.Resolve()
		Expect(err).ToNot(HaveOccurred())
		Expect(res).To(Equal([]kubernetes.Object{n3, n2, n1}))
	})

	It("should handle multiple nodes depend on the same node", func() {
		n1 := fakeObject("n1")
		n2 := fakeObject("n2")
		n3 := fakeObject("n3")
		resolver := kubernetes.NewDependencyResolver()
		resolver.AddNode(n1)
		resolver.AddNode(n2)
		resolver.AddNode(n3, n1, n2)

		res, err := resolver.Resolve()
		Expect(err).ToNot(HaveOccurred())
		Expect(objectIndex(n1, res)).To(BeNumerically("<", objectIndex(n3, res)))
		Expect(objectIndex(n2, res)).To(BeNumerically("<", objectIndex(n3, res)))
	})

	It("should handle multiple distinct dependency graphs", func() {
		n1 := fakeObject("n1")
		n2 := fakeObject("n2")
		n3 := fakeObject("n3")
		n4 := fakeObject("n4")
		resolver := kubernetes.NewDependencyResolver()
		resolver.AddNode(n1, n2)
		resolver.AddNode(n2)
		resolver.AddNode(n3, n4)
		resolver.AddNode(n4)

		res, err := resolver.Resolve()
		Expect(err).ToNot(HaveOccurred())
		Expect(objectIndex(n2, res)).To(BeNumerically("<", objectIndex(n1, res)))
		Expect(objectIndex(n4, res)).To(BeNumerically("<", objectIndex(n3, res)))
	})

	It("should return an error on a circular reference", func() {
		n1 := fakeObject("n1")
		n2 := fakeObject("n2")
		n3 := fakeObject("n3")
		resolver := kubernetes.NewDependencyResolver()
		resolver.AddNode(n1, n2)
		resolver.AddNode(n2, n3)
		resolver.AddNode(n3, n1)

		_, err := resolver.Resolve()
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("circular reference detected"))
	})

	It("should return an error when depending on a non-registered object", func() {
		n1 := fakeObject("n1")
		n2 := fakeObject("n2")
		resolver := kubernetes.NewDependencyResolver()
		resolver.AddNode(n1, n2)

		_, err := resolver.Resolve()
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("dependency not found"))
	})

	It("should panic if the same node is added multiple times", func() {
		n1 := fakeObject("n1")
		resolver := kubernetes.NewDependencyResolver()
		resolver.AddNode(n1)

		Expect(func() { resolver.AddNode(n1) }).To(Panic())
	})
})
