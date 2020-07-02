const { BasePage } = require('../../base')
import { clearFillTextInput } from '../../utils'

export class NewClusterPage extends BasePage {
  constructor(p, teamID) {
    super(p)
    this.pagePath = `/teams/${teamID}/clusters/new`
  }

  /**
   * Click the cloud selector, cloud can be one of gcp, aws or azure
   * @param cloud
   */
  async selectCloud(cloud) {
    expect(this.p).toClick(`#${cloud.toLowerCase()}`)
  }

  async selectPlan(plan) {
    expect(this.p).toClick('#cluster_options_plan span', { text: plan })
  }

  async populate({ cloud, plan, name }) {
    await this.selectCloud(cloud)
    await this.selectPlan(plan)
    await clearFillTextInput(this.p, 'cluster_options_clusterName', name)
  }

  async save() {
    await Promise.all([
      this.p.waitForNavigation(),
      await this.p.click('button#save')
    ])
  }

}
