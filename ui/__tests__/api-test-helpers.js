import nock from 'nock'
import config from '../config'
import fs from 'fs'
import path from 'path'
import url from 'url'

class ApiTestHelpers {
  static basePath = null

  static getScope = () => {
    const u = url.parse(config.koreApi.url)
    ApiTestHelpers.basePath = u.path
    const spec = JSON.parse((fs.readFileSync(path.join(__dirname, '../kore-api-swagger.json'))).toString())
  
    return nock(`${u.protocol}//${u.host}`)
      .get('/swagger.json')
      .optionally()
      .reply(200, spec, { 'Content-type': 'application/json'})
  }
}

export default ApiTestHelpers