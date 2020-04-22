const { BasePage } = require('./base')

export class LoginPage extends BasePage {

  constructor(page) {
    super(page)
    this.pagePath = '/login'
  }

  async localUserLogin() {
    await this.page.click('#local-user-login')
    await this.page.waitFor('#login-form')
    await this.page.type('#login_login', 'admin')
    await this.page.type('#login_password', 'password')
    const [response] = await Promise.all([
      this.page.waitForNavigation(),
      this.page.click('#submit'),
    ])
    expect(response.ok()).toBeTruthy()
  }

}
