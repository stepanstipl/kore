import React from 'react'
import Link from 'next/link'
import { Card, Col, Icon, Row, Statistic, Typography } from 'antd'
const { Paragraph, Text, Title } = Typography

import Breadcrumb from '../../../../../lib/components/layout/Breadcrumb'
import MonthlyCostTables from '../../../../../lib/prototype/components/costs/MonthlyCostTables'
import IconTooltip from '../../../../../lib/components/utils/IconTooltip'

class TeamCosts extends React.Component {

  render() {
    return (
      <>
        <Breadcrumb items={[{ text: 'Demo' }, { text: 'Team costs' }]}/>
        <Row gutter={16}>
          <Col span={10}>
            <Card bordered={false}>
              <Paragraph style={{ fontSize: '16px', marginBottom: 0 }} type="secondary">Current costs for June 2020</Paragraph>
              <Text strong style={{ fontSize: '60px', marginRight: '5px' }}>&pound;</Text><Text style={{ fontSize: '60px' }}>254.62</Text>
              <Paragraph>
                <Link href="/prototype/teams/demo/costs/history">
                  <a style={{ fontSize: '14px', textDecoration: 'underline' }}>See historical cost</a>
                </Link>
              </Paragraph>
            </Card>
          </Col>
          <Col span={14}>
            <Card bordered={false}>
              <Paragraph style={{ fontSize: '16px', marginBottom: 0 }} type="secondary">
                Total costs for June 2020 (predicted) <IconTooltip icon="info-circle" text="This figure is projected from the usage so far this month, it could increase or decrease with usage changes." />
              </Paragraph>
              <Text strong style={{ fontSize: '60px', marginRight: '5px' }}>&pound;</Text><Text style={{ fontSize: '60px' }}>734.72</Text>
              <Statistic
                title="compared to May 2020"
                value={5.5}
                precision={1}
                prefix={<Icon type="arrow-up" />}
                suffix="%"
                style={{ display: 'inline-block', marginLeft: '20px' }}
              />
            </Card>
          </Col>
        </Row>
        <Title level={3}>Costs breakdown for June 2020</Title>
        <Paragraph style={{ marginBottom: '20px' }} type="secondary">This shows the breakdown of the incurred costs for this month, the predicted total expenditure is not broken down here.</Paragraph>
        <MonthlyCostTables />
      </>
    )
  }
}

export default TeamCosts
