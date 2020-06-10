const { LoginPage } = require('./page-objects/login')
const { IndexPage } = require('./page-objects/index')
const { SetupKorePage } = require('./page-objects/setup-kore')
const { SetupKoreCloudAccessPage } = require('./page-objects/setup-kore-cloud-access')
const { SetupKoreCompletePage } = require('./page-objects/setup-kore-complete')

const page = global.page

describe('First time login and setup', () => {
  const loginPage = new LoginPage(page)
  const indexPage = new IndexPage(page)
  const setupKorePage = new SetupKorePage(page)
  const setupKoreCompletePage = new SetupKoreCompletePage(page)
  const setupKoreCloudAccessPage = new SetupKoreCloudAccessPage(page)

  beforeAll(async () => {
  })

  afterAll(() => {
  })

  test('admin user logs in and sets up GCP cloud access', async () => {
    // admin user login
    await loginPage.visitPage()
    await loginPage.localUserLogin()

    // If this is re-run, we don't go back to the setup. Check the URL and exit clean if
    // not auto-redirected to the setup page.
    // @TODO: Reset the environment as part of every test so this is re-runnable.
    const currUrl = await page.url()
    if (!currUrl.endsWith(setupKorePage.pagePath)) {
      return
    }

    // kore setup page
    setupKorePage.verifyPageURL()
    await setupKorePage.clickPrimaryButton()

    // cloud access setup
    setupKoreCloudAccessPage.verifyPageURL()
    await setupKoreCloudAccessPage.selectCloud('gcp')
    await setupKoreCloudAccessPage.selectKoreManagedGcpProjects()
    await setupKoreCloudAccessPage.addGcpOrganization()
    await setupKoreCloudAccessPage.nextStep()
    await setupKoreCloudAccessPage.selectKoreManagedProjects('custom')
    await setupKoreCloudAccessPage.setAutomatedProjectDefaults()
    await setupKoreCloudAccessPage.save()
    await setupKorePage.clickPrimaryButton()

    // kore setup complete
    setupKoreCompletePage.verifyPageURL()
    await setupKoreCompletePage.clickPrimaryButton()

    // main dashboard
    indexPage.verifyPageURL()
  })
})