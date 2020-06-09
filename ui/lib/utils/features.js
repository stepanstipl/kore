import getConfig from 'next/config'
const { publicRuntimeConfig } = getConfig()

export const KoreFeatures = {
  SERVICES: 'services',
  APPLICATION_SERVICES: 'application_services'
}

export function featureEnabled(feature) {
  return Boolean(publicRuntimeConfig.featureGates[feature])
}
