import redirect from './redirect'

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
    (res) => res.body,
    (err) => {
      // Handle not found as a null
      if (err.response && err.response.status === 404) {
        return null
      }
      // Handle 401 unauth:
      if (err.response && err.response.status === 401) {
        redirect(null, '/login/refresh', true)
      }

      throw err
    }
  )

  // @TODO: Auto-generate these?
  ListAllocations = (team, assigned = undefined) => this.apis.default.ListAllocations({team, assigned: assigned})
  ListPlans = () => this.apis.default.ListPlans()
  ListUsers = () => this.apis.default.ListUsers()
  GetTeam = (team) => this.apis.default.GetTeam({team})
  ListTeams = () => this.apis.default.ListTeams()
  ListTeamMembers = (team) => this.apis.default.ListTeamMembers({team})
  ListGKECredentials = (team) => this.apis.default.ListGKECredentials({team})
  ListGCPOrganizations = (team) => this.apis.default.findOrganizations({team})
  UpdateTeam = (team, teamSpec) => this.apis.default.UpdateTeam({team, body: teamSpec})
  AddTeamMember = (team, user) => this.apis.default.AddTeamMember({team, user})
  RemoveTeamMember = (team, user) => this.apis.default.RemoveTeamMember({team, user})
}

export default KoreApiClient