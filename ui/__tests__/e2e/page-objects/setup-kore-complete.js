const { BasePage } = require('./base')

export class SetupKoreCompletePage extends BasePage {
  constructor(page) {
    super(page)
    this.pagePath = '/setup/kore/complete'
  }
}
