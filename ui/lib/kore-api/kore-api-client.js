const redirect = require('../utils/redirect')
const inflect = require('inflect')

class KoreApiClient {
  spec = null
  apis = null
  client = null
  basePath = null
  constructor(api, basePath) {
    this.apis = api.apis
    this.spec = api.spec
    this.client = api
    this.basePath = basePath

    // This decorates every operation returned from the swagger with a function which unwraps the
    // returned object, making the usage of the api much cleaner in the rest of the code.
    // Also does global API error handling.
    Object.keys(this.apis).forEach(tagName =>
      Object.entries(this.apis[tagName]).forEach(([functionName, fnc]) =>
        this.apis[tagName][functionName] = (...args) => this._wrapFunc(fnc, ...args)
      )
    )
  }

  _wrapFunc = (fnc, ...args) => fnc(...args).then(
    // Unwrap body so normal usage 'just works':
    (res) => res.body,
    // Handle a few specific error cases (not found, auth errors):
    (err) => {
      // Handle not found as a null
      if (err.response && err.response.status === 404) {
        return null
      }
      // Handle 401 unauth, if running in a browser:
      if (err.response && err.response.status === 401 && process.browser) {
        redirect({
          path: '/login/refresh',
          ensureRefreshFromServer: true
        })
      }
      if (err.response && (err.response.status === 400 || err.response.status === 409) && err.response.body) {
        throw { ...err.response.body, statusCode: err.response.status }
      }
      // @TODO: Handle validation errors (400) and forbidden (403)
      throw err
    }
  )

  _mapResourceToOperation = (team, resource) => {
    const parts = resource.split('/')
    const resType = parts[0]
    const name = parts[1]
    let pathName = null
    let basePath = this.basePath
    if (process.browser) {
      basePath = '/apiproxy'
    }
    const resTypePlural = inflect.pluralize(resType).toLowerCase()
    pathName = `${basePath}/teams/{team}/${resTypePlural}/{name}`
    return {
      pathName: pathName,
      method: 'GET',
      parameters: {
        team: team,
        name: name
      }
    }
  }

  GetTeamResource = (team, resource) => this._wrapFunc((t, r) => this.client.execute(this._mapResourceToOperation(t, r)), team, resource)

  // @TODO: Auto-generate these?

  // IDP
  GetDefaultIDP = () => this.apis.default.GetDefaultIDP()
  UpdateIDP = (name, idp) => this.apis.default.UpdateIDP({ name, body: JSON.stringify(idp) })
  UpdateIDPClient = (name, idpClient) => this.apis.default.UpdateIDPClient({ name, body: JSON.stringify(idpClient) })

  WhoAmI = () => this.apis.default.WhoAmI()

  // Users
  ListUsers = () => this.apis.default.ListUsers()
  ListUserTeams = (user) => this.apis.default.ListUserTeams({ user })
  UpdateUser = (user, userSpec) => this.apis.default.UpdateUser({ user, body: JSON.stringify(userSpec) })

  // Plans
  ListPlans = (kind) => this.apis.default.ListPlans({ kind })
  UpdatePlan = (name, plan) => this.apis.default.UpdatePlan({ name, body: JSON.stringify(plan) })
  GetPlanSchema = (kind) => this.apis.default.GetPlanSchema({ kind })
  RemovePlan = (name) => this.apis.default.RemovePlan({ name })

  // Services
  ListServiceProviders = () => this.apis.default.ListServiceProviders()
  ListServiceKinds = () => this.apis.default.ListServiceKinds()
  ListServicePlans = (kind) => this.apis.default.ListServicePlans({ kind })
  GetServicePlan = (name) => this.apis.default.GetServicePlan({ name })
  GetServicePlanDetails = (name, team, cluster) => this.apis.default.GetServicePlanDetails({ name, team, cluster })
  UpdateServicePlan = (name, servicePlan) => this.apis.default.UpdateServicePlan({ name, body: JSON.stringify(servicePlan) })
  DeleteServicePlan = (name) => this.apis.default.DeleteServicePLan({ name })
  GetServiceKind = (name) => this.apis.default.GetServiceKind({ name })
  UpdateServiceKind = (name, kind) => this.apis.default.UpdateServiceKind({ name, body: kind })

  // Audit
  ListAuditEvents = () => this.apis.default.ListAuditEvents()

  // Security
  security = {
    GetSecurityOverview: () => this.apis.security.GetSecurityOverview(),
    GetSecurityScanForResource: (group, version, kind, namespace, name) => this.apis.security.GetSecurityScanForResource({ group, version, kind, namespace, name }),
    ListSecurityScansForResource: (group, version, kind, namespace, name) => this.apis.security.ListSecurityScansForResource({ group, version, kind, namespace, name }),
    GetSecurityScan: (id) => this.apis.security.GetSecurityScan({ id: id }),
    ListSecurityRules: () => this.apis.security.ListSecurityRules(),
    GetSecurityRule: (code) => this.apis.security.GetSecurityRule({ code })
  }

  // Metadata
  metadata = {
    GetCloudRegions: (cloud) => this.apis.metadata.GetCloudRegions({ cloud }),
    GetCloudRegionZones: (cloud, region) => this.apis.metadata.GetCloudRegionZones({ cloud, region }),
    GetCloudNodeTypes: (cloud, region) => this.apis.metadata.GetCloudNodeTypes({ cloud, region }),
    GetKubernetesRegions: (provider) => this.apis.metadata.GetKubernetesRegions({ provider }),
    GetKubernetesRegionZones: (provider, region) => this.apis.metadata.GetKubernetesRegionZones({ provider, region }),
    GetKubernetesNodeTypes: (provider, region) => this.apis.metadata.GetKubernetesNodeTypes({ provider, region }),
    GetKubernetesVersions: (provider, region) => this.apis.metadata.GetKubernetesVersions({ provider, region }),
  }

  // Costs
  costestimates = {
    EstimateClusterPlanCost: (plan) => this.apis.costestimates.EstimateClusterPlanCost({ body: JSON.stringify(plan) }),
    EstimateServicePlanCost: (plan) => this.apis.costestimates.EstimateServicePlanCost({ body: JSON.stringify(plan) })
  }

  // Policies
  ListPlanPolicies = (kind) => this.apis.default.ListPlanPolicies({ kind })
  GetPlanPolicy = (name) => this.apis.default.GetPlanPolicy({ name })
  UpdatePlanPolicy = (name, plan) => this.apis.default.UpdatePlanPolicy({ name, body: JSON.stringify(plan) })
  RemovePlanPolicy = (name) => this.apis.default.RemovePlanPolicy({ name })

  // Account management
  ListAccounts = (kind) => this.apis.default.ListAccounts({ kind })
  UpdateAccount = (name, account) => this.apis.default.UpdateAccount({ name, body: JSON.stringify(account) })
  RemoveAccount = (name) => this.apis.default.RemoveAccount({ name })

  // Credentials
  ListGKECredentials = (team) => this.apis.default.ListGKECredentials({ team })
  GetGKECredential = (team, name) => this.apis.default.GetGKECredential({ team, name })
  UpdateGKECredential = (team, name, resource) => this.apis.default.UpdateGKECredential({ team, name, body: JSON.stringify(resource) })
  RemoveGKECredential = (team, name) => this.apis.default.RemoveGKECredential({ team, name })
  ListGCPOrganizations = (team) => this.apis.default.ListGCPOrganizations({ team })
  GetGCPOrganization = (team, name) => this.apis.default.GetGCPOrganization({ team, name })
  UpdateGCPOrganization = (team, name, org) => this.apis.default.UpdateGCPOrganization({ team, name, body: JSON.stringify(org) })
  DeleteGCPOrganization = (team, name) => this.apis.default.DeleteGCPOrganization({ team, name })
  ListEKSCredentials = (team) => this.apis.default.ListEKSCredentials({ team })
  GetEKSCredentials = (team, name) => this.apis.default.GetEKSCredentials({ team, name })
  UpdateEKSCredentials = (team, name, resource) => this.apis.default.UpdateEKSCredentials({ team, name, body: JSON.stringify(resource) })
  DeleteEKSCredentials = (team, name) => this.apis.default.DeleteEKSCredentials({ team, name })

  // Invitation
  InvitationSubmit = (token) => this.apis.default.InvitationSubmit({ token })

  // Teams
  GetTeam = (team) => this.apis.default.GetTeam({ team })
  RemoveTeam = (team) => this.apis.default.RemoveTeam({ team })
  ListTeams = () => this.apis.default.ListTeams()
  UpdateTeam = (team, teamSpec) => this.apis.default.UpdateTeam({ team, body: JSON.stringify(teamSpec) })
  GenerateInviteLink = (team, expire) => this.apis.default.GenerateInviteLink({ team, expire })
  ListTeamMembers = (team) => this.apis.default.ListTeamMembers({ team })
  AddTeamMember = (team, user) => this.apis.default.AddTeamMember({ team, user })
  RemoveTeamMember = (team, user) => this.apis.default.RemoveTeamMember({ team, user })
  GetTeamSecurityOverview = (team) => this.apis.default.GetTeamSecurityOverview({ team })
  ListAllocations = (team, assigned = undefined) => this.apis.default.ListAllocations({ team, assigned })
  GetAllocation = (team, name) => this.apis.default.GetAllocation({ team, name })
  UpdateAllocation = (team, name, resource) => this.apis.default.UpdateAllocation({ team, name, body: JSON.stringify(resource) })
  RemoveAllocation = (team, name) => this.apis.default.RemoveAllocation({ team, name })
  ListClusters = (team) => this.apis.default.ListClusters({ team })
  UpdateCluster = (team, name, cluster) => this.apis.default.UpdateCluster({ team, name, body: JSON.stringify(cluster) })
  RemoveCluster = (team, name) => this.apis.default.RemoveCluster({ team, name })
  GetCluster = (team, name) => this.apis.default.GetCluster({ team, name })
  ListNamespaces = (team) => this.apis.default.ListNamespaces({ team })
  GetNamespace = (team, name) => this.apis.default.GetNamespace({ team, name })
  UpdateNamespace = (team, name, namespace) => this.apis.default.UpdateNamespace({ team, name, body: JSON.stringify(namespace) })
  RemoveNamespace = (team, name) => this.apis.default.RemoveNamespace({ team, name })
  GetTeamPlanDetails = (team, plan) => this.apis.default.GetTeamPlanDetails({ team, plan })
  UpdateTeamSecret = (team, name, secret) => this.apis.default.UpdateTeamSecret({ team, name, body: JSON.stringify(secret) })
  ListTeamAudit = (team) => this.apis.default.ListTeamAudit({ team })
  ListServices = (team) => this.apis.default.ListServices({ team })
  UpdateService = (team, name, service) => this.apis.default.UpdateService({ team, name, body: JSON.stringify(service) })
  GetService = (team, name) => this.apis.default.GetService({ team, name })
  DeleteService = (team, name) => this.apis.default.DeleteService({ team, name })
  GetServiceCredentials = (team, name) => this.apis.default.GetServiceCredentials({ team, name })
  UpdateServiceCredentials = (team, name, serviceCredential) => this.apis.default.UpdateServiceCredentials({ team, name, body: JSON.stringify(serviceCredential) })
  ListServiceCredentials = (team, cluster, service) => this.apis.default.ListServiceCredentials({ team, cluster, service })
  DeleteServiceCredentials = (team, name) => this.apis.default.DeleteServiceCredentials({ team, name })
}

module.exports = KoreApiClient
