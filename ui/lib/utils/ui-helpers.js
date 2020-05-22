module.exports = {
  statusColorMap: {
    'Success': 'green',
    'Pending': 'orange'
  },
  statusIconMap: {
    'Success': 'check-circle',
    'Pending': 'loading',
    'Failure': 'exclamation-circle',
    'Error': 'exclamation-circle',
  },
  clusterProviderIconSrcMap: {
    'GKE': '/static/images/GCP.png',
    'EKS': '/static/images/AWS.png'
  },
  verifiedStatusMessageMap: {
    'Success': 'Verified',
    'Failure': 'Not Verified',
    'Pending': 'Verifying'
  },
  inProgressStatusList: ['Pending', 'Deleting']
}
