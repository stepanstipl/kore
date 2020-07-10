import { ConfigureCloudClusterPlansBase } from '../cluster-plans-base'
import { clearFillTextInput, waitForDrawerOpenClose, setCascader } from '../../../utils'

export class ConfigureCloudAWSClusterPlans extends ConfigureCloudClusterPlansBase {
  constructor(p) {
    super(p)
    this.pagePath = '/configure/cloud/AWS/plans'
  }

  async openTab() {
    await this.selectCloud('aws')
    await this.selectSubTab('Cluster Plans', 'AWS/plans')
  }

  async populatePlan({ description, name, planDescription, region, version }) {
    await clearFillTextInput(this.p, 'plan_summary', description)
    await clearFillTextInput(this.p, 'plan_description', name)
    await clearFillTextInput(this.p, 'plan_input_description', planDescription)
    await setCascader(this.p, 'plan_input_region', region)
    await clearFillTextInput(this.p, 'plan_input_version', version)
  }

  async addNodeGroup() {
    await this.p.click('button#plan_nodegroup_add')
    await waitForDrawerOpenClose(this.p)
  }

  async viewEditNodeGroup(idx) {
    await this.p.click(`a#plan_nodegroup_${idx}_viewedit`)
    await waitForDrawerOpenClose(this.p)
  }

  async populateNodeGroup({ name, minSize, desiredSize, maxSize }) {
    await clearFillTextInput(this.p, 'plan_nodegroup_name', name)
    await clearFillTextInput(this.p, 'plan_nodegroup_minSize', minSize)
    await clearFillTextInput(this.p, 'plan_nodegroup_desiredSize', desiredSize)
    await clearFillTextInput(this.p, 'plan_nodegroup_maxSize', maxSize)
  }

  async closeNodeGroupDisabled() {
    return (await this.p.$('button#plan_nodegroup_close[disabled]')) !== null
  }

  async closeNodeGroup() {
    await this.p.click('button#plan_nodegroup_close')
    await waitForDrawerOpenClose(this.p)
  }
}
