const { BasePage } = require('../base')
const asyncForEach = require('../../../../lib/utils/async-foreach')
import { modalYes, popConfirmYes } from '../utils'

export class TeamPage extends BasePage {
  constructor(p, teamID) {
    super(p)
    this.pagePath = `/teams/${teamID}`
  }

  async newCluster() {
    await Promise.all([
      this.p.waitForNavigation(),
      expect(this.p).toClick('button', { text: 'New cluster' })
    ])
  }

  async checkForCluster({ name, plan, namespaces }) {
    await this.p.waitFor(`#cluster_${name}`)
    // check cluster details
    await expect(this.p).toMatchElement(`#cluster_${name}`)
    await expect(this.p).toMatchElement(`#cluster_plan_${name}`, { text: plan })
    await this.checkForClusterStatus(name, 'Success')
    const namespaceScope = await this.p.$(`#cluster_namespaces_${name}`)
    await asyncForEach(namespaces, async (namespace) => {
      await expect(namespaceScope).toMatch(namespace)
    })
  }

  async checkForClusterStatus(name, status) {
    await expect(this.p).toMatchElement(`#cluster_status_${name}`, { text: status } )
  }

  async deleteCluster(name) {
    expect(this.p).toClick(`#cluster_delete_${name}`)
  }

  async deleteClusterConfirm() {
    await popConfirmYes(this.p, 'Are you sure you want to delete this cluster?')
  }

  async waitForClusterDeleted(name, timeoutSeconds) {
    await this.p.waitFor(`#cluster_${name}`, { hidden: true, timeout: timeoutSeconds * 1000 })
  }

  async deleteTeam() {
    await expect(this.p).toClick('#team_settings')
    await expect(this.p).toClick('#delete_team')
  }

  async deleteTeamConfirm() {
    await Promise.all([
      this.p.waitForNavigation(),
      modalYes(this.p, 'Are you sure you want to delete this team?')
    ])
  }

}
