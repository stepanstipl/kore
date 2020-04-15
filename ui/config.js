const getOpenidURL = () => {
  const envValue = process.env.KORE_IDP_SERVER_URL
  if (!envValue) {
    return 'https://my-openid-domain.com'
  }
  // the URL is already correct
  if (envValue.indexOf('.well-known/openid-configuration') !== -1) {
    return envValue
  }
  // check if we need to add trailing slash to the URL before adding the required path
  let suffix = '.well-known/openid-configuration'
  if (envValue.lastIndexOf('/') !== envValue.length - 1) {
    suffix = `/${suffix}`
  }
  return `${envValue}${suffix}`
}

module.exports = {
  server: {
    port: process.env.PORT || '3000',
    session: {
      sessionSecret: process.env.SESSION_SECRET || 'sessionsecret',
      url: process.env.REDIS_URL || 'redis://127.0.0.1:6379',
      ttlInSeconds: 1200
    }
  },
  auth: {
    embedded: process.env.KORE_UI_USE_EMBEDDED_AUTH === 'true' || false,
    openid: {
      enabled: Boolean(process.env.KORE_IDP_SERVER_URL),
      url: getOpenidURL(),
      callbackURL: process.env.KORE_CALLBACK_URL || 'http://localhost:3000/auth/callback',
      clientID: process.env.KORE_IDP_CLIENT_ID || 'my-openid-client-id',
      clientSecret: process.env.KORE_IDP_CLIENT_SECRET || 'my-openid-client-secret',
      userClaimsOrder: process.env.KORE_IDP_USER_CLAIMS || 'preferred_username,email,name,username'
    }
  },
  kore: {
    baseUrl: process.env.KORE_BASE_URL || 'http://localhost:3000',
    koreAdminTeamName: 'kore-admin',
    ignoreTeams: ['kore-admin', 'kore-default'],
    gtmId: 'GTM-T9THH55'
  },
  koreApi: {
    url: process.env.KORE_API_URL || 'http://localhost:10080/api/v1alpha1',
    publicUrl: process.env.KORE_API_PUBLIC_URL || 'http://localhost:10080',
    token: process.env.KORE_API_TOKEN || 'password'
  }
}
