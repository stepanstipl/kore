const { ConfigureCloudPage } = require('./page-objects/configure/cloud/configure-cloud')
const { ConfigureCloudAWSAccounts } = require('./page-objects/configure/cloud/AWS/accounts')
const { ConfigureCloudAWSClusterPlans } = require('./page-objects/configure/cloud/AWS/cluster-plans')
const { waitForDrawerOpenClose } = require('./page-objects/utils')

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

    it('allows credentials to be deleted', async () => {
      await awsAccountsPage.delete(testCred.name)
      await awsAccountsPage.confirmDelete()
      await expect(page).toMatch('AWS account credentials deleted successfully')
    })

  })

  describe('Cluster Plans', () => {
    const awsClusterPlansPage = new ConfigureCloudAWSClusterPlans(page)

    const testPlan = {
      name: `testplan-${Math.random().toString(36).substr(2, 5)}`,
      description: 'Test plan for testing',
      planDescription: 'A plan description',
      region: 'eu-west-2',
      version: '1.15'
    }

    beforeAll(async () => {
      // In case of gruff from previous tests, wait a beat before starting.
      await cloudPage.closeAllNotifications()
      await waitForDrawerOpenClose(page)
      await awsClusterPlansPage.openTab()
      await awsClusterPlansPage.listLoaded()
    })

    beforeEach(async () => {
      await awsClusterPlansPage.closeAllNotifications()
    })

    it('has the correct URL', () => {
      awsClusterPlansPage.verifyPageURL()
    })

    it('views an existing plan', async () => {
      await awsClusterPlansPage.view('eks-development')
      // Check a random thing to ensure the plan is being displayed.
      await expect(page).toMatch('Default Team Role')
      await awsClusterPlansPage.closeDrawer()
    })

    it('does not allow editing of a read-only plan', async () => {
      await awsClusterPlansPage.edit('eks-development')
      await expect(page).toMatch('This plan is read-only')
    })

    it('does not allow deleting of a read-only plan', async () => {
      await awsClusterPlansPage.delete('eks-development')
      await expect(page).toMatch('This plan is read-only and cannot be deleted')
    })

    it('creates a new plan using default values', async () => {
      await awsClusterPlansPage.new()
      await expect(page).toMatch('New EKS plan')
      await awsClusterPlansPage.populatePlan(testPlan)
      await awsClusterPlansPage.addNodeGroup()
      await awsClusterPlansPage.populateNodeGroup({ name: 'default' })
      await awsClusterPlansPage.closeNodeGroup()
      await awsClusterPlansPage.save()
      await expect(page).toMatch('EKS plan created successfully')
    })

    it('edits an existing plan', async () => {
      await awsClusterPlansPage.edit(testPlan.name)
      await awsClusterPlansPage.populatePlan({ region: 'eu-west-1' })
      await awsClusterPlansPage.addNodeGroup()
      // No name, close should be disabled:
      expect(await awsClusterPlansPage.closeNodeGroupDisabled()).toBeTruthy()
      await awsClusterPlansPage.populateNodeGroup({ name: 'default' })
      // Name clashes with existing name, close should be disabled:
      expect(await awsClusterPlansPage.closeNodeGroupDisabled()).toBeTruthy()
      await awsClusterPlansPage.populateNodeGroup({ name: 'default2' })
      // Button should now be enabled:
      expect(await awsClusterPlansPage.closeNodeGroupDisabled()).not.toBeTruthy()
      await awsClusterPlansPage.closeNodeGroup()
      await awsClusterPlansPage.save()
      await expect(page).toMatch('EKS plan updated successfully')
    })

    it('edits an existing node group', async () => {
      await awsClusterPlansPage.edit(testPlan.name)
      await awsClusterPlansPage.viewEditNodeGroup(1)
      await awsClusterPlansPage.populateNodeGroup({ minSize: 2, desiredSize: 5, maxSize: 7 })
      await awsClusterPlansPage.closeNodeGroup()
      await awsClusterPlansPage.save()
      await expect(page).toMatch('EKS plan updated successfully')
    })

    it('allows deleting of a non-read-only plan', async () => {
      await awsClusterPlansPage.delete(testPlan.name)
      await awsClusterPlansPage.confirmDelete()
      await expect(page).toMatch(`${testPlan.name} plan deleted`)
    })
  })
})
