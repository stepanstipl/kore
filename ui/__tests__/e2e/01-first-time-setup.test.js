const { LoginPage } = require('./page-objects/login')
const { IndexPage } = require('./page-objects/index')
const { SetupKorePage } = require('./page-objects/setup-kore')
const { SetupKoreCloudAccessPage } = require('./page-objects/setup-kore-cloud-access')
const { SetupKoreCompletePage } = require('./page-objects/setup-kore-complete')
const { ConfigureCloudPage } = require('./page-objects/configure/cloud/configure-cloud')
const { ConfigureCloudGCPOrgs } = require('./page-objects/configure/cloud/GCP/organizations')

const page = global.page

describe('First time login and Kore setup - GCP', () => {
  const loginPage = new LoginPage(page)
  const indexPage = new IndexPage(page)
  const setupKorePage = new SetupKorePage(page)
  const setupKoreCompletePage = new SetupKoreCompletePage(page)
  const setupKoreCloudAccessPage = new SetupKoreCloudAccessPage(page)

  const testOrg = {
    // Randomise name of test org.
    name: `testorg-${Math.random().toString(36).substr(2, 5)}`,
    summary: 'Org used by automated test',
    parentID: '1234567890',
    billingAccount: 'BILL-1234-ABCD',
    json: 'kangaroo'
  }

  describe('Kore setup', () => {
    it('admin user logs in and sets up GCP cloud access', async () => {
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
      await setupKoreCloudAccessPage.addGcpOrganization(testOrg)
      await expect(page).toMatch('GCP organization created successfully')
      await setupKoreCloudAccessPage.nextStep()
      await setupKoreCloudAccessPage.selectKoreManagedAccounts('custom')
      await setupKoreCloudAccessPage.setAutomatedAccountDefaults()
      await setupKoreCloudAccessPage.save()
      await setupKorePage.clickPrimaryButton()

      // kore setup complete
      setupKoreCompletePage.verifyPageURL()
      await setupKoreCompletePage.clickPrimaryButton()

      // main dashboard
      indexPage.verifyPageURL()
    })
  })

  describe('Verifying Kore setup', () => {
    const cloudPage = new ConfigureCloudPage(page)
    const orgsPage = new ConfigureCloudGCPOrgs(page)

    beforeAll(async () => {
      await cloudPage.visitPage()
      cloudPage.verifyPageURL()
      await cloudPage.selectCloud('gcp')
      await orgsPage.openTab()
    })

    it('has the correct URL', () => {
      orgsPage.verifyPageURL()
    })

    it('allows the configured org to be deleted', async () => {
      const orgConfigured = await orgsPage.orgConfigured()
      if (!orgConfigured) {
        return
      }
      await orgsPage.checkOrgListed(testOrg.name)
      await orgsPage.delete(testOrg.name)
      await orgsPage.confirmDelete()
      await expect(page).toMatch('GCP organization deleted successfully')
    })
  })
})
