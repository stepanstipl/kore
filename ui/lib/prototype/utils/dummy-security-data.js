/* eslint-disable quotes */
const dummyOverview = {
  "kind": "SecurityOverview",
  "apiVersion": "security.kore.appvia.io/v1",
  "metadata": {
    "name": "overview",
    "creationTimestamp": null
  },
  "spec": {
    "openIssueCounts": {
      "Failure": 3,
      "Warning": 32,
      "Compliant": 20,
    },
    "resources": [
      {
        "resource": {
          "group": "clusters.compute.kore.appvia.io",
          "version": "v1",
          "kind": "Cluster",
          "namespace": "example-team",
          "name": "example-cluster-1"
        },
        "lastChecked": "2020-06-17T07:08:17Z",
        "overallStatus": "Failure",
        "openIssueCounts": {
          "Failure": 1,
          "Compliant": 2
        }
      },
      {
        "resource": {
          "group": "clusters.compute.kore.appvia.io",
          "version": "v1",
          "kind": "Cluster",
          "namespace": "example-team",
          "name": "example-cluster-2"
        },
        "lastChecked": "2020-06-12T15:06:34Z",
        "overallStatus": "Compliant",
        "openIssueCounts": {
          "Compliant": 15
        }
      },
      {
        "resource": {
          "group": "config.kore.appvia.io",
          "version": "v1",
          "kind": "Plan",
          "namespace": "kore",
          "name": "example-plan-1"
        },
        "lastChecked": "2020-06-17T05:46:43Z",
        "overallStatus": "Warning",
        "openIssueCounts": {
          "Warning": 12
        }
      },
      {
        "resource": {
          "group": "config.kore.appvia.io",
          "version": "v1",
          "kind": "Plan",
          "namespace": "kore",
          "name": "example-plan-2"
        },
        "lastChecked": "2020-06-17T05:46:43Z",
        "overallStatus": "Failure",
        "openIssueCounts": {
          "Failure": 2,
          "Warning": 20,
          "Compliant": 2
        }
      }
    ]
  }
}

const dummyRules = [
  { 
    'id': '1', 
    'ruleName': 'PSP-01-CheckPlanPSP',
    'ruleID': '/security/rules/1', 
    'description': 'Pod Security Policy is recommended to ensure that workloads within a cluster cannot be run as a privileged user. This check ensures that a plan is configured to enable PSP on all clusters created using that plan.',
    'details': ''
  },
]

class SecurityData {
  static overview = dummyOverview
  static rules = dummyRules
}

export default SecurityData