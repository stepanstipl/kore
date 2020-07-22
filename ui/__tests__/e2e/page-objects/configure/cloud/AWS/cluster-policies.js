import { ConfigureCloudClusterPoliciesBase } from '../cluster-policies-base'

export class ConfigureCloudAWSClusterPolicies extends ConfigureCloudClusterPoliciesBase {
  constructor(p) {
    super(p)
    this.pagePath = '/configure/cloud/AWS/policies'
  }

  async openTab() {
    await this.selectCloud('aws')
    await this.selectSubTab('Cluster policies', 'AWS/policies')
  }
}
