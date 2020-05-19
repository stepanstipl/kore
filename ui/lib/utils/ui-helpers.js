module.exports = {
  statusColorMap: {
    'Success': 'green',
    'Pending': 'orange'
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
