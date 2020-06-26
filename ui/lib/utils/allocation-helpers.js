import KoreApi from '../kore-api'
import config from '../../config' 

export default class AllocationHelpers {
  static getAllocationNameForResource = (resource) => {
    return `${resource.kind.toLowerCase()}-${resource.metadata.name}`
  }

  static getAllocationForResource = async (resource) => {
    if (!resource || !resource.metadata || !resource.metadata.name || !resource.kind) {
      return null
    }
    return await (await KoreApi.client()).GetAllocation(config.kore.koreAdminTeamName, AllocationHelpers.getAllocationNameForResource(resource))
  }

  /**
   * From a list of allocations (typically returned from api.ListAllocations), this will return the first allocation 
   * which has a resource kind, group, version, name and namespace matching those of the passed-in resource.
   */
  static findAllocationForResource = (allocationList, resource) => {
    if (!resource) {
      return null
    }
    const resGroupVersion = resource.apiVersion.split('/')
    return allocationList.items.find((a) => 
      a.spec.resource.kind === resource.kind && 
      a.spec.resource.group === resGroupVersion[0] && 
      a.spec.resource.version === resGroupVersion[1] && 
      a.spec.resource.name === resource.metadata.name && 
      a.spec.resource.namespace === resource.metadata.namespace)
  }

  static isAllTeams = (allocation) => {
    return allocation && allocation.spec && allocation.spec.teams && allocation.spec.teams.find((t) => t === '*')
  }

  static storeAllocation = async ({ resourceToAllocate, teams, name, summary }) => {
    const resourceName = AllocationHelpers.getAllocationNameForResource(resourceToAllocate)
    const allocationResource = KoreApi.resources().generateAllocationResource(resourceName, resourceToAllocate, teams, name, summary)

    return await (await KoreApi.client()).UpdateAllocation(config.kore.koreAdminTeamName, allocationResource.metadata.name, allocationResource)
  }

  static removeAllocation = async (resourceToDeallocate) => {
    return await (await KoreApi.client()).RemoveAllocation(config.kore.koreAdminTeamName, AllocationHelpers.getAllocationNameForResource(resourceToDeallocate))
  }
}