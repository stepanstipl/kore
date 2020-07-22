import { ConfigureCloudClusterPlansBase } from '../cluster-plans-base'
import { clearFillTextInput, setSelect, waitForDrawerOpenClose } from '../../../utils'

export class ConfigureCloudAzureClusterPlans extends ConfigureCloudClusterPlansBase {
  constructor(p) {
    super(p)
    this.pagePath = '/configure/cloud/Azure/plans'
  }

  async openTab() {
    await this.selectCloud('azure')
    await this.selectSubTab('Cluster plans', 'Azure/plans')
  }

  async populatePlan({ description, name, planDescription, region, version, dnsPrefix, networkPlugin }) {
    await this.viewPlanConfig()
    await clearFillTextInput(this.p, 'plan_summary', description)
    await clearFillTextInput(this.p, 'plan_description', name)
    await clearFillTextInput(this.p, 'plan_input_description', planDescription)
    await clearFillTextInput(this.p, 'plan_input_region', region)
    // wait for version control to be enabled after selecting the region
    await this.p.waitForSelector('#plan_input_version.ant-select-disabled', { hidden: true })
    await clearFillTextInput(this.p, 'plan_input_version', version)
    await clearFillTextInput(this.p, 'plan_input_dnsPrefix', dnsPrefix)
    await setSelect(this.p, 'plan_input_networkPlugin', networkPlugin)
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

  async populateNodePool({ name, mode, minSize, size, maxSize }) {
    await clearFillTextInput(this.p, 'plan_nodepool_name', name)
    await setSelect(this.p, 'plan_nodepool_mode', mode)
    await clearFillTextInput(this.p, 'plan_nodepool_minSize', minSize)
    await clearFillTextInput(this.p, 'plan_nodepool_size', size)
    await clearFillTextInput(this.p, 'plan_nodepool_maxSize', maxSize)
  }

  async closeNodePoolDisabled() {
    return (await this.p.$('button#plan_nodepool_close[disabled]')) !== null
  }

  async closeNodePool() {
    await this.p.click('button#plan_nodepool_close')
    await waitForDrawerOpenClose(this.p)
  }
}
