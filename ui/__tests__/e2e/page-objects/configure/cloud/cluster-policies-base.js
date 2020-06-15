import { ConfigureCloudPage } from './configure-cloud'
import { modalYes, waitForDrawerOpenClose, clearFillTextInput } from '../../utils'

export class ConfigureCloudClusterPoliciesBase extends ConfigureCloudPage {
  async listLoaded() {
    await this.p.waitForSelector('#policy_list', { timeout: 1000 })
  }

  async view(name) {
    await this.p.click(`a#policy_view_${name}`)
    await waitForDrawerOpenClose(this.p)
  }

  async edit(name) {
    await this.p.click(`a#policy_edit_${name}`)
    await waitForDrawerOpenClose(this.p)
  }

  async delete(name) {
    await this.p.click(`a#policy_delete_${name}`)
  }

  async confirmDelete() {
    await modalYes(this.p, 'Are you sure you want to delete the policy')
  }

  async new() {
    await this.p.click('button#add')
    await waitForDrawerOpenClose(this.p)
  }

  async populate({ description, name }) {
    await clearFillTextInput(this.p, 'policy_summary', name)
    await clearFillTextInput(this.p, 'policy_description', description)
  }

  async togglePolicyAllow(name) {
    await this.p.click(`#policy_${name}_allow`)
  }

  async togglePolicyDisallow(name) {
    await this.p.click(`#policy_${name}_disallow`)
  }

  static RESULT_EXPLICIT_DENY = 'Changes explicitly denied by this policy'
  static RESULT_EXPLICIT_ALLOW = 'Changes explicitly allowed by this policy'
  static RESULT_DEFAULT_DENY = 'Changes will be denied by default'
  static RESULT_DEFAULT_ALLOW = 'Changes will be allowed by default'

  async checkPolicyResult(name, expected) {
    await this.p.hover(`#policy_${name}_result`)
    const tt = await this.p.waitForSelector(`#policy_${name}_result_tt`)
    await expect(tt).toMatch(expected)
  }

  async save() {
    await this.p.click('button#policy_save')
    await waitForDrawerOpenClose(this.p)
  }
}
