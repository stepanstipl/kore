/**
 * Copyright (C) 2020 Appvia Ltd <info@appvia.io>
 *
 * This file is part of kore-apiserver.
 *
 * kore-apiserver is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 2 of the License, or
 * (at your option) any later version.
 *
 * kore-apiserver is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with kore-apiserver.  If not, see <http://www.gnu.org/licenses/>.
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
