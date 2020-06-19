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

// EKSInstanceTypes is the list of possible instance types for EKS clusters. Will be
// replaced by info sourced from cloudinfo shortly.
const EKSInstanceTypes = `{
    "products": [
        {
            "category": "Storage optimized",
            "type": "i3en.metal",
            "onDemandPrice": 12.624,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 3.7872
                },
                {
                    "zone": "eu-west-2b",
                    "price": 3.7872
                },
                {
                    "zone": "eu-west-2c",
                    "price": 12.624
                }
            ],
            "cpusPerVm": 96,
            "memPerVm": 768,
            "gpusPerVm": 0,
            "ntwPerf": "100 Gigabit",
            "ntwPerfCategory": "extra",
            "zones": null,
            "attributes": {
                "cpu": "96",
                "instanceTypeCategory": "Storage optimized",
                "memory": "768",
                "networkPerfCategory": "extra"
            },
            "currentGen": true
        },
        {
            "category": "General purpose",
            "type": "m5.8xlarge",
            "onDemandPrice": 1.776,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 0.5278
                },
                {
                    "zone": "eu-west-2b",
                    "price": 0.5408
                },
                {
                    "zone": "eu-west-2c",
                    "price": 0.5278
                }
            ],
            "cpusPerVm": 32,
            "memPerVm": 128,
            "gpusPerVm": 0,
            "ntwPerf": "10 Gigabit",
            "ntwPerfCategory": "high",
            "zones": null,
            "attributes": {
                "cpu": "32",
                "instanceTypeCategory": "General purpose",
                "memory": "128",
                "networkPerfCategory": "high"
            },
            "currentGen": true
        },
        {
            "category": "GPU instance",
            "type": "p3.2xlarge",
            "onDemandPrice": 3.589,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 1.0767
                },
                {
                    "zone": "eu-west-2b",
                    "price": 1.0767
                }
            ],
            "cpusPerVm": 8,
            "memPerVm": 61,
            "gpusPerVm": 1,
            "ntwPerf": "Up to 10 Gigabit",
            "ntwPerfCategory": "high",
            "zones": null,
            "attributes": {
                "cpu": "8",
                "instanceTypeCategory": "GPU instance",
                "memory": "61",
                "networkPerfCategory": "high"
            },
            "currentGen": true
        },
        {
            "category": "Compute optimized",
            "type": "c4.xlarge",
            "onDemandPrice": 0.237,
            "spotPrice": [
                {
                    "zone": "eu-west-2c",
                    "price": 0.06
                },
                {
                    "zone": "eu-west-2a",
                    "price": 0.0598
                },
                {
                    "zone": "eu-west-2b",
                    "price": 0.0598
                }
            ],
            "cpusPerVm": 4,
            "memPerVm": 7.5,
            "gpusPerVm": 0,
            "ntwPerf": "High",
            "ntwPerfCategory": "medium",
            "zones": null,
            "attributes": {
                "cpu": "4",
                "instanceTypeCategory": "Compute optimized",
                "memory": "7.5",
                "networkPerfCategory": "medium"
            },
            "currentGen": true
        },
        {
            "category": "GPU instance",
            "type": "g4dn.8xlarge",
            "onDemandPrice": 2.546,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 0.7638
                },
                {
                    "zone": "eu-west-2b",
                    "price": 0.7638
                },
                {
                    "zone": "eu-west-2c",
                    "price": 0.7638
                }
            ],
            "cpusPerVm": 32,
            "memPerVm": 128,
            "gpusPerVm": 1,
            "ntwPerf": "50 Gigabit",
            "ntwPerfCategory": "extra",
            "zones": null,
            "attributes": {
                "cpu": "32",
                "instanceTypeCategory": "GPU instance",
                "memory": "128",
                "networkPerfCategory": "extra"
            },
            "currentGen": true
        },
        {
            "category": "Compute optimized",
            "type": "c5d.2xlarge",
            "onDemandPrice": 0.46,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 0.1257
                },
                {
                    "zone": "eu-west-2b",
                    "price": 0.1257
                },
                {
                    "zone": "eu-west-2c",
                    "price": 0.1257
                }
            ],
            "cpusPerVm": 8,
            "memPerVm": 16,
            "gpusPerVm": 0,
            "ntwPerf": "Up to 10 Gigabit",
            "ntwPerfCategory": "high",
            "zones": null,
            "attributes": {
                "cpu": "8",
                "instanceTypeCategory": "Compute optimized",
                "memory": "16",
                "networkPerfCategory": "high"
            },
            "currentGen": true
        },
        {
            "category": "Memory optimized",
            "type": "r5a.2xlarge",
            "onDemandPrice": 0.532,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 0.1382
                },
                {
                    "zone": "eu-west-2b",
                    "price": 0.1382
                },
                {
                    "zone": "eu-west-2c",
                    "price": 0.1382
                }
            ],
            "cpusPerVm": 8,
            "memPerVm": 64,
            "gpusPerVm": 0,
            "ntwPerf": "10 Gigabit",
            "ntwPerfCategory": "high",
            "zones": null,
            "attributes": {
                "cpu": "8",
                "instanceTypeCategory": "Memory optimized",
                "memory": "64",
                "networkPerfCategory": "high"
            },
            "currentGen": true
        },
        {
            "category": "General purpose",
            "type": "m5.2xlarge",
            "onDemandPrice": 0.444,
            "spotPrice": [
                {
                    "zone": "eu-west-2b",
                    "price": 0.137
                },
                {
                    "zone": "eu-west-2c",
                    "price": 0.1571
                },
                {
                    "zone": "eu-west-2a",
                    "price": 0.1353
                }
            ],
            "cpusPerVm": 8,
            "memPerVm": 32,
            "gpusPerVm": 0,
            "ntwPerf": "Up to 10 Gigabit",
            "ntwPerfCategory": "high",
            "zones": null,
            "attributes": {
                "cpu": "8",
                "instanceTypeCategory": "General purpose",
                "memory": "32",
                "networkPerfCategory": "high"
            },
            "currentGen": true
        },
        {
            "category": "Compute optimized",
            "type": "c5d.12xlarge",
            "onDemandPrice": 2.76,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 0.754
                },
                {
                    "zone": "eu-west-2b",
                    "price": 0.754
                },
                {
                    "zone": "eu-west-2c",
                    "price": 0.754
                }
            ],
            "cpusPerVm": 48,
            "memPerVm": 96,
            "gpusPerVm": 0,
            "ntwPerf": "12 Gigabit",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": {
                "cpu": "48",
                "instanceTypeCategory": "Compute optimized",
                "memory": "96",
                "networkPerfCategory": ""
            },
            "currentGen": true
        },
        {
            "category": "Compute optimized",
            "type": "c4.4xlarge",
            "onDemandPrice": 0.95,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 0.245
                },
                {
                    "zone": "eu-west-2b",
                    "price": 0.2394
                },
                {
                    "zone": "eu-west-2c",
                    "price": 0.2434
                }
            ],
            "cpusPerVm": 16,
            "memPerVm": 30,
            "gpusPerVm": 0,
            "ntwPerf": "High",
            "ntwPerfCategory": "medium",
            "zones": null,
            "attributes": {
                "cpu": "16",
                "instanceTypeCategory": "Compute optimized",
                "memory": "30",
                "networkPerfCategory": "medium"
            },
            "currentGen": true
        },
        {
            "category": "General purpose",
            "type": "m5.metal",
            "onDemandPrice": 5.328,
            "spotPrice": [
                {
                    "zone": "eu-west-2b",
                    "price": 1.0557
                },
                {
                    "zone": "eu-west-2c",
                    "price": 1.0557
                },
                {
                    "zone": "eu-west-2a",
                    "price": 1.0557
                }
            ],
            "cpusPerVm": 96,
            "memPerVm": 384,
            "gpusPerVm": 0,
            "ntwPerf": "25 Gigabit",
            "ntwPerfCategory": "extra",
            "zones": null,
            "attributes": {
                "cpu": "96",
                "instanceTypeCategory": "General purpose",
                "memory": "384",
                "networkPerfCategory": "extra"
            },
            "currentGen": true
        },
        {
            "category": "General purpose",
            "type": "t3.nano",
            "onDemandPrice": 0.0059,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 0.0022
                },
                {
                    "zone": "eu-west-2b",
                    "price": 0.0023
                },
                {
                    "zone": "eu-west-2c",
                    "price": 0.0018
                }
            ],
            "cpusPerVm": 2,
            "memPerVm": 0.5,
            "gpusPerVm": 0,
            "ntwPerf": "Low",
            "ntwPerfCategory": "low",
            "zones": null,
            "attributes": {
                "cpu": "2",
                "instanceTypeCategory": "General purpose",
                "memory": "0.5",
                "networkPerfCategory": "low"
            },
            "currentGen": true,
            "burst": true
        },
        {
            "category": "General purpose",
            "type": "m4.4xlarge",
            "onDemandPrice": 0.928,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 0.2513
                },
                {
                    "zone": "eu-west-2b",
                    "price": 0.2513
                },
                {
                    "zone": "eu-west-2c",
                    "price": 0.2523
                }
            ],
            "cpusPerVm": 16,
            "memPerVm": 64,
            "gpusPerVm": 0,
            "ntwPerf": "High",
            "ntwPerfCategory": "medium",
            "zones": null,
            "attributes": {
                "cpu": "16",
                "instanceTypeCategory": "General purpose",
                "memory": "64",
                "networkPerfCategory": "medium"
            },
            "currentGen": true
        },
        {
            "category": "Storage optimized",
            "type": "i3en.3xlarge",
            "onDemandPrice": 1.578,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 0.4734
                },
                {
                    "zone": "eu-west-2b",
                    "price": 0.4734
                },
                {
                    "zone": "eu-west-2c",
                    "price": 0.4734
                }
            ],
            "cpusPerVm": 12,
            "memPerVm": 96,
            "gpusPerVm": 0,
            "ntwPerf": "Up to 25 Gigabit",
            "ntwPerfCategory": "extra",
            "zones": null,
            "attributes": {
                "cpu": "12",
                "instanceTypeCategory": "Storage optimized",
                "memory": "96",
                "networkPerfCategory": "extra"
            },
            "currentGen": true
        },
        {
            "category": "Compute optimized",
            "type": "c5.2xlarge",
            "onDemandPrice": 0.404,
            "spotPrice": [
                {
                    "zone": "eu-west-2c",
                    "price": 0.1337
                },
                {
                    "zone": "eu-west-2a",
                    "price": 0.1398
                },
                {
                    "zone": "eu-west-2b",
                    "price": 0.1266
                }
            ],
            "cpusPerVm": 8,
            "memPerVm": 16,
            "gpusPerVm": 0,
            "ntwPerf": "Up to 10 Gigabit",
            "ntwPerfCategory": "high",
            "zones": null,
            "attributes": {
                "cpu": "8",
                "instanceTypeCategory": "Compute optimized",
                "memory": "16",
                "networkPerfCategory": "high"
            },
            "currentGen": true
        },
        {
            "category": "Memory optimized",
            "type": "r5d.2xlarge",
            "onDemandPrice": 0.676,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 0.1382
                },
                {
                    "zone": "eu-west-2b",
                    "price": 0.1382
                },
                {
                    "zone": "eu-west-2c",
                    "price": 0.1382
                }
            ],
            "cpusPerVm": 8,
            "memPerVm": 64,
            "gpusPerVm": 0,
            "ntwPerf": "10 Gigabit",
            "ntwPerfCategory": "high",
            "zones": null,
            "attributes": {
                "cpu": "8",
                "instanceTypeCategory": "Memory optimized",
                "memory": "64",
                "networkPerfCategory": "high"
            },
            "currentGen": true
        },
        {
            "category": "GPU instance",
            "type": "g3.4xlarge",
            "onDemandPrice": 1.429,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 0.4287
                },
                {
                    "zone": "eu-west-2b",
                    "price": 0.4287
                }
            ],
            "cpusPerVm": 16,
            "memPerVm": 122,
            "gpusPerVm": 1,
            "ntwPerf": "Up to 10 Gigabit",
            "ntwPerfCategory": "high",
            "zones": null,
            "attributes": {
                "cpu": "16",
                "instanceTypeCategory": "GPU instance",
                "memory": "122",
                "networkPerfCategory": "high"
            },
            "currentGen": true
        },
        {
            "category": "General purpose",
            "type": "t3.micro",
            "onDemandPrice": 0.0118,
            "spotPrice": [
                {
                    "zone": "eu-west-2c",
                    "price": 0.0035
                },
                {
                    "zone": "eu-west-2a",
                    "price": 0.0035
                },
                {
                    "zone": "eu-west-2b",
                    "price": 0.0035
                }
            ],
            "cpusPerVm": 2,
            "memPerVm": 1,
            "gpusPerVm": 0,
            "ntwPerf": "Low to Moderate",
            "ntwPerfCategory": "low",
            "zones": null,
            "attributes": {
                "cpu": "2",
                "instanceTypeCategory": "General purpose",
                "memory": "1",
                "networkPerfCategory": "low"
            },
            "currentGen": true,
            "burst": true
        },
        {
            "category": "General purpose",
            "type": "m4.large",
            "onDemandPrice": 0.116,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 0.0314
                },
                {
                    "zone": "eu-west-2b",
                    "price": 0.0314
                },
                {
                    "zone": "eu-west-2c",
                    "price": 0.0315
                }
            ],
            "cpusPerVm": 2,
            "memPerVm": 8,
            "gpusPerVm": 0,
            "ntwPerf": "Moderate",
            "ntwPerfCategory": "medium",
            "zones": null,
            "attributes": {
                "cpu": "2",
                "instanceTypeCategory": "General purpose",
                "memory": "8",
                "networkPerfCategory": "medium"
            },
            "currentGen": true
        },
        {
            "category": "Compute optimized",
            "type": "c5d.metal",
            "onDemandPrice": 5.52,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 1.5081
                },
                {
                    "zone": "eu-west-2b",
                    "price": 5.52
                },
                {
                    "zone": "eu-west-2c",
                    "price": 1.5081
                }
            ],
            "cpusPerVm": 96,
            "memPerVm": 192,
            "gpusPerVm": 0,
            "ntwPerf": "25 Gigabit",
            "ntwPerfCategory": "extra",
            "zones": null,
            "attributes": {
                "cpu": "96",
                "instanceTypeCategory": "Compute optimized",
                "memory": "192",
                "networkPerfCategory": "extra"
            },
            "currentGen": true
        },
        {
            "category": "Memory optimized",
            "type": "r5d.24xlarge",
            "onDemandPrice": 8.112,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 1.6589
                },
                {
                    "zone": "eu-west-2b",
                    "price": 1.6589
                },
                {
                    "zone": "eu-west-2c",
                    "price": 1.6589
                }
            ],
            "cpusPerVm": 96,
            "memPerVm": 768,
            "gpusPerVm": 0,
            "ntwPerf": "25 Gigabit",
            "ntwPerfCategory": "extra",
            "zones": null,
            "attributes": {
                "cpu": "96",
                "instanceTypeCategory": "Memory optimized",
                "memory": "768",
                "networkPerfCategory": "extra"
            },
            "currentGen": true
        },
        {
            "category": "Memory optimized",
            "type": "r5d.8xlarge",
            "onDemandPrice": 2.704,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 0.553
                },
                {
                    "zone": "eu-west-2b",
                    "price": 0.553
                },
                {
                    "zone": "eu-west-2c",
                    "price": 0.553
                }
            ],
            "cpusPerVm": 32,
            "memPerVm": 256,
            "gpusPerVm": 0,
            "ntwPerf": "10 Gigabit",
            "ntwPerfCategory": "high",
            "zones": null,
            "attributes": {
                "cpu": "32",
                "instanceTypeCategory": "Memory optimized",
                "memory": "256",
                "networkPerfCategory": "high"
            },
            "currentGen": true
        },
        {
            "category": "Memory optimized",
            "type": "r4.16xlarge",
            "onDemandPrice": 4.992,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 1.0533
                },
                {
                    "zone": "eu-west-2b",
                    "price": 1.0533
                },
                {
                    "zone": "eu-west-2c",
                    "price": 1.0533
                }
            ],
            "cpusPerVm": 64,
            "memPerVm": 488,
            "gpusPerVm": 0,
            "ntwPerf": "20 Gigabit",
            "ntwPerfCategory": "extra",
            "zones": null,
            "attributes": {
                "cpu": "64",
                "instanceTypeCategory": "Memory optimized",
                "memory": "488",
                "networkPerfCategory": "extra"
            },
            "currentGen": true
        },
        {
            "category": "Memory optimized",
            "type": "r5ad.large",
            "onDemandPrice": 0.154,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 0.0346
                },
                {
                    "zone": "eu-west-2b",
                    "price": 0.0346
                },
                {
                    "zone": "eu-west-2c",
                    "price": 0.0346
                }
            ],
            "cpusPerVm": 2,
            "memPerVm": 16,
            "gpusPerVm": 0,
            "ntwPerf": "10 Gigabit",
            "ntwPerfCategory": "high",
            "zones": null,
            "attributes": {
                "cpu": "2",
                "instanceTypeCategory": "Memory optimized",
                "memory": "16",
                "networkPerfCategory": "high"
            },
            "currentGen": true
        },
        {
            "category": "Memory optimized",
            "type": "z1d.2xlarge",
            "onDemandPrice": 0.879,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 0.2637
                },
                {
                    "zone": "eu-west-2b",
                    "price": 0.2637
                },
                {
                    "zone": "eu-west-2c",
                    "price": 0.2637
                }
            ],
            "cpusPerVm": 8,
            "memPerVm": 64,
            "gpusPerVm": 0,
            "ntwPerf": "Up to 10 Gigabit",
            "ntwPerfCategory": "high",
            "zones": null,
            "attributes": {
                "cpu": "8",
                "instanceTypeCategory": "Memory optimized",
                "memory": "64",
                "networkPerfCategory": "high"
            },
            "currentGen": true
        },
        {
            "category": "General purpose",
            "type": "m5.4xlarge",
            "onDemandPrice": 0.888,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 0.2639
                },
                {
                    "zone": "eu-west-2b",
                    "price": 0.2639
                },
                {
                    "zone": "eu-west-2c",
                    "price": 0.2639
                }
            ],
            "cpusPerVm": 16,
            "memPerVm": 64,
            "gpusPerVm": 0,
            "ntwPerf": "Up to 10 Gigabit",
            "ntwPerfCategory": "high",
            "zones": null,
            "attributes": {
                "cpu": "16",
                "instanceTypeCategory": "General purpose",
                "memory": "64",
                "networkPerfCategory": "high"
            },
            "currentGen": true
        },
        {
            "category": "General purpose",
            "type": "t3.small",
            "onDemandPrice": 0.0236,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 0.0071
                },
                {
                    "zone": "eu-west-2b",
                    "price": 0.0071
                },
                {
                    "zone": "eu-west-2c",
                    "price": 0.0071
                }
            ],
            "cpusPerVm": 2,
            "memPerVm": 2,
            "gpusPerVm": 0,
            "ntwPerf": "Low to Moderate",
            "ntwPerfCategory": "low",
            "zones": null,
            "attributes": {
                "cpu": "2",
                "instanceTypeCategory": "General purpose",
                "memory": "2",
                "networkPerfCategory": "low"
            },
            "currentGen": true,
            "burst": true
        },
        {
            "category": "Storage optimized",
            "type": "i3.4xlarge",
            "onDemandPrice": 1.448,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 0.4344
                },
                {
                    "zone": "eu-west-2b",
                    "price": 1.448
                },
                {
                    "zone": "eu-west-2c",
                    "price": 0.4344
                }
            ],
            "cpusPerVm": 16,
            "memPerVm": 122,
            "gpusPerVm": 0,
            "ntwPerf": "Up to 10 Gigabit",
            "ntwPerfCategory": "high",
            "zones": null,
            "attributes": {
                "cpu": "16",
                "instanceTypeCategory": "Storage optimized",
                "memory": "122",
                "networkPerfCategory": "high"
            },
            "currentGen": true
        },
        {
            "category": "General purpose",
            "type": "m5.12xlarge",
            "onDemandPrice": 2.664,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 0.7918
                },
                {
                    "zone": "eu-west-2b",
                    "price": 0.7984
                },
                {
                    "zone": "eu-west-2c",
                    "price": 0.7918
                }
            ],
            "cpusPerVm": 48,
            "memPerVm": 192,
            "gpusPerVm": 0,
            "ntwPerf": "10 Gigabit",
            "ntwPerfCategory": "high",
            "zones": null,
            "attributes": {
                "cpu": "48",
                "instanceTypeCategory": "General purpose",
                "memory": "192",
                "networkPerfCategory": "high"
            },
            "currentGen": true
        },
        {
            "category": "Compute optimized",
            "type": "c5d.24xlarge",
            "onDemandPrice": 5.52,
            "spotPrice": [
                {
                    "zone": "eu-west-2c",
                    "price": 1.5081
                },
                {
                    "zone": "eu-west-2a",
                    "price": 1.5081
                },
                {
                    "zone": "eu-west-2b",
                    "price": 1.5081
                }
            ],
            "cpusPerVm": 96,
            "memPerVm": 192,
            "gpusPerVm": 0,
            "ntwPerf": "25 Gigabit",
            "ntwPerfCategory": "extra",
            "zones": null,
            "attributes": {
                "cpu": "96",
                "instanceTypeCategory": "Compute optimized",
                "memory": "192",
                "networkPerfCategory": "extra"
            },
            "currentGen": true
        },
        {
            "category": "General purpose",
            "type": "m5a.large",
            "onDemandPrice": 0.1,
            "spotPrice": [
                {
                    "zone": "eu-west-2b",
                    "price": 0.033
                },
                {
                    "zone": "eu-west-2c",
                    "price": 0.033
                },
                {
                    "zone": "eu-west-2a",
                    "price": 0.033
                }
            ],
            "cpusPerVm": 2,
            "memPerVm": 8,
            "gpusPerVm": 0,
            "ntwPerf": "Up to 10 Gigabit",
            "ntwPerfCategory": "high",
            "zones": null,
            "attributes": {
                "cpu": "2",
                "instanceTypeCategory": "General purpose",
                "memory": "8",
                "networkPerfCategory": "high"
            },
            "currentGen": true
        },
        {
            "category": "Compute optimized",
            "type": "c5.18xlarge",
            "onDemandPrice": 3.636,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 1.1311
                },
                {
                    "zone": "eu-west-2b",
                    "price": 1.1584
                },
                {
                    "zone": "eu-west-2c",
                    "price": 1.1311
                }
            ],
            "cpusPerVm": 72,
            "memPerVm": 144,
            "gpusPerVm": 0,
            "ntwPerf": "25 Gigabit",
            "ntwPerfCategory": "extra",
            "zones": null,
            "attributes": {
                "cpu": "72",
                "instanceTypeCategory": "Compute optimized",
                "memory": "144",
                "networkPerfCategory": "extra"
            },
            "currentGen": true
        },
        {
            "category": "General purpose",
            "type": "t3.2xlarge",
            "onDemandPrice": 0.3776,
            "spotPrice": [
                {
                    "zone": "eu-west-2b",
                    "price": 0.1133
                },
                {
                    "zone": "eu-west-2c",
                    "price": 0.1133
                },
                {
                    "zone": "eu-west-2a",
                    "price": 0.1133
                }
            ],
            "cpusPerVm": 8,
            "memPerVm": 32,
            "gpusPerVm": 0,
            "ntwPerf": "Moderate",
            "ntwPerfCategory": "medium",
            "zones": null,
            "attributes": {
                "cpu": "8",
                "instanceTypeCategory": "General purpose",
                "memory": "32",
                "networkPerfCategory": "medium"
            },
            "currentGen": true,
            "burst": true
        },
        {
            "category": "General purpose",
            "type": "t2.large",
            "onDemandPrice": 0.1056,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 0.0317
                },
                {
                    "zone": "eu-west-2b",
                    "price": 0.0317
                },
                {
                    "zone": "eu-west-2c",
                    "price": 0.0336
                }
            ],
            "cpusPerVm": 2,
            "memPerVm": 8,
            "gpusPerVm": 0,
            "ntwPerf": "Low to Moderate",
            "ntwPerfCategory": "low",
            "zones": null,
            "attributes": {
                "cpu": "2",
                "instanceTypeCategory": "General purpose",
                "memory": "8",
                "networkPerfCategory": "low"
            },
            "currentGen": true,
            "burst": true
        },
        {
            "category": "GPU instance",
            "type": "p3.16xlarge",
            "onDemandPrice": 28.712,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 8.6136
                },
                {
                    "zone": "eu-west-2b",
                    "price": 8.6136
                }
            ],
            "cpusPerVm": 64,
            "memPerVm": 488,
            "gpusPerVm": 8,
            "ntwPerf": "25 Gigabit",
            "ntwPerfCategory": "extra",
            "zones": null,
            "attributes": {
                "cpu": "64",
                "instanceTypeCategory": "GPU instance",
                "memory": "488",
                "networkPerfCategory": "extra"
            },
            "currentGen": true
        },
        {
            "category": "Memory optimized",
            "type": "z1d.3xlarge",
            "onDemandPrice": 1.318,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 0.3954
                },
                {
                    "zone": "eu-west-2b",
                    "price": 0.3954
                },
                {
                    "zone": "eu-west-2c",
                    "price": 0.3954
                }
            ],
            "cpusPerVm": 12,
            "memPerVm": 96,
            "gpusPerVm": 0,
            "ntwPerf": "Up to 10 Gigabit",
            "ntwPerfCategory": "high",
            "zones": null,
            "attributes": {
                "cpu": "12",
                "instanceTypeCategory": "Memory optimized",
                "memory": "96",
                "networkPerfCategory": "high"
            },
            "currentGen": true
        },
        {
            "category": "Memory optimized",
            "type": "r5ad.2xlarge",
            "onDemandPrice": 0.616,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 0.616
                },
                {
                    "zone": "eu-west-2b",
                    "price": 0.616
                },
                {
                    "zone": "eu-west-2c",
                    "price": 0.1382
                }
            ],
            "cpusPerVm": 8,
            "memPerVm": 64,
            "gpusPerVm": 0,
            "ntwPerf": "10 Gigabit",
            "ntwPerfCategory": "high",
            "zones": null,
            "attributes": {
                "cpu": "8",
                "instanceTypeCategory": "Memory optimized",
                "memory": "64",
                "networkPerfCategory": "high"
            },
            "currentGen": true
        },
        {
            "category": "GPU instance",
            "type": "g4dn.4xlarge",
            "onDemandPrice": 1.409,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 0.4227
                },
                {
                    "zone": "eu-west-2b",
                    "price": 0.4227
                },
                {
                    "zone": "eu-west-2c",
                    "price": 0.4227
                }
            ],
            "cpusPerVm": 16,
            "memPerVm": 64,
            "gpusPerVm": 1,
            "ntwPerf": "Up to 25 Gigabit",
            "ntwPerfCategory": "extra",
            "zones": null,
            "attributes": {
                "cpu": "16",
                "instanceTypeCategory": "GPU instance",
                "memory": "64",
                "networkPerfCategory": "extra"
            },
            "currentGen": true
        },
        {
            "category": "Compute optimized",
            "type": "c4.large",
            "onDemandPrice": 0.119,
            "spotPrice": [
                {
                    "zone": "eu-west-2b",
                    "price": 0.03
                },
                {
                    "zone": "eu-west-2c",
                    "price": 0.0532
                },
                {
                    "zone": "eu-west-2a",
                    "price": 0.0299
                }
            ],
            "cpusPerVm": 2,
            "memPerVm": 3.75,
            "gpusPerVm": 0,
            "ntwPerf": "Moderate",
            "ntwPerfCategory": "medium",
            "zones": null,
            "attributes": {
                "cpu": "2",
                "instanceTypeCategory": "Compute optimized",
                "memory": "3.75",
                "networkPerfCategory": "medium"
            },
            "currentGen": true
        },
        {
            "category": "Memory optimized",
            "type": "r5ad.8xlarge",
            "onDemandPrice": 2.464,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 0.553
                },
                {
                    "zone": "eu-west-2b",
                    "price": 0.553
                },
                {
                    "zone": "eu-west-2c",
                    "price": 0.553
                }
            ],
            "cpusPerVm": 32,
            "memPerVm": 256,
            "gpusPerVm": 0,
            "ntwPerf": "Up to 10 Gigabit",
            "ntwPerfCategory": "high",
            "zones": null,
            "attributes": {
                "cpu": "32",
                "instanceTypeCategory": "Memory optimized",
                "memory": "256",
                "networkPerfCategory": "high"
            },
            "currentGen": true
        },
        {
            "category": "Storage optimized",
            "type": "i3.large",
            "onDemandPrice": 0.181,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 0.0543
                },
                {
                    "zone": "eu-west-2b",
                    "price": 0.181
                },
                {
                    "zone": "eu-west-2c",
                    "price": 0.0543
                }
            ],
            "cpusPerVm": 2,
            "memPerVm": 15.25,
            "gpusPerVm": 0,
            "ntwPerf": "Up to 10 Gigabit",
            "ntwPerfCategory": "high",
            "zones": null,
            "attributes": {
                "cpu": "2",
                "instanceTypeCategory": "Storage optimized",
                "memory": "15.25",
                "networkPerfCategory": "high"
            },
            "currentGen": true
        },
        {
            "category": "Storage optimized",
            "type": "i3.16xlarge",
            "onDemandPrice": 5.792,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 1.7376
                },
                {
                    "zone": "eu-west-2b",
                    "price": 5.792
                },
                {
                    "zone": "eu-west-2c",
                    "price": 1.7376
                }
            ],
            "cpusPerVm": 64,
            "memPerVm": 488,
            "gpusPerVm": 0,
            "ntwPerf": "20 Gigabit",
            "ntwPerfCategory": "extra",
            "zones": null,
            "attributes": {
                "cpu": "64",
                "instanceTypeCategory": "Storage optimized",
                "memory": "488",
                "networkPerfCategory": "extra"
            },
            "currentGen": true
        },
        {
            "category": "Storage optimized",
            "type": "i3en.2xlarge",
            "onDemandPrice": 1.052,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 0.3156
                },
                {
                    "zone": "eu-west-2b",
                    "price": 1.052
                },
                {
                    "zone": "eu-west-2c",
                    "price": 1.052
                }
            ],
            "cpusPerVm": 8,
            "memPerVm": 64,
            "gpusPerVm": 0,
            "ntwPerf": "Up to 25 Gigabit",
            "ntwPerfCategory": "extra",
            "zones": null,
            "attributes": {
                "cpu": "8",
                "instanceTypeCategory": "Storage optimized",
                "memory": "64",
                "networkPerfCategory": "extra"
            },
            "currentGen": true
        },
        {
            "category": "Compute optimized",
            "type": "c5.12xlarge",
            "onDemandPrice": 2.424,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 0.754
                },
                {
                    "zone": "eu-west-2b",
                    "price": 0.754
                },
                {
                    "zone": "eu-west-2c",
                    "price": 0.754
                }
            ],
            "cpusPerVm": 48,
            "memPerVm": 96,
            "gpusPerVm": 0,
            "ntwPerf": "12 Gigabit",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": {
                "cpu": "48",
                "instanceTypeCategory": "Compute optimized",
                "memory": "96",
                "networkPerfCategory": ""
            },
            "currentGen": true
        },
        {
            "category": "Memory optimized",
            "type": "r5.8xlarge",
            "onDemandPrice": 2.368,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 0.5686
                },
                {
                    "zone": "eu-west-2b",
                    "price": 0.6658
                },
                {
                    "zone": "eu-west-2c",
                    "price": 0.906
                }
            ],
            "cpusPerVm": 32,
            "memPerVm": 256,
            "gpusPerVm": 0,
            "ntwPerf": "10 Gigabit",
            "ntwPerfCategory": "high",
            "zones": null,
            "attributes": {
                "cpu": "32",
                "instanceTypeCategory": "Memory optimized",
                "memory": "256",
                "networkPerfCategory": "high"
            },
            "currentGen": true
        },
        {
            "category": "Storage optimized",
            "type": "d2.xlarge",
            "onDemandPrice": 0.772,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 0.2316
                },
                {
                    "zone": "eu-west-2b",
                    "price": 0.2316
                },
                {
                    "zone": "eu-west-2c",
                    "price": 0.772
                }
            ],
            "cpusPerVm": 4,
            "memPerVm": 30.5,
            "gpusPerVm": 0,
            "ntwPerf": "Moderate",
            "ntwPerfCategory": "medium",
            "zones": null,
            "attributes": {
                "cpu": "4",
                "instanceTypeCategory": "Storage optimized",
                "memory": "30.5",
                "networkPerfCategory": "medium"
            },
            "currentGen": true
        },
        {
            "category": "General purpose",
            "type": "t3a.nano",
            "onDemandPrice": 0.0053,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 0.0016
                },
                {
                    "zone": "eu-west-2b",
                    "price": 0.0016
                },
                {
                    "zone": "eu-west-2c",
                    "price": 0.0016
                }
            ],
            "cpusPerVm": 2,
            "memPerVm": 0.5,
            "gpusPerVm": 0,
            "ntwPerf": "Low",
            "ntwPerfCategory": "low",
            "zones": null,
            "attributes": {
                "cpu": "2",
                "instanceTypeCategory": "General purpose",
                "memory": "0.5",
                "networkPerfCategory": "low"
            },
            "currentGen": true,
            "burst": true
        },
        {
            "category": "GPU instance",
            "type": "g3.16xlarge",
            "onDemandPrice": 5.716,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 1.7148
                },
                {
                    "zone": "eu-west-2b",
                    "price": 1.7148
                }
            ],
            "cpusPerVm": 64,
            "memPerVm": 488,
            "gpusPerVm": 4,
            "ntwPerf": "20 Gigabit",
            "ntwPerfCategory": "extra",
            "zones": null,
            "attributes": {
                "cpu": "64",
                "instanceTypeCategory": "GPU instance",
                "memory": "488",
                "networkPerfCategory": "extra"
            },
            "currentGen": true
        },
        {
            "category": "Memory optimized",
            "type": "r5a.12xlarge",
            "onDemandPrice": 3.192,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 3.192
                },
                {
                    "zone": "eu-west-2b",
                    "price": 0.8295
                },
                {
                    "zone": "eu-west-2c",
                    "price": 0.8295
                }
            ],
            "cpusPerVm": 48,
            "memPerVm": 384,
            "gpusPerVm": 0,
            "ntwPerf": "10 Gigabit",
            "ntwPerfCategory": "high",
            "zones": null,
            "attributes": {
                "cpu": "48",
                "instanceTypeCategory": "Memory optimized",
                "memory": "384",
                "networkPerfCategory": "high"
            },
            "currentGen": true
        },
        {
            "category": "Memory optimized",
            "type": "r5.xlarge",
            "onDemandPrice": 0.296,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 0.0851
                },
                {
                    "zone": "eu-west-2b",
                    "price": 0.0766
                },
                {
                    "zone": "eu-west-2c",
                    "price": 0.096
                }
            ],
            "cpusPerVm": 4,
            "memPerVm": 32,
            "gpusPerVm": 0,
            "ntwPerf": "10 Gigabit",
            "ntwPerfCategory": "high",
            "zones": null,
            "attributes": {
                "cpu": "4",
                "instanceTypeCategory": "Memory optimized",
                "memory": "32",
                "networkPerfCategory": "high"
            },
            "currentGen": true
        },
        {
            "category": "General purpose",
            "type": "m5d.metal",
            "onDemandPrice": 6.288,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 1.5835
                },
                {
                    "zone": "eu-west-2b",
                    "price": 1.5835
                },
                {
                    "zone": "eu-west-2c",
                    "price": 1.5835
                }
            ],
            "cpusPerVm": 96,
            "memPerVm": 384,
            "gpusPerVm": 0,
            "ntwPerf": "25 Gigabit",
            "ntwPerfCategory": "extra",
            "zones": null,
            "attributes": {
                "cpu": "96",
                "instanceTypeCategory": "General purpose",
                "memory": "384",
                "networkPerfCategory": "extra"
            },
            "currentGen": true
        },
        {
            "category": "Memory optimized",
            "type": "r4.8xlarge",
            "onDemandPrice": 2.496,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 0.5266
                },
                {
                    "zone": "eu-west-2b",
                    "price": 0.5297
                },
                {
                    "zone": "eu-west-2c",
                    "price": 0.5301
                }
            ],
            "cpusPerVm": 32,
            "memPerVm": 244,
            "gpusPerVm": 0,
            "ntwPerf": "10 Gigabit",
            "ntwPerfCategory": "high",
            "zones": null,
            "attributes": {
                "cpu": "32",
                "instanceTypeCategory": "Memory optimized",
                "memory": "244",
                "networkPerfCategory": "high"
            },
            "currentGen": true
        },
        {
            "category": "Memory optimized",
            "type": "r5d.metal",
            "onDemandPrice": 8.112,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 8.112
                },
                {
                    "zone": "eu-west-2b",
                    "price": 8.112
                },
                {
                    "zone": "eu-west-2c",
                    "price": 8.112
                }
            ],
            "cpusPerVm": 96,
            "memPerVm": 768,
            "gpusPerVm": 0,
            "ntwPerf": "25 Gigabit",
            "ntwPerfCategory": "extra",
            "zones": null,
            "attributes": {
                "cpu": "96",
                "instanceTypeCategory": "Memory optimized",
                "memory": "768",
                "networkPerfCategory": "extra"
            },
            "currentGen": true
        },
        {
            "category": "Compute optimized",
            "type": "c4.2xlarge",
            "onDemandPrice": 0.476,
            "spotPrice": [
                {
                    "zone": "eu-west-2c",
                    "price": 0.1335
                },
                {
                    "zone": "eu-west-2a",
                    "price": 0.1197
                },
                {
                    "zone": "eu-west-2b",
                    "price": 0.1197
                }
            ],
            "cpusPerVm": 8,
            "memPerVm": 15,
            "gpusPerVm": 0,
            "ntwPerf": "High",
            "ntwPerfCategory": "medium",
            "zones": null,
            "attributes": {
                "cpu": "8",
                "instanceTypeCategory": "Compute optimized",
                "memory": "15",
                "networkPerfCategory": "medium"
            },
            "currentGen": true
        },
        {
            "category": "General purpose",
            "type": "t2.nano",
            "onDemandPrice": 0.0066,
            "spotPrice": null,
            "cpusPerVm": 1,
            "memPerVm": 0.5,
            "gpusPerVm": 0,
            "ntwPerf": "Low",
            "ntwPerfCategory": "low",
            "zones": null,
            "attributes": {
                "cpu": "1",
                "instanceTypeCategory": "General purpose",
                "memory": "0.5",
                "networkPerfCategory": "low"
            },
            "currentGen": true,
            "burst": true
        },
        {
            "category": "GPU instance",
            "type": "g4dn.xlarge",
            "onDemandPrice": 0.615,
            "spotPrice": [
                {
                    "zone": "eu-west-2b",
                    "price": 0.2057
                },
                {
                    "zone": "eu-west-2c",
                    "price": 0.1845
                },
                {
                    "zone": "eu-west-2a",
                    "price": 0.1871
                }
            ],
            "cpusPerVm": 4,
            "memPerVm": 16,
            "gpusPerVm": 1,
            "ntwPerf": "Up to 25 Gigabit",
            "ntwPerfCategory": "extra",
            "zones": null,
            "attributes": {
                "cpu": "4",
                "instanceTypeCategory": "GPU instance",
                "memory": "16",
                "networkPerfCategory": "extra"
            },
            "currentGen": true
        },
        {
            "category": "Storage optimized",
            "type": "i3en.xlarge",
            "onDemandPrice": 0.526,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 0.1578
                },
                {
                    "zone": "eu-west-2b",
                    "price": 0.1578
                },
                {
                    "zone": "eu-west-2c",
                    "price": 0.1578
                }
            ],
            "cpusPerVm": 4,
            "memPerVm": 32,
            "gpusPerVm": 0,
            "ntwPerf": "Up to 25 Gigabit",
            "ntwPerfCategory": "extra",
            "zones": null,
            "attributes": {
                "cpu": "4",
                "instanceTypeCategory": "Storage optimized",
                "memory": "32",
                "networkPerfCategory": "extra"
            },
            "currentGen": true
        },
        {
            "category": "Storage optimized",
            "type": "i3.xlarge",
            "onDemandPrice": 0.362,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 0.1086
                },
                {
                    "zone": "eu-west-2b",
                    "price": 0.362
                },
                {
                    "zone": "eu-west-2c",
                    "price": 0.1086
                }
            ],
            "cpusPerVm": 4,
            "memPerVm": 30.5,
            "gpusPerVm": 0,
            "ntwPerf": "Up to 10 Gigabit",
            "ntwPerfCategory": "high",
            "zones": null,
            "attributes": {
                "cpu": "4",
                "instanceTypeCategory": "Storage optimized",
                "memory": "30.5",
                "networkPerfCategory": "high"
            },
            "currentGen": true
        },
        {
            "category": "Memory optimized",
            "type": "r5d.xlarge",
            "onDemandPrice": 0.338,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 0.0691
                },
                {
                    "zone": "eu-west-2b",
                    "price": 0.0691
                },
                {
                    "zone": "eu-west-2c",
                    "price": 0.0691
                }
            ],
            "cpusPerVm": 4,
            "memPerVm": 32,
            "gpusPerVm": 0,
            "ntwPerf": "10 Gigabit",
            "ntwPerfCategory": "high",
            "zones": null,
            "attributes": {
                "cpu": "4",
                "instanceTypeCategory": "Memory optimized",
                "memory": "32",
                "networkPerfCategory": "high"
            },
            "currentGen": true
        },
        {
            "category": "GPU instance",
            "type": "g3s.xlarge",
            "onDemandPrice": 0.94,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 0.282
                },
                {
                    "zone": "eu-west-2b",
                    "price": 0.282
                }
            ],
            "cpusPerVm": 4,
            "memPerVm": 30.5,
            "gpusPerVm": 1,
            "ntwPerf": "10 Gigabit",
            "ntwPerfCategory": "high",
            "zones": null,
            "attributes": {
                "cpu": "4",
                "instanceTypeCategory": "GPU instance",
                "memory": "30.5",
                "networkPerfCategory": "high"
            },
            "currentGen": true
        },
        {
            "category": "General purpose",
            "type": "m5a.24xlarge",
            "onDemandPrice": 4.8,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 1.5835
                },
                {
                    "zone": "eu-west-2b",
                    "price": 1.5835
                },
                {
                    "zone": "eu-west-2c",
                    "price": 4.8
                }
            ],
            "cpusPerVm": 96,
            "memPerVm": 384,
            "gpusPerVm": 0,
            "ntwPerf": "20 Gigabit",
            "ntwPerfCategory": "extra",
            "zones": null,
            "attributes": {
                "cpu": "96",
                "instanceTypeCategory": "General purpose",
                "memory": "384",
                "networkPerfCategory": "extra"
            },
            "currentGen": true
        },
        {
            "category": "General purpose",
            "type": "t3a.medium",
            "onDemandPrice": 0.0425,
            "spotPrice": [
                {
                    "zone": "eu-west-2c",
                    "price": 0.0128
                },
                {
                    "zone": "eu-west-2a",
                    "price": 0.0128
                },
                {
                    "zone": "eu-west-2b",
                    "price": 0.0128
                }
            ],
            "cpusPerVm": 2,
            "memPerVm": 4,
            "gpusPerVm": 0,
            "ntwPerf": "Low to Moderate",
            "ntwPerfCategory": "low",
            "zones": null,
            "attributes": {
                "cpu": "2",
                "instanceTypeCategory": "General purpose",
                "memory": "4",
                "networkPerfCategory": "low"
            },
            "currentGen": true,
            "burst": true
        },
        {
            "category": "Memory optimized",
            "type": "r4.2xlarge",
            "onDemandPrice": 0.624,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 0.1317
                },
                {
                    "zone": "eu-west-2b",
                    "price": 0.1317
                },
                {
                    "zone": "eu-west-2c",
                    "price": 0.1449
                }
            ],
            "cpusPerVm": 8,
            "memPerVm": 61,
            "gpusPerVm": 0,
            "ntwPerf": "Up to 10 Gigabit",
            "ntwPerfCategory": "high",
            "zones": null,
            "attributes": {
                "cpu": "8",
                "instanceTypeCategory": "Memory optimized",
                "memory": "61",
                "networkPerfCategory": "high"
            },
            "currentGen": true
        },
        {
            "category": "General purpose",
            "type": "t3.large",
            "onDemandPrice": 0.0944,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 0.0322
                },
                {
                    "zone": "eu-west-2b",
                    "price": 0.0289
                },
                {
                    "zone": "eu-west-2c",
                    "price": 0.0283
                }
            ],
            "cpusPerVm": 2,
            "memPerVm": 8,
            "gpusPerVm": 0,
            "ntwPerf": "Low to Moderate",
            "ntwPerfCategory": "low",
            "zones": null,
            "attributes": {
                "cpu": "2",
                "instanceTypeCategory": "General purpose",
                "memory": "8",
                "networkPerfCategory": "low"
            },
            "currentGen": true,
            "burst": true
        },
        {
            "category": "General purpose",
            "type": "m5d.16xlarge",
            "onDemandPrice": 4.192,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 1.0557
                },
                {
                    "zone": "eu-west-2b",
                    "price": 1.0557
                },
                {
                    "zone": "eu-west-2c",
                    "price": 1.0557
                }
            ],
            "cpusPerVm": 64,
            "memPerVm": 256,
            "gpusPerVm": 0,
            "ntwPerf": "20 Gigabit",
            "ntwPerfCategory": "extra",
            "zones": null,
            "attributes": {
                "cpu": "64",
                "instanceTypeCategory": "General purpose",
                "memory": "256",
                "networkPerfCategory": "extra"
            },
            "currentGen": true
        },
        {
            "category": "Memory optimized",
            "type": "r5a.4xlarge",
            "onDemandPrice": 1.064,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 0.2765
                },
                {
                    "zone": "eu-west-2b",
                    "price": 0.2765
                },
                {
                    "zone": "eu-west-2c",
                    "price": 0.2765
                }
            ],
            "cpusPerVm": 16,
            "memPerVm": 128,
            "gpusPerVm": 0,
            "ntwPerf": "10 Gigabit",
            "ntwPerfCategory": "high",
            "zones": null,
            "attributes": {
                "cpu": "16",
                "instanceTypeCategory": "Memory optimized",
                "memory": "128",
                "networkPerfCategory": "high"
            },
            "currentGen": true
        },
        {
            "category": "General purpose",
            "type": "m5d.12xlarge",
            "onDemandPrice": 3.144,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 0.7918
                },
                {
                    "zone": "eu-west-2b",
                    "price": 0.7918
                },
                {
                    "zone": "eu-west-2c",
                    "price": 0.7918
                }
            ],
            "cpusPerVm": 48,
            "memPerVm": 192,
            "gpusPerVm": 0,
            "ntwPerf": "10 Gigabit",
            "ntwPerfCategory": "high",
            "zones": null,
            "attributes": {
                "cpu": "48",
                "instanceTypeCategory": "General purpose",
                "memory": "192",
                "networkPerfCategory": "high"
            },
            "currentGen": true
        },
        {
            "category": "Memory optimized",
            "type": "r4.4xlarge",
            "onDemandPrice": 1.248,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 0.2633
                },
                {
                    "zone": "eu-west-2b",
                    "price": 0.2633
                },
                {
                    "zone": "eu-west-2c",
                    "price": 0.2633
                }
            ],
            "cpusPerVm": 16,
            "memPerVm": 122,
            "gpusPerVm": 0,
            "ntwPerf": "Up to 10 Gigabit",
            "ntwPerfCategory": "high",
            "zones": null,
            "attributes": {
                "cpu": "16",
                "instanceTypeCategory": "Memory optimized",
                "memory": "122",
                "networkPerfCategory": "high"
            },
            "currentGen": true
        },
        {
            "category": "General purpose",
            "type": "m5a.8xlarge",
            "onDemandPrice": 1.6,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 0.5278
                },
                {
                    "zone": "eu-west-2b",
                    "price": 0.5278
                },
                {
                    "zone": "eu-west-2c",
                    "price": 1.6
                }
            ],
            "cpusPerVm": 32,
            "memPerVm": 128,
            "gpusPerVm": 0,
            "ntwPerf": "Up to 10 Gigabit",
            "ntwPerfCategory": "high",
            "zones": null,
            "attributes": {
                "cpu": "32",
                "instanceTypeCategory": "General purpose",
                "memory": "128",
                "networkPerfCategory": "high"
            },
            "currentGen": true
        },
        {
            "category": "Memory optimized",
            "type": "r4.large",
            "onDemandPrice": 0.156,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 0.0329
                },
                {
                    "zone": "eu-west-2b",
                    "price": 0.033
                },
                {
                    "zone": "eu-west-2c",
                    "price": 0.0329
                }
            ],
            "cpusPerVm": 2,
            "memPerVm": 15.25,
            "gpusPerVm": 0,
            "ntwPerf": "Up to 10 Gigabit",
            "ntwPerfCategory": "high",
            "zones": null,
            "attributes": {
                "cpu": "2",
                "instanceTypeCategory": "Memory optimized",
                "memory": "15.25",
                "networkPerfCategory": "high"
            },
            "currentGen": true
        },
        {
            "category": "General purpose",
            "type": "m5a.xlarge",
            "onDemandPrice": 0.2,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 0.0673
                },
                {
                    "zone": "eu-west-2b",
                    "price": 0.066
                },
                {
                    "zone": "eu-west-2c",
                    "price": 0.2
                }
            ],
            "cpusPerVm": 4,
            "memPerVm": 16,
            "gpusPerVm": 0,
            "ntwPerf": "Up to 10 Gigabit",
            "ntwPerfCategory": "high",
            "zones": null,
            "attributes": {
                "cpu": "4",
                "instanceTypeCategory": "General purpose",
                "memory": "16",
                "networkPerfCategory": "high"
            },
            "currentGen": true
        },
        {
            "category": "Memory optimized",
            "type": "r5d.large",
            "onDemandPrice": 0.169,
            "spotPrice": [
                {
                    "zone": "eu-west-2c",
                    "price": 0.0346
                },
                {
                    "zone": "eu-west-2a",
                    "price": 0.0346
                },
                {
                    "zone": "eu-west-2b",
                    "price": 0.0346
                }
            ],
            "cpusPerVm": 2,
            "memPerVm": 16,
            "gpusPerVm": 0,
            "ntwPerf": "10 Gigabit",
            "ntwPerfCategory": "high",
            "zones": null,
            "attributes": {
                "cpu": "2",
                "instanceTypeCategory": "Memory optimized",
                "memory": "16",
                "networkPerfCategory": "high"
            },
            "currentGen": true
        },
        {
            "category": "Compute optimized",
            "type": "c5.large",
            "onDemandPrice": 0.101,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 0.0328
                },
                {
                    "zone": "eu-west-2b",
                    "price": 0.0356
                },
                {
                    "zone": "eu-west-2c",
                    "price": 0.0322
                }
            ],
            "cpusPerVm": 2,
            "memPerVm": 4,
            "gpusPerVm": 0,
            "ntwPerf": "Up to 10 Gigabit",
            "ntwPerfCategory": "high",
            "zones": null,
            "attributes": {
                "cpu": "2",
                "instanceTypeCategory": "Compute optimized",
                "memory": "4",
                "networkPerfCategory": "high"
            },
            "currentGen": true
        },
        {
            "category": "Memory optimized",
            "type": "r5a.8xlarge",
            "onDemandPrice": 2.128,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 0.553
                },
                {
                    "zone": "eu-west-2b",
                    "price": 0.553
                },
                {
                    "zone": "eu-west-2c",
                    "price": 0.553
                }
            ],
            "cpusPerVm": 32,
            "memPerVm": 256,
            "gpusPerVm": 0,
            "ntwPerf": "Up to 10 Gigabit",
            "ntwPerfCategory": "high",
            "zones": null,
            "attributes": {
                "cpu": "32",
                "instanceTypeCategory": "Memory optimized",
                "memory": "256",
                "networkPerfCategory": "high"
            },
            "currentGen": true
        },
        {
            "category": "Memory optimized",
            "type": "z1d.metal",
            "onDemandPrice": 5.273,
            "spotPrice": [
                {
                    "zone": "eu-west-2c",
                    "price": 1.5819
                },
                {
                    "zone": "eu-west-2a",
                    "price": 1.5819
                },
                {
                    "zone": "eu-west-2b",
                    "price": 1.5819
                }
            ],
            "cpusPerVm": 48,
            "memPerVm": 384,
            "gpusPerVm": 0,
            "ntwPerf": "25 Gigabit",
            "ntwPerfCategory": "extra",
            "zones": null,
            "attributes": {
                "cpu": "48",
                "instanceTypeCategory": "Memory optimized",
                "memory": "384",
                "networkPerfCategory": "extra"
            },
            "currentGen": true
        },
        {
            "category": "General purpose",
            "type": "m5ad.2xlarge",
            "onDemandPrice": 0.48,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 0.132
                },
                {
                    "zone": "eu-west-2b",
                    "price": 0.132
                },
                {
                    "zone": "eu-west-2c",
                    "price": 0.132
                }
            ],
            "cpusPerVm": 8,
            "memPerVm": 32,
            "gpusPerVm": 0,
            "ntwPerf": "Up to 10 Gigabit",
            "ntwPerfCategory": "high",
            "zones": null,
            "attributes": {
                "cpu": "8",
                "instanceTypeCategory": "General purpose",
                "memory": "32",
                "networkPerfCategory": "high"
            },
            "currentGen": true
        },
        {
            "category": "Memory optimized",
            "type": "r5a.xlarge",
            "onDemandPrice": 0.266,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 0.0691
                },
                {
                    "zone": "eu-west-2b",
                    "price": 0.0691
                },
                {
                    "zone": "eu-west-2c",
                    "price": 0.0691
                }
            ],
            "cpusPerVm": 4,
            "memPerVm": 32,
            "gpusPerVm": 0,
            "ntwPerf": "10 Gigabit",
            "ntwPerfCategory": "high",
            "zones": null,
            "attributes": {
                "cpu": "4",
                "instanceTypeCategory": "Memory optimized",
                "memory": "32",
                "networkPerfCategory": "high"
            },
            "currentGen": true
        },
        {
            "category": "Storage optimized",
            "type": "i3en.24xlarge",
            "onDemandPrice": 12.624,
            "spotPrice": [
                {
                    "zone": "eu-west-2c",
                    "price": 3.7872
                },
                {
                    "zone": "eu-west-2a",
                    "price": 3.7872
                },
                {
                    "zone": "eu-west-2b",
                    "price": 3.7872
                }
            ],
            "cpusPerVm": 96,
            "memPerVm": 768,
            "gpusPerVm": 0,
            "ntwPerf": "100 Gigabit",
            "ntwPerfCategory": "extra",
            "zones": null,
            "attributes": {
                "cpu": "96",
                "instanceTypeCategory": "Storage optimized",
                "memory": "768",
                "networkPerfCategory": "extra"
            },
            "currentGen": true
        },
        {
            "category": "Compute optimized",
            "type": "c5d.large",
            "onDemandPrice": 0.115,
            "spotPrice": [
                {
                    "zone": "eu-west-2b",
                    "price": 0.0314
                },
                {
                    "zone": "eu-west-2c",
                    "price": 0.0314
                },
                {
                    "zone": "eu-west-2a",
                    "price": 0.0314
                }
            ],
            "cpusPerVm": 2,
            "memPerVm": 4,
            "gpusPerVm": 0,
            "ntwPerf": "Up to 10 Gigabit",
            "ntwPerfCategory": "high",
            "zones": null,
            "attributes": {
                "cpu": "2",
                "instanceTypeCategory": "Compute optimized",
                "memory": "4",
                "networkPerfCategory": "high"
            },
            "currentGen": true
        },
        {
            "category": "General purpose",
            "type": "t3a.2xlarge",
            "onDemandPrice": 0.3398,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 0.1019
                },
                {
                    "zone": "eu-west-2b",
                    "price": 0.1019
                },
                {
                    "zone": "eu-west-2c",
                    "price": 0.1019
                }
            ],
            "cpusPerVm": 8,
            "memPerVm": 32,
            "gpusPerVm": 0,
            "ntwPerf": "Moderate",
            "ntwPerfCategory": "medium",
            "zones": null,
            "attributes": {
                "cpu": "8",
                "instanceTypeCategory": "General purpose",
                "memory": "32",
                "networkPerfCategory": "medium"
            },
            "currentGen": true,
            "burst": true
        },
        {
            "category": "Memory optimized",
            "type": "z1d.xlarge",
            "onDemandPrice": 0.439,
            "spotPrice": [
                {
                    "zone": "eu-west-2c",
                    "price": 0.1317
                },
                {
                    "zone": "eu-west-2a",
                    "price": 0.1317
                },
                {
                    "zone": "eu-west-2b",
                    "price": 0.1317
                }
            ],
            "cpusPerVm": 4,
            "memPerVm": 32,
            "gpusPerVm": 0,
            "ntwPerf": "Up to 10 Gigabit",
            "ntwPerfCategory": "high",
            "zones": null,
            "attributes": {
                "cpu": "4",
                "instanceTypeCategory": "Memory optimized",
                "memory": "32",
                "networkPerfCategory": "high"
            },
            "currentGen": true
        },
        {
            "category": "Memory optimized",
            "type": "r5ad.xlarge",
            "onDemandPrice": 0.308,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 0.0691
                },
                {
                    "zone": "eu-west-2b",
                    "price": 0.0691
                },
                {
                    "zone": "eu-west-2c",
                    "price": 0.0691
                }
            ],
            "cpusPerVm": 4,
            "memPerVm": 32,
            "gpusPerVm": 0,
            "ntwPerf": "10 Gigabit",
            "ntwPerfCategory": "high",
            "zones": null,
            "attributes": {
                "cpu": "4",
                "instanceTypeCategory": "Memory optimized",
                "memory": "32",
                "networkPerfCategory": "high"
            },
            "currentGen": true
        },
        {
            "category": "General purpose",
            "type": "m5.large",
            "onDemandPrice": 0.111,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 0.0338
                },
                {
                    "zone": "eu-west-2b",
                    "price": 0.035
                },
                {
                    "zone": "eu-west-2c",
                    "price": 0.0347
                }
            ],
            "cpusPerVm": 2,
            "memPerVm": 8,
            "gpusPerVm": 0,
            "ntwPerf": "Up to 10 Gigabit",
            "ntwPerfCategory": "high",
            "zones": null,
            "attributes": {
                "cpu": "2",
                "instanceTypeCategory": "General purpose",
                "memory": "8",
                "networkPerfCategory": "high"
            },
            "currentGen": true
        },
        {
            "category": "General purpose",
            "type": "m4.xlarge",
            "onDemandPrice": 0.232,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 0.0655
                },
                {
                    "zone": "eu-west-2b",
                    "price": 0.0628
                },
                {
                    "zone": "eu-west-2c",
                    "price": 0.0628
                }
            ],
            "cpusPerVm": 4,
            "memPerVm": 16,
            "gpusPerVm": 0,
            "ntwPerf": "High",
            "ntwPerfCategory": "medium",
            "zones": null,
            "attributes": {
                "cpu": "4",
                "instanceTypeCategory": "General purpose",
                "memory": "16",
                "networkPerfCategory": "medium"
            },
            "currentGen": true
        },
        {
            "category": "Compute optimized",
            "type": "c5d.4xlarge",
            "onDemandPrice": 0.92,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 0.2513
                },
                {
                    "zone": "eu-west-2b",
                    "price": 0.2513
                },
                {
                    "zone": "eu-west-2c",
                    "price": 0.2513
                }
            ],
            "cpusPerVm": 16,
            "memPerVm": 32,
            "gpusPerVm": 0,
            "ntwPerf": "Up to 10 Gigabit",
            "ntwPerfCategory": "high",
            "zones": null,
            "attributes": {
                "cpu": "16",
                "instanceTypeCategory": "Compute optimized",
                "memory": "32",
                "networkPerfCategory": "high"
            },
            "currentGen": true
        },
        {
            "category": "Memory optimized",
            "type": "r5d.4xlarge",
            "onDemandPrice": 1.352,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 0.2765
                },
                {
                    "zone": "eu-west-2b",
                    "price": 0.2765
                },
                {
                    "zone": "eu-west-2c",
                    "price": 0.2765
                }
            ],
            "cpusPerVm": 16,
            "memPerVm": 128,
            "gpusPerVm": 0,
            "ntwPerf": "10 Gigabit",
            "ntwPerfCategory": "high",
            "zones": null,
            "attributes": {
                "cpu": "16",
                "instanceTypeCategory": "Memory optimized",
                "memory": "128",
                "networkPerfCategory": "high"
            },
            "currentGen": true
        },
        {
            "category": "Memory optimized",
            "type": "z1d.large",
            "onDemandPrice": 0.22,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 0.066
                },
                {
                    "zone": "eu-west-2b",
                    "price": 0.066
                },
                {
                    "zone": "eu-west-2c",
                    "price": 0.083
                }
            ],
            "cpusPerVm": 2,
            "memPerVm": 16,
            "gpusPerVm": 0,
            "ntwPerf": "Up to 10 Gigabit",
            "ntwPerfCategory": "high",
            "zones": null,
            "attributes": {
                "cpu": "2",
                "instanceTypeCategory": "Memory optimized",
                "memory": "16",
                "networkPerfCategory": "high"
            },
            "currentGen": true
        },
        {
            "category": "General purpose",
            "type": "m5a.16xlarge",
            "onDemandPrice": 3.2,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 1.0557
                },
                {
                    "zone": "eu-west-2b",
                    "price": 1.0557
                },
                {
                    "zone": "eu-west-2c",
                    "price": 3.2
                }
            ],
            "cpusPerVm": 64,
            "memPerVm": 256,
            "gpusPerVm": 0,
            "ntwPerf": "12 Gigabit",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": {
                "cpu": "64",
                "instanceTypeCategory": "General purpose",
                "memory": "256",
                "networkPerfCategory": ""
            },
            "currentGen": true
        },
        {
            "category": "General purpose",
            "type": "t2.micro",
            "onDemandPrice": 0.0132,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 0.004
                },
                {
                    "zone": "eu-west-2b",
                    "price": 0.004
                },
                {
                    "zone": "eu-west-2c",
                    "price": 0.004
                }
            ],
            "cpusPerVm": 1,
            "memPerVm": 1,
            "gpusPerVm": 0,
            "ntwPerf": "Low to Moderate",
            "ntwPerfCategory": "low",
            "zones": null,
            "attributes": {
                "cpu": "1",
                "instanceTypeCategory": "General purpose",
                "memory": "1",
                "networkPerfCategory": "low"
            },
            "currentGen": true,
            "burst": true
        },
        {
            "category": "General purpose",
            "type": "m4.16xlarge",
            "onDemandPrice": 3.712,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 1.0054
                },
                {
                    "zone": "eu-west-2b",
                    "price": 1.0054
                },
                {
                    "zone": "eu-west-2c",
                    "price": 1.0269
                }
            ],
            "cpusPerVm": 64,
            "memPerVm": 256,
            "gpusPerVm": 0,
            "ntwPerf": "20 Gigabit",
            "ntwPerfCategory": "extra",
            "zones": null,
            "attributes": {
                "cpu": "64",
                "instanceTypeCategory": "General purpose",
                "memory": "256",
                "networkPerfCategory": "extra"
            },
            "currentGen": true
        },
        {
            "category": "General purpose",
            "type": "t3a.xlarge",
            "onDemandPrice": 0.1699,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 0.051
                },
                {
                    "zone": "eu-west-2b",
                    "price": 0.051
                },
                {
                    "zone": "eu-west-2c",
                    "price": 0.051
                }
            ],
            "cpusPerVm": 4,
            "memPerVm": 16,
            "gpusPerVm": 0,
            "ntwPerf": "Moderate",
            "ntwPerfCategory": "medium",
            "zones": null,
            "attributes": {
                "cpu": "4",
                "instanceTypeCategory": "General purpose",
                "memory": "16",
                "networkPerfCategory": "medium"
            },
            "currentGen": true,
            "burst": true
        },
        {
            "category": "Memory optimized",
            "type": "z1d.12xlarge",
            "onDemandPrice": 5.273,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 1.5819
                },
                {
                    "zone": "eu-west-2b",
                    "price": 1.5819
                },
                {
                    "zone": "eu-west-2c",
                    "price": 1.5819
                }
            ],
            "cpusPerVm": 48,
            "memPerVm": 384,
            "gpusPerVm": 0,
            "ntwPerf": "25 Gigabit",
            "ntwPerfCategory": "extra",
            "zones": null,
            "attributes": {
                "cpu": "48",
                "instanceTypeCategory": "Memory optimized",
                "memory": "384",
                "networkPerfCategory": "extra"
            },
            "currentGen": true
        },
        {
            "category": "General purpose",
            "type": "m5a.4xlarge",
            "onDemandPrice": 0.8,
            "spotPrice": [
                {
                    "zone": "eu-west-2b",
                    "price": 0.2639
                },
                {
                    "zone": "eu-west-2c",
                    "price": 0.8
                },
                {
                    "zone": "eu-west-2a",
                    "price": 0.2639
                }
            ],
            "cpusPerVm": 16,
            "memPerVm": 64,
            "gpusPerVm": 0,
            "ntwPerf": "Up to 10 Gigabit",
            "ntwPerfCategory": "high",
            "zones": null,
            "attributes": {
                "cpu": "16",
                "instanceTypeCategory": "General purpose",
                "memory": "64",
                "networkPerfCategory": "high"
            },
            "currentGen": true
        },
        {
            "category": "General purpose",
            "type": "t2.2xlarge",
            "onDemandPrice": 0.4224,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 0.1267
                },
                {
                    "zone": "eu-west-2b",
                    "price": 0.1267
                },
                {
                    "zone": "eu-west-2c",
                    "price": 0.1267
                }
            ],
            "cpusPerVm": 8,
            "memPerVm": 32,
            "gpusPerVm": 0,
            "ntwPerf": "Moderate",
            "ntwPerfCategory": "medium",
            "zones": null,
            "attributes": {
                "cpu": "8",
                "instanceTypeCategory": "General purpose",
                "memory": "32",
                "networkPerfCategory": "medium"
            },
            "currentGen": true,
            "burst": true
        },
        {
            "category": "GPU instance",
            "type": "g4dn.16xlarge",
            "onDemandPrice": 5.092,
            "spotPrice": [
                {
                    "zone": "eu-west-2c",
                    "price": 1.5276
                },
                {
                    "zone": "eu-west-2a",
                    "price": 1.5276
                },
                {
                    "zone": "eu-west-2b",
                    "price": 1.5276
                }
            ],
            "cpusPerVm": 64,
            "memPerVm": 256,
            "gpusPerVm": 1,
            "ntwPerf": "50 Gigabit",
            "ntwPerfCategory": "extra",
            "zones": null,
            "attributes": {
                "cpu": "64",
                "instanceTypeCategory": "GPU instance",
                "memory": "256",
                "networkPerfCategory": "extra"
            },
            "currentGen": true
        },
        {
            "category": "Compute optimized",
            "type": "c5d.18xlarge",
            "onDemandPrice": 4.14,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 1.1311
                },
                {
                    "zone": "eu-west-2b",
                    "price": 1.1311
                },
                {
                    "zone": "eu-west-2c",
                    "price": 1.1311
                }
            ],
            "cpusPerVm": 72,
            "memPerVm": 144,
            "gpusPerVm": 0,
            "ntwPerf": "25 Gigabit",
            "ntwPerfCategory": "extra",
            "zones": null,
            "attributes": {
                "cpu": "72",
                "instanceTypeCategory": "Compute optimized",
                "memory": "144",
                "networkPerfCategory": "extra"
            },
            "currentGen": true
        },
        {
            "category": "Memory optimized",
            "type": "r5ad.4xlarge",
            "onDemandPrice": 1.232,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 1.232
                },
                {
                    "zone": "eu-west-2b",
                    "price": 0.2765
                },
                {
                    "zone": "eu-west-2c",
                    "price": 0.2765
                }
            ],
            "cpusPerVm": 16,
            "memPerVm": 128,
            "gpusPerVm": 0,
            "ntwPerf": "10 Gigabit",
            "ntwPerfCategory": "high",
            "zones": null,
            "attributes": {
                "cpu": "16",
                "instanceTypeCategory": "Memory optimized",
                "memory": "128",
                "networkPerfCategory": "high"
            },
            "currentGen": true
        },
        {
            "category": "General purpose",
            "type": "t3.medium",
            "onDemandPrice": 0.0472,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 0.0142
                },
                {
                    "zone": "eu-west-2b",
                    "price": 0.0142
                },
                {
                    "zone": "eu-west-2c",
                    "price": 0.0142
                }
            ],
            "cpusPerVm": 2,
            "memPerVm": 4,
            "gpusPerVm": 0,
            "ntwPerf": "Low to Moderate",
            "ntwPerfCategory": "low",
            "zones": null,
            "attributes": {
                "cpu": "2",
                "instanceTypeCategory": "General purpose",
                "memory": "4",
                "networkPerfCategory": "low"
            },
            "currentGen": true,
            "burst": true
        },
        {
            "category": "General purpose",
            "type": "t2.medium",
            "onDemandPrice": 0.052,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 0.0156
                },
                {
                    "zone": "eu-west-2b",
                    "price": 0.0156
                },
                {
                    "zone": "eu-west-2c",
                    "price": 0.0156
                }
            ],
            "cpusPerVm": 2,
            "memPerVm": 4,
            "gpusPerVm": 0,
            "ntwPerf": "Low to Moderate",
            "ntwPerfCategory": "low",
            "zones": null,
            "attributes": {
                "cpu": "2",
                "instanceTypeCategory": "General purpose",
                "memory": "4",
                "networkPerfCategory": "low"
            },
            "currentGen": true,
            "burst": true
        },
        {
            "category": "Compute optimized",
            "type": "c5.metal",
            "onDemandPrice": 4.848,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 1.5081
                },
                {
                    "zone": "eu-west-2b",
                    "price": 1.5081
                },
                {
                    "zone": "eu-west-2c",
                    "price": 1.5081
                }
            ],
            "cpusPerVm": 96,
            "memPerVm": 192,
            "gpusPerVm": 0,
            "ntwPerf": "25 Gigabit",
            "ntwPerfCategory": "extra",
            "zones": null,
            "attributes": {
                "cpu": "96",
                "instanceTypeCategory": "Compute optimized",
                "memory": "192",
                "networkPerfCategory": "extra"
            },
            "currentGen": true
        },
        {
            "category": "Storage optimized",
            "type": "d2.8xlarge",
            "onDemandPrice": 6.174,
            "spotPrice": [
                {
                    "zone": "eu-west-2b",
                    "price": 1.8522
                },
                {
                    "zone": "eu-west-2c",
                    "price": 6.174
                },
                {
                    "zone": "eu-west-2a",
                    "price": 1.8522
                }
            ],
            "cpusPerVm": 36,
            "memPerVm": 244,
            "gpusPerVm": 0,
            "ntwPerf": "10 Gigabit",
            "ntwPerfCategory": "high",
            "zones": null,
            "attributes": {
                "cpu": "36",
                "instanceTypeCategory": "Storage optimized",
                "memory": "244",
                "networkPerfCategory": "high"
            },
            "currentGen": true
        },
        {
            "category": "General purpose",
            "type": "t3.xlarge",
            "onDemandPrice": 0.1888,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 0.0566
                },
                {
                    "zone": "eu-west-2b",
                    "price": 0.0569
                },
                {
                    "zone": "eu-west-2c",
                    "price": 0.0566
                }
            ],
            "cpusPerVm": 4,
            "memPerVm": 16,
            "gpusPerVm": 0,
            "ntwPerf": "Moderate",
            "ntwPerfCategory": "medium",
            "zones": null,
            "attributes": {
                "cpu": "4",
                "instanceTypeCategory": "General purpose",
                "memory": "16",
                "networkPerfCategory": "medium"
            },
            "currentGen": true,
            "burst": true
        },
        {
            "category": "Memory optimized",
            "type": "r5a.24xlarge",
            "onDemandPrice": 6.384,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 1.6589
                },
                {
                    "zone": "eu-west-2b",
                    "price": 1.6589
                },
                {
                    "zone": "eu-west-2c",
                    "price": 1.6589
                }
            ],
            "cpusPerVm": 96,
            "memPerVm": 768,
            "gpusPerVm": 0,
            "ntwPerf": "20 Gigabit",
            "ntwPerfCategory": "extra",
            "zones": null,
            "attributes": {
                "cpu": "96",
                "instanceTypeCategory": "Memory optimized",
                "memory": "768",
                "networkPerfCategory": "extra"
            },
            "currentGen": true
        },
        {
            "category": "General purpose",
            "type": "m5.16xlarge",
            "onDemandPrice": 3.552,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 1.0557
                },
                {
                    "zone": "eu-west-2b",
                    "price": 1.0583
                },
                {
                    "zone": "eu-west-2c",
                    "price": 1.0557
                }
            ],
            "cpusPerVm": 64,
            "memPerVm": 256,
            "gpusPerVm": 0,
            "ntwPerf": "20 Gigabit",
            "ntwPerfCategory": "extra",
            "zones": null,
            "attributes": {
                "cpu": "64",
                "instanceTypeCategory": "General purpose",
                "memory": "256",
                "networkPerfCategory": "extra"
            },
            "currentGen": true
        },
        {
            "category": "GPU instance",
            "type": "p3.8xlarge",
            "onDemandPrice": 14.356,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 4.3068
                },
                {
                    "zone": "eu-west-2b",
                    "price": 4.3068
                }
            ],
            "cpusPerVm": 32,
            "memPerVm": 244,
            "gpusPerVm": 4,
            "ntwPerf": "10 Gigabit",
            "ntwPerfCategory": "high",
            "zones": null,
            "attributes": {
                "cpu": "32",
                "instanceTypeCategory": "GPU instance",
                "memory": "244",
                "networkPerfCategory": "high"
            },
            "currentGen": true
        },
        {
            "category": "Compute optimized",
            "type": "c5d.9xlarge",
            "onDemandPrice": 2.07,
            "spotPrice": [
                {
                    "zone": "eu-west-2b",
                    "price": 0.5655
                },
                {
                    "zone": "eu-west-2c",
                    "price": 0.567
                },
                {
                    "zone": "eu-west-2a",
                    "price": 0.5655
                }
            ],
            "cpusPerVm": 36,
            "memPerVm": 72,
            "gpusPerVm": 0,
            "ntwPerf": "10 Gigabit",
            "ntwPerfCategory": "high",
            "zones": null,
            "attributes": {
                "cpu": "36",
                "instanceTypeCategory": "Compute optimized",
                "memory": "72",
                "networkPerfCategory": "high"
            },
            "currentGen": true
        },
        {
            "category": "General purpose",
            "type": "m5d.24xlarge",
            "onDemandPrice": 6.288,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 1.5835
                },
                {
                    "zone": "eu-west-2b",
                    "price": 1.5835
                },
                {
                    "zone": "eu-west-2c",
                    "price": 1.5835
                }
            ],
            "cpusPerVm": 96,
            "memPerVm": 384,
            "gpusPerVm": 0,
            "ntwPerf": "25 Gigabit",
            "ntwPerfCategory": "extra",
            "zones": null,
            "attributes": {
                "cpu": "96",
                "instanceTypeCategory": "General purpose",
                "memory": "384",
                "networkPerfCategory": "extra"
            },
            "currentGen": true
        },
        {
            "category": "General purpose",
            "type": "m5ad.4xlarge",
            "onDemandPrice": 0.96,
            "spotPrice": [
                {
                    "zone": "eu-west-2b",
                    "price": 0.2639
                },
                {
                    "zone": "eu-west-2c",
                    "price": 0.2639
                },
                {
                    "zone": "eu-west-2a",
                    "price": 0.2639
                }
            ],
            "cpusPerVm": 16,
            "memPerVm": 64,
            "gpusPerVm": 0,
            "ntwPerf": "Up to 10 Gigabit",
            "ntwPerfCategory": "high",
            "zones": null,
            "attributes": {
                "cpu": "16",
                "instanceTypeCategory": "General purpose",
                "memory": "64",
                "networkPerfCategory": "high"
            },
            "currentGen": true
        },
        {
            "category": "General purpose",
            "type": "t2.xlarge",
            "onDemandPrice": 0.2112,
            "spotPrice": [
                {
                    "zone": "eu-west-2b",
                    "price": 0.0634
                },
                {
                    "zone": "eu-west-2c",
                    "price": 0.0634
                },
                {
                    "zone": "eu-west-2a",
                    "price": 0.0634
                }
            ],
            "cpusPerVm": 4,
            "memPerVm": 16,
            "gpusPerVm": 0,
            "ntwPerf": "Moderate",
            "ntwPerfCategory": "medium",
            "zones": null,
            "attributes": {
                "cpu": "4",
                "instanceTypeCategory": "General purpose",
                "memory": "16",
                "networkPerfCategory": "medium"
            },
            "currentGen": true,
            "burst": true
        },
        {
            "category": "Memory optimized",
            "type": "r5a.large",
            "onDemandPrice": 0.133,
            "spotPrice": [
                {
                    "zone": "eu-west-2c",
                    "price": 0.0346
                },
                {
                    "zone": "eu-west-2a",
                    "price": 0.0346
                },
                {
                    "zone": "eu-west-2b",
                    "price": 0.0346
                }
            ],
            "cpusPerVm": 2,
            "memPerVm": 16,
            "gpusPerVm": 0,
            "ntwPerf": "10 Gigabit",
            "ntwPerfCategory": "high",
            "zones": null,
            "attributes": {
                "cpu": "2",
                "instanceTypeCategory": "Memory optimized",
                "memory": "16",
                "networkPerfCategory": "high"
            },
            "currentGen": true
        },
        {
            "category": "Compute optimized",
            "type": "c5.24xlarge",
            "onDemandPrice": 4.848,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 1.5081
                },
                {
                    "zone": "eu-west-2b",
                    "price": 1.5081
                },
                {
                    "zone": "eu-west-2c",
                    "price": 1.5081
                }
            ],
            "cpusPerVm": 96,
            "memPerVm": 192,
            "gpusPerVm": 0,
            "ntwPerf": "25 Gigabit",
            "ntwPerfCategory": "extra",
            "zones": null,
            "attributes": {
                "cpu": "96",
                "instanceTypeCategory": "Compute optimized",
                "memory": "192",
                "networkPerfCategory": "extra"
            },
            "currentGen": true
        },
        {
            "category": "Memory optimized",
            "type": "r5d.12xlarge",
            "onDemandPrice": 4.056,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 0.8295
                },
                {
                    "zone": "eu-west-2b",
                    "price": 0.8295
                },
                {
                    "zone": "eu-west-2c",
                    "price": 0.8295
                }
            ],
            "cpusPerVm": 48,
            "memPerVm": 384,
            "gpusPerVm": 0,
            "ntwPerf": "10 Gigabit",
            "ntwPerfCategory": "high",
            "zones": null,
            "attributes": {
                "cpu": "48",
                "instanceTypeCategory": "Memory optimized",
                "memory": "384",
                "networkPerfCategory": "high"
            },
            "currentGen": true
        },
        {
            "category": "General purpose",
            "type": "m5ad.8xlarge",
            "onDemandPrice": 1.92,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 0.5278
                },
                {
                    "zone": "eu-west-2b",
                    "price": 0.5278
                },
                {
                    "zone": "eu-west-2c",
                    "price": 0.5278
                }
            ],
            "cpusPerVm": 32,
            "memPerVm": 128,
            "gpusPerVm": 0,
            "ntwPerf": "Up to 10 Gigabit",
            "ntwPerfCategory": "high",
            "zones": null,
            "attributes": {
                "cpu": "32",
                "instanceTypeCategory": "General purpose",
                "memory": "128",
                "networkPerfCategory": "high"
            },
            "currentGen": true
        },
        {
            "category": "General purpose",
            "type": "m5d.xlarge",
            "onDemandPrice": 0.262,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 0.0662
                },
                {
                    "zone": "eu-west-2b",
                    "price": 0.066
                },
                {
                    "zone": "eu-west-2c",
                    "price": 0.066
                }
            ],
            "cpusPerVm": 4,
            "memPerVm": 16,
            "gpusPerVm": 0,
            "ntwPerf": "Up to 10 Gigabit",
            "ntwPerfCategory": "high",
            "zones": null,
            "attributes": {
                "cpu": "4",
                "instanceTypeCategory": "General purpose",
                "memory": "16",
                "networkPerfCategory": "high"
            },
            "currentGen": true
        },
        {
            "category": "Memory optimized",
            "type": "r5.12xlarge",
            "onDemandPrice": 3.552,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 0.835
                },
                {
                    "zone": "eu-west-2b",
                    "price": 0.9718
                },
                {
                    "zone": "eu-west-2c",
                    "price": 0.971
                }
            ],
            "cpusPerVm": 48,
            "memPerVm": 384,
            "gpusPerVm": 0,
            "ntwPerf": "10 Gigabit",
            "ntwPerfCategory": "high",
            "zones": null,
            "attributes": {
                "cpu": "48",
                "instanceTypeCategory": "Memory optimized",
                "memory": "384",
                "networkPerfCategory": "high"
            },
            "currentGen": true
        },
        {
            "category": "Storage optimized",
            "type": "d2.4xlarge",
            "onDemandPrice": 3.087,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 3.087
                },
                {
                    "zone": "eu-west-2b",
                    "price": 3.087
                },
                {
                    "zone": "eu-west-2c",
                    "price": 3.087
                }
            ],
            "cpusPerVm": 16,
            "memPerVm": 122,
            "gpusPerVm": 0,
            "ntwPerf": "High",
            "ntwPerfCategory": "medium",
            "zones": null,
            "attributes": {
                "cpu": "16",
                "instanceTypeCategory": "Storage optimized",
                "memory": "122",
                "networkPerfCategory": "medium"
            },
            "currentGen": true
        },
        {
            "category": "GPU instance",
            "type": "g3.8xlarge",
            "onDemandPrice": 2.858,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 0.8574
                },
                {
                    "zone": "eu-west-2b",
                    "price": 0.8574
                }
            ],
            "cpusPerVm": 32,
            "memPerVm": 244,
            "gpusPerVm": 2,
            "ntwPerf": "10 Gigabit",
            "ntwPerfCategory": "high",
            "zones": null,
            "attributes": {
                "cpu": "32",
                "instanceTypeCategory": "GPU instance",
                "memory": "244",
                "networkPerfCategory": "high"
            },
            "currentGen": true
        },
        {
            "category": "Compute optimized",
            "type": "c4.8xlarge",
            "onDemandPrice": 1.902,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 0.4788
                },
                {
                    "zone": "eu-west-2b",
                    "price": 0.4788
                },
                {
                    "zone": "eu-west-2c",
                    "price": 1.902
                }
            ],
            "cpusPerVm": 36,
            "memPerVm": 60,
            "gpusPerVm": 0,
            "ntwPerf": "10 Gigabit",
            "ntwPerfCategory": "high",
            "zones": null,
            "attributes": {
                "cpu": "36",
                "instanceTypeCategory": "Compute optimized",
                "memory": "60",
                "networkPerfCategory": "high"
            },
            "currentGen": true
        },
        {
            "category": "General purpose",
            "type": "m4.2xlarge",
            "onDemandPrice": 0.464,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 0.1276
                },
                {
                    "zone": "eu-west-2b",
                    "price": 0.1319
                },
                {
                    "zone": "eu-west-2c",
                    "price": 0.1327
                }
            ],
            "cpusPerVm": 8,
            "memPerVm": 32,
            "gpusPerVm": 0,
            "ntwPerf": "High",
            "ntwPerfCategory": "medium",
            "zones": null,
            "attributes": {
                "cpu": "8",
                "instanceTypeCategory": "General purpose",
                "memory": "32",
                "networkPerfCategory": "medium"
            },
            "currentGen": true
        },
        {
            "category": "Storage optimized",
            "type": "i3en.6xlarge",
            "onDemandPrice": 3.156,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 0.9468
                },
                {
                    "zone": "eu-west-2b",
                    "price": 0.9468
                },
                {
                    "zone": "eu-west-2c",
                    "price": 0.9468
                }
            ],
            "cpusPerVm": 24,
            "memPerVm": 192,
            "gpusPerVm": 0,
            "ntwPerf": "25 Gigabit",
            "ntwPerfCategory": "extra",
            "zones": null,
            "attributes": {
                "cpu": "24",
                "instanceTypeCategory": "Storage optimized",
                "memory": "192",
                "networkPerfCategory": "extra"
            },
            "currentGen": true
        },
        {
            "category": "Memory optimized",
            "type": "r5.2xlarge",
            "onDemandPrice": 0.592,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 0.1473
                },
                {
                    "zone": "eu-west-2b",
                    "price": 0.1382
                },
                {
                    "zone": "eu-west-2c",
                    "price": 0.1465
                }
            ],
            "cpusPerVm": 8,
            "memPerVm": 64,
            "gpusPerVm": 0,
            "ntwPerf": "10 Gigabit",
            "ntwPerfCategory": "high",
            "zones": null,
            "attributes": {
                "cpu": "8",
                "instanceTypeCategory": "Memory optimized",
                "memory": "64",
                "networkPerfCategory": "high"
            },
            "currentGen": true
        },
        {
            "category": "Compute optimized",
            "type": "c5d.xlarge",
            "onDemandPrice": 0.23,
            "spotPrice": [
                {
                    "zone": "eu-west-2b",
                    "price": 0.0628
                },
                {
                    "zone": "eu-west-2c",
                    "price": 0.0628
                },
                {
                    "zone": "eu-west-2a",
                    "price": 0.0628
                }
            ],
            "cpusPerVm": 4,
            "memPerVm": 8,
            "gpusPerVm": 0,
            "ntwPerf": "Up to 10 Gigabit",
            "ntwPerfCategory": "high",
            "zones": null,
            "attributes": {
                "cpu": "4",
                "instanceTypeCategory": "Compute optimized",
                "memory": "8",
                "networkPerfCategory": "high"
            },
            "currentGen": true
        },
        {
            "category": "Memory optimized",
            "type": "z1d.6xlarge",
            "onDemandPrice": 2.636,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 0.7908
                },
                {
                    "zone": "eu-west-2b",
                    "price": 0.7908
                },
                {
                    "zone": "eu-west-2c",
                    "price": 0.7908
                }
            ],
            "cpusPerVm": 24,
            "memPerVm": 192,
            "gpusPerVm": 0,
            "ntwPerf": "10 Gigabit",
            "ntwPerfCategory": "high",
            "zones": null,
            "attributes": {
                "cpu": "24",
                "instanceTypeCategory": "Memory optimized",
                "memory": "192",
                "networkPerfCategory": "high"
            },
            "currentGen": true
        },
        {
            "category": "GPU instance",
            "type": "g4dn.12xlarge",
            "onDemandPrice": 4.577,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 1.3731
                },
                {
                    "zone": "eu-west-2b",
                    "price": 1.3731
                },
                {
                    "zone": "eu-west-2c",
                    "price": 4.577
                }
            ],
            "cpusPerVm": 48,
            "memPerVm": 192,
            "gpusPerVm": 4,
            "ntwPerf": "50 Gigabit",
            "ntwPerfCategory": "extra",
            "zones": null,
            "attributes": {
                "cpu": "48",
                "instanceTypeCategory": "GPU instance",
                "memory": "192",
                "networkPerfCategory": "extra"
            },
            "currentGen": true
        },
        {
            "category": "Memory optimized",
            "type": "r5d.16xlarge",
            "onDemandPrice": 5.408,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 4.4071
                },
                {
                    "zone": "eu-west-2b",
                    "price": 1.1059
                },
                {
                    "zone": "eu-west-2c",
                    "price": 1.1059
                }
            ],
            "cpusPerVm": 64,
            "memPerVm": 512,
            "gpusPerVm": 0,
            "ntwPerf": "20 Gigabit",
            "ntwPerfCategory": "extra",
            "zones": null,
            "attributes": {
                "cpu": "64",
                "instanceTypeCategory": "Memory optimized",
                "memory": "512",
                "networkPerfCategory": "extra"
            },
            "currentGen": true
        },
        {
            "category": "Memory optimized",
            "type": "r5ad.24xlarge",
            "onDemandPrice": 7.392,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 1.6589
                },
                {
                    "zone": "eu-west-2b",
                    "price": 1.6589
                },
                {
                    "zone": "eu-west-2c",
                    "price": 1.6589
                }
            ],
            "cpusPerVm": 96,
            "memPerVm": 768,
            "gpusPerVm": 0,
            "ntwPerf": "20 Gigabit",
            "ntwPerfCategory": "extra",
            "zones": null,
            "attributes": {
                "cpu": "96",
                "instanceTypeCategory": "Memory optimized",
                "memory": "768",
                "networkPerfCategory": "extra"
            },
            "currentGen": true
        },
        {
            "category": "GPU instance",
            "type": "g4dn.2xlarge",
            "onDemandPrice": 0.88,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 0.264
                },
                {
                    "zone": "eu-west-2b",
                    "price": 0.264
                },
                {
                    "zone": "eu-west-2c",
                    "price": 0.264
                }
            ],
            "cpusPerVm": 8,
            "memPerVm": 32,
            "gpusPerVm": 1,
            "ntwPerf": "Up to 25 Gigabit",
            "ntwPerfCategory": "extra",
            "zones": null,
            "attributes": {
                "cpu": "8",
                "instanceTypeCategory": "GPU instance",
                "memory": "32",
                "networkPerfCategory": "extra"
            },
            "currentGen": true
        },
        {
            "category": "Storage optimized",
            "type": "d2.2xlarge",
            "onDemandPrice": 1.544,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 1.544
                },
                {
                    "zone": "eu-west-2b",
                    "price": 1.544
                },
                {
                    "zone": "eu-west-2c",
                    "price": 1.544
                }
            ],
            "cpusPerVm": 8,
            "memPerVm": 61,
            "gpusPerVm": 0,
            "ntwPerf": "High",
            "ntwPerfCategory": "medium",
            "zones": null,
            "attributes": {
                "cpu": "8",
                "instanceTypeCategory": "Storage optimized",
                "memory": "61",
                "networkPerfCategory": "medium"
            },
            "currentGen": true
        },
        {
            "category": "General purpose",
            "type": "m5ad.xlarge",
            "onDemandPrice": 0.24,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 0.0667
                },
                {
                    "zone": "eu-west-2b",
                    "price": 0.066
                },
                {
                    "zone": "eu-west-2c",
                    "price": 0.066
                }
            ],
            "cpusPerVm": 4,
            "memPerVm": 16,
            "gpusPerVm": 0,
            "ntwPerf": "Up to 10 Gigabit",
            "ntwPerfCategory": "high",
            "zones": null,
            "attributes": {
                "cpu": "4",
                "instanceTypeCategory": "General purpose",
                "memory": "16",
                "networkPerfCategory": "high"
            },
            "currentGen": true
        },
        {
            "category": "General purpose",
            "type": "m5d.large",
            "onDemandPrice": 0.131,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 0.033
                },
                {
                    "zone": "eu-west-2b",
                    "price": 0.033
                },
                {
                    "zone": "eu-west-2c",
                    "price": 0.0332
                }
            ],
            "cpusPerVm": 2,
            "memPerVm": 8,
            "gpusPerVm": 0,
            "ntwPerf": "Up to 10 Gigabit",
            "ntwPerfCategory": "high",
            "zones": null,
            "attributes": {
                "cpu": "2",
                "instanceTypeCategory": "General purpose",
                "memory": "8",
                "networkPerfCategory": "high"
            },
            "currentGen": true
        },
        {
            "category": "General purpose",
            "type": "m5.xlarge",
            "onDemandPrice": 0.222,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 0.067
                },
                {
                    "zone": "eu-west-2b",
                    "price": 0.066
                },
                {
                    "zone": "eu-west-2c",
                    "price": 0.066
                }
            ],
            "cpusPerVm": 4,
            "memPerVm": 16,
            "gpusPerVm": 0,
            "ntwPerf": "Up to 10 Gigabit",
            "ntwPerfCategory": "high",
            "zones": null,
            "attributes": {
                "cpu": "4",
                "instanceTypeCategory": "General purpose",
                "memory": "16",
                "networkPerfCategory": "high"
            },
            "currentGen": true
        },
        {
            "category": "Compute optimized",
            "type": "c5.4xlarge",
            "onDemandPrice": 0.808,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 0.2531
                },
                {
                    "zone": "eu-west-2b",
                    "price": 0.2582
                },
                {
                    "zone": "eu-west-2c",
                    "price": 0.2679
                }
            ],
            "cpusPerVm": 16,
            "memPerVm": 32,
            "gpusPerVm": 0,
            "ntwPerf": "Up to 10 Gigabit",
            "ntwPerfCategory": "high",
            "zones": null,
            "attributes": {
                "cpu": "16",
                "instanceTypeCategory": "Compute optimized",
                "memory": "32",
                "networkPerfCategory": "high"
            },
            "currentGen": true
        },
        {
            "category": "Memory optimized",
            "type": "r5a.16xlarge",
            "onDemandPrice": 4.256,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 1.1059
                },
                {
                    "zone": "eu-west-2b",
                    "price": 1.1059
                },
                {
                    "zone": "eu-west-2c",
                    "price": 1.1059
                }
            ],
            "cpusPerVm": 64,
            "memPerVm": 512,
            "gpusPerVm": 0,
            "ntwPerf": "12 Gigabit",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": {
                "cpu": "64",
                "instanceTypeCategory": "Memory optimized",
                "memory": "512",
                "networkPerfCategory": ""
            },
            "currentGen": true
        },
        {
            "category": "General purpose",
            "type": "m4.10xlarge",
            "onDemandPrice": 2.32,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 0.6284
                },
                {
                    "zone": "eu-west-2b",
                    "price": 0.6284
                },
                {
                    "zone": "eu-west-2c",
                    "price": 0.6318
                }
            ],
            "cpusPerVm": 40,
            "memPerVm": 160,
            "gpusPerVm": 0,
            "ntwPerf": "10 Gigabit",
            "ntwPerfCategory": "high",
            "zones": null,
            "attributes": {
                "cpu": "40",
                "instanceTypeCategory": "General purpose",
                "memory": "160",
                "networkPerfCategory": "high"
            },
            "currentGen": true
        },
        {
            "category": "Memory optimized",
            "type": "r4.xlarge",
            "onDemandPrice": 0.312,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 0.0666
                },
                {
                    "zone": "eu-west-2b",
                    "price": 0.0658
                },
                {
                    "zone": "eu-west-2c",
                    "price": 0.0658
                }
            ],
            "cpusPerVm": 4,
            "memPerVm": 30.5,
            "gpusPerVm": 0,
            "ntwPerf": "Up to 10 Gigabit",
            "ntwPerfCategory": "high",
            "zones": null,
            "attributes": {
                "cpu": "4",
                "instanceTypeCategory": "Memory optimized",
                "memory": "30.5",
                "networkPerfCategory": "high"
            },
            "currentGen": true
        },
        {
            "category": "General purpose",
            "type": "m5.24xlarge",
            "onDemandPrice": 5.328,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 1.5835
                },
                {
                    "zone": "eu-west-2b",
                    "price": 1.5897
                },
                {
                    "zone": "eu-west-2c",
                    "price": 1.5835
                }
            ],
            "cpusPerVm": 96,
            "memPerVm": 384,
            "gpusPerVm": 0,
            "ntwPerf": "25 Gigabit",
            "ntwPerfCategory": "extra",
            "zones": null,
            "attributes": {
                "cpu": "96",
                "instanceTypeCategory": "General purpose",
                "memory": "384",
                "networkPerfCategory": "extra"
            },
            "currentGen": true
        },
        {
            "category": "Storage optimized",
            "type": "i3en.12xlarge",
            "onDemandPrice": 6.312,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 1.8936
                },
                {
                    "zone": "eu-west-2b",
                    "price": 1.8936
                },
                {
                    "zone": "eu-west-2c",
                    "price": 1.8936
                }
            ],
            "cpusPerVm": 48,
            "memPerVm": 384,
            "gpusPerVm": 0,
            "ntwPerf": "50 Gigabit",
            "ntwPerfCategory": "extra",
            "zones": null,
            "attributes": {
                "cpu": "48",
                "instanceTypeCategory": "Storage optimized",
                "memory": "384",
                "networkPerfCategory": "extra"
            },
            "currentGen": true
        },
        {
            "category": "Memory optimized",
            "type": "r5.large",
            "onDemandPrice": 0.148,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 0.0401
                },
                {
                    "zone": "eu-west-2b",
                    "price": 0.0449
                },
                {
                    "zone": "eu-west-2c",
                    "price": 0.0442
                }
            ],
            "cpusPerVm": 2,
            "memPerVm": 16,
            "gpusPerVm": 0,
            "ntwPerf": "10 Gigabit",
            "ntwPerfCategory": "high",
            "zones": null,
            "attributes": {
                "cpu": "2",
                "instanceTypeCategory": "Memory optimized",
                "memory": "16",
                "networkPerfCategory": "high"
            },
            "currentGen": true
        },
        {
            "category": "Compute optimized",
            "type": "c5.xlarge",
            "onDemandPrice": 0.202,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 0.0637
                },
                {
                    "zone": "eu-west-2b",
                    "price": 0.0643
                },
                {
                    "zone": "eu-west-2c",
                    "price": 0.0714
                }
            ],
            "cpusPerVm": 4,
            "memPerVm": 8,
            "gpusPerVm": 0,
            "ntwPerf": "Up to 10 Gigabit",
            "ntwPerfCategory": "high",
            "zones": null,
            "attributes": {
                "cpu": "4",
                "instanceTypeCategory": "Compute optimized",
                "memory": "8",
                "networkPerfCategory": "high"
            },
            "currentGen": true
        },
        {
            "category": "Memory optimized",
            "type": "x1.16xlarge",
            "onDemandPrice": 8.403,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 2.5209
                },
                {
                    "zone": "eu-west-2b",
                    "price": 2.5209
                },
                {
                    "zone": "eu-west-2c",
                    "price": 2.5209
                }
            ],
            "cpusPerVm": 64,
            "memPerVm": 976,
            "gpusPerVm": 0,
            "ntwPerf": "High",
            "ntwPerfCategory": "medium",
            "zones": null,
            "attributes": {
                "cpu": "64",
                "instanceTypeCategory": "Memory optimized",
                "memory": "976",
                "networkPerfCategory": "medium"
            },
            "currentGen": true
        },
        {
            "category": "General purpose",
            "type": "t3a.micro",
            "onDemandPrice": 0.0106,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 0.0032
                },
                {
                    "zone": "eu-west-2b",
                    "price": 0.0032
                },
                {
                    "zone": "eu-west-2c",
                    "price": 0.0032
                }
            ],
            "cpusPerVm": 2,
            "memPerVm": 1,
            "gpusPerVm": 0,
            "ntwPerf": "Low to Moderate",
            "ntwPerfCategory": "low",
            "zones": null,
            "attributes": {
                "cpu": "2",
                "instanceTypeCategory": "General purpose",
                "memory": "1",
                "networkPerfCategory": "low"
            },
            "currentGen": true,
            "burst": true
        },
        {
            "category": "General purpose",
            "type": "t3a.small",
            "onDemandPrice": 0.0212,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 0.0064
                },
                {
                    "zone": "eu-west-2b",
                    "price": 0.0064
                },
                {
                    "zone": "eu-west-2c",
                    "price": 0.0064
                }
            ],
            "cpusPerVm": 2,
            "memPerVm": 2,
            "gpusPerVm": 0,
            "ntwPerf": "Low to Moderate",
            "ntwPerfCategory": "low",
            "zones": null,
            "attributes": {
                "cpu": "2",
                "instanceTypeCategory": "General purpose",
                "memory": "2",
                "networkPerfCategory": "low"
            },
            "currentGen": true,
            "burst": true
        },
        {
            "category": "Storage optimized",
            "type": "i3.8xlarge",
            "onDemandPrice": 2.896,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 0.8688
                },
                {
                    "zone": "eu-west-2b",
                    "price": 2.896
                },
                {
                    "zone": "eu-west-2c",
                    "price": 0.8688
                }
            ],
            "cpusPerVm": 32,
            "memPerVm": 244,
            "gpusPerVm": 0,
            "ntwPerf": "10 Gigabit",
            "ntwPerfCategory": "high",
            "zones": null,
            "attributes": {
                "cpu": "32",
                "instanceTypeCategory": "Storage optimized",
                "memory": "244",
                "networkPerfCategory": "high"
            },
            "currentGen": true
        },
        {
            "category": "Memory optimized",
            "type": "r5.metal",
            "onDemandPrice": 7.104,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 1.6589
                },
                {
                    "zone": "eu-west-2b",
                    "price": 1.6589
                },
                {
                    "zone": "eu-west-2c",
                    "price": 1.6589
                }
            ],
            "cpusPerVm": 96,
            "memPerVm": 768,
            "gpusPerVm": 0,
            "ntwPerf": "25 Gigabit",
            "ntwPerfCategory": "extra",
            "zones": null,
            "attributes": {
                "cpu": "96",
                "instanceTypeCategory": "Memory optimized",
                "memory": "768",
                "networkPerfCategory": "extra"
            },
            "currentGen": true
        },
        {
            "category": "Storage optimized",
            "type": "i3en.large",
            "onDemandPrice": 0.263,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 0.0789
                },
                {
                    "zone": "eu-west-2b",
                    "price": 0.0789
                },
                {
                    "zone": "eu-west-2c",
                    "price": 0.0789
                }
            ],
            "cpusPerVm": 2,
            "memPerVm": 16,
            "gpusPerVm": 0,
            "ntwPerf": "Up to 25 Gigabit",
            "ntwPerfCategory": "extra",
            "zones": null,
            "attributes": {
                "cpu": "2",
                "instanceTypeCategory": "Storage optimized",
                "memory": "16",
                "networkPerfCategory": "extra"
            },
            "currentGen": true
        },
        {
            "category": "General purpose",
            "type": "m5d.2xlarge",
            "onDemandPrice": 0.524,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 0.132
                },
                {
                    "zone": "eu-west-2b",
                    "price": 0.1397
                },
                {
                    "zone": "eu-west-2c",
                    "price": 0.1345
                }
            ],
            "cpusPerVm": 8,
            "memPerVm": 32,
            "gpusPerVm": 0,
            "ntwPerf": "Up to 10 Gigabit",
            "ntwPerfCategory": "high",
            "zones": null,
            "attributes": {
                "cpu": "8",
                "instanceTypeCategory": "General purpose",
                "memory": "32",
                "networkPerfCategory": "high"
            },
            "currentGen": true
        },
        {
            "category": "Memory optimized",
            "type": "r5.4xlarge",
            "onDemandPrice": 1.184,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 0.3384
                },
                {
                    "zone": "eu-west-2b",
                    "price": 0.3194
                },
                {
                    "zone": "eu-west-2c",
                    "price": 0.2908
                }
            ],
            "cpusPerVm": 16,
            "memPerVm": 128,
            "gpusPerVm": 0,
            "ntwPerf": "10 Gigabit",
            "ntwPerfCategory": "high",
            "zones": null,
            "attributes": {
                "cpu": "16",
                "instanceTypeCategory": "Memory optimized",
                "memory": "128",
                "networkPerfCategory": "high"
            },
            "currentGen": true
        },
        {
            "category": "General purpose",
            "type": "m5ad.12xlarge",
            "onDemandPrice": 2.88,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 0.7918
                },
                {
                    "zone": "eu-west-2b",
                    "price": 0.7918
                },
                {
                    "zone": "eu-west-2c",
                    "price": 0.7918
                }
            ],
            "cpusPerVm": 48,
            "memPerVm": 192,
            "gpusPerVm": 0,
            "ntwPerf": "10 Gigabit",
            "ntwPerfCategory": "high",
            "zones": null,
            "attributes": {
                "cpu": "48",
                "instanceTypeCategory": "General purpose",
                "memory": "192",
                "networkPerfCategory": "high"
            },
            "currentGen": true
        },
        {
            "category": "Compute optimized",
            "type": "c5.9xlarge",
            "onDemandPrice": 1.818,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 0.5655
                },
                {
                    "zone": "eu-west-2b",
                    "price": 0.5655
                },
                {
                    "zone": "eu-west-2c",
                    "price": 0.5655
                }
            ],
            "cpusPerVm": 36,
            "memPerVm": 72,
            "gpusPerVm": 0,
            "ntwPerf": "10 Gigabit",
            "ntwPerfCategory": "high",
            "zones": null,
            "attributes": {
                "cpu": "36",
                "instanceTypeCategory": "Compute optimized",
                "memory": "72",
                "networkPerfCategory": "high"
            },
            "currentGen": true
        },
        {
            "category": "General purpose",
            "type": "m5d.8xlarge",
            "onDemandPrice": 2.096,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 0.5278
                },
                {
                    "zone": "eu-west-2b",
                    "price": 0.534
                },
                {
                    "zone": "eu-west-2c",
                    "price": 0.5278
                }
            ],
            "cpusPerVm": 32,
            "memPerVm": 128,
            "gpusPerVm": 0,
            "ntwPerf": "10 Gigabit",
            "ntwPerfCategory": "high",
            "zones": null,
            "attributes": {
                "cpu": "32",
                "instanceTypeCategory": "General purpose",
                "memory": "128",
                "networkPerfCategory": "high"
            },
            "currentGen": true
        },
        {
            "category": "General purpose",
            "type": "t2.small",
            "onDemandPrice": 0.026,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 0.0078
                },
                {
                    "zone": "eu-west-2b",
                    "price": 0.0078
                },
                {
                    "zone": "eu-west-2c",
                    "price": 0.0078
                }
            ],
            "cpusPerVm": 1,
            "memPerVm": 2,
            "gpusPerVm": 0,
            "ntwPerf": "Low to Moderate",
            "ntwPerfCategory": "low",
            "zones": null,
            "attributes": {
                "cpu": "1",
                "instanceTypeCategory": "General purpose",
                "memory": "2",
                "networkPerfCategory": "low"
            },
            "currentGen": true,
            "burst": true
        },
        {
            "category": "General purpose",
            "type": "m5ad.large",
            "onDemandPrice": 0.12,
            "spotPrice": [
                {
                    "zone": "eu-west-2c",
                    "price": 0.033
                },
                {
                    "zone": "eu-west-2a",
                    "price": 0.033
                },
                {
                    "zone": "eu-west-2b",
                    "price": 0.033
                }
            ],
            "cpusPerVm": 2,
            "memPerVm": 8,
            "gpusPerVm": 0,
            "ntwPerf": "Up to 10 Gigabit",
            "ntwPerfCategory": "high",
            "zones": null,
            "attributes": {
                "cpu": "2",
                "instanceTypeCategory": "General purpose",
                "memory": "8",
                "networkPerfCategory": "high"
            },
            "currentGen": true
        },
        {
            "category": "Memory optimized",
            "type": "r5.24xlarge",
            "onDemandPrice": 7.104,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 1.6589
                },
                {
                    "zone": "eu-west-2b",
                    "price": 1.6589
                },
                {
                    "zone": "eu-west-2c",
                    "price": 1.6589
                }
            ],
            "cpusPerVm": 96,
            "memPerVm": 768,
            "gpusPerVm": 0,
            "ntwPerf": "25 Gigabit",
            "ntwPerfCategory": "extra",
            "zones": null,
            "attributes": {
                "cpu": "96",
                "instanceTypeCategory": "Memory optimized",
                "memory": "768",
                "networkPerfCategory": "extra"
            },
            "currentGen": true
        },
        {
            "category": "General purpose",
            "type": "m5a.2xlarge",
            "onDemandPrice": 0.4,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 0.1349
                },
                {
                    "zone": "eu-west-2b",
                    "price": 0.1377
                },
                {
                    "zone": "eu-west-2c",
                    "price": 0.166
                }
            ],
            "cpusPerVm": 8,
            "memPerVm": 32,
            "gpusPerVm": 0,
            "ntwPerf": "Up to 10 Gigabit",
            "ntwPerfCategory": "high",
            "zones": null,
            "attributes": {
                "cpu": "8",
                "instanceTypeCategory": "General purpose",
                "memory": "32",
                "networkPerfCategory": "high"
            },
            "currentGen": true
        },
        {
            "category": "Memory optimized",
            "type": "r5.16xlarge",
            "onDemandPrice": 4.736,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 1.1059
                },
                {
                    "zone": "eu-west-2b",
                    "price": 1.1273
                },
                {
                    "zone": "eu-west-2c",
                    "price": 1.1234
                }
            ],
            "cpusPerVm": 64,
            "memPerVm": 512,
            "gpusPerVm": 0,
            "ntwPerf": "20 Gigabit",
            "ntwPerfCategory": "extra",
            "zones": null,
            "attributes": {
                "cpu": "64",
                "instanceTypeCategory": "Memory optimized",
                "memory": "512",
                "networkPerfCategory": "extra"
            },
            "currentGen": true
        },
        {
            "category": "Memory optimized",
            "type": "x1.32xlarge",
            "onDemandPrice": 16.806,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 5.0418
                },
                {
                    "zone": "eu-west-2b",
                    "price": 5.0418
                },
                {
                    "zone": "eu-west-2c",
                    "price": 5.0418
                }
            ],
            "cpusPerVm": 128,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "High",
            "ntwPerfCategory": "medium",
            "zones": null,
            "attributes": {
                "cpu": "128",
                "instanceTypeCategory": "Memory optimized",
                "memory": "1,952",
                "networkPerfCategory": "medium"
            },
            "currentGen": true
        },
        {
            "category": "Storage optimized",
            "type": "i3.2xlarge",
            "onDemandPrice": 0.724,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 0.2172
                },
                {
                    "zone": "eu-west-2b",
                    "price": 0.724
                },
                {
                    "zone": "eu-west-2c",
                    "price": 0.2172
                }
            ],
            "cpusPerVm": 8,
            "memPerVm": 61,
            "gpusPerVm": 0,
            "ntwPerf": "Up to 10 Gigabit",
            "ntwPerfCategory": "high",
            "zones": null,
            "attributes": {
                "cpu": "8",
                "instanceTypeCategory": "Storage optimized",
                "memory": "61",
                "networkPerfCategory": "high"
            },
            "currentGen": true
        },
        {
            "category": "General purpose",
            "type": "t3a.large",
            "onDemandPrice": 0.085,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 0.0255
                },
                {
                    "zone": "eu-west-2b",
                    "price": 0.0255
                },
                {
                    "zone": "eu-west-2c",
                    "price": 0.0255
                }
            ],
            "cpusPerVm": 2,
            "memPerVm": 8,
            "gpusPerVm": 0,
            "ntwPerf": "Low to Moderate",
            "ntwPerfCategory": "low",
            "zones": null,
            "attributes": {
                "cpu": "2",
                "instanceTypeCategory": "General purpose",
                "memory": "8",
                "networkPerfCategory": "low"
            },
            "currentGen": true,
            "burst": true
        },
        {
            "category": "Memory optimized",
            "type": "r5ad.16xlarge",
            "onDemandPrice": 4.928,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 1.1059
                },
                {
                    "zone": "eu-west-2b",
                    "price": 1.1059
                },
                {
                    "zone": "eu-west-2c",
                    "price": 1.1059
                }
            ],
            "cpusPerVm": 64,
            "memPerVm": 512,
            "gpusPerVm": 0,
            "ntwPerf": "12 Gigabit",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": {
                "cpu": "64",
                "instanceTypeCategory": "Memory optimized",
                "memory": "512",
                "networkPerfCategory": ""
            },
            "currentGen": true
        },
        {
            "category": "General purpose",
            "type": "m5a.12xlarge",
            "onDemandPrice": 2.4,
            "spotPrice": [
                {
                    "zone": "eu-west-2b",
                    "price": 0.7918
                },
                {
                    "zone": "eu-west-2c",
                    "price": 2.4
                },
                {
                    "zone": "eu-west-2a",
                    "price": 0.7918
                }
            ],
            "cpusPerVm": 48,
            "memPerVm": 192,
            "gpusPerVm": 0,
            "ntwPerf": "10 Gigabit",
            "ntwPerfCategory": "high",
            "zones": null,
            "attributes": {
                "cpu": "48",
                "instanceTypeCategory": "General purpose",
                "memory": "192",
                "networkPerfCategory": "high"
            },
            "currentGen": true
        },
        {
            "category": "Memory optimized",
            "type": "r5ad.12xlarge",
            "onDemandPrice": 3.696,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 0.8295
                },
                {
                    "zone": "eu-west-2b",
                    "price": 0.8295
                },
                {
                    "zone": "eu-west-2c",
                    "price": 0.8295
                }
            ],
            "cpusPerVm": 48,
            "memPerVm": 384,
            "gpusPerVm": 0,
            "ntwPerf": "10 Gigabit",
            "ntwPerfCategory": "high",
            "zones": null,
            "attributes": {
                "cpu": "48",
                "instanceTypeCategory": "Memory optimized",
                "memory": "384",
                "networkPerfCategory": "high"
            },
            "currentGen": true
        },
        {
            "category": "General purpose",
            "type": "m5ad.16xlarge",
            "onDemandPrice": 3.84,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 1.0557
                },
                {
                    "zone": "eu-west-2b",
                    "price": 1.0557
                },
                {
                    "zone": "eu-west-2c",
                    "price": 1.0557
                }
            ],
            "cpusPerVm": 64,
            "memPerVm": 256,
            "gpusPerVm": 0,
            "ntwPerf": "12 Gigabit",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": {
                "cpu": "64",
                "instanceTypeCategory": "General purpose",
                "memory": "256",
                "networkPerfCategory": ""
            },
            "currentGen": true
        },
        {
            "category": "General purpose",
            "type": "m5d.4xlarge",
            "onDemandPrice": 1.048,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 0.2639
                },
                {
                    "zone": "eu-west-2b",
                    "price": 0.2639
                },
                {
                    "zone": "eu-west-2c",
                    "price": 0.2639
                }
            ],
            "cpusPerVm": 16,
            "memPerVm": 64,
            "gpusPerVm": 0,
            "ntwPerf": "Up to 10 Gigabit",
            "ntwPerfCategory": "high",
            "zones": null,
            "attributes": {
                "cpu": "16",
                "instanceTypeCategory": "General purpose",
                "memory": "64",
                "networkPerfCategory": "high"
            },
            "currentGen": true
        },
        {
            "category": "General purpose",
            "type": "m5ad.24xlarge",
            "onDemandPrice": 5.76,
            "spotPrice": [
                {
                    "zone": "eu-west-2a",
                    "price": 1.5835
                },
                {
                    "zone": "eu-west-2b",
                    "price": 1.5835
                },
                {
                    "zone": "eu-west-2c",
                    "price": 1.5835
                }
            ],
            "cpusPerVm": 96,
            "memPerVm": 384,
            "gpusPerVm": 0,
            "ntwPerf": "20 Gigabit",
            "ntwPerfCategory": "extra",
            "zones": null,
            "attributes": {
                "cpu": "96",
                "instanceTypeCategory": "General purpose",
                "memory": "384",
                "networkPerfCategory": "extra"
            },
            "currentGen": true
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.2,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.2,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.2,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.2,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.2,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.2,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.2,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.2,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.2,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.2,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.2,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.2,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.2,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.2,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.2,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.2,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.2,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.2,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.2,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.2,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.2,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.2,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.2,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.2,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.2,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.2,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.2,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.2,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.2,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.2,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.2,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.2,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.2,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.2,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.2,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.2,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.2,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.2,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.2,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.2,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.2,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.2,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.2,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.2,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.2,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.2,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.2,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.2,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.2,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.2,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.2,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.2,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.2,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.2,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.2,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.2,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.2,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.2,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.2,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.2,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.2,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.2,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.2,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.2,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.2,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.2,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.2,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.2,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.2,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.2,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.2,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.2,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.2,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.2,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.2,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.2,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.2,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.2,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.2,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.2,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.2,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.2,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.2,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.2,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.1,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.1,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.1,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.1,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.1,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.1,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.1,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.1,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.1,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.1,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.1,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.1,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.1,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.1,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.1,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.1,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.1,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.1,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.1,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.1,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.1,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.1,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.1,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.1,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.1,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.1,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.1,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.1,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.1,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.1,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.1,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.1,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.1,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.1,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.1,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.1,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.1,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.1,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.1,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.1,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.1,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.1,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.1,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.1,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.1,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.1,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.1,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.1,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.1,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.1,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.1,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.1,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.1,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.1,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.1,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.1,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.1,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.1,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.1,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.1,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.1,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.1,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.1,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.1,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.1,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.1,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.1,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.1,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.1,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.1,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.1,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.1,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.1,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.1,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.1,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.1,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.1,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.1,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.1,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        },
        {
            "category": "",
            "type": "EKS Control Plane",
            "onDemandPrice": 0.1,
            "spotPrice": null,
            "cpusPerVm": 0,
            "memPerVm": 0,
            "gpusPerVm": 0,
            "ntwPerf": "",
            "ntwPerfCategory": "",
            "zones": null,
            "attributes": null,
            "currentGen": false
        }
    ],
    "scrapingTime": "1592310039108"
}`
