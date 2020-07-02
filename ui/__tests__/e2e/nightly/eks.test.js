const { TeamCluster } = require('./team-cluster')

describe('EKS end-to-end', () => {
  new TeamCluster({
    provider: 'EKS',
    plan: 'EKS Development Cluster',
    timeouts: {
      // 30 minutes
      create: 30 * 60 * 1000,
      // 20 minutes
      delete: 20 * 60 * 1000,
    }
  }).run()
})
