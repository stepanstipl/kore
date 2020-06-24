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

package kubernetes

import (
	"github.com/appvia/kore/pkg/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// EnsureOwnerReference ensures the given owner is present in the object's owner references list
func EnsureOwnerReference(object Object, owner Object, blockOwnerDeletion bool) {
	if HasOwnerReference(object, owner) {
		return
	}

	object.SetOwnerReferences(append(object.GetOwnerReferences(), metav1.OwnerReference{
		APIVersion:         owner.GetObjectKind().GroupVersionKind().GroupVersion().String(),
		Kind:               owner.GetObjectKind().GroupVersionKind().Kind,
		Name:               owner.GetName(),
		UID:                owner.GetUID(),
		BlockOwnerDeletion: utils.BoolPtr(blockOwnerDeletion),
	}))
}

func HasOwnerReference(object Object, owner Object) bool {
	for _, o := range object.GetOwnerReferences() {
		if o.UID == owner.GetUID() {
			return true
		}
	}

	return false
}

func HasOwnerReferenceWithKind(object Object, gvk schema.GroupVersionKind) bool {
	for _, o := range object.GetOwnerReferences() {
		if o.APIVersion == gvk.GroupVersion().String() && o.Kind == gvk.Kind {
			return true
		}
	}

	return false
}
