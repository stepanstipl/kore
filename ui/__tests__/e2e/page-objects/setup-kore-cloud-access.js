const { BasePage } = require('./base')

export class SetupKoreCloudAccessPage extends BasePage {
  constructor(page) {
    super(page)
    this.pagePath = '/setup/kore/cloud-access'
  }

  async selectCloud(name) {
    await this.page.click(`#${name}`)
  }

  async selectKoreManagedGcpProjects() {
    await this.page.click('.use-kore-managed-projects')
  }

  async selectKoreManagedProjects(type) {
    // cluster or custom
    await this.page.click(`.automated-projects-${type}`)
  }

  async addGcpOrganization() {
    await this.page.click('.new-gcp-organization')
    await this.page.waitFor('#gcp_organization_parentID')
    await this.page.type('#gcp_organization_parentID', 'test-org')
    await this.page.type('#gcp_organization_billingAccount', 'BILL-1234')
    await this.page.type('#gcp_organization_account', 'invalid service account JSON')
    await this.page.type('#gcp_organization_name', 'Test org')
    await this.page.type('#gcp_organization_summary', 'Org used by automated test')
    await Promise.all([
      this.page.waitFor('#gcp_organization_parentID', { hidden: true }),
      this.page.click('#save')
    ])
  }

  async nextStep() {
    await this.page.click('.steps-action .ant-btn-primary')
  }

  async save() {
    await this.page.click('.steps-action .ant-btn-primary')
    await this.page.waitFor('.kore-managed-setup-complete')
  }

  async setAutomatedProjectDefaults() {
    await this.page.waitFor('.set-kore-defaults')
    await this.page.click('.set-kore-defaults')
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
