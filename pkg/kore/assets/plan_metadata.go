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

package assets

// EKSContinents is the list of supported regions for EKS clusters, organised by continent. Will be
// replaced by info sourced from cloudinfo shortly.
const EKSContinents = `[
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
        "name": "Asia",
        "regions": [
            {
                "id": "ap-east-1",
                "name": "Asia Pacific (Hong Kong)"
            },
            {
                "id": "ap-northeast-2",
                "name": "Asia Pacific (Seoul)"
            },
            {
                "id": "ap-southeast-1",
                "name": "Asia Pacific (Singapore)"
            },
            {
                "id": "ap-northeast-1",
                "name": "Asia Pacific (Tokyo)"
            },
            {
                "id": "ap-south-1",
                "name": "Asia Pacific (Mumbai)"
            },
            {
                "id": "me-south-1",
                "name": "Middle East (Bahrain)"
            }
        ]
    },
    {
        "name": "Europe",
        "regions": [
            {
                "id": "eu-west-2",
                "name": "EU (London)"
            },
            {
                "id": "eu-north-1",
                "name": "EU (Stockholm)"
            },
            {
                "id": "eu-west-3",
                "name": "EU (Paris)"
            },
            {
                "id": "eu-central-1",
                "name": "EU (Frankfurt)"
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
                "id": "us-west-2",
                "name": "US West (Oregon)"
            },
            {
                "id": "us-east-2",
                "name": "US East (Ohio)"
            }
        ]
    }
]`

// GKEContinents is the list of supported regions for GKE clusters, organised by continent. Will be
// replaced by info sourced from cloudinfo shortly.
const GKEContinents = `[
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
        "name": "Europe",
        "regions": [
            {
                "id": "europe-west1",
                "name": "EU (Belgium)"
            },
            {
                "id": "europe-west3",
                "name": "EU (Frankfurt)"
            },
            {
                "id": "europe-west4",
                "name": "EU (Netherlands)"
            },
            {
                "id": "europe-north1",
                "name": "EU (Finland)"
            },
            {
                "id": "europe-west2",
                "name": "EU (London)"
            }
        ]
    },
    {
        "name": "North America",
        "regions": [
            {
                "id": "us-east4",
                "name": "US East (Northern Virginia)"
            },
            {
                "id": "us-west2",
                "name": "US West (Los Angeles)"
            },
            {
                "id": "us-west1",
                "name": "US West (Oregon)"
            },
            {
                "id": "northamerica-northeast1",
                "name": "Canada (Montréal)"
            },
            {
                "id": "us-east1",
                "name": "US East (South Carolina)"
            },
            {
                "id": "us-central1",
                "name": "US Central (Iowa)"
            }
        ]
    },
    {
        "name": "Asia",
        "regions": [
            {
                "id": "asia-south1",
                "name": "Asia Pacific (Mumbai)"
            },
            {
                "id": "asia-east2",
                "name": "Asia Pacific (Hong Kong)"
            },
            {
                "id": "asia-east1",
                "name": "Asia Pacific (Taiwan)"
            },
            {
                "id": "asia-northeast1",
                "name": "Asia Pacific (Tokyo)"
            },
            {
                "id": "asia-southeast1",
                "name": "Asia Pacific (Singapore)"
            }
        ]
    },
    {
        "name": "South America",
        "regions": [
            {
                "id": "southamerica-east1",
                "name": "South America (São Paulo)"
            }
        ]
    }
]`
