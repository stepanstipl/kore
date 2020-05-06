import AllocationHelpers from '../../../../lib/utils/allocation-helpers'
import config from '../../../../config'

describe('AllocationHelpers', () => {
  describe('#generateAllocation', () => {
    const resource = {
      apiVersion: 'testgroup.appvia.io/v1alpha1',
      kind: 'TestKind',
      metadata: {
        name: 'test-resource',
        namespace: 'test-team'
      },
    }

    it('should return an allocation when called with a resource', () => {
      const allocation = AllocationHelpers.generateAllocation({ resourceToAllocate: resource, teams: [] })
      expect(allocation).toBeDefined()
      expect(allocation.metadata.name).toEqual(`${resource.kind.toLowerCase()}-${resource.metadata.name}`)
      expect(allocation.metadata.namespace).toEqual(config.kore.koreAdminTeamName)
      expect(allocation.spec.resource.group).toEqual('testgroup.appvia.io')
      expect(allocation.spec.resource.version).toEqual('v1alpha1')
      expect(allocation.spec.resource.kind).toEqual(resource.kind)
      expect(allocation.spec.resource.namespace).toEqual(resource.metadata.namespace)
    })

    it('should set the teams to all if no teams provided', () => {
      const allocation = AllocationHelpers.generateAllocation({ resourceToAllocate: resource, teams: [] })
      expect(allocation.spec.teams).toEqual(['*'])
    })

    it('should set the teams to specific value if teams provided', () => {
      const allocation = AllocationHelpers.generateAllocation({ resourceToAllocate: resource, teams: ['team1','team2'] })
      expect(allocation.spec.teams).toEqual(['team1','team2'])
    })

    it('should set the spec name and summary if provided', () => {
      const allocation = AllocationHelpers.generateAllocation({ resourceToAllocate: resource, teams: ['team1','team2'], name: 'horse', summary: 'cheese' })
      expect(allocation.spec.name).toEqual('horse')
      expect(allocation.spec.summary).toEqual('cheese')
    })

    it('should set defaults for the spec name and description if not provided', () => {
      const allocation = AllocationHelpers.generateAllocation({ resourceToAllocate: resource, teams: ['team1','team2'] })
      expect(allocation.spec.name).toEqual('test-resource')
      expect(allocation.spec.summary).toEqual('Allocation of test-resource')
    })
  })
})
