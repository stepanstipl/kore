const config = require('../../config')

module.exports = {
  apiVersion: 'core.kore.appvia.io/v1',
  kind: 'IDPClient',
  metadata: {
    name: 'default'
  },
  spec: {
    displayName: 'Kore UI',
    secret: config.auth.openid.clientSecret,
    id: config.auth.openid.clientID,
    redirectURIs: [`${config.kore.baseUrl}/auth/callback`]
  }
}
