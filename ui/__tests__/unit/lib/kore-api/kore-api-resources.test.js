import KoreApiResources from '../../../../lib/kore-api/kore-api-resources'
import publicConfig from '../../../../config.public'

describe('KoreApiResources', () => {
  let koreApiResources

  beforeEach(() => {
    koreApiResources = new KoreApiResources()
  })

  describe('#generatePlanResource', () => {
    it('generates expected resource', () => {
      const values = {
        name: 'my-plan',
        summary: 'Plan summary',
        description: 'Plan description',
        configuration: {
          a: 1, b: 2
        }
      }
      const plan = koreApiResources.generatePlanResource('plan-kind', values)
      expect(plan).toBeDefined()
      expect(plan.apiVersion).toBe('config.kore.appvia.io/v1')
      expect(plan.kind).toBe('Plan')
      expect(plan.metadata.name).toBe(values.name)
      expect(plan.metadata.namespace).toBe(undefined)
      expect(plan.spec.kind).toBe('plan-kind')
      expect(plan.spec.summary).toBe(values.summary)
      expect(plan.spec.description).toBe(values.description)
      expect(plan.spec.configuration).toEqual(values.configuration)
    })
  })

  describe('#generatePolicyResource', () => {
    it('generates expected resource', () => {
      const values = {
        name: 'my-policy',
        summary: 'Policy summary',
        description: 'Policy description',
        properties: {
          a: 1, b: 2
        }
      }
      const policy = koreApiResources.generatePolicyResource('policy-kind', values)
      expect(policy).toBeDefined()
      expect(policy.apiVersion).toBe('config.kore.appvia.io/v1')
      expect(policy.kind).toBe('PlanPolicy')
      expect(policy.metadata.name).toBe(values.name)
      expect(policy.metadata.namespace).toBe(publicConfig.koreAdminTeamName)
      expect(policy.spec.kind).toBe('policy-kind')
      expect(policy.spec.summary).toBe(values.summary)
      expect(policy.spec.description).toBe(values.description)
      expect(policy.spec.properties).toEqual(values.properties)
    })
  })

  describe('#generateServicePlanResource', () => {
    it('generates expected resource', () => {
      const values = {
        name: 'My service plan',
        summary: 'Policy summary',
        description: 'Policy description',
        configuration: {
          a: 1, b: 2
        }
      }
      const resourceName = 'my-service-plan'
      const servicePlan = koreApiResources.generateServicePlanResource('service-kind', resourceName, values)
      expect(servicePlan).toBeDefined()
      expect(servicePlan.apiVersion).toBe('services.kore.appvia.io/v1')
      expect(servicePlan.kind).toBe('ServicePlan')
      expect(servicePlan.metadata.name).toBe(resourceName)
      expect(servicePlan.metadata.namespace).toBe('kore')
      expect(servicePlan.spec.kind).toBe('service-kind')
      expect(servicePlan.spec.summary).toBe(values.summary)
      expect(servicePlan.spec.description).toBe(values.description)
      expect(servicePlan.spec.configuration).toEqual(values.configuration)
    })
  })

  describe('#generateAccountManagementResource', () => {
    const gcpOrg = {
      apiVersion: 'gcp.compute.kore.appvia.io/v1alpha1',
      kind: 'Organization',
      metadata: {
        name: 'my-gcp',
        namespace: 'team1'
      }
    }
    const provider = 'GKE'

    it('generates expected resource, with GCP org only', () => {
      const accountMgt = koreApiResources.generateAccountManagementResource(provider, gcpOrg)
      expect(accountMgt).toBeDefined()
      expect(accountMgt.apiVersion).toBe('accounts.kore.appvia.io/v1beta1')
      expect(accountMgt.kind).toBe('AccountManagement')
      expect(accountMgt.metadata.name).toBe(`am-${gcpOrg.metadata.name}`)
      expect(accountMgt.metadata.namespace).toBe(publicConfig.koreAdminTeamName)
      expect(accountMgt.spec.provider).toBe('GKE')

      const [ group, version ] = gcpOrg.apiVersion.split('/')
      expect(accountMgt.spec.organization).toEqual({
        group,
        version,
        kind: gcpOrg.kind,
        name: gcpOrg.metadata.name,
        namespace: gcpOrg.metadata.namespace
      })
    })

    it('generates expected resource with project rules', () => {
      const gcpProjectList = [{
        name: 'Not prod',
        description: 'Not prod account',
        prefix: 'kore',
        suffix: 'notprod',
        plans: ['gke-development']
      }, {
        name: 'Prod',
        description: 'Prod account',
        prefix: 'kore',
        suffix: 'prod',
        plans: ['gke-production']
      }]
      const accountMgt = koreApiResources.generateAccountManagementResource(provider, gcpOrg, gcpProjectList)
      expect(accountMgt).toBeDefined()
      expect(accountMgt.spec.rules).toHaveLength(2)
      expect(accountMgt.spec.rules).toEqual(gcpProjectList)
    })

    it('sets resourceVersion if specified', () => {
      const accountMgt = koreApiResources.generateAccountManagementResource(provider, gcpOrg, [], '123')
      expect(accountMgt).toBeDefined()
      expect(accountMgt.metadata.resourceVersion).toBe('123')
    })

    it('sets provider', () => {
      const provider = 'EKS'
      const accountMgt = koreApiResources.generateAccountManagementResource(provider, gcpOrg)
      expect(accountMgt).toBeDefined()
      expect(accountMgt.spec.provider).toBe(provider)
    })
  })

  describe('#generateAllocationResource', () => {
    const resourceToAllocate = {
      apiVersion: 'test.kore.appvia.io/v1alpha1',
      kind: 'AllocatedKind',
      metadata: {
        name: 'my-resource',
        namespace: 'team1'
      }
    }

    it('generates expected resource', () => {
      const resourceName = 'allocatedkind-my-resource'
      const teams = ['*']
      const name = 'My resource'
      const summary = 'This is my resource'
      const allocation = koreApiResources.generateAllocationResource(resourceName, resourceToAllocate, teams, name, summary)
      expect(allocation).toBeDefined()
      expect(allocation.apiVersion).toBe('config.kore.appvia.io/v1')
      expect(allocation.kind).toBe('Allocation')
      expect(allocation.metadata.name).toBe(resourceName)
      expect(allocation.metadata.namespace).toBe(resourceToAllocate.metadata.namespace)
      expect(allocation.spec.name).toBe(name)
      expect(allocation.spec.summary).toBe(summary)
      expect(allocation.spec.teams).toBe(teams)
      const [ group, version ] = resourceToAllocate.apiVersion.split('/')
      expect(allocation.spec.resource).toEqual({
        group,
        version,
        kind: resourceToAllocate.kind,
        name: resourceToAllocate.metadata.name,
        namespace: resourceToAllocate.metadata.namespace
      })
    })
  })

  describe('#generateTeamResource', () => {
    it('generates expected resource', () => {
      const values = {
        teamName: 'Team 1',
        teamDescription: 'This is the first team'
      }
      const team = koreApiResources.generateTeamResource(values)
      expect(team).toBeDefined()
      expect(team.apiVersion).toBe('config.kore.appvia.io/v1')
      expect(team.kind).toBe('Team')
      expect(team.metadata.name).toBe('team-1')
      expect(team.metadata.namespace).toBe(undefined)
      expect(team.spec.summary).toBe(values.teamName)
      expect(team.spec.description).toBe(values.teamDescription)
    })
  })

  describe('#generateSecretResource', () => {
    const data = {
      accessKeyID: 'access-key',
      secretAccessKey: 'secret-key'
    }

    it('generates expected resource', () => {
      const secret = koreApiResources.generateSecretResource('team1', 'my-secret', 'secret-type', 'Super secret', data)
      expect(secret).toBeDefined()
      expect(secret.apiVersion).toBe('config.kore.appvia.io')
      expect(secret.kind).toBe('Secret')
      expect(secret.metadata.name).toBe('my-secret')
      expect(secret.metadata.namespace).toBe('team1')
      expect(secret.spec.type).toBe('secret-type')
      expect(secret.spec.description).toBe('Super secret')
      expect(secret.spec.data).toEqual(data)
    })

    it('works when using the team wrapper', () => {
      const secret = koreApiResources.team('team1').generateSecretResource('my-secret', 'secret-type', 'Super secret', data)
      expect(secret).toBeDefined()
    })
  })

  describe('#generateEKSCredentialsResource', () => {
    const values = {
      name: 'eks',
      accountID: '1234567890'
    }

    it('generates expected resource', () => {
      const eksCredential = koreApiResources.generateEKSCredentialsResource('team1', values, 'my-secret')
      expect(eksCredential).toBeDefined()
      expect(eksCredential.apiVersion).toBe('aws.compute.kore.appvia.io/v1alpha1')
      expect(eksCredential.kind).toBe('EKSCredentials')
      expect(eksCredential.metadata.name).toBe(values.name)
      expect(eksCredential.metadata.namespace).toBe('team1')
      expect(eksCredential.spec.accountID).toBe(values.accountID)
      expect(eksCredential.spec.credentialsRef).toEqual({ name: 'my-secret', namespace: 'team1' })
    })

    it('works when using the team wrapper', () => {
      const eksCredential = koreApiResources.team('team1').generateEKSCredentialsResource('team1', values, 'my-secret')
      expect(eksCredential).toBeDefined()
    })
  })

  describe('#generateGKECredentialsResource', () => {
    const values = { project: 'my-project' }

    it('generates expected resource', () => {
      const gkeCredential = koreApiResources.generateGKECredentialsResource('team1', values, 'my-secret')
      expect(gkeCredential).toBeDefined()
      expect(gkeCredential.apiVersion).toBe('gke.compute.kore.appvia.io/v1alpha1')
      expect(gkeCredential.kind).toBe('GKECredentials')
      expect(gkeCredential.metadata.name).toBe(values.name)
      expect(gkeCredential.metadata.namespace).toBe('team1')
      expect(gkeCredential.spec.accountID).toBe(values.accountID)
      expect(gkeCredential.spec.credentialsRef).toEqual({ name: 'my-secret', namespace: 'team1' })
    })

    it('works when using the team wrapper', () => {
      const gkeCredential = koreApiResources.team('team1').generateGKECredentialsResource('team1', values, 'my-secret')
      expect(gkeCredential).toBeDefined()
    })
  })

  describe('#generateGCPOrganizationResource', () => {
    const values = {
      name: 'gcp',
      parentID: '1234567890',
      billingAccount: 'ABC-124'
    }

    it('generates expected resource', () => {
      const gcpOrg = koreApiResources.generateGCPOrganizationResource('team1', values, 'my-secret')
      expect(gcpOrg).toBeDefined()
      expect(gcpOrg.apiVersion).toBe('gcp.compute.kore.appvia.io/v1alpha1')
      expect(gcpOrg.kind).toBe('Organization')
      expect(gcpOrg.metadata.name).toBe(values.name)
      expect(gcpOrg.metadata.namespace).toBe('team1')
      expect(gcpOrg.spec.parentID).toBe(values.parentID)
      expect(gcpOrg.spec.billingAccount).toBe(values.billingAccount)
      expect(gcpOrg.spec.credentialsRef).toEqual({ name: 'my-secret', namespace: 'team1' })
    })

    it('works when using the team wrapper', () => {
      const gcpOrg = koreApiResources.team('team1').generateGCPOrganizationResource('team1', values, 'my-secret')
      expect(gcpOrg).toBeDefined()
    })
  })

  describe('#generateClusterResource', () => {
    const user = 'user@appvia.io'
    const values = {
      clusterName: 'example-cluster'
    }
    const plan = {
      metadata: { name: 'dev-plan' },
      spec: { kind: 'dev' }
    }
    const planValues = {
      a: 1, b: 2
    }
    const credentials = 'some-creds'

    it('generates expected resource', () => {
      const cluster = koreApiResources.generateClusterResource('team1', user, values, plan, planValues, credentials)
      expect(cluster).toBeDefined()
      expect(cluster.apiVersion).toBe('clusters.compute.kore.appvia.io/v1')
      expect(cluster.kind).toBe('Cluster')
      expect(cluster.metadata.name).toBe(values.clusterName)
      expect(cluster.metadata.namespace).toBe('team1')
      expect(cluster.spec.kind).toBe(plan.spec.kind)
      expect(cluster.spec.plan).toBe(plan.metadata.name)
      expect(cluster.spec.configuration).toBe(planValues)
      expect(cluster.spec.credentials).toEqual(credentials)
    })

    it('works when using the team wrapper', () => {
      const cluster = koreApiResources.team('team1').generateClusterResource(user, values, plan, planValues, credentials)
      expect(cluster).toBeDefined()
    })
  })

  describe('#generateNamespaceClaimResource', () => {
    const cluster = {
      apiVersion: 'clusters.compute.kore.appvia.io/v1',
      kind: 'Cluster',
      metadata: { name: 'my-cluster', namespace: 'team1' },
      spec: { kind: 'dev' }
    }
    const resourceName = 'my-cluster-my-namespace'
    const values = {
      name: 'my-namespace'
    }

    it('generates expected resource', () => {
      const namespaceClaim = koreApiResources.generateNamespaceClaimResource('team1', cluster, resourceName, values)
      expect(namespaceClaim).toBeDefined()
      expect(namespaceClaim.apiVersion).toBe('namespaceclaims.clusters.compute.kore.appvia.io/v1alpha1')
      expect(namespaceClaim.kind).toBe('NamespaceClaim')
      expect(namespaceClaim.metadata.name).toBe(resourceName)
      expect(namespaceClaim.metadata.namespace).toBe('team1')
      expect(namespaceClaim.spec.name).toBe(values.name)

      const [ group, version ] = cluster.apiVersion.split('/')
      expect(namespaceClaim.spec.cluster).toEqual({
        group,
        version,
        kind: cluster.kind,
        name: cluster.metadata.name,
        namespace: cluster.metadata.namespace
      })

    })

    it('works when using the team wrapper', () => {
      const namespaceClaim = koreApiResources.team('team1').generateNamespaceClaimResource(cluster, resourceName, values)
      expect(namespaceClaim).toBeDefined()
    })
  })

  describe('#generateServiceResource', () => {
    const cluster = {
      apiVersion: 'clusters.compute.kore.appvia.io/v1',
      kind: 'Cluster',
      metadata: { name: 'my-cluster', namespace: 'team1' },
      spec: { kind: 'dev' }
    }
    const values = {
      serviceName: 'my-service',
      namespace: 'my-namespace'
    }
    const plan = {
      metadata: { name: 'dev-plan' },
      spec: { kind: 'dev' }
    }
    const planValues = {
      a: 1, b: 2
    }

    it('generates expected resource', () => {
      const service = koreApiResources.generateServiceResource('team1', cluster, values, plan, planValues)
      expect(service).toBeDefined()
      expect(service.apiVersion).toBe('services.compute.kore.appvia.io/v1')
      expect(service.kind).toBe('Service')
      expect(service.metadata.name).toBe(values.serviceName)
      expect(service.metadata.namespace).toBe('team1')
      expect(service.spec.kind).toBe(plan.spec.kind)
      expect(service.spec.plan).toBe(plan.metadata.name)
      expect(service.spec.configuration).toEqual(planValues)
      expect(service.spec.clusterNamespace).toBe(values.namespace)

      const [ group, version ] = cluster.apiVersion.split('/')
      expect(service.spec.cluster).toEqual({
        group,
        version,
        kind: cluster.kind,
        name: cluster.metadata.name,
        namespace: cluster.metadata.namespace
      })
    })

    it('sets cluster namespace as a new namespace, if specified', () => {
      const values2 = { ...values, createNamespace: 'new-namespace' }
      const service = koreApiResources.generateServiceResource('team1', cluster, values2, plan, planValues)
      expect(service).toBeDefined()
      expect(service.spec.clusterNamespace).toBe(values2.createNamespace)
    })

    it('works when using the team wrapper', () => {
      const service = koreApiResources.team('team1').generateServiceResource(cluster, values, plan, planValues)
      expect(service).toBeDefined()
    })
  })

  describe('#generateServiceCredentialsResource', () => {
    const name = 'test-service-access'
    const secretName = 'test-secret'
    const config = {
      a: 1, b: 2
    }
    const cluster = {
      apiVersion: 'clusters.compute.kore.appvia.io/v1',
      kind: 'Cluster',
      metadata: { name: 'my-cluster', namespace: 'team1' },
      spec: { kind: 'dev' }
    }
    const service = {
      apiVersion: 'services.compute.kore.appvia.io/v1',
      kind: 'Service',
      metadata: { name: 'my-service', namespace: 'team1' },
      spec: { kind: 'a-service' }
    }
    const clusterNamespace = 'my-namespace'

    it('generates expected resource', () => {
      const serviceCred = koreApiResources.generateServiceCredentialsResource('team1', name, secretName, config, service, cluster, clusterNamespace)
      expect(serviceCred).toBeDefined()
      expect(serviceCred.apiVersion).toBe('servicecredentials.services.kore.appvia.io/v1')
      expect(serviceCred.kind).toBe('ServiceCredentials')
      expect(serviceCred.metadata.name).toBe(name)
      expect(serviceCred.metadata.namespace).toBe('team1')
      expect(serviceCred.spec.kind).toBe(service.spec.kind)
      expect(serviceCred.spec.secretName).toBe(secretName)
      expect(serviceCred.spec.configuration).toBe(config)

      const [ sGroup, sVersion ] = service.apiVersion.split('/')
      expect(serviceCred.spec.service).toEqual({
        group: sGroup,
        version: sVersion,
        kind: service.kind,
        name: service.metadata.name,
        namespace: service.metadata.namespace
      })

      const [ cGroup, cVersion ] = cluster.apiVersion.split('/')
      expect(serviceCred.spec.cluster).toEqual({
        group: cGroup,
        version: cVersion,
        kind: cluster.kind,
        name: cluster.metadata.name,
        namespace: cluster.metadata.namespace
      })

    })

    it('works when using the team wrapper', () => {
      const serviceCred = koreApiResources.team('team1').generateServiceCredentialsResource(name, secretName, config, service, cluster, clusterNamespace)
      expect(serviceCred).toBeDefined()
    })
  })

})
