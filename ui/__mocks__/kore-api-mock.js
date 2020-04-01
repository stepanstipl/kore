import Swagger from 'swagger-client'
import fs from 'fs'
import path from 'path'

class KoreApiMock {
  static spec = null
  static basePath = '/api/v1alpha1'
  static swaggerUrl = 'http://localhost/swagger.json'
  static mocks = {}
  static calls = {}
  static allCallCount = 0

  static reset = () => {
    KoreApiMock.calls = {}
    KoreApiMock.allCallCount = 0
  }

  static client = async (ctx = null) => {
    const spec = JSON.parse((await fs.promises.readFile(path.join(__dirname, 'kore-api-swagger.json'))).toString())
    const api = await Swagger(
      KoreApiMock.swaggerUrl,
      {
        spec: spec,
        requestInterceptor: (req) => {
          if (process.browser) {
            req.url = req.url.replace(KoreApiMock.basePath, '/apiproxy')
          } else {
            req.headers['Authorization'] = `Bearer ${ctx.req.session.passport.user.id_token}`
          }
          return req
        }
      }
    )
    return KoreApiMock._decorateFunctions(api).default
  }

  static registerMock = (tagName, functionName, func) => {
    if (!(KoreApiMock.mocks[tagName])) {
      KoreApiMock.mocks[tagName] = {}
    }
    KoreApiMock.mocks[tagName][functionName] = func
  }

  static _decorateFunctions = (api) => {
    let apis = api.apis
    Object.keys(apis).forEach(tagName => {
      Object.entries(apis[tagName]).forEach(([functionName]) => {
        if (!(KoreApiMock.calls[tagName])) {
          KoreApiMock.calls[tagName] = {}
        }
        if (!(KoreApiMock.calls[tagName][functionName])) {
          KoreApiMock.calls[tagName][functionName] = []
        }
        apis[tagName][functionName] = (...args) => {
          KoreApiMock.calls[tagName][functionName].push(...args)
          KoreApiMock.allCallCount++
          // Call the mock if we have one, otherwise return a dummy
          // empty response.
          if (KoreApiMock.mocks[tagName] && KoreApiMock.mocks[tagName][functionName]) {
            return KoreApiMock.mocks[tagName][functionName](...args)
          } else {
            return Promise.resolve({})
          }
        }
      })
    })
    return apis
  }
}

export default KoreApiMock