const { BasePage } = require('./base')

export class SetupKoreCloudAccessPage extends BasePage {
  constructor(p) {
    super(p)
    this.pagePath = '/setup/kore/cloud-access'
  }

  async selectCloud(name) {
    await this.p.click(`#${name}`)
  }

  async selectKoreManagedGcpProjects() {
    await this.p.click('.use-kore-managed-projects')
  }

  async selectKoreManagedProjects(type) {
    // cluster or custom
    await this.p.click(`.automated-projects-${type}`)
  }

  async addGcpOrganization() {
    await this.p.waitFor('.new-gcp-organization')
    await this.p.click('.new-gcp-organization')
    await this.p.waitFor('#gcp_organization_parentID')
    await this.p.type('#gcp_organization_parentID', 'test-org')
    await this.p.type('#gcp_organization_billingAccount', 'BILL-1234')
    await this.p.type('#gcp_organization_account', 'invalid service account JSON')
    await this.p.type('#gcp_organization_name', 'Test org')
    await this.p.type('#gcp_organization_summary', 'Org used by automated test')
    await Promise.all([
      this.p.waitFor('#gcp_organization_parentID', { hidden: true }),
      this.p.click('#save')
    ])
  }

  async nextStep() {
    await this.p.click('.steps-action .ant-btn-primary')
  }

  async save() {
    await this.p.click('.steps-action .ant-btn-primary')
    await this.p.waitFor('.kore-managed-setup-complete')
  }

  async setAutomatedProjectDefaults() {
    await this.p.waitFor('.set-kore-defaults')
    await this.p.click('.set-kore-defaults')
  }

  async enterGkeCredentials() {
    await this.p.type('#gke_credentials_project', 'test-project')
    await this.p.type('#gke_credentials_account', 'invalid service account JSON')
    await this.p.type('#gke_credentials_name', 'Test project')
    await this.p.type('#gke_credentials_summary', 'Project used by automated test')
  }

  async continueWithoutVerification() {
    await this.p.waitFor('#continue-without-verification')
    await Promise.all([
      this.p.waitForNavigation(),
      await this.p.click('#continue-without-verification')
    ])
  }
}
