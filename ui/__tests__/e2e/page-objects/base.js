const { testUrl } = require('../config')

export class BasePage {

  constructor(page) {
    this.page = page
  }

  getFileName() {
    if (this.pagePath === '/') {
      return 'index'
    }
    return this.pagePath.replace('/', '').replace(/\//g,'-')
  }

  async screenshot(filename) {
    await this.page.screenshot({
      fullPage: true,
      path: `__tests__/e2e/screenshots/${filename}.png`
    })
  }

  verifyPage() {
    expect(this.page.url()).toBe(`${testUrl}${this.pagePath}`)
  }

  async visitPage(query = '') {
    try {
      console.log(`goto [${testUrl}${this.pagePath}${query}]`)
      await this.page.goto(`${testUrl}${this.pagePath}${query}`)
      await this.page.waitForSelector('body')
    } catch (error) {
      const filename = this.getFileName()
      console.log(`Exception caught in visitPage, taking screenshot ${filename}.png. Error is: ${error}`)
      await this.screenshot(`failed-visit-${filename}`)
    }
  }

  async getHeading() {
    try {
      return await this.page.$eval('h1', el => el.innerHTML)
    } catch (error) {
      const filename = `failed-getHeading-${this.getFileName()}`
      console.log(`Exception caught in getHeading on ${this.page.url()}, taking screenshot ${filename}.png. Error is: ${error}`)
      await this.screenshot(filename)
    }
  }

  async clickPrimaryButton(options) {
    options = options || { waitForNav: true }
    if (options.waitForNav) {
      await Promise.all([
        this.page.waitForNavigation(),
        this.page.click('.ant-btn-primary')
      ])
    } else {
      await this.page.click('.ant-btn-primary')
    }
  }

}
