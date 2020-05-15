import PlanOptionClusterUsers from './PlanOptionClusterUsers'

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