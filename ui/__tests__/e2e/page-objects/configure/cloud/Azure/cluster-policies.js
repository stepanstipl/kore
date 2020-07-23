import { ConfigureCloudClusterPoliciesBase } from '../cluster-policies-base'

export class ConfigureCloudAzureClusterPolicies extends ConfigureCloudClusterPoliciesBase {
  constructor(p) {
    super(p)
    this.pagePath = '/configure/cloud/Azure/policies'
  }

  async openTab() {
    await this.selectCloud('azure')
    await this.selectSubTab('Cluster policies', 'Azure/policies')
  }
}
