import { ConfigureCloudPage } from '../configure-cloud'
import { clearFillTextInput, modalYes, waitForDrawerOpenClose } from '../../../utils'

export class ConfigureCloudGCPProjects extends ConfigureCloudPage {
  constructor(p) {
    super(p)
    this.pagePath = '/configure/cloud/GCP/projects'
  }

  async openTab() {
    await this.selectCloud('gcp')
    await this.selectSubTab('Project credentials', 'GCP/projects')
  }

  /**
   * Checks credential for project listed in list of credentials
   */
  async checkCredentialListed(name) {
    await expect(this.p).toMatchElement(`#gkecreds_${name}`)
  }

  async add() {
    await expect(this.p).toClick('button', { text: '+ New' })
    await waitForDrawerOpenClose(this.p)
    await expect(this.p).toMatch('New GCP project')
  }

  async edit(name, project) {
    await this.p.click(`a#gkecreds_edit_${name}`)
    await waitForDrawerOpenClose(this.p)
    await expect(this.p).toMatch(`GCP project: ${project}`)
  }

  async populate({ name, summary, project, json }) {
    await clearFillTextInput(this.p, 'gke_credentials_name', name)
    await clearFillTextInput(this.p, 'gke_credentials_summary', summary)
    await clearFillTextInput(this.p, 'gke_credentials_project', project)
    if (json !== undefined) {
      await this.p.type('textarea#gke_credentials_account', json)
    }
  }

  async replaceKey(json) {
    await this.p.type('input#gke_credentials_replace_key', ' ')
    // Wait for service account text field to be shown:
    await expect(this.p).toMatch('Service Account JSON')
    await this.p.type('textarea#gke_credentials_account', json)
  }

  async save() {
    await this.p.click('button#save')
    await waitForDrawerOpenClose(this.p)
  }

  async delete(name) {
    await this.p.click(`a#gkecreds_del_${name}`)
  }

  async confirmDelete() {
    await modalYes(this.p, 'Are you sure you want to delete the credentials')
  }
}