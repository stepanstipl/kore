import React from 'react'
import { Tabs } from 'antd'
const { TabPane } = Tabs

import Breadcrumb from '../../../../../lib/components/layout/Breadcrumb'
import EventNotificationsTab from '../../../../../lib/prototype/components/teams/settings/EventNotificationsTab'
import NotificationIntegrationsTab from '../../../../../lib/prototype/components/teams/settings/NotificationIntegrationsTab'

class TeamNotificationsSettings extends React.Component {

  render() {
    return (
      <>
        <Breadcrumb items={[{ text: 'Proto', link: '/prototype/teams/proto', href: '/prototype/teams/proto' }, { text: 'Notification settings' }]}/>

        <Tabs defaultActiveKey="notifications" tabBarStyle={{ marginBottom: '20px' }}>
          <TabPane key="notifications" tab="Notifications" forceRender={true}>
            <EventNotificationsTab />
          </TabPane>
          <TabPane key="integrations" tab="Integrations" forceRender={true}>
            <NotificationIntegrationsTab />
          </TabPane>
        </Tabs>
      </>
    )
  }
}

export default TeamNotificationsSettings
