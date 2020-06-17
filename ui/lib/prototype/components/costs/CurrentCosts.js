import PropTypes from 'prop-types'
import Link from 'next/link'
import { Card, Col, Icon, Row, Statistic, Typography } from 'antd'
const { Paragraph, Text } = Typography

import IconTooltip from '../../../components/utils/IconTooltip'

const CurrentCosts = ({ team, currentCost, predictedCost, predictedCostChangePercent }) => (
  <Row gutter={16}>
    <Col span={10}>
      <Card bordered={false}>
        <Paragraph style={{ fontSize: '16px', marginBottom: 0 }} type="secondary">Current costs for June 2020</Paragraph>
        <Text strong style={{ fontSize: '60px', marginRight: '5px' }}>&pound;</Text><Text style={{ fontSize: '60px' }}>{currentCost}</Text>

        {team ? (
          <Paragraph>
            <Link href={`/prototype/teams/${team}/costs/history`}>
              <a style={{ fontSize: '14px', textDecoration: 'underline' }}>See historical cost</a>
            </Link>
          </Paragraph>
        ) : null}

      </Card>
    </Col>
    <Col span={14}>
      <Card bordered={false}>
        <Paragraph style={{ fontSize: '16px', marginBottom: 0 }} type="secondary">
          Total costs for June 2020 (predicted) <IconTooltip icon="info-circle" text="This figure is projected from the usage so far this month, it could increase or decrease with usage changes." />
        </Paragraph>
        <Text strong style={{ fontSize: '60px', marginRight: '5px' }}>&pound;</Text><Text style={{ fontSize: '60px' }}>{predictedCost}</Text>
        <Statistic
          title="compared to May 2020"
          value={predictedCostChangePercent}
          precision={1}
          prefix={predictedCostChangePercent >= 0 ? <Icon type="arrow-up" /> : <Icon type="arrow-down" />}
          suffix="%"
          style={{ display: 'inline-block', marginLeft: '20px' }}
        />
      </Card>
    </Col>
  </Row>
)

CurrentCosts.propTypes = {
  team: PropTypes.string,
  currentCost: PropTypes.number.isRequired,
  predictedCost: PropTypes.number.isRequired,
  predictedCostChangePercent: PropTypes.number.isRequired
}

export default CurrentCosts
