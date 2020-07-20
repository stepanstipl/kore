import getConfig from 'next/config'
const { publicRuntimeConfig } = getConfig()

export function getProviderCloudInfo(planKind) {
  return publicRuntimeConfig.clusterProviderCloudMap[planKind]
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
