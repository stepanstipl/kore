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
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// CopyWithMeta will copy pertenant metadata objects
func CopyWithMeta(objIn runtime.Object) (runtime.Object, error) {
	objCopy := objIn.DeepCopyObject()
	objCopyMeta, err := meta.Accessor(objCopy)
	if err != nil {
		return objCopy, err
	}
	objInMeta, err := meta.Accessor(objIn)
	if err != nil {
		return objCopy, err
	}
	// Update all metadata on the copied object
	objCopyMeta.SetName(objInMeta.GetName())
	objCopyMeta.SetNamespace(objInMeta.GetNamespace())
	objCopyMeta.SetLabels(objInMeta.GetLabels())
	return objCopy, nil
}

// GetMeta will get a metadata object from a runtime.object
func GetMeta(obj runtime.Object) (metav1.ObjectMeta, error) {
	objMeta := metav1.ObjectMeta{}
	accessor, err := meta.Accessor(obj)
	if err != nil {
		return objMeta, err
	}
	objMeta.Name = accessor.GetName()
	objMeta.Namespace = accessor.GetNamespace()
	objMeta.Labels = accessor.GetLabels()
	return objMeta, nil
}
