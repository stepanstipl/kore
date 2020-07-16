import getConfig from 'next/config'
const { publicRuntimeConfig } = getConfig()

export function getProviderCloudInfo(planKind) {
  return publicRuntimeConfig.clusterProviderCloudMap[planKind]
}
