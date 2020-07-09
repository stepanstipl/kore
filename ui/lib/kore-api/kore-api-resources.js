import * as models from './model/*'
import canonical from '../utils/canonical'
import getConfig from 'next/config'
const { publicRuntimeConfig } = getConfig()

export default class KoreApiResources {
  newV1ObjectMeta(name, namespace, resourceVersion) {
    const meta = new models.V1ObjectMeta()
    meta.setName(name)
    if (namespace) {
      meta.setNamespace(namespace)
    }
    if (resourceVersion) {
      meta.setResourceVersion(resourceVersion)
    }
    return meta
  }

  newResource({ type, apiVersion, kind, name, namespace, resourceVersion }) {
    const resource = new models[type]()
    resource.setApiVersion(apiVersion)
    resource.setKind(kind)
    resource.setMetadata(this.newV1ObjectMeta(name, namespace, resourceVersion))
    return resource
  }

  newV1SecretReference(name, namespace) {
    const secretRef = new models.V1SecretReference()
    secretRef.setName(name)
    secretRef.setNamespace(namespace)
    return secretRef
  }

  newV1Ownership({ group, version, kind, name, namespace }) {
    const ownership = new models.V1Ownership()
    ownership.setGroup(group)
    ownership.setVersion(version)
    ownership.setKind(kind)
    ownership.setName(name)
    ownership.setNamespace(namespace)
    return ownership
  }

  generatePlanResource(kind, values) {
    const resource = this.newResource({
      type: 'V1Plan',
      apiVersion: 'config.kore.appvia.io/v1',
      kind: 'Plan',
      name: values.name
    })

    const spec = new models.V1PlanSpec()
    spec.setKind(kind)
    spec.setDescription(values.description)
    spec.setSummary(values.summary)
    spec.setConfiguration(values.configuration)
    resource.setSpec(spec)

    return resource
  }

  generatePolicyResource(kind, values) {
    const resource = this.newResource({
      type: 'V1PlanPolicy',
      apiVersion: 'config.kore.appvia.io/v1',
      kind: 'PlanPolicy',
      name: values.name,
      namespace: publicRuntimeConfig.koreAdminTeamName
    })

    const spec = new models.V1PlanPolicySpec()
    spec.setKind(kind)
    spec.setDescription(values.description)
    spec.setSummary(values.summary)
    spec.setProperties(values.properties)
    resource.setSpec(spec)

    return resource
  }

  generateServicePlanResource(kind, resourceName, values) {
    const resource = this.newResource({
      type: 'V1ServicePlan',
      apiVersion: 'services.kore.appvia.io/v1',
      kind: 'ServicePlan',
      name: resourceName,
      namespace: 'kore'
    })

    const spec = new models.V1ServicePlanSpec()
    spec.setKind(kind)
    spec.setDescription(values.description)
    spec.setSummary(values.summary)
    spec.setConfiguration(values.configuration)
    resource.setSpec(spec)

    return resource
  }

  generateAccountManagementResource(resourceName, provider, orgResource, accountList, resourceVersion) {
    const resource = this.newResource({
      type: 'V1beta1AccountManagement',
      apiVersion: 'accounts.kore.appvia.io/v1beta1',
      kind: 'AccountManagement',
      name: resourceName,
      namespace: publicRuntimeConfig.koreAdminTeamName,
      resourceVersion
    })

    const spec = new models.V1beta1AccountManagementSpec()
    spec.setProvider(provider)

    const [ group, version ] = orgResource.apiVersion.split('/')
    spec.setOrganization(this.newV1Ownership({
      group,
      version,
      kind: orgResource.kind,
      name: orgResource.metadata.name,
      namespace: orgResource.metadata.namespace
    }))

    if (accountList) {
      const rules = accountList.map(project => {
        const rule = new models.V1beta1AccountsRule()
        rule.setName(project.name)
        rule.setDescription(project.description)
        rule.setPrefix(project.prefix)
        rule.setSuffix(project.suffix)
        rule.setPlans(project.plans)
        return rule
      })
      spec.setRules(rules)
    }

    resource.setSpec(spec)
    return resource
  }

  generateAllocationResource(resourceName, resourceToAllocate, teams, name, summary) {
    const resource = this.newResource({
      type: 'V1Allocation',
      apiVersion: 'config.kore.appvia.io/v1',
      kind: 'Allocation',
      name: resourceName,
      namespace: resourceToAllocate.metadata.namespace
    })

    const spec = new models.V1AllocationSpec()
    spec.setName(name ? name : resourceToAllocate.metadata.name)
    spec.setSummary(summary ? summary : `Allocation of ${resourceToAllocate.metadata.name}`)
    spec.setTeams(teams && teams.length > 0 ? teams : ['*'])

    const [group, version] = resourceToAllocate.apiVersion.split('/')
    spec.setResource(this.newV1Ownership({
      group,
      version,
      kind: resourceToAllocate.kind,
      name: resourceToAllocate.metadata.name,
      namespace: resourceToAllocate.metadata.namespace
    }))

    resource.setSpec(spec)
    return resource
  }

  generateTeamResource(values) {
    const resource = this.newResource({
      type: 'V1Team',
      apiVersion: 'config.kore.appvia.io/v1',
      kind: 'Team',
      name: canonical(values.teamName)
    })

    const spec = new models.V1TeamSpec()
    spec.setSummary(values.teamName.trim())
    spec.setDescription(values.teamDescription.trim())
    resource.setSpec(spec)
    return resource
  }

  team(team) {
    return {
      generateSecretResource: (name, secretType, description, data) => this.generateSecretResource(team, name, secretType, description, data),
      generateEKSCredentialsResource: (values, secretName) => this.generateEKSCredentialsResource(team, values, secretName),
      generateGKECredentialsResource: (values, secretName) => this.generateGKECredentialsResource(team, values, secretName),
      generateGCPOrganizationResource: (values, secretName) => this.generateGCPOrganizationResource(team, values, secretName),
      generateAWSOrganizationResource: (values, secretName) => this.generateAWSOrganizationResource(team, values, secretName),
      generateClusterResource: (user, values, plan, planValues, credentials) => this.generateClusterResource(team, user, values, plan, planValues, credentials),
      generateNamespaceClaimResource: (cluster, resourceName, values) => this.generateNamespaceClaimResource(team, cluster, resourceName, values),
      generateServiceResource: (cluster, values, plan, planValues) => this.generateServiceResource(team, cluster, values, plan, planValues),
      generateServiceCredentialsResource: (name, secretName, config, service, cluster, namespaceClaim) => this.generateServiceCredentialsResource(team, name, secretName, config, service, cluster, namespaceClaim)
    }
  }

  generateSecretResource(team, name, secretType, description, data) {
    const resource = this.newResource({
      type: 'V1Secret',
      apiVersion: 'config.kore.appvia.io',
      kind: 'Secret',
      name,
      namespace: team
    })

    const spec = new models.V1SecretSpec()
    spec.setType(secretType)
    spec.setDescription(description)
    spec.setData(data)
    resource.setSpec(spec)

    return resource
  }

  generateEKSCredentialsResource(team, values, secretName) {
    const resource = this.newResource({
      type: 'V1alpha1EKSCredentials',
      apiVersion: 'aws.compute.kore.appvia.io/v1alpha1',
      kind: 'EKSCredentials',
      name: values.name,
      namespace: team
    })

    const spec = new models.V1alpha1EKSCredentialsSpec()
    spec.setAccountID(values.accountID)
    spec.setCredentialsRef(this.newV1SecretReference(secretName, team))
    resource.setSpec(spec)

    return resource
  }

  generateGKECredentialsResource(team, values, secretName) {
    const resource = this.newResource({
      type: 'V1alpha1GKECredentials',
      apiVersion: 'gke.compute.kore.appvia.io/v1alpha1',
      kind: 'GKECredentials',
      name: values.name,
      namespace: team
    })

    const spec = new models.V1alpha1GKECredentialsSpec()
    spec.setProject(values.project)
    spec.setCredentialsRef(this.newV1SecretReference(secretName, team))
    resource.setSpec(spec)

    return resource
  }

  generateGCPOrganizationResource(team, values, secretName) {
    const resource = this.newResource({
      type: 'V1alpha1Organization',
      apiVersion: 'gcp.compute.kore.appvia.io/v1alpha1',
      kind: 'Organization',
      name: values.name,
      namespace: team
    })

    const spec = new models.V1alpha1OrganizationSpec()
    spec.setParentType('organization')
    spec.setParentID(values.parentID)
    spec.setBillingAccount(values.billingAccount)
    spec.setServiceAccount('kore')
    spec.setCredentialsRef(this.newV1SecretReference(secretName, team))
    resource.setSpec(spec)

    return resource
  }

  // TODO: complete and test
  generateAWSOrganizationResource(team, values, secretName) {
    const resource = this.newResource({
      type: 'V1alpha1AWSOrganization',
      apiVersion: 'aws.org.kore.appvia.io/v1alpha1',
      kind: 'AWSOrganization',
      name: values.name,
      namespace: team
    })

    const spec = new models.V1alpha1OrganizationSpec()
    spec.setOuName(values.ouName)
    spec.setRegion(values.region)
    spec.setRoleARN(values.roleARN)
    spec.setSsoUser({
      firstName: values.ssoUserFirstName,
      lastName: values.ssoUserLastName,
      email: values.ssoUserEmailAddress,
    })
    spec.setCredentialsRef(this.newV1SecretReference(secretName, team))
    resource.setSpec(spec)

    return resource
  }

  generateClusterResource(team, user, values, plan, planValues, credentials) {
    const resource = this.newResource({
      type: 'V1Cluster',
      apiVersion: 'clusters.compute.kore.appvia.io/v1',
      kind: 'Cluster',
      name: values.clusterName,
      namespace: team
    })

    const spec = new models.V1ClusterSpec()
    spec.setKind(plan.spec.kind)
    spec.setPlan(plan.metadata.name)
    spec.setConfiguration(planValues)
    spec.setCredentials(credentials)

    if (!(spec.configuration['clusterUsers'])) {
      spec.configuration['clusterUsers'] = [
        {
          username: user,
          roles: ['cluster-admin']
        }
      ]
    }

    resource.setSpec(spec)
    return resource
  }

  generateNamespaceClaimResource(team, cluster, resourceName, values) {
    const resource = this.newResource({
      type: 'V1NamespaceClaim',
      apiVersion: 'namespaceclaims.clusters.compute.kore.appvia.io/v1alpha1',
      kind: 'NamespaceClaim',
      name: resourceName,
      namespace: team
    })

    const spec = new models.V1NamespaceClaimSpec()
    spec.setName(values.name)

    const [ group, version ] = cluster.apiVersion.split('/')
    spec.setCluster(this.newV1Ownership({
      group,
      version,
      kind: cluster.kind,
      name: cluster.metadata.name,
      namespace: team
    }))
    resource.setSpec(spec)

    return resource
  }

  generateServiceResource(team, cluster, values, plan, planValues) {
    const resource = this.newResource({
      type: 'V1Service',
      apiVersion: 'services.compute.kore.appvia.io/v1',
      kind: 'Service',
      name: values.serviceName,
      namespace: team
    })

    const serviceSpec = new models.V1ServiceSpec()
    serviceSpec.setKind(plan.spec.kind)
    serviceSpec.setPlan(plan.metadata.name)
    serviceSpec.setConfiguration(planValues)
    serviceSpec.setClusterNamespace(values.createNamespace || values.namespace)

    const [ group, version ] = cluster.apiVersion.split('/')
    serviceSpec.setCluster(this.newV1Ownership({
      group,
      version,
      kind: cluster.kind,
      name: cluster.metadata.name,
      namespace: team
    }))

    resource.setSpec(serviceSpec)
    return resource
  }

  generateServiceCredentialsResource(team, name, secretName, config, service, cluster, clusterNamespace) {
    const resource = this.newResource({
      type: 'V1ServiceCredentials',
      apiVersion: 'servicecredentials.services.kore.appvia.io/v1',
      kind: 'ServiceCredentials',
      name: name,
      namespace: team
    })

    const spec = new models.V1ServiceCredentialsSpec()
    spec.setKind(service.spec.kind)
    const [ sGroup, sVersion ] = service.apiVersion.split('/')
    spec.setService(this.newV1Ownership({
      group: sGroup,
      version: sVersion,
      kind: service.kind,
      name: service.metadata.name,
      namespace: team
    }))
    const [ cGroup, cVersion ] = cluster.apiVersion.split('/')
    spec.setCluster(this.newV1Ownership({
      group: cGroup,
      version: cVersion,
      kind: cluster.kind,
      name: cluster.metadata.name,
      namespace: team
    }))
    spec.setClusterNamespace(clusterNamespace)
    spec.setSecretName(secretName)
    spec.setConfiguration(config)

    resource.setSpec(spec)
    return resource
  }

}

