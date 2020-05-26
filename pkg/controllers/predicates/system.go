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

package predicates

import (
	"github.com/appvia/kore/pkg/kore"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/event"
)

type SystemResourcePredicate struct {
}

func (s SystemResourcePredicate) Create(e event.CreateEvent) bool {
	return s.shouldProcess(e.Meta)
}

func (s SystemResourcePredicate) Delete(e event.DeleteEvent) bool {
	return s.shouldProcess(e.Meta)
}

func (s SystemResourcePredicate) Update(e event.UpdateEvent) bool {
	return s.shouldProcess(e.MetaOld)
}

func (s SystemResourcePredicate) Generic(e event.GenericEvent) bool {
	return s.shouldProcess(e.Meta)
}

func (s SystemResourcePredicate) shouldProcess(meta metav1.Object) bool {
	annotations := meta.GetAnnotations()
	if annotations[kore.AnnotationSystem] == kore.AnnotationValueTrue {
		return false
	}

	return true
}
