import Link from 'next/link'
import { Typography, List, Button } from 'antd'
const { Title, Text } = Typography

const prototypeList = [{
  title: 'Setup wizard',
  description: 'Setup Kore cloud access and team project automation settings for GCP',
  path: '/prototype/setup/kore'
}, {
  title: 'Project automation settings',
  description: 'Configure the team project automation settings, within the configure cloud page',
  path: '/prototype/configure/cloud'
}, {
  title: 'Team costs',
  description: 'Page showing cloud costs for a team, select "Team costs" from the settings downdown in the top-right corner',
  path: '/prototype/teams/costs-demo'
}, {
  title: 'Organisation reports',
  description: 'Adding a reports section for org-level security and costs reports. This would be accessed from a "Reports" link on left-side menu.',
  path: '/prototype/reports'
}]

const PrototypeIndex = () => (
  <>
    <Title>Prototypes</Title>
    <List
      dataSource={prototypeList}
      renderItem={item => (
        <List.Item actions={[<Link key="view" href={item.path}><Button type="primary">View</Button></Link>]}>
          <List.Item.Meta
            title={<Text style={{ fontSize: '20px', fontWeight: '600' }}>{item.title}</Text>}
            description={item.description}
          />
        </List.Item>
      )}
    />
  </>
)

PrototypeIndex.staticProps = {
  title: 'Kore prototypes',
  hideSider: true
}

export default PrototypeIndex
