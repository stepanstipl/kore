import { ConfigureCloudPage } from '../configure-cloud'
import { waitForDrawerOpenClose, clearFillTextInput, modalYes } from '../../../utils'

export class ConfigureCloudAzureSubscriptions extends ConfigureCloudPage {
  constructor(p) {
    super(p)
    this.pagePath = '/configure/cloud/Azure/subscriptions'
  }

  async openTab() {
    await this.selectCloud('azure')
    await this.selectSubTab('Subscription credentials', 'Azure/subscriptions')
  }

  async checkSubscriptionListed(name) {
    await expect(this.p).toMatchElement(`#akscreds_${name}`)
  }

  async add() {
    await expect(this.p).toClick('button', { text: '+ New' })
    await waitForDrawerOpenClose(this.p)
    await expect(this.p).toMatch('New Azure subscription')
  }

  async edit(name, subscriptionID) {
    await this.p.click(`a#akscreds_edit_${name}`)
    await waitForDrawerOpenClose(this.p)
    await expect(this.p).toMatch(`Azure subscription: ${subscriptionID}`)
  }

  async populate({ subscriptionID, tenantID, clientID, clientSecret, name, summary }) {
    await clearFillTextInput(this.p, 'aks_credentials_subscriptionID', subscriptionID)
    await clearFillTextInput(this.p, 'aks_credentials_tenantID', tenantID)
    await clearFillTextInput(this.p, 'aks_credentials_clientID', clientID)
    await clearFillTextInput(this.p, 'aks_credentials_clientSecret', clientSecret)
    await clearFillTextInput(this.p, 'aks_credentials_name', name)
    await clearFillTextInput(this.p, 'aks_credentials_summary', summary)
  }

  async replacePassword(clientSecret) {
    await this.p.type('input#aks_credentials_replace_key', ' ')
    // Wait for password field to be shown:
    await expect(this.p).toMatch('Password')
    await clearFillTextInput(this.p, 'aks_credentials_clientSecret', clientSecret)
  }

  async saveButtonDisabled() {
    return (await this.p.$('button#save[disabled]')) !== null
  }

  async save() {
    await this.p.click('button#save')
    await waitForDrawerOpenClose(this.p)
  }

  async delete(name) {
    await this.p.click(`a#akscreds_del_${name}`)
  }

  async confirmDelete() {
    await modalYes(this.p, 'Are you sure you want to delete the credentials')
  }
}
