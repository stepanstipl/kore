const redirect = ({ res, router, forceSSR = false, path }) => {
  if (res) {
    res.redirect(path)
    res.end()
    return
  }
  if (forceSSR) {
    window.location.href = path
    return
  }
  router.push(path)
}

module.exports = redirect
