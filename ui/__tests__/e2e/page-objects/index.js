const { BasePage } = require('./base')

export class IndexPage extends BasePage {
  constructor(page) {
    super(page)
    this.pagePath = '/'
  }
}
