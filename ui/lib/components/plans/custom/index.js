import PlanOptionClusterUsers from './PlanOptionClusterUsers'
import PlanOptionEKSNodeGroups from './PlanOptionEKSNodeGroups'

export default class CustomPlanOptionRegistry {
  static controls = {
    'cluster': {
      'GKE': {
        'clusterUsers': function clusterUsers(props) {
          return <PlanOptionClusterUsers {...props} />
        }
      },
      'EKS': {
        'clusterUsers': function clusterUsers(props) { 
          return <PlanOptionClusterUsers {...props} /> 
        },
        'nodeGroups': function nodeGroups(props) { 
          return <PlanOptionEKSNodeGroups {...props} /> 
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