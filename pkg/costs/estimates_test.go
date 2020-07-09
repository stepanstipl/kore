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

package costs_test

import (
	"encoding/json"

	configv1 "github.com/appvia/kore/pkg/apis/config/v1"
	costsv1 "github.com/appvia/kore/pkg/apis/costs/v1beta1"
	"github.com/appvia/kore/pkg/costs"
	"github.com/appvia/kore/pkg/costs/costsfakes"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func eksSamplePlan() map[string]interface{} {
	return map[string]interface{}{
		"region": "eu-west-2",
		"nodeGroups": []map[string]interface{}{
			{
				"instanceType":     "t3.medium",
				"diskSize":         10,
				"name":             "group1",
				"desiredSize":      6,
				"minSize":          3,
				"maxSize":          30,
				"enableAutoscaler": true,
			},
		},
	}
}

func gkeSamplePlan() map[string]interface{} {
	return map[string]interface{}{
		"region": "europe-west2",
		"nodePools": []map[string]interface{}{
			{
				"machineType":      "n1-standard-2",
				"diskSize":         10,
				"name":             "compute1",
				"size":             5,
				"minSize":          2,
				"maxSize":          7,
				"enableAutoscaler": true,
			},
		},
	}
}

func encodePlan(config map[string]interface{}) []byte {
	e, _ := json.Marshal(config)
	return e
}

var _ = Describe("Costs - Estimates", func() {
	var e costs.Estimates
	var m costsfakes.FakeMetadata
	BeforeEach(func() {
		m = costsfakes.FakeMetadata{}
		m.KubernetesControlPlaneCostReturns(100000, nil)
		m.KubernetesExposedServiceCostReturns(25000, nil)
		e = costs.NewEstimates(&m)
	})
	When("GetClusterEstimate is called", func() {
		It("should return validation error if kind is not a supported cloud", func() {
			est, err := e.GetClusterEstimate(&configv1.PlanSpec{
				Kind: "unmanaged",
			})
			Expect(err).To(HaveOccurred())
			Expect(est).To(BeNil())
		})

		It("should return validation error if region is not specified on the plan", func() {
			config := eksSamplePlan()
			config["region"] = ""

			est, err := e.GetClusterEstimate(&configv1.PlanSpec{
				Kind: "EKS",
				Configuration: v1beta1.JSON{
					Raw: encodePlan(config),
				},
			})
			Expect(err).To(HaveOccurred())
			Expect(est).To(BeNil())
		})

		It("should return an estimate for a single EKS node pool cluster", func() {
			config := eksSamplePlan()
			m.InstanceTypeReturns(&costsv1.InstanceType{
				Prices: map[costsv1.PriceType]int64{
					costsv1.PriceTypeOnDemand: 25000,
				},
			}, nil)

			est, err := e.GetClusterEstimate(&configv1.PlanSpec{
				Kind: "EKS",
				Configuration: v1beta1.JSON{
					Raw: encodePlan(config),
				},
			})
			Expect(err).ToNot(HaveOccurred())
			Expect(est).ToNot(BeNil())
			Expect(est.CostElements[2].Name).To(Equal("Node Pool group1"))
			Expect(est.CostElements[2].MinCost).To(Equal(int64(3 * 25000)))
			Expect(est.CostElements[2].MaxCost).To(Equal(int64(30 * 25000)))
			Expect(est.CostElements[2].TypicalCost).To(Equal(int64(6 * 25000)))
			// Expect overall cost to be node group cost + exposed service cost + k8s control plane cost
			Expect(est.MinCost).To(Equal(int64(100000 + 25000 + (3 * 25000))))
			Expect(est.MaxCost).To(Equal(int64(100000 + 25000 + (30 * 25000))))
			Expect(est.TypicalCost).To(Equal(int64(100000 + 25000 + (6 * 25000))))
		})

		It("should return an estimate for a single GKE node pool cluster", func() {
			config := gkeSamplePlan()
			m.InstanceTypeReturns(&costsv1.InstanceType{
				Prices: map[costsv1.PriceType]int64{
					costsv1.PriceTypeOnDemand: 35000,
				},
			}, nil)
			m.RegionZonesReturns([]string{"az1", "az2", "az3"}, nil)

			est, err := e.GetClusterEstimate(&configv1.PlanSpec{
				Kind: "GKE",
				Configuration: v1beta1.JSON{
					Raw: encodePlan(config),
				},
			})
			Expect(err).ToNot(HaveOccurred())
			Expect(est).ToNot(BeNil())
			Expect(est.CostElements[2].Name).To(Equal("Node Pool compute1"))
			// GKE deploys the count into each AZ, we deploy to all AZs, we've returned 3 AZs
			// above, hence multiply node pool costs by 3:
			Expect(est.CostElements[2].MinCost).To(Equal(int64(2 * 35000 * 3)))
			Expect(est.CostElements[2].MaxCost).To(Equal(int64(7 * 35000 * 3)))
			Expect(est.CostElements[2].TypicalCost).To(Equal(int64(5 * 35000 * 3)))

			// Expect overall cost to be node group cost + exposed service cost + k8s control plane cost
			Expect(est.MinCost).To(Equal(int64(100000 + 25000 + (2 * 35000 * 3))))
			Expect(est.MaxCost).To(Equal(int64(100000 + 25000 + (7 * 35000 * 3))))
			Expect(est.TypicalCost).To(Equal(int64(100000 + 25000 + (5 * 35000 * 3))))
		})

	})
})
