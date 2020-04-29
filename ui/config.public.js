const config = require('./config')

// pick vales from the config which can be exposed to the client
// should contain nothing secret
module.exports = {
  koreApiPublicUrl: config.api.publicUrl,
  koreBaseUrl: config.kore.baseUrl,
  koreAdminTeamName: config.kore.koreAdminTeamName,
  ignoreTeams: config.kore.ignoreTeams,
  sessionTtlInSeconds: config.session.ttlInSeconds,
  authOpenidUrl: config.auth.openid.url,
  authOpenidCallbackUrl: config.auth.openid.callbackURL,
  gtmId: config.kore.gtmId,
  showPrototypes: config.kore.showPrototypes,
  featureGates: config.kore.featureGates
}
