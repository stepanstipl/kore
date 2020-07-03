import PlanOptionClusterUsers from './PlanOptionClusterUsers'
import PlanOptionEKSNodeGroups from './PlanOptionEKSNodeGroups'
import PlanOptionGKENodePools from './PlanOptionGKENodePools'
import PlanOptionGKEReleaseChannel from './PlanOptionGKEReleaseChannel'
import PlanOptionGKEVersion from './PlanOptionGKEVersion'
import PlanOptionClusterRegion from './PlanOptionClusterRegion'

export default class CustomPlanOptionRegistry {
  static controls = {
    'cluster': {
      'GKE': {
        'clusterUsers': function clusterUsers(props) {
          return <PlanOptionClusterUsers {...props} />
        },
        'nodePools': function nodePools(props) {
          return <PlanOptionGKENodePools {...props} />
        },
        'releaseChannel': function releaseChannel(props) {
          return <PlanOptionGKEReleaseChannel {...props} />
        },
        'version': function releaseChannel(props) {
          return <PlanOptionGKEVersion {...props} />
        },
        'region': function region(props) {
          return <PlanOptionClusterRegion {...props} />
        }
      },
      'EKS': {
        'clusterUsers': function clusterUsers(props) { 
          return <PlanOptionClusterUsers {...props} /> 
        },
        'nodeGroups': function nodeGroups(props) { 
          return <PlanOptionEKSNodeGroups {...props} /> 
        },
        'region': function region(props) {
          return <PlanOptionClusterRegion {...props} />
        }
      }
    }
  }

  static getCustomPlanOption = (planType, planKind, fieldName, props) => {
    if (!CustomPlanOptionRegistry.controls[planType] || 
      !CustomPlanOptionRegistry.controls[planType][planKind] || 
      !CustomPlanOptionRegistry.controls[planType][planKind][fieldName]) {
      return null
    }
    return CustomPlanOptionRegistry.controls[planType][planKind][fieldName](props)
  }
}