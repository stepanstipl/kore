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
  title: 'Security',
  description: 'Review the security posture of all Kore-provisioned clusters and plans',
  path: '/prototype/security'
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
