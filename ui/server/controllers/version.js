const Router = require('express').Router

function version() {
  return async (_, res) => {
    return res.json({ version: process.env.UI_VERSION })
  }
}

function initRouter() {
  const router = Router()
  router.use('/version', version())
  return router
}

module.exports = {
  initRouter
}
