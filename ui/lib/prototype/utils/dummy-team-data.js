const team = {
  'kind': 'Team',
  'apiVersion': 'org.kore.appvia.io/v1',
  'metadata': {
    'name': 'demo',
    'namespace': 'kore',
    'creationTimestamp': '2020-05-05T09:31:14Z'
  },
  'spec': {
    'summary': 'Demo team',
    'description': 'This is a demo team'
  },
  'status': {
    'conditions': [],
    'status': 'Success'
  }
}

const members = {
  'kind': 'List',
  'apiVersion': 'v1',
  'metadata': {},
  'items': [
    'dave.thompson@appvia.io',
    'jon.shanks@appvia.io'
  ]
}

const clusters = {
  'metadata': {
    'selfLink': '/apis/clusters.compute.kore.appvia.io/v1/namespaces/demo/clusters',
    'resourceVersion': '486395'
  },
  'items': [
    {
      'kind': 'Cluster',
      'apiVersion': 'clusters.compute.kore.appvia.io/v1',
      'metadata': {
        'name': 'demo-notprod',
        'namespace': 'demo'
      },
      'spec': {
        'kind': 'EKS',
        'plan': 'eks-development',
      },
      'status': {
        'status': 'Success'
      }
    }
  ]
}

const namespaceClaims = {
  'metadata': {
    'selfLink': '/apis/clusters.compute.kore.appvia.io/v1/namespaces/demo/namespaceclaims',
    'resourceVersion': '487255'
  },
  'items': [
    {
      'kind': 'NamespaceClaim',
      'apiVersion': 'clusters.compute.kore.appvia.io/v1',
      'metadata': {
        'name': 'demo-notprod-dev',
        'namespace': 'demo',
      },
      'spec': {
        'cluster': {
          'group': 'clusters.compute.kore.appvia.io',
          'version': 'v1',
          'kind': 'Cluster',
          'namespace': 'demo',
          'name': 'demo-notprod'
        },
        'name': 'dev'
      },
      'status': {
        'status': 'Success'
      }
    },
    {
      'kind': 'NamespaceClaim',
      'apiVersion': 'clusters.compute.kore.appvia.io/v1',
      'metadata': {
        'name': 'demo-notprod-qa',
        'namespace': 'demo',
      },
      'spec': {
        'cluster': {
          'group': 'clusters.compute.kore.appvia.io',
          'version': 'v1',
          'kind': 'Cluster',
          'namespace': 'demo',
          'name': 'demo-notprod'
        },
        'name': 'qa'
      },
      'status': {
        'status': 'Success'
      }
    }
  ]
}

const allocations = {
  'metadata': {},
  'items': []
}

class TeamData {
  static team = team
  static members = members
  static clusters = clusters
  static namespaceClaims = namespaceClaims
  static allocations = allocations
}

export default TeamData
