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

package persistence_test

import (
	"testing"

	"github.com/appvia/kore/pkg/persistence"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var ginkgoTestContext *testing.T

func TestPersistence(t *testing.T) {
	ginkgoTestContext = t
	RegisterFailHandler(Fail)
	RunSpecs(t, "Persistence Suite")
}

func getTestStore() persistence.Interface {
	return makeTestStore(ginkgoTestContext)
}
