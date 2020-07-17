const { ConfigureCloudPage } = require('./page-objects/configure/cloud/configure-cloud')
const { ConfigureCloudAzureSubscriptions } = require('./page-objects/configure/cloud/Azure/subscriptions')
const { ConfigureCloudAzureClusterPlans } = require('./page-objects/configure/cloud/Azure/cluster-plans')
const { ConfigureCloudAzureClusterPolicies } = require('./page-objects/configure/cloud/Azure/cluster-policies')
const { ConfigureCloudClusterPoliciesBase } = require('./page-objects/configure/cloud/cluster-policies-base')

const page = global.page

// This test assumes it runs after a test which performs a login as local admin (e.g. 01-first-time-setup.test.js)
describe('Configure Cloud - Azure', () => {
  const cloudPage = new ConfigureCloudPage(page)

  beforeAll(async () => {
    await cloudPage.visitPage()
    cloudPage.verifyPageURL()
    await cloudPage.selectCloud('azure')
  })

  describe('Subscription Credentials', () => {
    const azureSubscriptionsPage = new ConfigureCloudAzureSubscriptions(page)
    const testCred = {
      // Randomise name of test cred.
      name: `testsub-${Math.random().toString(36).substr(2, 5)}`,
      summary: 'AKS test summary',
      subscriptionID: '1234-abcd-5678',
      tenantID: 'my-aks-tenant-01',
      clientID: 'some-client-id',
      clientSecret: 'super-secret-password'
    }

    beforeAll(async () => {
      await azureSubscriptionsPage.openTab()
    })

    beforeEach(async () => {
      await azureSubscriptionsPage.closeAllNotifications()
    })

    it('has the correct URL', () => {
      azureSubscriptionsPage.verifyPageURL()
    })

    it('adds a new subscriptions credential', async () => {
      await azureSubscriptionsPage.add()
      await expect(page).toMatch('New Azure subscription')
      expect(await azureSubscriptionsPage.saveButtonDisabled()).toBeTruthy()
      await azureSubscriptionsPage.populate(testCred)
      expect(await azureSubscriptionsPage.saveButtonDisabled()).not.toBeTruthy()
      await azureSubscriptionsPage.save()
    })

    it('shows subscription credentials', async () => {
      await azureSubscriptionsPage.checkSubscriptionListed(testCred.name)
    })

    it('edits a credential with a new description', async () => {
      await azureSubscriptionsPage.edit(testCred.name, testCred.subscriptionID)
      await azureSubscriptionsPage.populate({ summary: 'summary2' })
      await azureSubscriptionsPage.save()
      await expect(page).toMatch('Azure subscription credentials updated successfully')
    })

    it('edits a credential with a new key', async () => {
      await azureSubscriptionsPage.edit(testCred.name, testCred.subscriptionID)
      await azureSubscriptionsPage.replacePassword('pqrstuvwxAB')
      await azureSubscriptionsPage.save()
      await expect(page).toMatch('Azure subscription credentials updated successfully')
    })

    it('allows credentials to be deleted', async () => {
      await azureSubscriptionsPage.delete(testCred.name)
      await azureSubscriptionsPage.confirmDelete()
      await expect(page).toMatch('Azure subscription credentials deleted successfully')
    })

  })

  describe('Cluster Plans', () => {
    const azureClusterPlansPage = new ConfigureCloudAzureClusterPlans(page)

    const testPlan = {
      name: `testplan-${Math.random().toString(36).substr(2, 5)}`,
      description: 'AKS plan for testing',
      planDescription: 'A plan description',
      region: 'uksouth',
      version: '1.16.10',
      dnsPrefix: 'example',
      networkPlugin: 'azure'
    }

    beforeAll(async () => {
      // In case of gruff from previous tests, wait a beat before starting.
      await cloudPage.closeAllNotifications()
      await azureClusterPlansPage.openTab()
      await azureClusterPlansPage.listLoaded()
    })

    beforeEach(async () => {
      await azureClusterPlansPage.closeAllNotifications()
    })

    it('has the correct URL', () => {
      azureClusterPlansPage.verifyPageURL()
    })

    it('views an existing plan', async () => {
      await azureClusterPlansPage.view('aks-development')
      // Check a random thing to ensure the plan is being displayed.
      await expect(page).toMatch('Default Team Role')
      await azureClusterPlansPage.closeDrawer()
    })

    it('does not allow editing of a read-only plan', async () => {
      await azureClusterPlansPage.edit('aks-development')
      await expect(page).toMatch('This plan is read-only')
    })

    it('does not allow deleting of a read-only plan', async () => {
      await azureClusterPlansPage.delete('aks-development')
      await expect(page).toMatch('This plan is read-only and cannot be deleted')
    })

    it('creates a new plan using default values', async () => {
      await azureClusterPlansPage.new()
      await expect(page).toMatch('New AKS plan')
      await azureClusterPlansPage.populatePlan(testPlan)
      await azureClusterPlansPage.addNodePool()
      await azureClusterPlansPage.populateNodePool({ name: 'default', mode: 'System', minSize: '1', size: '1', maxSize: '5' })
      await azureClusterPlansPage.closeNodePool()
      await azureClusterPlansPage.save()
      await expect(page).toMatch('AKS plan created successfully')
    })

    it('edits an existing plan', async () => {
      await azureClusterPlansPage.edit(testPlan.name)
      await azureClusterPlansPage.populatePlan({ dnsPrefix: 'example2' })
      await azureClusterPlansPage.addNodePool()
      // No name, close should be disabled:
      expect(await azureClusterPlansPage.closeNodePoolDisabled()).toBeTruthy()
      await azureClusterPlansPage.populateNodePool({ name: 'default', mode: 'System' })
      // Name clashes with existing name, close should be disabled:
      expect(await azureClusterPlansPage.closeNodePoolDisabled()).toBeTruthy()
      await azureClusterPlansPage.populateNodePool({ name: 'default2' })
      // Button should now be enabled:
      expect(await azureClusterPlansPage.closeNodePoolDisabled()).not.toBeTruthy()
      await azureClusterPlansPage.closeNodePool()
      await azureClusterPlansPage.save()
      await expect(page).toMatch('AKS plan updated successfully')
    })

    it('edits an existing node pool', async () => {
      await azureClusterPlansPage.edit(testPlan.name)
      await azureClusterPlansPage.viewEditNodePool(1)
      await azureClusterPlansPage.populateNodePool({ minSize: 2, size: 5, maxSize: 7 })
      await azureClusterPlansPage.closeNodePool()
      await azureClusterPlansPage.save()
      await expect(page).toMatch('AKS plan updated successfully')
    })

    it('allows deleting of a non-read-only plan', async () => {
      await azureClusterPlansPage.delete(testPlan.name)
      await azureClusterPlansPage.confirmDelete()
      await expect(page).toMatch(`${testPlan.name} plan deleted`)
    })
  })

  describe('Cluster Policies', () => {
    const policiesPage = new ConfigureCloudAzureClusterPolicies(page)
    const testPolicy = {
      name: `testpolicy-${Math.random().toString(36).substr(2, 5)}`,
      description: 'AKS policy for testing'
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

    it('has the correct URL', () => {
      policiesPage.verifyPageURL()
    })

    it('views an existing policy', async () => {
      await policiesPage.view('default-aks')
      // Check a random thing to ensure the plan is being displayed.
      await expect(page).toMatch('Authorized Master Networks')
      await policiesPage.closeDrawer()
    })

    it('does not allow editing of a read-only policy', async () => {
      await policiesPage.edit('default-aks')
      await expect(page).toMatch('This policy is read-only')
    })

    it('does not allow deleting of a read-only policy', async () => {
      await policiesPage.delete('default-aks')
      await expect(page).toMatch('This policy is read-only and cannot be deleted')
    })

    it('creates a new policy', async () => {
      await policiesPage.new()
      await expect(page).toMatch('New AKS policy')
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

    it('updates an existing policy', async () => {
      await policiesPage.edit(testPolicy.name)
      await policiesPage.populate({ description: 'Updated Policy Description' })
      await policiesPage.togglePolicyAllow('clusterUsers')
      await policiesPage.checkPolicyResult('clusterUsers', ConfigureCloudClusterPoliciesBase.RESULT_DEFAULT_DENY)

      await policiesPage.save()
      await expect(page).toMatch('Policy saved successfully')
    })

    it('views the updated policy', async () => {
      await policiesPage.view(testPolicy.name)
      await policiesPage.checkPolicyResult('domain', ConfigureCloudClusterPoliciesBase.RESULT_EXPLICIT_DENY)
      await policiesPage.checkPolicyResult('description', ConfigureCloudClusterPoliciesBase.RESULT_DEFAULT_DENY)
      await policiesPage.checkPolicyResult('authorizedMasterNetworks', ConfigureCloudClusterPoliciesBase.RESULT_EXPLICIT_DENY)
      await policiesPage.checkPolicyResult('clusterUsers', ConfigureCloudClusterPoliciesBase.RESULT_DEFAULT_DENY)
      await policiesPage.closeDrawer()
    })

    it('allows deleting of a non-read-only policy', async () => {
      await policiesPage.delete(testPolicy.name)
      await policiesPage.confirmDelete()
      await expect(page).toMatch('Policy Updated Policy Description deleted')
    })

  })

})
