const { ConfigureCloudPage } = require('./page-objects/configure/cloud/configure-cloud')
const { ConfigureCloudGCPProjects } = require('./page-objects/configure/cloud/GCP/projects')
const { ConfigureCloudGCPClusterPlans } = require('./page-objects/configure/cloud/GCP/cluster-plans')
const { waitForDrawerOpenClose } = require('./page-objects/utils')

const page = global.page

// This test assumes it runs after a test which performs a login as local admin (e.g. 01-first-time-setup.test.js)
describe('Configure Cloud - GCP', () => {
  const cloudPage = new ConfigureCloudPage(page)

  beforeAll(async () => {
    await cloudPage.visitPage()
    cloudPage.verifyPageURL()
    await cloudPage.selectCloud('gcp')
  })

  describe('Project Credentials', () => {
    const gkeProjCredsPage = new ConfigureCloudGCPProjects(page)
    const testCred = {
      // Randomise name of test cred.
      name: `testproj-${Math.random().toString(36).substr(2, 5)}`,
      summary: 'a summary',
      project: 'project001',
      json: 'horse'
    }

    beforeAll(async () => {
      await gkeProjCredsPage.openTab()
    })

    beforeEach(async () => {
      await gkeProjCredsPage.closeAllNotifications()
    })

    it('has the correct URL', () => {
      gkeProjCredsPage.verifyPageURL()
    })

    it('adds a new project credential', async () => {
      await gkeProjCredsPage.add()
      await gkeProjCredsPage.populate(testCred)
      await gkeProjCredsPage.save()
      await expect(page).toMatch('GCP project credentials created successfully')
    })

    it('shows project credentials', async () => {
      await gkeProjCredsPage.checkCredentialListed(testCred.name)
    })

    it('edits a project credential with a new description', async () => {
      await gkeProjCredsPage.edit(testCred.name, testCred.project)
      await gkeProjCredsPage.populate({ summary: 'summary2' })
      await gkeProjCredsPage.save()
      await expect(page).toMatch('GCP project credentials updated successfully')
    })

    it('edits a project credential with a new key', async () => {
      await gkeProjCredsPage.edit(testCred.name, testCred.project)
      await gkeProjCredsPage.replaceKey('sheep')
      await gkeProjCredsPage.save()
      await expect(page).toMatch('GCP project credentials updated successfully')
    })

    it('allows credentials to be deleted', async () => {
      await gkeProjCredsPage.delete(testCred.name)
      await gkeProjCredsPage.confirmDelete()
      await expect(page).toMatch('GCP project credentials deleted successfully')
    })

  })

  describe('Cluster Plans', () => {
    const gkeClusterPlansPage = new ConfigureCloudGCPClusterPlans(page)

    const testPlan = {
      name: `testplan-${Math.random().toString(36).substr(2, 5)}`,
      description: 'Test plan for testing',
      planDescription: 'A plan description',
      region: 'europe-west2'
    }

    beforeAll(async () => {
      // In case of gruff from previous tests, wait a beat before starting.
      await cloudPage.closeAllNotifications()
      await waitForDrawerOpenClose(page)
      await gkeClusterPlansPage.openTab()
      await gkeClusterPlansPage.listLoaded()
    })

    beforeEach(async () => {
      await gkeClusterPlansPage.closeAllNotifications()
    })

    it('has the correct URL', () => {
      gkeClusterPlansPage.verifyPageURL()
    })

    it('views an existing plan', async () => {
      await gkeClusterPlansPage.view('gke-development')
      // Check a random thing to ensure the plan is being displayed.
      await expect(page).toMatch('Authorized Master Networks')
      await gkeClusterPlansPage.closeDrawer()
    })

    it('does not allow editing of a read-only plan', async () => {
      await gkeClusterPlansPage.edit('gke-development')
      await expect(page).toMatch('This plan is read-only')
    })

    it('does not allow deleting of a read-only plan', async () => {
      await gkeClusterPlansPage.delete('gke-development')
      await expect(page).toMatch('This plan is read-only and cannot be deleted')
    })

    it('creates a new plan using default values', async () => {
      await gkeClusterPlansPage.new()
      await expect(page).toMatch('New GKE plan')
      await gkeClusterPlansPage.populatePlan(testPlan)
      await gkeClusterPlansPage.addNodePool()
      await gkeClusterPlansPage.populateNodePool({ name: 'compute' })
      await gkeClusterPlansPage.closeNodePool()
      await gkeClusterPlansPage.save()
      await expect(page).toMatch('GKE plan created successfully')
    })

    it('edits an existing plan', async () => {
      await gkeClusterPlansPage.edit(testPlan.name)
      await gkeClusterPlansPage.populatePlan({ region: 'europe-west1' })
      await gkeClusterPlansPage.addNodePool()
      // No node pool name, close should be disabled:
      expect(await gkeClusterPlansPage.closeNodePoolDisabled()).toBeTruthy()
      await gkeClusterPlansPage.populateNodePool({ name: 'compute' })
      // Node pool name clashes with existing name, close should be disabled:
      expect(await gkeClusterPlansPage.closeNodePoolDisabled()).toBeTruthy()
      await gkeClusterPlansPage.populateNodePool({ name: 'compute2' })
      // Button should now be enabled:
      expect(await gkeClusterPlansPage.closeNodePoolDisabled()).not.toBeTruthy()
      await gkeClusterPlansPage.closeNodePool()
      await gkeClusterPlansPage.save()
      await expect(page).toMatch('GKE plan updated successfully')
    })

    it('edits an existing node pool', async () => {
      await gkeClusterPlansPage.edit(testPlan.name)
      await gkeClusterPlansPage.viewEditNodePool(1)
      await gkeClusterPlansPage.populateNodePool({ enableAutoscaler: false, size: 10 })
      await gkeClusterPlansPage.closeNodePool()
      await gkeClusterPlansPage.save()
      await expect(page).toMatch('GKE plan updated successfully')
    })

    it('allows deleting of a non-read-only plan', async () => {
      await gkeClusterPlansPage.delete(testPlan.name)
      await gkeClusterPlansPage.confirmDelete()
      await expect(page).toMatch(`${testPlan.name} plan deleted`)
    })

  })

})
