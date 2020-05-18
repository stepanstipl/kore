const puppeteer = require('puppeteer')
const { LoginPage } = require('./page-objects/login')
const { IndexPage } = require('./page-objects/index')
const { SetupKorePage } = require('./page-objects/setup-kore')
const { SetupKoreCloudAccessPage } = require('./page-objects/setup-kore-cloud-access')
const { SetupKoreCompletePage } = require('./page-objects/setup-kore-complete')

const headless = process.env.SHOW_BROWSER !== 'true'

describe('First time login and setup', () => {
  let browser
  let page
  let loginPage
  let indexPage
  let setupKorePage
  let setupKoreCompletePage
  let setupKoreCloudAccessPage

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
    setupKoreCloudAccessPage = new SetupKoreCloudAccessPage(page)
  })

  afterAll(() => {
    console.log('closing browser')
    browser.close()
  })

  test('admin user logs in and sets up GCP cloud access', async () => {
    // admin user login
    await loginPage.visitPage()
    await loginPage.localUserLogin()

    // kore setup page
    setupKorePage.verifyPage()
    await setupKorePage.clickPrimaryButton()

    // cloud access setup
    setupKoreCloudAccessPage.verifyPage()
    await setupKoreCloudAccessPage.selectCloud('gcp')
    await setupKoreCloudAccessPage.selectKoreManagedGcpProjects()
    await setupKoreCloudAccessPage.addGcpOrganization()
    await setupKoreCloudAccessPage.nextStep()
    await setupKoreCloudAccessPage.selectKoreManagedProjects('custom')
    await setupKoreCloudAccessPage.setAutomatedProjectDefaults()
    await setupKoreCloudAccessPage.save()
    await setupKorePage.clickPrimaryButton()

    // kore setup complete
    setupKoreCompletePage.verifyPage()
    await setupKoreCompletePage.clickPrimaryButton()

    // main dashboard
    indexPage.verifyPage()
  }, 60000)
})
