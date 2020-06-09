const { BasePage } = require('./base')

export class SetupKorePage extends BasePage {
  constructor(p) {
    super(p)
    this.pagePath = '/setup/kore'
  }
}
