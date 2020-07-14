module.exports = (ctx) => (req) => {
  // If we're running on the server, we need to layer in the user's identity to
  // the request. When running in the browser, this is handled by the cookie-based
  // session.
  if (!process.browser && !ctx) {
    throw new Error('KoreApi client requires ctx containing id_token OR passport session data')
  }
  if (!process.browser && ctx) {
    if (ctx.id_token || ctx.req) {
      req.headers['Authorization'] = `Bearer ${ctx.id_token || ctx.req.session.passport.user.id_token}`
    }
    // optionally override the auth with a username/password
    // used when authenticating local users
    if (ctx.basicAuth) {
      req.headers['Authorization'] = `Basic ${ctx.basicAuth}`
    }
  }
  return req
}
