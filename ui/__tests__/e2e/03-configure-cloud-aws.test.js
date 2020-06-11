const { ConfigureCloudPage } = require('./page-objects/configure/cloud/configure-cloud')
const { ConfigureCloudAWSAccounts } = require('./page-objects/configure/cloud/AWS/accounts')

const page = global.page

// This test assumes it runs after a test which performs a login as local admin (e.g. 01-first-time-setup.test.js)
describe('Configure Cloud - AWS', () => {
  const cloudPage = new ConfigureCloudPage(page)

  beforeAll(async () => {
    await cloudPage.visitPage()
    cloudPage.verifyPageURL()
    await cloudPage.selectCloud('aws')
  })

  describe('Account Credentials', () => {
    const awsAccountsPage = new ConfigureCloudAWSAccounts(page)
    const testCred = {
      // Randomise name of test cred.
      name: `testproj-${Math.random().toString(36).substr(2, 5)}`,
      summary: 'a summary',
      accountID: '123456',
      accessKeyID: 'abcdef',
      secretAccessKey: 'pqrstuvwx'
    }

    beforeAll(async () => {
      await awsAccountsPage.openTab()
    })

    beforeEach(async () => {
      await awsAccountsPage.closeAllNotifications()
    })

    it('has the correct URL', () => {
      awsAccountsPage.verifyPageURL()
    })

    it('adds a new account credential', async () => {
      await awsAccountsPage.add()
      await expect(page).toMatch('New AWS account')
      await awsAccountsPage.populate(testCred)
      expect(await awsAccountsPage.saveButtonDisabled()).not.toBeTruthy()
      await awsAccountsPage.save()
    })

    it('shows account credentials', async () => {
      await awsAccountsPage.checkAccountListed(testCred.name)
    })

    it('edits a credential with a new description', async () => {
      await awsAccountsPage.edit(testCred.name, testCred.accountID)
      await awsAccountsPage.populate({ summary: 'summary2' })
      await awsAccountsPage.save()
      await expect(page).toMatch('AWS account credentials updated successfully')
    })

    it('edits a credential with a new key', async () => {
      await awsAccountsPage.edit(testCred.name, testCred.accountID)
      await awsAccountsPage.replaceKey('abcdefAB', 'pqrstuvwxAB')
      await awsAccountsPage.save()
      await expect(page).toMatch('AWS account credentials updated successfully')
    })

    // Re-instate this once deleting EKS credentials implemented:
    // it('allows credentials to be deleted', async () => {
    //   await awsAccountsPage.delete(testCred.name)
    //   await awsAccountsPage.confirmDelete()
    //   await expect(page).toMatch('AWS account credentials deleted successfully')
    // })
  })

})
