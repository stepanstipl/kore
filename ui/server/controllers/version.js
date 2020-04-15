const Router = require('express').Router
const fs = require('fs')
const path = require('path')

function version() {
  return async (_, res) => {
    const v = await fs.promises.readFile(path.resolve(__dirname, '../../VERSION'), 'utf8')
    return res.json({ version: v })
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
