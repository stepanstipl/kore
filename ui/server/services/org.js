const User = require('../../lib/crd/User')
const config = require('../../config')

class OrgService {
  constructor(KoreApi) {
    this.KoreApi = KoreApi
  }

  async getApiClient(id_token) {
    const api = await this.KoreApi.client({ id_token })
    return api
  }

  async getOrCreateUser(user) {
    try {
      const userResource = await User(user)
      console.log(`*** putting user ${user.id}`, userResource)
      const api = await this.getApiClient(config.api.token)
      const userResult = await api.UpdateUser(user.id, userResource)
      const adminTeamMembers = await this.getTeamMembers(config.kore.koreAdminTeamName, config.api.token)
      if (adminTeamMembers.length === 1) {
        await this.addUserToTeam(config.kore.koreAdminTeamName, user.id, config.api.token)
      }
      userResult.teams = await this.getUserTeams(user)
      userResult.isAdmin = this.isAdmin(userResult.teams.userTeams)
      return userResult
    } catch (err) {
      console.error('Error in getOrCreateUser from API', err)
      return Promise.reject(err)
    }
  }

  /* eslint-disable require-atomic-updates */
  async refreshUser(user) {
    user.teams = await this.getUserTeams(user)
    user.isAdmin = this.isAdmin(user.teams.userTeams)
  }
  /* eslint-enable require-atomic-updates */

  isAdmin(userTeams) {
    return (userTeams || []).filter(t => t.metadata && t.metadata.name === config.kore.koreAdminTeamName).length > 0
  }

  async getTeamMembers(team, requestingIdToken) {
    try {
      const api = await this.getApiClient(requestingIdToken)
      const result = await api.ListTeamMembers(team)
      console.log(`*** found team members for team: ${team}`, result.items)
      return result.items
    } catch (err) {
      console.error('Error getting team members from API', err)
      return Promise.reject(err)
    }
  }

  async addUserToTeam(team, username, requestingIdToken) {
    console.log(`*** adding user ${username} to team ${team}`)
    try {
      const api = await this.getApiClient(requestingIdToken)
      await api.AddTeamMember(team, username)
    } catch (err) {
      console.error('Error adding user to team', err)
      return Promise.reject(err)
    }
  }

  async getUserTeams(user) {
    try {
      const api = await this.getApiClient(user.id_token)

      const userTeamsList = await api.ListUserTeams(user.id)
      const userTeams = userTeamsList.items

      if (this.isAdmin(userTeams)) {
        const allTeamsList = await api.ListTeams()
        const userTeamIdList = userTeams.map(t => t.metadata.name)
        return {
          userTeams,
          otherTeams: allTeamsList.items.filter(at => !userTeamIdList.includes(at.metadata.name))
        }
      }
      return { userTeams }
    } catch (err) {
      console.error('Error getting teams for user', err)
      return Promise.reject(err)
    }
  }

  async getTeamGkeCredentials(team, requestingIdToken) {
    try {
      const api = await this.getApiClient(requestingIdToken)
      const result = await api.ListGKECredentials(team)
      return result
    } catch (err) {
      if (err.response && err.response.status === 404) {
        return null
      }
      console.error('Error getting team GKE credentials from API', err)
      return Promise.reject(err)
    }
  }

  async hasTeamCredentials(team, requestingIdToken) {
    try {
      const api = await this.getApiClient(requestingIdToken)
      const [ gkeCredentialsList, gcpOrgList, eksCredentialsList ] = await Promise.all([
        api.ListGKECredentials(team),
        api.ListGCPOrganizations(team),
        api.ListEKSCredentials(team)
      ])
      return gkeCredentialsList.items.length !== 0 ||
        gcpOrgList.items.length !== 0 ||
        eksCredentialsList.items.length !== 0
    } catch (err) {
      console.error('Error checking for team credentials from API', err)
      return Promise.reject(err)
    }
  }
}

module.exports = OrgService
