const Router = require('express').Router

function getSessionUser(orgService) {
  return async(req, res) => {
    const user = req.session.passport.user
    await orgService.refreshUser(user)
    return res.json(user)
  }
}

function getConfig(koreApi) {
  return (req, res) => {
    res.json({ apiUrl: koreApi.url })
  }
}

function initRouter({ ensureAuthenticated, ensureUserCurrent, persistRequestedPath, orgService, koreApi }) {
  const router = Router()
  router.get('/session/user', ensureAuthenticated, ensureUserCurrent, persistRequestedPath, getSessionUser(orgService))
  router.get('/session/config', ensureAuthenticated, ensureUserCurrent, getConfig(koreApi))
  return router
}

module.exports = {
  initRouter
}
