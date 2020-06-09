const { BasePage } = require('./base')

export class SetupKoreCompletePage extends BasePage {
  constructor(p) {
    super(p)
    this.pagePath = '/setup/kore/complete'
  }
}
