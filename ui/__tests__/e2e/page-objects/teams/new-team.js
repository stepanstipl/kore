const { BasePage } = require('../base')
import { clearFillTextInput } from '../utils'

export class NewTeamPage extends BasePage {
  constructor(p) {
    super(p)
    this.pagePath = '/teams/new'
  }

  async populate({ name, description }) {
    await clearFillTextInput(this.p, 'new_team_teamName', name)
    await clearFillTextInput(this.p, 'new_team_teamDescription', description)
  }

  async checkTeamID(id) {
    await expect(this.p).toMatchElement('#team_id', { text: id })
  }

  async save() {
    await this.p.click('button#save')
    await expect(this.p).toMatchElement('#created_team')
  }

  async skipToTeamDashboard() {
    await Promise.all([
      this.p.waitForNavigation(),
      expect(this.p).toClick('button', { text: 'Skip to team dashboard' })
    ])
  }

}
