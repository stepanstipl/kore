import { ConfigureCloudPage } from '../configure-cloud'
import { clearFillTextInput, modalYes, waitForDrawerOpenClose } from '../../../utils'

export class ConfigureCloudGCPOrgs extends ConfigureCloudPage {
  constructor(p) {
    super(p)
    this.pagePath = '/configure/cloud/GCP/orgs'
  }

  async openTab() {
    await this.selectCloud('gcp')
    await this.selectSubTab('Organization credentials', 'GCP/orgs')
  }

  verifyPageURL() {
    expect([this.pagePath, '/configure/cloud'].includes(this.p.url()))
  }

  /**
   * checks if any org is configured
   */
  async orgConfigured() {
    try {
      await this.p.waitForSelector('.ant-list-items', { timeout: 2000 })
      return true
    } catch(err) {
      return false
    }
  }

  /**
   * Checks if a specific org is listed
   */
  async checkOrgListed(name) {
    await expect(this.p).toMatchElement(`#gcporg_${name}`)
  }

  async add() {
    await this.p.waitFor('.new-gcp-organization')
    await expect(this.p).toClick('button', { text: 'Configure' })
    await waitForDrawerOpenClose(this.p)
    await expect(this.p).toMatch('New GCP organization')
  }

  async edit(name, parentID) {
    await this.p.click(`a#gcporg_edit_${name}`)
    await waitForDrawerOpenClose(this.p)
    await expect(this.p).toMatch(`GCP Organization: ${parentID}`)
  }

  async populate({ name, summary, parentID, billingAccount, json }) {
    await clearFillTextInput(this.p, 'gcp_organization_name', name)
    await clearFillTextInput(this.p, 'gcp_organization_summary', summary)
    await clearFillTextInput(this.p, 'gcp_organization_parentID', parentID)
    await clearFillTextInput(this.p, 'gcp_organization_billingAccount', billingAccount)
    if (json !== undefined) {
      await this.p.type('textarea#gcp_organization_account', json)
    }
  }

  async replaceKey(json) {
    await this.p.type('input#gcporg_replace_key', ' ')
    // Wait for service account text field to be shown:
    await expect(this.p).toMatch('Service Account JSON')
    await this.p.type('textarea#gcp_organization_account', json)
  }

  async save() {
    await this.p.click('button#save')
    await waitForDrawerOpenClose(this.p)
  }

  async delete(name) {
    await this.p.click(`a#gcporg_del_${name}`)
  }

  async confirmDelete() {
    await modalYes(this.p, 'Are you sure you want to delete the GCP Organization')
  }
}