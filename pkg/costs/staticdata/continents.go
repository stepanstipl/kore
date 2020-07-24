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

package staticdata

func Continents(cloud string) string {
	switch cloud {
	case "gcp":
		return gcpContinents
	case "aws":
		return awsContinents
	case "azure":
		return azureContinents
	}
	return ""
}

// gcp_continents source https://banzaicloud.com/cloudinfo/api/v1/providers/google/services/gke/continents
const gcpContinents = `[
	 {
		 "name": "South America",
		 "regions": [
			 {
				 "id": "southamerica-east1",
				 "name": "South America (São Paulo)"
			 }
		 ]
	 },
	 {
		 "name": "Australia",
		 "regions": [
			 {
				 "id": "australia-southeast1",
				 "name": "Asia Pacific (Sydney)"
			 }
		 ]
	 },
	 {
		 "name": "Asia",
		 "regions": [
			 {
				 "id": "asia-east2",
				 "name": "Asia Pacific (Hong Kong)"
			 },
			 {
				 "id": "asia-northeast1",
				 "name": "Asia Pacific (Tokyo)"
			 },
			 {
				 "id": "asia-south1",
				 "name": "Asia Pacific (Mumbai)"
			 },
			 {
				 "id": "asia-southeast1",
				 "name": "Asia Pacific (Singapore)"
			 },
			 {
				 "id": "asia-east1",
				 "name": "Asia Pacific (Taiwan)"
			 }
		 ]
	 },
	 {
		 "name": "North America",
		 "regions": [
			 {
				 "id": "northamerica-northeast1",
				 "name": "Canada (Montréal)"
			 },
			 {
				 "id": "us-east4",
				 "name": "US East (Northern Virginia)"
			 },
			 {
				 "id": "us-west1",
				 "name": "US West (Oregon)"
			 },
			 {
				 "id": "us-west2",
				 "name": "US West (Los Angeles)"
			 },
			 {
				 "id": "us-central1",
				 "name": "US Central (Iowa)"
			 },
			 {
				 "id": "us-east1",
				 "name": "US East (South Carolina)"
			 }
		 ]
	 },
	 {
		 "name": "Europe",
		 "regions": [
			 {
				 "id": "europe-west2",
				 "name": "EU (London)"
			 },
			 {
				 "id": "europe-west4",
				 "name": "EU (Netherlands)"
			 },
			 {
				 "id": "europe-west3",
				 "name": "EU (Frankfurt)"
			 },
			 {
				 "id": "europe-north1",
				 "name": "EU (Finland)"
			 },
			 {
				 "id": "europe-west1",
				 "name": "EU (Belgium)"
			 }
		 ]
	 }
 ]`

// aws_continents source https://banzaicloud.com/cloudinfo/api/v1/providers/amazon/services/eks/continents
const awsContinents = `[
	 {
		 "name": "Asia",
		 "regions": [
			 {
				 "id": "ap-northeast-2",
				 "name": "Asia Pacific (Seoul)"
			 },
			 {
				 "id": "ap-south-1",
				 "name": "Asia Pacific (Mumbai)"
			 },
			 {
				 "id": "ap-northeast-1",
				 "name": "Asia Pacific (Tokyo)"
			 },
			 {
				 "id": "ap-east-1",
				 "name": "Asia Pacific (Hong Kong)"
			 },
			 {
				 "id": "ap-southeast-1",
				 "name": "Asia Pacific (Singapore)"
			 },
			 {
				 "id": "me-south-1",
				 "name": "Middle East (Bahrain)"
			 }
		 ]
	 },
	 {
		 "name": "Australia",
		 "regions": [
			 {
				 "id": "ap-southeast-2",
				 "name": "Asia Pacific (Sydney)"
			 }
		 ]
	 },
	 {
		 "name": "Europe",
		 "regions": [
			 {
				 "id": "eu-central-1",
				 "name": "EU (Frankfurt)"
			 },
			 {
				 "id": "eu-west-2",
				 "name": "EU (London)"
			 },
			 {
				 "id": "eu-west-3",
				 "name": "EU (Paris)"
			 },
			 {
				 "id": "eu-north-1",
				 "name": "EU (Stockholm)"
			 },
			 {
				 "id": "eu-west-1",
				 "name": "EU (Ireland)"
			 }
		 ]
	 },
	 {
		 "name": "North America",
		 "regions": [
			 {
				 "id": "us-east-1",
				 "name": "US East (N. Virginia)"
			 },
			 {
				 "id": "us-east-2",
				 "name": "US East (Ohio)"
			 },
			 {
				 "id": "us-west-2",
				 "name": "US West (Oregon)"
			 }
		 ]
	 }
 ]`

const azureContinents = `[{"name":"Australia","regions":[{"id":"australiasoutheast","name":"Australia Southeast"},{"id":"australiaeast","name":"Australia East"}]},{"name":"South America","regions":[{"id":"brazilsouth","name":"Brazil South"}]},{"name":"Africa","regions":[{"id":"southafricanorth","name":"South Africa North"}]},{"name":"Asia","regions":[{"id":"koreasouth","name":"Korea South"},{"id":"eastasia","name":"East Asia"},{"id":"southeastasia","name":"Southeast Asia"},{"id":"southindia","name":"South India"},{"id":"centralindia","name":"Central India"},{"id":"japanwest","name":"Japan West"},{"id":"japaneast","name":"Japan East"},{"id":"uaenorth","name":"UAE North"},{"id":"koreacentral","name":"Korea Central"}]},{"name":"North America","regions":[{"id":"westcentralus","name":"West Central US"},{"id":"southcentralus","name":"South Central US"},{"id":"canadacentral","name":"Canada Central"},{"id":"centralus","name":"Central US"},{"id":"westus2","name":"West US 2"},{"id":"eastus","name":"East US"},{"id":"westus","name":"West US"},{"id":"canadaeast","name":"Canada East"},{"id":"eastus2","name":"East US 2"},{"id":"northcentralus","name":"North Central US"}]},{"name":"Europe","regions":[{"id":"westeurope","name":"West Europe"},{"id":"norwayeast","name":"Norway East"},{"id":"germanynorth","name":"Germany North"},{"id":"northeurope","name":"North Europe"},{"id":"francecentral","name":"France Central"},{"id":"uksouth","name":"UK South"},{"id":"germanywestcentral","name":"Germany West Central"},{"id":"switzerlandnorth","name":"Switzerland North"},{"id":"switzerlandwest","name":"Switzerland West"},{"id":"norwaywest","name":"Norway West"},{"id":"ukwest","name":"UK West"}]}]`
