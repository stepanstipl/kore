const { BasePage } = require('./base')

export class IndexPage extends BasePage {
  constructor(p) {
    super(p)
    this.pagePath = '/'
  }

  async checkForMenuLink(text) {
    await expect(this.p).toMatchElement('#sider a', { text })
  }

  async clickMenuLink(text) {
    await expect(this.p).toClick('#sider a', { text })
  }
}
