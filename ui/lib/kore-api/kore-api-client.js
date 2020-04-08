const redirect = require('../utils/redirect')

class KoreApiClient {
  spec = null
  apis = null
  constructor(api) {
    this.apis = api.apis
    this.spec = api.spec

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
      // @TODO: Handle validation errors (400) and forbidden (403)
      throw err
    }
  )

  // @TODO: Auto-generate these?

  // Users
  ListUsers = () => this.apis.default.ListUsers()
  ListUserTeams = (user) => this.apis.default.ListUserTeams({ user })
  UpdateUser = (user, userSpec) => this.apis.default.UpdateUser({ user, body: JSON.stringify(userSpec) })

  // Plans
  ListPlans = (kind) => this.apis.default.ListPlans({ kind })

  // Teams
  GetTeam = (team) => this.apis.default.GetTeam({ team })
  ListTeams = () => this.apis.default.ListTeams()
  UpdateTeam = (team, teamSpec) => this.apis.default.UpdateTeam({ team, body: JSON.stringify(teamSpec) })
  ListTeamMembers = (team) => this.apis.default.ListTeamMembers({ team })
  AddTeamMember = (team, user) => this.apis.default.AddTeamMember({ team, user })
  RemoveTeamMember = (team, user) => this.apis.default.RemoveTeamMember({ team, user })
  ListGKECredentials = (team) => this.apis.default.ListGKECredentials({ team })
  ListGCPOrganizations = (team) => this.apis.default.findOrganizations({ team })
  ListAllocations = (team, assigned = undefined) => this.apis.default.ListAllocations({ team, assigned: assigned })
  ListClusters = (team) => this.apis.default.ListClusters({ team })
  UpdateCluster = (team, name, cluster) => this.apis.default.UpdateCluster({ team, name, body: JSON.stringify(cluster) })
  ListNamespaces = (team) => this.apis.default.ListNamespaces({ team })
}

module.exports = KoreApiClient
