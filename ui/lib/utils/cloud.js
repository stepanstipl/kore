import getConfig from 'next/config'
const { publicRuntimeConfig } = getConfig()

/**
 * Get cloud info for the cluster provider
 * @param provider The cluster provider, eg GKE, EKS, AKS
 */
export function getProviderCloudInfo(provider) {
  return publicRuntimeConfig.clusterProviderCloudMap[provider]
}

/**
 * Get cloud info for the cloud, via the cluster provider
 * @param cloud The cloud, eg GCP, AWS, Azure
 * @returns {*}
 */
export function getCloudInfo(cloud) {
  return getProviderCloudInfo(publicRuntimeConfig.clusterProviderMap[cloud.toUpperCase()])
}

/**
 * Remove redundant rules and rule plans from the automated cloud account list
 * @param cloudAccountList the list of automated account rules for a provider eg GKE
 * @param plans Plans specific to the provider eg GKE
 */
export function filterCloudAccountList(cloudAccountList, plans) {
  return cloudAccountList
    // remove rules with no plans associated
    .filter(cloudAccount => cloudAccount.plans.length > 0)
    // remove non-existent plans from rules
    .map(cloudAccount => {
      return { ...cloudAccount, plans: cloudAccount.plans.filter(p => plans.map(p => p.metadata.name).includes(p)) }
    })
}
