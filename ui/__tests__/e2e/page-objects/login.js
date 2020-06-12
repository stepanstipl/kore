const { BasePage } = require('./base')
const config = require('../config')

export class LoginPage extends BasePage {
  constructor(p) {
    super(p)
    this.pagePath = '/login'
  }

  async localUserLogin() {
    await this.p.click('#local-user-login')
    await this.p.waitFor('#login-form')
    await this.p.type('#login_login', 'admin')
    await this.p.type('#login_password', config.adminPass)
    const [response] = await Promise.all([
      this.p.waitForNavigation(),
      this.p.click('#submit'),
    ])
    expect(response.ok()).toBeTruthy()
  }

}
