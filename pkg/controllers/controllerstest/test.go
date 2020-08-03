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

package controllerstest

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/appvia/kore/pkg/kore"

	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/appvia/kore/pkg/controllers"
	"github.com/appvia/kore/pkg/controllers/controllersfakes"
	"github.com/appvia/kore/pkg/kore/korefakes"
	"github.com/appvia/kore/pkg/utils/kubernetes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
)

const LabelGetError = "testing.kore.appvia.io/get-error"

type Test struct {
	context      context.Context
	cancel       context.CancelFunc
	Context      kore.Context
	Client       *controllersfakes.FakeClient
	StatusClient *controllersfakes.FakeStatusWriter
	Objects      []kubernetes.Object
	Controller   *controllersfakes.FakeController
	Kore         *korefakes.FakeInterface
	Logger       *logrus.Logger
}

func NewTest(ctx context.Context) *Test {
	test := &Test{}

	test.context, test.cancel = context.WithCancel(ctx)

	test.StatusClient = &controllersfakes.FakeStatusWriter{}
	test.Client = &controllersfakes.FakeClient{}
	test.Client.StatusReturns(test.StatusClient)
	test.Client.GetStub = func(_ context.Context, name types.NamespacedName, object runtime.Object) error {
		for _, o := range test.Objects {
			if reflect.TypeOf(o) == reflect.TypeOf(object) {
				if o.GetName() == name.Name && o.GetNamespace() == name.Namespace {
					if o.GetLabels()[LabelGetError] != "" {
						return errors.New(o.GetLabels()[LabelGetError])
					}
					res := o.DeepCopyObject()

					// The runtime client doesn't set the GVK on the result object, so we shouldn't either
					res.GetObjectKind().SetGroupVersionKind(schema.EmptyObjectKind.GroupVersionKind())

					reflect.ValueOf(object).Elem().Set(reflect.ValueOf(res).Elem())
					return nil
				}
			}
		}

		gr := schema.GroupResource{
			Group:    object.GetObjectKind().GroupVersionKind().Group,
			Resource: object.GetObjectKind().GroupVersionKind().Kind,
		}
		return kerrors.NewNotFound(gr, name.Name)
	}
	test.Controller = &controllersfakes.FakeController{}
	test.Kore = &korefakes.FakeInterface{}
	test.Logger = logrus.New()
	test.Logger.Out = GinkgoWriter

	test.Context = kore.NewContext(test.context, test.Logger, test.Client, test.Kore)

	return test
}

func (t *Test) Initialize(c controllers.Interface2) {
	err := c.Initialize(kore.NewContext(t.Context, t.Logger, t.Client, t.Kore), t.Controller)
	Expect(err).ToNot(HaveOccurred())
}

func (t *Test) Stop() {
	t.cancel()
}

func (t *Test) ExpectCreate(i int, obj interface{}) {
	t.checkCallCount(i, t.Client.CreateCallCount(), "create")
	_, createdObj, _ := t.Client.CreateArgsForCall(i)
	Expect(obj).To(BeAssignableToTypeOf(obj))
	reflect.ValueOf(obj).Elem().Set(reflect.ValueOf(createdObj).Elem())
}

func (t *Test) ExpectUpdate(i int, obj interface{}) {
	t.checkCallCount(i, t.Client.UpdateCallCount(), "update")
	_, updatedObj, _ := t.Client.UpdateArgsForCall(i)
	Expect(obj).To(BeAssignableToTypeOf(obj))
	reflect.ValueOf(obj).Elem().Set(reflect.ValueOf(updatedObj).Elem())
}

func (t *Test) ExpectDelete(i int, obj interface{}) {
	t.checkCallCount(i, t.Client.DeleteCallCount(), "delete")
	_, deletedObj, _ := t.Client.DeleteArgsForCall(i)
	Expect(obj).To(BeAssignableToTypeOf(obj))
	reflect.ValueOf(obj).Elem().Set(reflect.ValueOf(deletedObj).Elem())
}

func (t *Test) ExpectPatch(i int, obj interface{}) client.Patch {
	t.checkCallCount(i, t.Client.PatchCallCount(), "patch")
	_, patchedObj, patch, _ := t.Client.PatchArgsForCall(i)
	Expect(obj).To(BeAssignableToTypeOf(obj))
	reflect.ValueOf(obj).Elem().Set(reflect.ValueOf(patchedObj).Elem())
	return patch
}

func (t *Test) ExpectStatusUpdate(i int, obj interface{}) {
	t.checkCallCount(i, t.StatusClient.UpdateCallCount(), "statusUpdate")
	_, updatedObj, _ := t.StatusClient.UpdateArgsForCall(0)
	Expect(obj).To(BeAssignableToTypeOf(obj))
	reflect.ValueOf(obj).Elem().Set(reflect.ValueOf(updatedObj).Elem())
}

func (t *Test) ExpectStatusPatch(i int, obj interface{}) client.Patch {
	t.checkCallCount(i, t.StatusClient.PatchCallCount(), "statusPatch")
	_, patchedObj, patch, _ := t.StatusClient.PatchArgsForCall(0)
	Expect(obj).To(BeAssignableToTypeOf(obj))
	reflect.ValueOf(obj).Elem().Set(reflect.ValueOf(patchedObj).Elem())
	return patch
}

func (t *Test) checkCallCount(i, callCount int, method string) {
	if callCount < i+1 {
		Fail(fmt.Sprintf("less than %d %s call(s) happened", i+1, method))
	}
}

func (t *Test) ExpectRequeue(res reconcile.Result, err error) {
	// We shouldn't use the HaveOccurred() matcher because it returns true if the value is a nil concrete value as an interface
	if err != nil {
		Fail(fmt.Sprintf("reconcile error is not nil: (%T) %v", err, err))
	}
	if !res.Requeue && res.RequeueAfter == 0 {
		Fail("was expecting Requeue = true or RequeueAfter > 0")
	}
}
