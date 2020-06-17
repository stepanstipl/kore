import { Table } from 'antd'

const columns = (cost) => [
  {
    title: 'Type',
    dataIndex: 'type',
    key: 'type',
  },
  {
    title: 'Name',
    dataIndex: 'name',
    key: 'name',
  },
  {
    title: 'Account / project',
    dataIndex: 'account',
    key: 'account',
  },
  {
    title: cost,
    dataIndex: 'cost',
    key: 'cost',
    align: 'right'
  }
]

const clustersData = [
  { key: '1', type: 'EKS cluster', name: 'eks-demo-notprod', account: 'kore-demo-notprod', cost: '£54.62', children: [
    { key: '1a', type: 'EKS control plane', name: '', account: 'kore-demo-notprod', cost: '£20.00' },
    { key: '1b', type: 'EKS node group', name: 'eks-demo-notprod-default', account: 'kore-demo-notprod', cost: '£25.00' },
    { key: '1c', type: 'Namespace', name: 'development', account: 'kore-demo-notprod', cost: '£4.62', children: [
      { key: '1c1', type: 'Amazon SQS', name: 'dev-sqs', account: 'kore-demo-notprod', cost: '£2.02' },
      { key: '1c2', type: 'Amazon RDS for MySQL', name: 'dev-rds', account: 'kore-demo-notprod', cost: '£2.60' }
    ] },
    { key: '1d', type: 'Namespace', name: 'qa', account: 'kore-demo-notprod', cost: '£5.00' },
  ] },
  { key: '2', type: 'EKS cluster', name: 'eks-demo-prod', account: 'kore-demo-prod', cost: '£150.00', children: [
    { key: '2a', type: 'EKS control plane', name: '', account: 'kore-demo-prod', cost: '£20.00' },
    { key: '2b', type: 'EKS node group', name: 'eks-demo-prod-default', account: 'kore-demo-prod', cost: '£110.00' },
    { key: '2c', type: 'Namespace', name: 'production', account: 'kore-demo-prod', cost: '£20.00', children: [
      { key: '2c1', type: 'Amazon SQS', name: 'prod-sqs', account: 'kore-demo-prod', cost: '£5.00' },
      { key: '2c2', type: 'Amazon RDS for MySQL', name: 'prod-rds', account: 'kore-demo-prod', cost: '£15.00' }
    ] }
  ] }
]

const TeamMonthlyCostTable = () => (
  <Table style={{ marginTop: '10px', marginBottom: '20px' }} showHeader={true} pagination={false} columns={columns('£204.62')} dataSource={clustersData} />
)

export default TeamMonthlyCostTable
