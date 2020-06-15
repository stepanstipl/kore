import { ConfigureCloudClusterPoliciesBase } from '../cluster-policies-base'

export class ConfigureCloudGCPClusterPolicies extends ConfigureCloudClusterPoliciesBase {
  constructor(p) {
    super(p)
    this.pagePath = '/configure/cloud/GCP/policies'
  }

  async openTab() {
    await this.selectCloud('gcp')
    await this.selectSubTab('Cluster Policies', 'GCP/policies')
  }
}
