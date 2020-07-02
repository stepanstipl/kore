const { TeamCluster } = require('./team-cluster')

describe('GKE end-to-end', () => {
  new TeamCluster({
    provider: 'GKE',
    plan: 'GKE Development Cluster',
    timeouts: {
      // 15 minutes
      create: 15 * 60 * 1000,
      // 10 minutes
      delete: 10 * 60 * 1000,
    }
  }).run()
})
