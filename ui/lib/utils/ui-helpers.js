module.exports = {
  statusColorMap: {
    'Success': 'green',
    'Pending': 'orange'
  },
  clusterProviderIconSrcMap: {
    'GKE': '/static/images/GCP.png',
    'EKS': '/static/images/AWS.png',
    'Kore': '/static/images/appvia-colour.svg'
  },
  verifiedStatusMessageMap: {
    'Success': 'Verified',
    'Failure': 'Not Verified',
    'Pending': 'Verifying'
  },
  inProgressStatusList: ['Pending', 'Deleting']
}
