import getConfig from 'next/config'
const { publicRuntimeConfig } = getConfig()

export function getPlanCloudInfo(planKind) {
  return publicRuntimeConfig.clusterProviderCloudMap[planKind]
}
