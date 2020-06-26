const { ConfigureCloudPage } = require('./page-objects/configure/cloud/configure-cloud')
const { ConfigureCloudGCPOrgs } = require('./page-objects/configure/cloud/GCP/organizations')
const { ConfigureCloudGCPProjects } = require('./page-objects/configure/cloud/GCP/projects')
const { ConfigureCloudGCPClusterPlans } = require('./page-objects/configure/cloud/GCP/cluster-plans')
const { ConfigureCloudGCPClusterPolicies } = require('./page-objects/configure/cloud/GCP/cluster-policies')
const { ConfigureCloudClusterPoliciesBase } = require('./page-objects/configure/cloud/cluster-policies-base')

const page = global.page

// This test assumes it runs after a test which performs a login as local admin (e.g. 01-first-time-setup.test.js)
describe('Configure Cloud - GCP', () => {
  const cloudPage = new ConfigureCloudPage(page)

  beforeAll(async () => {
    await cloudPage.visitPage()
    cloudPage.verifyPageURL()
    await cloudPage.selectCloud('gcp')
  })

  describe('Organizations', () => {
    const orgsPage = new ConfigureCloudGCPOrgs(page)
    const testOrg = {
      // Randomise name of test org.
      name: `testorg-${Math.random().toString(36).substr(2, 5)}`,
      summary: 'a summary',
      parentID: '1234567890',
      billingAccount: 'BILL-1234-ABCD',
      json: 'crocodile'
    }

    beforeAll(async () => {
      await orgsPage.openTab()
    })

    beforeEach(async () => {
      await orgsPage.closeAllNotifications()
    })

    test('has the correct URL', () => {
      orgsPage.verifyPageURL()
    })

    test('adds a new organization', async () => {
      await orgsPage.add()
      await orgsPage.populate(testOrg)
      await orgsPage.save()
      await expect(page).toMatch('GCP organization created successfully')
    })

    test('shows the organization', async () => {
      await orgsPage.checkOrgListed(testOrg.name)
    })

    test('edits the organization with a new description', async () => {
      await orgsPage.edit(testOrg.name, testOrg.parentID)
      await orgsPage.populate({ summary: 'summary2' })
      await orgsPage.save()
      await expect(page).toMatch('GCP organization updated successfully')
    })

    test('edits a project credential with a new key', async () => {
      await orgsPage.edit(testOrg.name, testOrg.parentID)
      await orgsPage.replaceKey('chicken')
      await orgsPage.save()
      await expect(page).toMatch('GCP organization updated successfully')
    })

    test('allows the organization to be deleted', async () => {
      await orgsPage.delete(testOrg.name)
      await orgsPage.confirmDelete()
      await expect(page).toMatch('GCP organization deleted successfully')
    })

  })

  describe('Project Credentials', () => {
    const projCredsPage = new ConfigureCloudGCPProjects(page)
    const testCred = {
      // Randomise name of test cred.
      name: `testproj-${Math.random().toString(36).substr(2, 5)}`,
      summary: 'a summary',
      project: 'project001',
      json: 'horse'
    }

    beforeAll(async () => {
      await projCredsPage.openTab()
    })

    beforeEach(async () => {
      await projCredsPage.closeAllNotifications()
    })

    test('has the correct URL', () => {
      projCredsPage.verifyPageURL()
    })

    test('adds a new project credential', async () => {
      await projCredsPage.add()
      await projCredsPage.populate(testCred)
      await projCredsPage.save()
      await expect(page).toMatch('GCP project credentials created successfully')
    })

    test('shows project credentials', async () => {
      await projCredsPage.checkCredentialListed(testCred.name)
    })

    test('edits a project credential with a new description', async () => {
      await projCredsPage.edit(testCred.name, testCred.project)
      await projCredsPage.populate({ summary: 'summary2' })
      await projCredsPage.save()
      await expect(page).toMatch('GCP project credentials updated successfully')
    })

    test('edits a project credential with a new key', async () => {
      await projCredsPage.edit(testCred.name, testCred.project)
      await projCredsPage.replaceKey('sheep')
      await projCredsPage.save()
      await expect(page).toMatch('GCP project credentials updated successfully')
    })

    test('allows credentials to be deleted', async () => {
      await projCredsPage.delete(testCred.name)
      await projCredsPage.confirmDelete()
      await expect(page).toMatch('GCP project credentials deleted successfully')
    })

  })

  describe('Cluster Plans', () => {
    const clusterPlansPage = new ConfigureCloudGCPClusterPlans(page)

    const testPlan = {
      name: `testplan-${Math.random().toString(36).substr(2, 5)}`,
      description: 'Test plan for testing',
      planDescription: 'A plan description',
      region: 'europe-west2'
    }

    beforeAll(async () => {
      // In case of gruff from previous tests, wait a beat before starting.
      await cloudPage.closeAllNotifications()
      await clusterPlansPage.openTab()
      await clusterPlansPage.listLoaded()
    })

    beforeEach(async () => {
      await clusterPlansPage.closeAllNotifications()
    })

    test('has the correct URL', () => {
      clusterPlansPage.verifyPageURL()
    })

    test('views an existing plan', async () => {
      await clusterPlansPage.view('gke-development')
      // Check a random thing to ensure the plan is being displayed.
      await expect(page).toMatch('Authorized Master Networks')
      await clusterPlansPage.closeDrawer()
    })

    test('does not allow editing of a read-only plan', async () => {
      await clusterPlansPage.edit('gke-development')
      await expect(page).toMatch('This plan is read-only')
    })

    test('does not allow deleting of a read-only plan', async () => {
      await clusterPlansPage.delete('gke-development')
      await expect(page).toMatch('This plan is read-only and cannot be deleted')
    })

    test('creates a new plan using default values', async () => {
      await clusterPlansPage.new()
      await expect(page).toMatch('New GKE plan')
      await clusterPlansPage.populatePlan(testPlan)
      await clusterPlansPage.addNodePool()
      await clusterPlansPage.populateNodePool({ name: 'compute' })
      await clusterPlansPage.closeNodePool()
      await clusterPlansPage.save()
      await expect(page).toMatch('GKE plan created successfully')
    })

    test('edits an existing plan', async () => {
      await clusterPlansPage.edit(testPlan.name)
      await clusterPlansPage.populatePlan({ region: 'europe-west1' })
      await clusterPlansPage.addNodePool()
      // No node pool name, close should be disabled:
      expect(await clusterPlansPage.closeNodePoolDisabled()).toBeTruthy()
      await clusterPlansPage.populateNodePool({ name: 'compute' })
      // Node pool name clashes with existing name, close should be disabled:
      expect(await clusterPlansPage.closeNodePoolDisabled()).toBeTruthy()
      await clusterPlansPage.populateNodePool({ name: 'compute2' })
      // Button should now be enabled:
      expect(await clusterPlansPage.closeNodePoolDisabled()).not.toBeTruthy()
      await clusterPlansPage.closeNodePool()
      await clusterPlansPage.save()
      await expect(page).toMatch('GKE plan updated successfully')
    })

    test('edits an existing node pool', async () => {
      await clusterPlansPage.edit(testPlan.name)
      await clusterPlansPage.viewEditNodePool(1)
      await clusterPlansPage.populateNodePool({ enableAutoscaler: false, size: 10 })
      await clusterPlansPage.closeNodePool()
      await clusterPlansPage.save()
      await expect(page).toMatch('GKE plan updated successfully')
    })

    test('allows deleting of a non-read-only plan', async () => {
      await clusterPlansPage.delete(testPlan.name)
      await clusterPlansPage.confirmDelete()
      await expect(page).toMatch(`${testPlan.name} plan deleted`)
    })

  })

  describe('Cluster Policies', () => {
    const policiesPage = new ConfigureCloudGCPClusterPolicies(page)
    const testPolicy = {
      name: `testpolicy-${Math.random().toString(36).substr(2, 5)}`,
      description: 'Test policy for testing'
    }

    beforeAll(async () => {
      // In case of gruff from previous tests, wait a beat before starting.
      await cloudPage.closeAllNotifications()
      await policiesPage.openTab()
      await policiesPage.listLoaded()
    })

    beforeEach(async () => {
      await policiesPage.closeAllNotifications()
    })

    test('has the correct URL', () => {
      policiesPage.verifyPageURL()
    })

    test('views an existing policy', async () => {
      await policiesPage.view('default-gke')
      // Check a random thing to ensure the plan is being displayed.
      await expect(page).toMatch('Authorized Master Networks')
      await policiesPage.closeDrawer()
    })

    test('does not allow editing of a read-only policy', async () => {
      await policiesPage.edit('default-gke')
      await expect(page).toMatch('This policy is read-only')
    })

    test('does not allow deleting of a read-only policy', async () => {
      await policiesPage.delete('default-gke')
      await expect(page).toMatch('This policy is read-only and cannot be deleted')
    })

    test('creates a new policy', async () => {
      await policiesPage.new()
      await expect(page).toMatch('New GKE policy')
      await policiesPage.populate(testPolicy)
      await policiesPage.togglePolicyAllow('clusterUsers')
      await policiesPage.togglePolicyDisallow('domain')
      await policiesPage.checkPolicyResult('clusterUsers', ConfigureCloudClusterPoliciesBase.RESULT_EXPLICIT_ALLOW)
      await policiesPage.checkPolicyResult('domain', ConfigureCloudClusterPoliciesBase.RESULT_EXPLICIT_DENY)
      await policiesPage.checkPolicyResult('description', ConfigureCloudClusterPoliciesBase.RESULT_DEFAULT_DENY)

      // Both deny + allow = deny:
      await policiesPage.togglePolicyAllow('authorizedMasterNetworks')
      await policiesPage.togglePolicyDisallow('authorizedMasterNetworks')
      await policiesPage.checkPolicyResult('authorizedMasterNetworks', ConfigureCloudClusterPoliciesBase.RESULT_EXPLICIT_DENY)

      await policiesPage.save()
      await expect(page).toMatch('Policy created successfully')
    })

    test('updates an existing policy', async () => {
      await policiesPage.edit(testPolicy.name)
      await policiesPage.populate({ description: 'Updated Policy Description' })
      await policiesPage.togglePolicyAllow('clusterUsers')
      await policiesPage.checkPolicyResult('clusterUsers', ConfigureCloudClusterPoliciesBase.RESULT_DEFAULT_DENY)

      await policiesPage.save()
      await expect(page).toMatch('Policy saved successfully')
    })

    test('views the updated policy', async () => {
      await policiesPage.view(testPolicy.name)
      await policiesPage.checkPolicyResult('domain', ConfigureCloudClusterPoliciesBase.RESULT_EXPLICIT_DENY)
      await policiesPage.checkPolicyResult('description', ConfigureCloudClusterPoliciesBase.RESULT_DEFAULT_DENY)
      await policiesPage.checkPolicyResult('authorizedMasterNetworks', ConfigureCloudClusterPoliciesBase.RESULT_EXPLICIT_DENY)
      await policiesPage.checkPolicyResult('clusterUsers', ConfigureCloudClusterPoliciesBase.RESULT_DEFAULT_DENY)
      await policiesPage.closeDrawer()
    })

    test('allows deleting of a non-read-only policy', async () => {
      await policiesPage.delete(testPolicy.name)
      await policiesPage.confirmDelete()
      await expect(page).toMatch('Policy Updated Policy Description deleted')
    })

  })

})
