const { BasePage } = require('./base')

export class SetupKoreCloudProvidersPage extends BasePage {
  constructor(page) {
    super(page)
    this.pagePath = '/setup/kore/cloud-providers'
  }

  async selectCloud(name) {
    await this.page.click(`#${name}`)
  }

  async enterGkeCredentials() {
    await this.page.type('#gke_credentials_project', 'test-project')
    await this.page.type('#gke_credentials_account', 'invalid service account JSON')
    await this.page.type('#gke_credentials_name', 'Test project')
    await this.page.type('#gke_credentials_summary', 'Project used by automated test')
  }

  async continueWithoutVerification() {
    await this.page.waitFor('#continue-without-verification')
    await Promise.all([
      this.page.waitForNavigation(),
      await this.page.click('#continue-without-verification')
    ])
  }
}
