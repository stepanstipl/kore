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

package kore

// LabelKoreIdentifier is the label for the immutable unique identifier for this instance
// of Kore. This will always be set by Kore and any user-provided values will be ignored.
var LabelKoreIdentifier = Label("koreid")

// LabelClusterIdentifier is the label for the immutable unique identifier for a cluster.
// This should be left blank for new clusters (Kore will auto-assign a new value) or
// populated with the identifier of a previously-deleted cluster which this cluster
// replaces. No other values will be permitted, and read-only after the cluster is created.
var LabelClusterIdentifier = Label("clusterid")

// LabelTeamIdentifier is the label for the immutable unique identifier for a team that
// owns this resource. This should be left blank (Kore will auto-populate with the correct
// identifier for the team) or populated with the team's identifier for a resource. No
// other values permitted, and read-only after a cluster or team is created.
var LabelTeamIdentifier = Label("teamid")
