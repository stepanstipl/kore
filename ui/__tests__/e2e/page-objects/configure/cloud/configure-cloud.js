const { BasePage } = require('../../base')

export class ConfigureCloudPage extends BasePage {
  constructor(p) {
    super(p)
    this.pagePath = '/configure/cloud'
  }

  async selectCloud(name) {
    // Check if we're already on the right page (or the default page is what we want)
    const currUrl = await this.p.url()
    if (currUrl.toLowerCase().indexOf(`/configure/cloud/${name.toLowerCase()}/`) > -1 || (name.toLowerCase() === 'gcp' && currUrl.endsWith('/configure/cloud'))) {
      return
    }
    await this.p.waitFor(100)
    await Promise.all([
      this.p.waitForNavigation(),
      this.p.click(`#tab-${name}`)
    ])
  }

  async selectSubTab(name, urlName) {
    const subTab = await this.p.$x(`//div[@id='cloud_subtabs']//div[@role='tab' and text()='${name}']`)
    if (subTab.length === 0) {
      throw new Error(`No sub-tab exists with name ${name}`)
    }
    // Check if we're already on the right page/tab (or the default page/tab is what we want)
    const currUrl = await this.p.url()
    if (currUrl.endsWith(`/configure/cloud/${urlName}`) || (urlName === 'GCP/orgs' && currUrl.endsWith('/configure/cloud'))) {
      return
    }
    await this.p.waitFor(100)
    await Promise.all([
      this.p.waitForNavigation(),
      subTab[0].click()
    ])
  }
}
