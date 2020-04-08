package controllers_test

import (
	"errors"

	"github.com/appvia/kore/pkg/controllers"
	"github.com/appvia/kore/pkg/controllers/controllersfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
)

var _ = Describe("CompositeReconciler", func() {
	var c *controllers.CompositeReconciler
	var c1, c2 *controllersfakes.FakeComponent
	var components []controllers.Component
	var requeue bool
	var reconcileErr error

	BeforeEach(func() {
		logger := logrus.New()
		logger.Out = GinkgoWriter
		c = controllers.NewCompositeReconciler(logger)
		components = nil
	})

	JustBeforeEach(func() {
		for _, comp := range components {
			c.RegisterComponent(comp)
		}
	})

	Context("Reconcile", func() {
		JustBeforeEach(func() {
			requeue, reconcileErr = c.Reconcile()
		})

		When("there is an unfinished component", func() {
			BeforeEach(func() {
				c1 = &controllersfakes.FakeComponent{}
				c1.NameReturns("c1")
				c1.ReconcileReturns(true, nil)

				c2 = &controllersfakes.FakeComponent{}
				c2.NameReturns("c2")
				c2.ReconcileReturns(false, nil)

				components = append(components, c1)
			})

			It("should call reconcile", func() {
				Expect(c1.ReconcileCallCount()).To(Equal(1))
			})

			It("should requeue and no error", func() {
				Expect(requeue).To(BeTrue())
				Expect(reconcileErr).ToNot(HaveOccurred())
			})
		})

		When("all components finished", func() {
			BeforeEach(func() {
				c1 = &controllersfakes.FakeComponent{}
				c1.NameReturns("c1")
				c1.ReconcileReturns(false, nil)

				c2 = &controllersfakes.FakeComponent{}
				c2.NameReturns("c2")
				c2.ReconcileReturns(false, nil)

				components = append(components, c1, c2)
			})

			It("should not requeue and no error", func() {
				Expect(requeue).To(BeFalse())
				Expect(reconcileErr).ToNot(HaveOccurred())
			})
		})

		When("there is a critical error", func() {
			BeforeEach(func() {
				c1 = &controllersfakes.FakeComponent{}
				c1.NameReturns("c1")
				c1.ReconcileReturns(false, controllers.NewCriticalError(errors.New("critical error")))

				c2 = &controllersfakes.FakeComponent{}
				c2.NameReturns("c2")
				c2.ReconcileReturns(false, nil)

				components = append(components, c1, c2)
			})

			It("should not requeue and return the error", func() {
				Expect(requeue).To(BeFalse())
				Expect(reconcileErr).To(MatchError("critical error"))
			})
		})

		When("there are errors", func() {
			BeforeEach(func() {
				c1 = &controllersfakes.FakeComponent{}
				c1.NameReturns("c1")
				c1.ReconcileReturns(false, errors.New("error 1"))

				c2 = &controllersfakes.FakeComponent{}
				c2.NameReturns("c2")
				c2.ReconcileReturns(false, errors.New("error 2"))

				components = append(components, c1, c2)
			})

			It("should requeue and return the errors", func() {
				Expect(requeue).To(BeTrue())
				Expect(reconcileErr).To(HaveOccurred())
				Expect(reconcileErr.Error()).To(ContainSubstring("error 1"))
				Expect(reconcileErr.Error()).To(ContainSubstring("error 2"))
			})
		})

		When("there is a component depending on an other", func() {
			BeforeEach(func() {
				c1 = &controllersfakes.FakeComponent{}
				c1.NameReturns("c1")
				c1.ReconcileReturns(true, nil)

				c2 = &controllersfakes.FakeComponent{}
				c2.NameReturns("c2")
				c1.ReconcileReturns(true, nil)
				c2.DependenciesReturns([]string{"c1"})

				components = append(components, c1, c2)
			})

			It("should call reconcile only on c1", func() {
				Expect(c1.ReconcileCallCount()).To(Equal(1))
				Expect(c2.ReconcileCallCount()).To(Equal(0))
			})
		})

		When("there is a component depending on a finished component", func() {
			BeforeEach(func() {
				c1 = &controllersfakes.FakeComponent{}
				c1.NameReturns("c1")
				c1.ReconcileReturns(false, nil)

				c2 = &controllersfakes.FakeComponent{}
				c2.NameReturns("c2")
				c2.ReconcileReturns(true, nil)
				c2.DependenciesReturns([]string{"c1"})

				components = append(components, c1, c2)
			})

			JustBeforeEach(func() {
				_, _ = c.Reconcile()
			})

			It("should call c2 too", func() {
				Expect(c1.ReconcileCallCount()).To(Equal(1))
				Expect(c2.ReconcileCallCount()).To(BeNumerically(">=", 1))
			})
		})
	})

})
