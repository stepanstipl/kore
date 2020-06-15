const Router = require('express').Router
const router = Router()
const config = require('../config')
const KoreApi = require('../lib/kore-api')
const AuthService = require('./services/auth')
const OrgService = require('./services/org')
const OpenIdClient = require('./lib/openid-client')
const ensureAuthenticated = require('./middleware/ensure-authenticated')
const ensureUserCurrent = require('./middleware/ensure-user-current')
const persistRequestedPath = require('./middleware/persist-requested-path')

const koreConfig = config.kore
const embeddedAuth = config.auth.embedded
const userClaimsOrder = config.auth.openid.userClaimsOrder.split(',')
const koreApiUrl = config.api.url
const baseUrl = config.kore.baseUrl

const authService = new AuthService(KoreApi)
const orgService = new OrgService(KoreApi)
const authCallback = require('./middleware/auth-callback')(orgService, authService, koreConfig, userClaimsOrder, embeddedAuth)

const openIdClient = new OpenIdClient(baseUrl, config.auth.openid, embeddedAuth, authService)
openIdClient.init()
  .then(() => {})
  .catch(err => {
    console.error('Unexpected error occurred during openIdClient initialisation', err)
  })
const ensureOpenIdClient = require('./middleware/ensure-openid-client')(openIdClient)
const ensureAuthenticatedRedirect = ensureAuthenticated({ redirect: true })
const ensureAuthenticated401 = ensureAuthenticated({ redirect: false })

router.use(require('./controllers/auth-local').initRouter({ authService, ensureAuthenticated: ensureAuthenticatedRedirect, authCallback }))
router.use(require('./controllers/auth-openid').initRouter({ authService, ensureOpenIdClient, persistRequestedPath, embeddedAuth, authCallback }))

router.use(require('./controllers/swagger').initRouter({ koreApiUrl }))

// other routes must have an authenticated user
router.use(require('./controllers/session').initRouter({ ensureAuthenticated: ensureAuthenticated401, ensureUserCurrent, persistRequestedPath, orgService }))
router.use(require('./controllers/apiproxy').initRouter({ ensureAuthenticated: ensureAuthenticatedRedirect, ensureUserCurrent, koreApiUrl }))
router.use(require('./controllers/process').initRouter({ ensureAuthenticated: ensureAuthenticatedRedirect, ensureUserCurrent, KoreApi }))
router.use(require('./controllers/version').initRouter())

module.exports = router
