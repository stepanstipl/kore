const integrations = {
  items: [
    {
      apiVersion: 'integrations.kore.appvia.io/v1',
      kind: 'SlackIntegration',
      metadata: {
        name: 'proto-slack-integration',
        namespace: 'proto',
      },
      spec: {
        enabled: true,
        displayName: 'Example slack',
        webhookURL: 'https://ff3fsdfdfdf.slack.com/webhooks'
      }
    },
    {
      apiVersion: 'integrations.kore.appvia.io/v1',
      kind: 'EmailIntegration',
      metadata: {
        name: 'proto-email-integration',
        namespace: 'proto',
      },
      spec: {
        enabled: true
      }
    }
  ]
}

const notifications = {
  items: [
    {
      metadata: {
        name: 'cluster-created'
      },
      spec: {
        event: 'CLUSTER_CREATED',
        channel: 'my-notifications',
        integration: {
          group: 'integrations.kore.appvia.io',
          version: 'v1',
          kind: 'SlackIntegration',
          namespace: 'proto',
          name: 'proto-slack-integration'
        }
      }
    },
    {
      metadata: {
        name: 'cluster-created-2'
      },
      spec: {
        event: 'CLUSTER_CREATED',
        integration: {
          group: 'integrations.kore.appvia.io',
          version: 'v1',
          kind: 'EmailIntegration',
          namespace: 'proto',
          name: 'proto-email-integration'
        }
      }
    },
    {
      metadata: {
        name: 'cluster-deleted'
      },
      spec: {
        event: 'CLUSTER_DELETED'
      }
    },
    {
      metadata: {
        name: 'service-created'
      },
      spec: {
        event: 'SERVICE_CREATED',
        emailAddressList: ['bob@appvia.io', 'alice@appvia.io'],
        integration: {
          group: 'integrations.kore.appvia.io',
          version: 'v1',
          kind: 'EmailIntegration',
          namespace: 'proto',
          name: 'proto-email-integration'
        }
      }
    }
  ]
}

class TeamNotificationData {
  static notifications = notifications
  static integrations = integrations
}

export default TeamNotificationData