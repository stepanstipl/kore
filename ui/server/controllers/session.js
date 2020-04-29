const Router = require('express').Router

function getSessionUser(orgService) {
  return async(req, res) => {
    const user = req.session.passport.user
    try {
      await orgService.refreshUser(user)
      return res.json(user)
    } catch (err) {
      console.log('Failed to refresh user in /session/user', err)
      return res.status(err.statusCode || 500).send()
    }
  }
}

function initRouter({ ensureAuthenticated, ensureUserCurrent, persistRequestedPath, orgService }) {
  const router = Router()
  router.get('/session/user', ensureAuthenticated, ensureUserCurrent, persistRequestedPath, getSessionUser(orgService))
  return router
}

module.exports = {
  initRouter
}
