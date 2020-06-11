import { ConfigureCloudPage } from '../configure-cloud'
import { clearFillTextInput, modalYes, waitForDrawerOpenClose, setSwitch } from '../../../utils'

export class ConfigureCloudGCPClusterPlans extends ConfigureCloudPage {
  constructor(p) {
    super(p)
    this.pagePath = '/configure/cloud/GCP/plans'
  }

  async openTab() {
    await this.selectCloud('gcp')
    await this.selectSubTab('Cluster Plans', 'GCP/plans')
  }

  async view(name) {
    await this.p.click(`a#gkeplans_view_${name}`)
    await waitForDrawerOpenClose(this.p)
  }

  async edit(name) {
    await this.p.click(`a#gkeplans_edit_${name}`)
    await waitForDrawerOpenClose(this.p)
  }

  async delete(name) {
    await this.p.click(`a#gkeplans_delete_${name}`)
  }

  async confirmDelete() {
    await modalYes(this.p, 'Are you sure you want to delete the plan')
  }

  async new() {
    await expect(this.p).toClick('button', { text: '+ New' })
    await waitForDrawerOpenClose(this.p)
    await expect(this.p).toMatch('New GKE plan')
  }

  async populatePlan({ description, name, planDescription, region }) {
    await clearFillTextInput(this.p, 'plan_summary', description)
    await clearFillTextInput(this.p, 'plan_description', name)
    await clearFillTextInput(this.p, 'plan_input_description', planDescription)
    await clearFillTextInput(this.p, 'plan_input_region', region)
  }

  async addNodePool() {
    await this.p.click('button#plan_nodepool_add')
    await waitForDrawerOpenClose(this.p)
  }

  async viewEditNodePool(idx) {
    await this.p.click(`a#plan_nodepool_${idx}_viewedit`)
    await waitForDrawerOpenClose(this.p)
  }

  async populateNodePool({ name, enableAutoscaler, minSize, size, maxSize }) {
    await clearFillTextInput(this.p, 'plan_nodepool_name', name)
    await setSwitch(this.p, 'plan_nodepool_enableAutoscaler', enableAutoscaler)
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

  async save() {
    await this.p.click('button#plan_save')
    await waitForDrawerOpenClose(this.p)
  }
}
