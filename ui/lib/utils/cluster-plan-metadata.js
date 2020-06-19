export default class ClusterPlanMetadataService {
  static async getRegions(kind) {
    const resp = await fetch({ url: `/pricing/${kind}-continents.json`})
    const details = await resp.json()
    return details
  }
  static async getProducts(kind, region) {
    console.log("Region not implemented, returning for default region London", region)
    const resp = await fetch({ url: `/pricing/${kind}-products.json`})
    const details = await resp.json()
    return details
  }
  static async getVersions(kind, region) {
    console.log("Region not implemented, returning for default region London", region)
    const resp = await fetch({ url: `/pricing/${kind}-versions.json`})
    const details = await resp.json()
    return details
  }
}