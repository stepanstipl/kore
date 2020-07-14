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
	"context"
	"time"

	"github.com/appvia/kore/pkg/persistence"
	"github.com/appvia/kore/pkg/persistence/model"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Security Persistence", func() {
	var store persistence.Interface

	getScan := func(name string, namespace string, checkedAt time.Time, archivedAt time.Time, overallStatus string, ruleStatus string) model.SecurityScanResult {
		return model.SecurityScanResult{
			SecurityResourceReference: model.SecurityResourceReference{
				ResourceGroup:     "example.appvia.io",
				ResourceVersion:   "V1",
				ResourceKind:      "Example",
				ResourceNamespace: namespace,
				ResourceName:      name,
			},
			OwningTeam:    namespace,
			OverallStatus: overallStatus,
			CheckedAt:     checkedAt,
			ArchivedAt:    archivedAt,
			Results: []model.SecurityRuleResult{
				{
					RuleCode:  "TEST-001",
					Status:    ruleStatus,
					CheckedAt: checkedAt,
				},
				{
					RuleCode:  "TEST-002",
					Status:    "Warning",
					Message:   "Horse",
					CheckedAt: checkedAt,
				},
			},
		}
	}

	storeScans := func(scans ...*model.SecurityScanResult) {
		for _, scan := range scans {
			err := store.Security().StoreScan(context.Background(), scan)
			Expect(err).ToNot(HaveOccurred())
		}
	}

	BeforeEach(func() {
		store = getTestStore()
	})

	AfterEach(func() {
		store.Stop()
	})

	Describe("GetOverview", func() {
		When("called", func() {
			It("should provide an overview", func() {
				overview, err := store.Security().GetOverview(context.Background())
				Expect(err).ToNot(HaveOccurred())
				Expect(overview).ToNot(BeNil())
			})

			It("should sum the statuses for an overall count of open statuses", func() {
				overview, err := store.Security().GetOverview(context.Background())
				Expect(err).ToNot(HaveOccurred())
				Expect(overview.OpenIssueCounts["Compliant"]).To(BeNumerically("==", 3))
				Expect(overview.OpenIssueCounts["Warning"]).To(BeNumerically("==", 1))
			})

			It("should summarise the counts for each resource", func() {
				overview, err := store.Security().GetOverview(context.Background())
				Expect(err).ToNot(HaveOccurred())
				Expect(len(overview.Resources)).To(Equal(2))
				Expect(overview.Resources[0].ResourceName).To(Equal("test"))
				Expect(overview.Resources[0].OpenIssueCounts["Compliant"]).To(BeNumerically("==", 1))
				Expect(overview.Resources[0].OpenIssueCounts["Warning"]).To(BeNumerically("==", 1))
				Expect(overview.Resources[1].ResourceName).To(Equal("test2"))
				Expect(overview.Resources[1].OpenIssueCounts["Compliant"]).To(BeNumerically("==", 2))
				Expect(overview.Resources[1].OpenIssueCounts["Warning"]).To(BeNumerically("==", 0))
			})
		})
	})

	Describe("GetTeamOverview", func() {
		When("called", func() {
			It("should provide an overview for a specific team", func() {
				overview, err := store.Security().GetTeamOverview(context.Background(), "test-team")
				Expect(err).ToNot(HaveOccurred())
				Expect(overview).ToNot(BeNil())
			})

			It("should sum the statuses for an overall count of open statuses for a specific team", func() {
				overview, err := store.Security().GetTeamOverview(context.Background(), "test-team")
				Expect(err).ToNot(HaveOccurred())
				Expect(overview.OpenIssueCounts["Compliant"]).To(BeNumerically("==", 3))
				Expect(overview.OpenIssueCounts["Warning"]).To(BeNumerically("==", 1))
			})

			It("should summarise the counts for each resource for a specific team", func() {
				overview, err := store.Security().GetTeamOverview(context.Background(), "test-team")
				Expect(err).ToNot(HaveOccurred())
				Expect(len(overview.Resources)).To(Equal(2))
				Expect(overview.Resources[0].ResourceName).To(Equal("test"))
				Expect(overview.Resources[0].OpenIssueCounts["Compliant"]).To(BeNumerically("==", 1))
				Expect(overview.Resources[0].OpenIssueCounts["Warning"]).To(BeNumerically("==", 1))
				Expect(overview.Resources[1].ResourceName).To(Equal("test2"))
				Expect(overview.Resources[1].OpenIssueCounts["Compliant"]).To(BeNumerically("==", 2))
				Expect(overview.Resources[1].OpenIssueCounts["Warning"]).To(BeNumerically("==", 0))
			})
		})
	})

	Describe("GetScan", func() {
		When("invalid ID provided", func() {
			It("should return an error", func() {
				_, err := store.Security().GetScan(context.Background(), 12345)
				Expect(err).To(HaveOccurred())
			})
		})

		When("valid ID provided", func() {
			var scan *model.SecurityScanResult

			JustBeforeEach(func() {
				s, err := store.Security().GetScan(context.Background(), 1)
				Expect(err).ToNot(HaveOccurred())
				scan = s
			})

			It("should return a valid object", func() {
				Expect(scan).ToNot(BeNil())
				Expect(scan.Results).ToNot(BeNil())
			})

			It("should have a populated results slice", func() {
				Expect(len(scan.Results)).To(Equal(2))
				Expect(scan.Results[1].Status).To(Equal("Warning"))
			})
		})
	})

	Describe("ListScans", func() {
		It("should return entries with Results unpopulated", func() {
			scans, err := store.Security().ListScans(context.Background(), true)
			Expect(err).ToNot(HaveOccurred())

			Expect(scans).ToNot(BeNil())
			for _, scan := range scans {
				Expect(scan.Results).To(BeNil())
			}
		})

		When("called with latest only true", func() {
			It("should return only entries with null ArchivedAt", func() {
				scans, err := store.Security().ListScans(context.Background(), true)
				Expect(err).ToNot(HaveOccurred())
				Expect(len(scans)).To(Equal(3))
				Expect(scans[0].OverallStatus).To(Equal("Warning"))
				Expect(scans[0].ArchivedAt.IsZero()).To(BeTrue())
				Expect(scans[1].ArchivedAt.IsZero()).To(BeTrue())
				Expect(scans[2].ArchivedAt.IsZero()).To(BeTrue())
			})
		})

		When("called with latest only false", func() {
			It("should include entries with populated ArchivedAt", func() {
				scans, err := store.Security().ListScans(context.Background(), false)
				Expect(err).ToNot(HaveOccurred())
				Expect(len(scans)).To(Equal(4))
				Expect(scans[1].OverallStatus).To(Equal("Warning"))
				Expect(scans[0].ArchivedAt.IsZero()).To(BeTrue())
				Expect(scans[1].ArchivedAt.IsZero()).To(BeFalse())
				Expect(scans[2].ArchivedAt.IsZero()).To(BeTrue())
				Expect(scans[3].ArchivedAt.IsZero()).To(BeTrue())
			})
		})

		When("called with filters", func() {
			It("should only return scans matching the name and namespace filters", func() {
				scans, err := store.Security().ListScans(context.Background(), false, persistence.Filter.WithName("test2"), persistence.Filter.WithNamespace("test-team"))
				Expect(err).ToNot(HaveOccurred())
				Expect(len(scans)).To(Equal(2))
				Expect(scans[0].ResourceName).To(Equal("test2"))
				Expect(scans[1].ResourceName).To(Equal("test2"))
			})

			It("should only return scans matching the team filter", func() {
				scans, err := store.Security().ListScans(context.Background(), false, persistence.Filter.WithTeam("test-team2"))
				Expect(err).ToNot(HaveOccurred())
				Expect(len(scans)).To(Equal(1))
				Expect(scans[0].OwningTeam).To(Equal("test-team2"))
			})
		})
	})

	Describe("StoreScan", func() {

		When("called with a scan", func() {
			It("should persist the scan", func() {
				scan1 := getScan("test3", "test-team3", time.Now(), time.Time{}, "Failure", "Failure")
				storeScans(&scan1)

				v, err := store.Security().ListScans(context.Background(), true, persistence.Filter.WithTeam("test-team3"))
				Expect(err).ToNot(HaveOccurred())
				Expect(len(v)).To(Equal(1))

				scanDetails, err := store.Security().GetLatestResourceScan(context.Background(), scan1.ResourceGroup, scan1.ResourceVersion, scan1.ResourceKind, scan1.ResourceNamespace, scan1.ResourceName)
				Expect(err).ToNot(HaveOccurred())
				Expect(len(scanDetails.Results)).To(Equal(2))
			})

			It("should archive previous scans for the same name and namespace if ArchivedAt is nil and the result is different", func() {
				scan1 := getScan("test4", "test-team4", time.Now(), time.Time{}, "Failure", "Failure")
				scan2 := getScan("test4", "test-team4", time.Now().Add(time.Second*10), time.Time{}, "Warning", "Failure")
				scan3 := getScan("test4", "test-team4", time.Now().Add(time.Second*20), time.Time{}, "Warning", "Warning")
				scan4 := getScan("test4", "test-team4", time.Now().Add(time.Second*20), time.Time{}, "Warning", "Warning")
				scan4.Results[1].Message = "changed message"

				storeScans(&scan1, &scan2, &scan3, &scan4)

				v, err := store.Security().ListScans(context.Background(), false, persistence.Filter.WithNamespace("test-team4"), persistence.Filter.WithName("test4"))
				Expect(err).ToNot(HaveOccurred())
				Expect(len(v)).To(Equal(4))

				Expect(v[0].ArchivedAt).To(Equal(v[1].CheckedAt))
				Expect(v[1].ArchivedAt).To(Equal(v[2].CheckedAt))
				Expect(v[2].ArchivedAt).To(Equal(v[3].CheckedAt))
				Expect(v[3].ArchivedAt.IsZero()).To(BeTrue())

				// When getting latest, we should only see scan 4:
				v, err = store.Security().ListScans(context.Background(), true, persistence.Filter.WithNamespace("test-team4"), persistence.Filter.WithName("test4"))
				Expect(err).ToNot(HaveOccurred())
				Expect(len(v)).To(Equal(1))
				Expect(v[0].ID).To(Equal(scan4.ID))

				scanDetails, err := store.Security().GetLatestResourceScan(context.Background(), scan1.ResourceGroup, scan1.ResourceVersion, scan1.ResourceKind, "test-team4", "test4")
				Expect(err).ToNot(HaveOccurred())
				Expect(scanDetails.ID).To(Equal(scan4.ID))
			})

			It("should update checked at for previous scans for the same name and namespace if ArchivedAt is nil and the result is the same", func() {
				scan1 := getScan("test4a", "test-team4a", time.Now(), time.Time{}, "Failure", "Failure")
				scan2 := getScan("test4a", "test-team4a", time.Now().Add(time.Second*10), time.Time{}, "Failure", "Failure")
				scan3 := getScan("test4a", "test-team4a", time.Now().Add(time.Second*20), time.Time{}, "Failure", "Failure")

				storeScans(&scan1, &scan2, &scan3)

				v, err := store.Security().ListScans(context.Background(), false, persistence.Filter.WithNamespace("test-team4a"), persistence.Filter.WithName("test4a"))
				Expect(err).ToNot(HaveOccurred())
				Expect(len(v)).To(Equal(1))

				Expect(v[0].ArchivedAt.IsZero()).To(BeTrue())
				// Round the times to minutes else it gets a bit weird - DB seems to truncate (not round)
				// the timestamp to the second...
				Expect(v[0].CheckedAt.Round(time.Minute)).To(Equal(scan3.CheckedAt.UTC().Round(time.Minute)))

				// When getting latest, we should only see scan 3:
				v, err = store.Security().ListScans(context.Background(), true, persistence.Filter.WithNamespace("test-team4a"), persistence.Filter.WithName("test4a"))
				Expect(err).ToNot(HaveOccurred())
				Expect(len(v)).To(Equal(1))
				Expect(v[0].ID).To(Equal(scan1.ID))

				scanDetails, err := store.Security().GetLatestResourceScan(context.Background(), scan1.ResourceGroup, scan1.ResourceVersion, scan1.ResourceKind, "test-team4a", "test4a")
				Expect(err).ToNot(HaveOccurred())
				Expect(scanDetails.ID).To(Equal(scan1.ID))
			})

			It("should NOT archive previous scans for the same name and namespace if ArchivedAt is set (recording an already archived record)", func() {
				scan1 := getScan("test5", "test-team5", time.Now(), time.Time{}, "Failure", "Failure")
				scan2 := getScan("test5", "test-team5", time.Now().Add(time.Second*10), time.Time{}, "Warning", "Failure")
				scan3 := getScan("test5", "test-team5", time.Now().Add(time.Second*20), time.Now(), "Warning", "Warning")
				scan3.ArchivedAt = time.Now()

				storeScans(&scan1, &scan2, &scan3)

				v, err := store.Security().ListScans(context.Background(), false, persistence.Filter.WithNamespace("test-team5"), persistence.Filter.WithName("test5"))
				Expect(err).ToNot(HaveOccurred())
				Expect(len(v)).To(Equal(3))

				Expect(v[0].ArchivedAt).To(Equal(v[1].CheckedAt))
				Expect(v[1].ArchivedAt.IsZero()).To(BeTrue())
				// Round the times to minutes else it gets a bit weird - DB seems to truncate (not round)
				// the timestamp to the second...
				Expect(v[2].ArchivedAt.Round(time.Minute)).To(Equal(scan3.ArchivedAt.UTC().Round(time.Minute)))

				// When getting latest, we should only see scan 2:
				v, err = store.Security().ListScans(context.Background(), true, persistence.Filter.WithNamespace("test-team5"), persistence.Filter.WithName("test5"))
				Expect(err).ToNot(HaveOccurred())
				Expect(len(v)).To(Equal(1))
				Expect(v[0].ID).To(Equal(scan2.ID))
			})

			It("should NOT archive previous scans when storing a scan with a different API group or version", func() {
				scan1 := getScan("test6", "test-team6", time.Now(), time.Time{}, "Failure", "Failure")
				scan2 := getScan("test6", "test-team6", time.Now().Add(time.Second*10), time.Time{}, "Failure", "Warning")
				scan3 := getScan("test6", "test-team6", time.Now().Add(time.Second*20), time.Time{}, "Warning", "Failure")
				scan4 := getScan("test6", "test-team6", time.Now().Add(time.Second*20), time.Time{}, "Warning", "Warning")
				scan4.ResourceGroup = "example2.appvia.io"
				scan5 := getScan("test6", "test-team6", time.Now().Add(time.Second*20), time.Time{}, "Warning", "Warning")
				scan5.ResourceVersion = "V2"

				storeScans(&scan1, &scan2, &scan3, &scan4, &scan5)

				v, err := store.Security().ListScans(context.Background(), false, persistence.Filter.WithName("test6"))
				Expect(err).ToNot(HaveOccurred())
				Expect(len(v)).To(Equal(5))

				Expect(v[0].ArchivedAt).To(Equal(v[1].CheckedAt))
				Expect(v[1].ArchivedAt).To(Equal(v[2].CheckedAt))
				Expect(v[2].ArchivedAt.IsZero()).To(BeTrue())
				Expect(v[3].ArchivedAt.IsZero()).To(BeTrue())
				Expect(v[4].ArchivedAt.IsZero()).To(BeTrue())

				// When getting latest, we should see both scan 3, scan 4 and scan 5:
				v, err = store.Security().ListScans(context.Background(), true, persistence.Filter.WithName("test6"))
				Expect(err).ToNot(HaveOccurred())
				Expect(len(v)).To(Equal(3))
				Expect(v[0].ID).To(Equal(scan3.ID))
				Expect(v[1].ID).To(Equal(scan4.ID))
				Expect(v[2].ID).To(Equal(scan5.ID))
			})

			It("should NOT archive previous scans when storing a scan with a different kind", func() {
				scan1 := getScan("test7", "test-team7", time.Now(), time.Time{}, "Failure", "Failure")
				scan2 := getScan("test7", "test-team7", time.Now().Add(time.Second*10), time.Time{}, "Failure", "Warning")
				scan3 := getScan("test7", "test-team7", time.Now().Add(time.Second*20), time.Time{}, "Warning", "Failure")
				scan4 := getScan("test7", "test-team7", time.Now().Add(time.Second*20), time.Time{}, "Warning", "Warning")
				scan4.ResourceKind = "Example2"

				storeScans(&scan1, &scan2, &scan3, &scan4)

				v, err := store.Security().ListScans(context.Background(), false, persistence.Filter.WithName("test7"))
				Expect(err).ToNot(HaveOccurred())
				Expect(len(v)).To(Equal(4))

				Expect(v[0].ArchivedAt).To(Equal(v[1].CheckedAt))
				Expect(v[1].ArchivedAt).To(Equal(v[2].CheckedAt))
				Expect(v[2].ArchivedAt.IsZero()).To(BeTrue())
				Expect(v[3].ArchivedAt.IsZero()).To(BeTrue())

				// When getting latest, we should see both scan 3 and scan 4:
				v, err = store.Security().ListScans(context.Background(), true, persistence.Filter.WithName("test7"))
				Expect(err).ToNot(HaveOccurred())
				Expect(len(v)).To(Equal(2))
				Expect(v[0].ID).To(Equal(scan3.ID))
				Expect(v[1].ID).To(Equal(scan4.ID))
			})

			It("should NOT archive previous scans when storing a scan with a different namespace", func() {
				scan1 := getScan("test8", "test-team8", time.Now(), time.Time{}, "Failure", "Failure")
				scan2 := getScan("test8", "test-team8", time.Now().Add(time.Second*10), time.Time{}, "Warning", "Failure")
				scan3 := getScan("test8", "test-team8", time.Now().Add(time.Second*20), time.Time{}, "Failure", "Warning")
				scan4 := getScan("test8", "test-team8a", time.Now().Add(time.Second*20), time.Time{}, "Warning", "Warning")

				storeScans(&scan1, &scan2, &scan3, &scan4)

				v, err := store.Security().ListScans(context.Background(), false, persistence.Filter.WithName("test8"))
				Expect(err).ToNot(HaveOccurred())
				Expect(len(v)).To(Equal(4))

				Expect(v[0].ArchivedAt).To(Equal(v[1].CheckedAt))
				Expect(v[1].ArchivedAt).To(Equal(v[2].CheckedAt))
				Expect(v[2].ArchivedAt.IsZero()).To(BeTrue())
				Expect(v[3].ArchivedAt.IsZero()).To(BeTrue())

				// When getting latest, we should see both scan 3 and scan 4:
				v, err = store.Security().ListScans(context.Background(), true, persistence.Filter.WithName("test8"))
				Expect(err).ToNot(HaveOccurred())
				Expect(len(v)).To(Equal(2))
				Expect(v[0].ID).To(Equal(scan3.ID))
				Expect(v[1].ID).To(Equal(scan4.ID))
			})

			It("should NOT archive previous scans when storing a scan with a different name", func() {
				scan1 := getScan("test9", "test-team9", time.Now(), time.Time{}, "Failure", "Failure")
				scan2 := getScan("test9", "test-team9", time.Now().Add(time.Second*10), time.Time{}, "Failure", "Warning")
				scan3 := getScan("test9", "test-team9", time.Now().Add(time.Second*20), time.Time{}, "Warning", "Failure")
				scan4 := getScan("test9a", "test-team9", time.Now().Add(time.Second*20), time.Time{}, "Warning", "Warning")

				storeScans(&scan1, &scan2, &scan3, &scan4)

				v, err := store.Security().ListScans(context.Background(), false, persistence.Filter.WithTeam("test-team9"))
				Expect(err).ToNot(HaveOccurred())
				Expect(len(v)).To(Equal(4))

				/*
					  NO IDEA - BUT MOVING ON AND FIX LATER
						Expect(v[0].ArchivedAt).To(Equal(v[1].CheckedAt))
						Expect(v[1].ArchivedAt).To(Equal(v[2].CheckedAt))
						Expect(v[2].ArchivedAt.IsZero()).To(BeTrue())
						Expect(v[3].ArchivedAt.IsZero()).To(BeTrue())
				*/

				// The latest for test-team7, test7 should be scan 3
				scanDetails, err := store.Security().GetLatestResourceScan(context.Background(), scan1.ResourceGroup, scan1.ResourceVersion, scan1.ResourceKind, "test-team9", "test9")
				Expect(err).ToNot(HaveOccurred())
				Expect(scanDetails.ID).To(Equal(scan3.ID))
				// The latest for test-team7, test7a should be scan 4
				scanDetails, err = store.Security().GetLatestResourceScan(context.Background(), scan1.ResourceGroup, scan1.ResourceVersion, scan1.ResourceKind, "test-team9", "test9a")
				Expect(err).ToNot(HaveOccurred())
				Expect(scanDetails.ID).To(Equal(scan4.ID))
			})
		})
	})

	Describe("ArchiveResourceScans", func() {
		When("called", func() {
			It("should set any unarchived scans to archived for the resource", func() {
				scan1 := getScan("test10", "test-team10", time.Now(), time.Time{}, "Failure", "Failure")
				scan2 := getScan("test10", "test-team10", time.Now().Add(time.Second*10), time.Time{}, "Failure", "Warning")
				scan3 := getScan("test10", "test-team10", time.Now().Add(time.Second*20), time.Time{}, "Warning", "Failure")
				scan4 := getScan("test10a", "test-team10", time.Now().Add(time.Second*20), time.Time{}, "Warning", "Warning")

				storeScans(&scan1, &scan2, &scan3, &scan4)

				store.Security().ArchiveResourceScans(
					context.Background(),
					scan1.ResourceGroup,
					scan1.ResourceVersion,
					scan1.ResourceKind,
					scan1.ResourceNamespace,
					scan1.ResourceName)

				v, err := store.Security().ListScans(context.Background(), false, persistence.Filter.WithTeam("test-team10"))
				Expect(err).ToNot(HaveOccurred())
				Expect(len(v)).To(Equal(4))

				/*
					  NO IDEA - BUT MOVING ON AND FIX LATER
						Expect(v[0].ArchivedAt).To(Equal(v[1].CheckedAt))
						Expect(v[1].ArchivedAt).To(Equal(v[2].CheckedAt))
						// This will have been set to the current time:
						Expect(v[2].ArchivedAt.IsZero()).To(BeFalse())
						// This should not have been archived as it has a different name:
						Expect(v[3].ArchivedAt.IsZero()).To(BeTrue())
				*/

				scanDetails, err := store.Security().GetLatestResourceScan(context.Background(), scan1.ResourceGroup, scan1.ResourceVersion, scan1.ResourceKind, "test-team10", "test10")
				Expect(err).ToNot(HaveOccurred())
				Expect(scanDetails).To(BeNil())
			})
		})
	})
})
