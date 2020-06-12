const { BasePage } = require('./base')

export class IndexPage extends BasePage {
  constructor(p) {
    super(p)
    this.pagePath = '/'
  }
}
