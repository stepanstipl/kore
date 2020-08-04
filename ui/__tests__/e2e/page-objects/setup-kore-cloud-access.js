import { BasePage } from './base'
import { ConfigureCloudGCPOrgs } from './configure/cloud/GCP/organizations'

export class SetupKoreCloudAccessPage extends BasePage {
  constructor(p) {
    super(p)
    this.pagePath = '/setup/kore/cloud-access'
    this.orgsPage = new ConfigureCloudGCPOrgs(p)
  }

  async selectCloud(name) {
    await this.p.click(`#${name}`)
  }

  async selectKoreManagedGcpProjects() {
    await this.p.click('.use-kore-managed-projects')
  }

  async addGcpOrganization(testOrg) {
    await this.orgsPage.add()
    await this.orgsPage.populate(testOrg)
    await this.orgsPage.save()
  }

  async nextStep() {
    await this.p.click('.steps-action .ant-btn-primary')
  }

  async save() {
    await this.p.click('.steps-action .ant-btn-primary')
    await this.p.waitFor('.kore-managed-setup-complete')
  }

  async setAutomatedAccountDefaults() {
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
