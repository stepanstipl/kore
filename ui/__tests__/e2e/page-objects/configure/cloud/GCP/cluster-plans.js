import { ConfigureCloudClusterPlansBase } from '../cluster-plans-base'
import { clearFillTextInput, waitForDrawerOpenClose, setSwitch, setCascader } from '../../../utils'

export class ConfigureCloudGCPClusterPlans extends ConfigureCloudClusterPlansBase {
  constructor(p) {
    super(p)
    this.pagePath = '/configure/cloud/GCP/plans'
  }

  async openTab() {
    await this.selectCloud('gcp')
    await this.selectSubTab('Cluster plans', 'GCP/plans')
  }

  async populatePlan({ description, name, planDescription, region }) {
    await this.viewPlanConfig()
    await clearFillTextInput(this.p, 'plan_summary', description)
    await clearFillTextInput(this.p, 'plan_description', name)
    await clearFillTextInput(this.p, 'plan_input_description', planDescription)
    await setCascader(this.p, 'plan_input_region', region)
  }

  async addNodePool() {
    await this.p.click('button#plan_nodepool_add')
    await waitForDrawerOpenClose(this.p)
  }

  async viewEditNodePool(idx) {
    await this.viewPlanConfig()
    await waitForDrawerOpenClose(this.p)
    await this.p.click(`a#plan_nodepool_${idx}_viewedit`)
    await waitForDrawerOpenClose(this.p)
  }

  async populateNodePool({ name, enableAutoscaler, minSize, size, maxSize }) {
    await clearFillTextInput(this.p, 'plan_nodepool_name', name)
    await setSwitch(this.p, 'plan_nodepool_enableAutoscaler', enableAutoscaler)
    if (enableAutoscaler) {
      await clearFillTextInput(this.p, 'plan_nodepool_minSize', minSize)
    }
    await clearFillTextInput(this.p, 'plan_nodepool_size', size)
    if (enableAutoscaler) {
      await clearFillTextInput(this.p, 'plan_nodepool_maxSize', maxSize)
    }
  }

  async closeNodePoolDisabled() {
    return (await this.p.$('button#plan_nodepool_close[disabled]')) !== null
  }

  async closeNodePool() {
    await this.p.click('button#plan_nodepool_close')
    await waitForDrawerOpenClose(this.p)
  }
}
