import getConfig from 'next/config'
const { publicRuntimeConfig } = getConfig()

export const KoreFeatures = {
  APPLICATION_SERVICES: 'application_services',
  COSTS: 'costs',
  MONITORING_SERVICES: 'monitoring_services',
  SERVICES: 'services'
}

export function featureEnabled(feature) {
  return Boolean(publicRuntimeConfig.featureGates[feature])
}
