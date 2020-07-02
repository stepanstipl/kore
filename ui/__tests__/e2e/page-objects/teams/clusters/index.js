const { BasePage } = require('../../base')
import { clearFillTextInput, popConfirmYes, waitForDrawerOpenClose } from '../../utils'

export class ClusterPage extends BasePage {
  constructor(p, teamID, clusterName) {
    super(p)
    this.teamID = teamID
    this.clusterName = clusterName
    this.pagePath = `/teams/${teamID}/clusters/${clusterName}`
  }

  async checkClusterName() {
    await expect(this.p).toMatchElement('#cluster_name', { text: this.clusterName, timeout: 5000 })
  }

  async checkForNamespace(name, status) {
    await this.p.waitFor('.namespace')
    await expect(this.p).toMatchElement(`#namespace_${name}`)
    if (status) {
      await expect(this.p).toMatchElement(`#namespace_status_${name}`, { text: status })
    }
  }

  async checkForClusterStatus(status) {
    await expect(this.p).toMatchElement('#cluster_status', { text: status })
  }

  async waitForClusterStatus(statusList, timeoutSeconds) {
    const timeout = timeoutSeconds * 1000
    await Promise.race(statusList.map(s => expect(this.p).toMatchElement('#cluster_status', { text: s, timeout })))
  }

  async waitForNamespaceStatus(name, statusList, timeoutSeconds) {
    const timeout = timeoutSeconds * 1000
    await Promise.race(statusList.map(s => expect(this.p).toMatchElement(`#namespace_status_${name}`, { text: s, timeout })))
  }

  async addNamespace() {
    await expect(this.p).toClick('button', { text: 'New namespace' })
    await waitForDrawerOpenClose(this.p)
    await expect(this.p).toMatchElement('.ant-drawer-title', { text: 'New namespace' })
  }

  async populateNamespace(name) {
    await clearFillTextInput(this.p, 'namespace_claim_name', name)
  }

  async saveNamespace() {
    await this.p.click('button#save')
    await waitForDrawerOpenClose(this.p)
  }

  async deleteNamespace(name) {
    await expect(this.p).toClick(`#namespace_delete_${name}`)
  }

  async deleteNamespaceConfirm() {
    await popConfirmYes(this.p, 'Are you sure you want to delete this namespace?')
  }

  async waitForNamespaceDeleted(namespaceName, timeoutSeconds) {
    await this.p.waitFor(`#namespace_${namespaceName}`, { hidden: true, timeout: timeoutSeconds * 1000 })
  }

}
