const dummyEvents = [
  { 
    'id': '1', 
    'level': 'info', 
    'headline': 'Plan \'gke-development\' has PSP enabled', 
    'ruleSummary': 'PSP-01-CheckPlanPSP',
    'ruleID': '/security/rules/1', 
    'description': 'Pod Security Policy is recommended to ensure that workloads within a cluster cannot be run as a privileged user. This check ensures that a plan is configured to enable PSP on all clusters created using that plan.',
    'resourceType': 'plan', 
    'resourceID': 'example-plan' 
  },
  { 
    'id': '2', 
    'level': 'warn', 
    'headline': 'Plan \'gke-production\' does not specify PSP enabled', 
    'ruleSummary': 'PSP-01-CheckPlanPSP',
    'ruleID': '/security/rules/1', 
    'description': 'Pod Security Policy is recommended to ensure that workloads within a cluster cannot be run as a privileged user. This check ensures that a plan is configured to enable PSP on all clusters created using that plan.',
    'resourceType': 'plan', 
    'resourceID': 'example-plan' 
  },
  { 
    'id': '3', 
    'level': 'critical', 
    'headline': 'Plan \'eks-development\' has PSP disabled', 
    'ruleSummary': 'PSP-01-CheckPlanPSP',
    'ruleID': '/security/rules/1', 
    'description': 'Pod Security Policy is recommended to ensure that workloads within a cluster cannot be run as a privileged user. This check ensures that a plan is configured to enable PSP on all clusters created using that plan.',
    'resourceType': 'plan', 
    'resourceID': 'example-plan' 
  },
]

const dummyOverview = {
  'status': {
    'critical': 3,
    'warn': 32,
    'info': 20
  },
  'overallStatus': 'critical',
  'teamSummary': [
    { 
      'name': 'Example Team 1', 
      'overallStatus': 'critical',   
      'status': {
        'critical': 1,
        'warn': 0,
        'info': 2
      }
    },
    { 
      'name': 'Example Team 2', 
      'overallStatus': 'info',   
      'status': {
        'critical': 0,
        'warn': 0,
        'info': 15
      }
    },
  ],
  'planSummary': [
    {
      'name': 'Example Plan 1',
      'overallStatus': 'warn',   
      'status': {
        'critical': 0,
        'warn': 12,
        'info': 0
      }
    },
    {
      'name': 'Example Plan 2',
      'overallStatus': 'critical',   
      'status': {
        'critical': 2,
        'warn': 20,
        'info': 3
      }
    },
  ]
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
  static events = dummyEvents
}

export default SecurityData