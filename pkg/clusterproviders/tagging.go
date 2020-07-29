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

package clusterproviders

import (
	clustersv1 "github.com/appvia/kore/pkg/apis/clusters/v1"
	"github.com/appvia/kore/pkg/kore"
)

// SetCommonTags ensures that the common tags (such as team and cluster identifiers)
// are set on the provided tag map
func SetCommonTags(tags *map[string]string, cluster *clustersv1.Cluster) {
	if *tags == nil {
		*tags = map[string]string{}
	}
	if cluster.Labels[kore.LabelKoreIdentifier] != "" {
		(*tags)["kore-instance"] = cluster.Labels[kore.LabelKoreIdentifier]
	}
	if cluster.Labels[kore.LabelTeamIdentifier] != "" {
		(*tags)["kore-team"] = cluster.Labels[kore.LabelTeamIdentifier]
	}
	if cluster.Labels[kore.LabelClusterIdentifier] != "" {
		(*tags)["kore-cluster"] = cluster.Labels[kore.LabelClusterIdentifier]
	}
}
