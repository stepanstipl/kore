const { BasePage } = require('../../base')

export class ConfigureCloudPage extends BasePage {
  constructor(p) {
    super(p)
    this.pagePath = '/configure/cloud'
  }

  async selectCloud(name) {
    await this.p.click(`#tab-${name}`)
  }

  async selectSubTab(name) {
    const subTab = await this.p.$x(`//div[@id='cloud_subtabs']//div[@role='tab' and text()='${name}']`)
    if (subTab.length === 0) {
      throw new Error(`No sub-tab exists with name ${name}`)
    }
    await Promise.all([
      this.p.waitForNavigation(),
      subTab[0].click()
    ])
  }
}
