const puppeteer = require('puppeteer')
const { LoginPage } = require('./page-objects/login')
const { IndexPage } = require('./page-objects/index')
const { SetupKorePage } = require('./page-objects/setup-kore')
const { SetupKoreCloudProvidersPage } = require('./page-objects/setup-kore-cloud-providers')
const { SetupKoreCompletePage } = require('./page-objects/setup-kore-complete')

const headless = process.env.SHOW_BROWSER !== 'true'

describe('First time login and setup', () => {
  let browser
  let page
  let loginPage
  let indexPage
  let setupKorePage
  let setupKoreCompletePage
  let setupKoreCloudProvidersPage

  beforeAll(async () => {
    console.log('creating browser')
    browser = await puppeteer.launch({
      args: ['--no-sandbox', '--start-maximized'],
      headless
    })
    page = await browser.newPage()
    loginPage = new LoginPage(page)
    indexPage = new IndexPage(page)
    setupKorePage = new SetupKorePage(page)
    setupKoreCompletePage = new SetupKoreCompletePage(page)
    setupKoreCloudProvidersPage = new SetupKoreCloudProvidersPage(page)
  })

  afterAll(() => {
    console.log('closing browser')
    browser.close()
  })

  test('admin user logs in and sets up a cloud provider', async () => {
    // admin user login
    await loginPage.visitPage()
    await loginPage.localUserLogin()

    // kore setup page
    setupKorePage.verifyPage()
    await setupKorePage.clickPrimaryButton()

    // cloud provider setup
    setupKoreCloudProvidersPage.verifyPage()
    await setupKoreCloudProvidersPage.selectCloud('gcp')
    await setupKoreCloudProvidersPage.enterGkeCredentials()
    await setupKoreCloudProvidersPage.clickPrimaryButton({ waitForNav: false })
    await setupKoreCloudProvidersPage.continueWithoutVerification()

    // kore setup complete
    setupKoreCompletePage.verifyPage()
    await setupKoreCompletePage.clickPrimaryButton()

    // main dashboard
    indexPage.verifyPage()
  }, 60000)
})
