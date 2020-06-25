const moment = require('moment')

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

const eventNotifications = {
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

const notifications = {
  items: [
    {
      event: 'CLUSTER_DELETED',
      detail: 'example-cluster-1',
      creationTimestamp: moment().subtract(10, 'minutes').format(),
      acknowledged: false
    },
    {
      event: 'CLUSTER_CREATED',
      detail: 'example-cluster-2',
      creationTimestamp: moment().subtract(31, 'minutes').format(),
      acknowledged: false
    },
    {
      event: 'SERVICE_CREATED',
      detail: 'Amazon SQS proto-message-queue',
      creationTimestamp: moment().subtract(45, 'minutes').format(),
      acknowledged: true
    },
    {
      event: 'SERVICE_CREATED',
      detail: 'Amazon S3 proto-bucket',
      creationTimestamp: moment().subtract(1, 'hour').format(),
      acknowledged: false
    },
    {
      event: 'SERVICE_CREATED',
      detail: 'Amazon SQS proto-message-queue2',
      creationTimestamp: moment().subtract(2, 'hours').format(),
      acknowledged: true
    }
  ]
}

class TeamNotificationData {
  static eventNotifications = eventNotifications
  static integrations = integrations
  static notifications = notifications
}

export default TeamNotificationData