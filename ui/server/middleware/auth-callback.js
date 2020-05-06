module.exports = (orgService, authService, koreConfig, userClaimsOrder, embeddedAuth) => {
  return async (req, res) => {
    const session = req.session
    const user = session.passport.user
    if (session.localUser) {
      user.id = user.username
    } else {
      userClaimsOrder.some(c => {
        user.id = user[c]
        return user.id
      })
    }
    try {
      const userInfo = await orgService.getOrCreateUser(user)

      /* eslint-disable require-atomic-updates */
      user.teams = userInfo.teams || []
      user.isAdmin = userInfo.isAdmin
      /* eslint-enable require-atomic-updates */
      if (session.requestedPath) {
        return res.redirect(session.requestedPath)
      }

      let redirectPath = '/'
      if (user.isAdmin) {
        if (embeddedAuth) {
          const authProvider = await authService.getDefaultConfiguredIdp()
          if (!authProvider) {
            redirectPath = '/setup/authentication'
          }
        }
        const setupComplete = await orgService.hasTeamCredentials(koreConfig.koreAdminTeamName, user.id_token)
        if (!setupComplete) {
          /* eslint-disable-next-line require-atomic-updates */
          redirectPath = '/setup/kore'
        }
      }
      res.redirect(redirectPath)

    } catch (err) {
      /* eslint-disable-next-line require-atomic-updates */
      req.session.loginError = 500
      return res.redirect('/login')
    }
  }
}
