const { BasePage } = require('./base')

export class SetupKorePage extends BasePage {
  constructor(page) {
    super(page)
    this.pagePath = '/setup/kore'
  }
}
