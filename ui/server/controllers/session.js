const Router = require('express').Router

function getSessionUser(orgService, config) {
  return async(req, res) => {
    const user = req.session.passport.user
    try {
      await orgService.refreshUser(user)
      const clientConfig = {
        apiUrl: config.koreApi.publicUrl,
        featureGates: config.kore.featureGates
      }
      return res.json({ user, config: clientConfig })
    } catch (err) {
      console.log('Failed to refresh user in /session/user', err)
      return res.status(err.statusCode || 500).send()
    }
  }
}

function getConfig(config) {
  return (req, res) => {
    res.json({ apiUrl: config.koreApi.publicUrl })
  }
}

function initRouter({ ensureAuthenticated, ensureUserCurrent, persistRequestedPath, orgService, config }) {
  const router = Router()
  router.get('/session/user', ensureAuthenticated, ensureUserCurrent, persistRequestedPath, getSessionUser(orgService, config))
  router.get('/session/config', ensureAuthenticated, ensureUserCurrent, getConfig(config))
  return router
}

module.exports = {
  initRouter
}
