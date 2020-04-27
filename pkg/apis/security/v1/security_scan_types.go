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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// RuleStatus values represent the possible status of compliance with a security rule.
type RuleStatus string

func (r RuleStatus) String() string {
	return string(r)
}

const (
	// Compliant indicates that this target is fully compliant with the specified rule.
	Compliant RuleStatus = "Compliant"
	// Warning indicates that this target is uncompliant in such a way that
	// consideration should be made as to whether this should be remediated. This would
	// typically be used for best practice considerations, where not being compliant
	// isn't necessarily a critical issue.
	Warning RuleStatus = "Warning"
	// Failure indicates that this target is uncompliant in a significant way and
	// should be mitigated. This would typically be used for rules where compliance is
	// considered to be vital to a well-run cluster.
	Failure RuleStatus = "Failure"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ScanResult contains the result of a scan against all registered rules
// +k8s:openapi-gen=false
type ScanResult struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              ScanResultSpec `json:"spec,omitempty"`
}

// ScanResultSpec shows the overall result of a scan against all registered rules
// +k8s:openapi-gen=false
type ScanResultSpec struct {
	// ID is the ID of this scan result in the data store
	ID uint64 `json:"id,omitempty"`
	// ResourceAPIVersion is the group and version of the resource scanned by this scan
	ResourceAPIVersion string `json:"resourceApiVersion,omitempty"`
	// ResourceKind is the kind of the resource scanned by this scan
	ResourceKind string `json:"resourceKind,omitempty"`
	// ResourceNamespace is the namespace of the resource scanned by this scan
	ResourceNamespace string `json:"resourceNamespace,omitempty"`
	// ResourceName is the name of the resource scanned by this scan
	ResourceName string `json:"resourceName,omitempty"`
	// OwningTeam is the name of the Kore team that owns this resource, will be empty if it is a non-team resource.
	OwningTeam string `json:"owningTeam,omitempty"`
	// CheckedAt is the timestamp this result was determined
	CheckedAt metav1.Time `json:"checkedAt,omitempty"`
	// ArchivedAt is the timestamp this result was superceded by a later scan - if ArchivedAt.IsZero() is true this is the most recent scan.
	ArchivedAt metav1.Time `json:"archivedAt,omitempty"`
	// OverallStatus indicates the worst-case status of the rules checked in this scan
	OverallStatus RuleStatus `json:"overallStatus,omitempty"`
	// Results are the underlying results of the individual rules run as part of this scan
	Results []RuleResult `json:"results,omitempty"`
}

// RuleResult represents the compliance status of a target with respect to a
// specific security rule.
// +k8s:openapi-gen=false
type RuleResult struct {
	// RuleCode indicates the rule that this result relates to
	RuleCode string `json:"ruleCode,omitempty"`
	// Status indicates the compliance of the target with this rule
	Status RuleStatus `json:"status,omitempty"`
	// Message provides additional information about the status of this rule on this
	// target, if applicable
	Message string `json:"message,omitempty"`
	// CheckedAt is the timestamp this result was determined
	CheckedAt metav1.Time `json:"checkedAt,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ScanResultList contains a list of scan results event
type ScanResultList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ScanResult `json:"items"`
}
