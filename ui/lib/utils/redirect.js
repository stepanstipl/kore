const redirect = ({ res, router, ensureRefreshFromServer = false, path }) => {
  if (res) {
    res.redirect(path)
    res.end()
    return
  }
  if (ensureRefreshFromServer) {
    window.location.href = path
    return
  }
  router.push(path)
}

module.exports = redirect
