import { ConfigureCloudPage } from './configure-cloud'
import { modalYes, waitForDrawerOpenClose } from '../../utils'

export class ConfigureCloudClusterPlansBase extends ConfigureCloudPage {
  async listLoaded() {
    await this.p.waitForSelector('#plans_list', { timeout: 1000 })
  }

  async view(name) {
    await this.p.click(`a#plans_view_${name}`)
    await waitForDrawerOpenClose(this.p)
  }

  async edit(name) {
    await this.p.click(`a#plans_edit_${name}`)
    await waitForDrawerOpenClose(this.p)
  }

  async delete(name) {
    await this.p.click(`a#plans_delete_${name}`)
  }

  async confirmDelete() {
    await modalYes(this.p, 'Are you sure you want to delete the plan')
  }

  async new() {
    await this.p.click('button#add')
    await waitForDrawerOpenClose(this.p)
  }

  async save() {
    await this.p.click('button#plan_save')
    await waitForDrawerOpenClose(this.p)
  }

  async viewPlanConfig() {
    await this.p.evaluate(() => {
      document.querySelector('#plan_config').scrollIntoView()
    })
  }
}
