const Router = require('express').Router

function processTeamInvitation(KoreApi) {
  return async (req, res) => {
    const token = req.params.token
    try {
      const api = await KoreApi.client({ id_token: req.session.passport.user.id_token })
      const invitationResponse = await api.InvitationSubmit(token)
      let redirectTo = '/'
      if (invitationResponse.team) {
        redirectTo = `/teams/${invitationResponse.team}?invitation=true`
      }
      return res.redirect(redirectTo)
    } catch (err) {
      const status = (err.response && err.response.status) || 500
      const message = (err.response && err.response.data && err.response.data.message) || err.message
      console.error(`Error processing team invitation link with token ${token}`, status, message, err)
      const invitationLink = req.protocol + '://' + req.get('host') + req.originalUrl
      return res.redirect(`/invalid-team-invitation?link=${invitationLink}`)
    }
  }
}

function persistPath(req, res, next) {
  req.session.requestedPath = req.path
  next()
}

function initRouter({ ensureAuthenticated, ensureUserCurrent, KoreApi }) {
  const router = Router()
  router.get('/process/teams/invitation/:token', persistPath, ensureAuthenticated, ensureUserCurrent, processTeamInvitation(KoreApi))
  return router
}

module.exports = {
  initRouter
}
