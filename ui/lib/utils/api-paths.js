module.exports = {
  users: '/users',
  teams: '/teams',
  plans: '/plans',
  useTeamInvitation: token => `/teams/invitation/${token}`,
  user: id => ({
    self: `/users/${id}`,
    teams: `/users/${id}/teams`
  }),
  team: id => ({
    self: `/teams/${id}`,
    members: `/teams/${id}/members`,
    gkes: `/teams/${id}/gkes`,
    clusters: `/teams/${id}/clusters`,
    namespaceClaims: `/teams/${id}/namespaceclaims`,
    secrets: `/teams/${id}/secrets`,
    gcpOrganizations: `/teams/${id}/organizations`,
    gkeCredentials: `/teams/${id}/gkecredentials`,
    allocations: `/teams/${id}/allocations`,
    generateInviteLink: `/teams/${id}/invites/generate`,
    audit: `/teams/${id}/audit`
  }),
  idp: {
    default: '/idp/default',
    clients: '/idp/clients',
    configured: '/idp/configured'
  },
  audit: '/audit',
  whoAmI: '/whoami'
}
