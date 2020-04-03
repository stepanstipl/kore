const Router = require('next/router')

const redirect = (res, path, forceSSR = false) => {
  if (res) {
    res.redirect(path)
    res.end()
    return {}
  }
  if (forceSSR) {
    window.location.href = path
  } else {
    Router.push(path)
  }
  return {}
}

module.exports = redirect
