import KoreApi from '../kore-api'
import V1Allocation from '../kore-api/model/V1Allocation'
import V1AllocationSpec from '../kore-api/model/V1AllocationSpec'
import V1ObjectMeta from '../kore-api/model/V1ObjectMeta'
import V1Ownership from '../kore-api/model/V1Ownership'
import config from '../../config' 

export default class AllocationHelpers {
  static getAllocationNameForResource = (resource) => {
    // @TODO: Include kind in allocation name so they're predictabley-named but don't clash
    return resource.metadata.name
  }

  static generateAllocation = ({ resourceToAllocate, teams, name, summary }) => {
    if (!resourceToAllocate) {
      throw 'Must specify resourceToAllocate'
    }
    const resource = new V1Allocation()
    resource.setApiVersion('config.kore.appvia.io/v1')
    resource.setKind('Allocation')

    const meta = new V1ObjectMeta()
    meta.setName(AllocationHelpers.getAllocationNameForResource(resourceToAllocate))
    meta.setNamespace(config.kore.koreAdminTeamName)
    resource.setMetadata(meta)

    const spec = new V1AllocationSpec()
    spec.setName(name ? name : resourceToAllocate.metadata.name)
    spec.setSummary(summary ? summary : `Allocation of ${resourceToAllocate.metadata.name}`)
    spec.setTeams(teams && teams.length > 0 ? teams : ['*'])

    const resGroupVersion = resourceToAllocate.apiVersion.split('/')
    const owner = new V1Ownership()
    owner.setGroup(resGroupVersion[0])
    owner.setVersion(resGroupVersion[1])
    owner.setKind(resourceToAllocate.kind)
    owner.setName(resourceToAllocate.metadata.name)
    owner.setNamespace(resourceToAllocate.metadata.namespace)
    spec.setResource(owner)

    resource.setSpec(spec)
    return resource
  }

  static getAllocationForResource = async (resource) => {
    if (!resource || !resource.metadata || !resource.metadata.name) {
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
    const allocation = AllocationHelpers.generateAllocation({ resourceToAllocate, teams, name, summary })
    return await (await KoreApi.client()).UpdateAllocation(config.kore.koreAdminTeamName, allocation.metadata.name, allocation)
  }

  static removeAllocation = async (resourceToDeallocate) => {
    return await (await KoreApi.client()).RemoveAllocation(config.kore.koreAdminTeamName, AllocationHelpers.getAllocationNameForResource(resourceToDeallocate))
  }
}