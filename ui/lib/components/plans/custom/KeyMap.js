import * as React from 'react'
import PropTypes from 'prop-types'
import { Table, Input } from 'antd'

export default class KeyMap extends React.Component {
  static propTypes = {
    value: PropTypes.object,
    property: PropTypes.object.isRequired,
    editable: PropTypes.bool.isRequired,
    onChange: PropTypes.func.isRequired,
  }

  onKeyChange = (key, newKey, value) => {
    const newVal = this.props.value ? { ...this.props.value } : {}
    if (newKey && newKey.length > 0) {
      newVal[newKey] = value
    }

    // If this was an existing key, remove it:
    if (key && key.length > 0) {
      delete newVal[key]
    }

    this.props.onChange(newVal)
  }

  onValueChange = (key, value) => {
    const newVal = this.props.value ? { ...this.props.value, [key]: value } : { [key]: value }
    this.props.onChange(newVal)
  }

  render() {
    const { value, editable } = this.props
    const keys  = value ? Object.keys(value) : []

    if (!editable && keys.length === 0) {
      return <>None set</>
    }

    const columns = [
      { title: 'Key', dataIndex: 'key', key: 'key', width: '35%', render: function renderAction(_,kv) { 
        return <Input value={kv.key} readOnly={!editable} onChange={(e) => this.onKeyChange(kv.key, e.target.value, kv.value)} /> 
      }.bind(this) },
      { title: 'Value', dataIndex: 'value', key: 'value', width: '65%', render: function renderAction(_,kv) { 
        return <Input value={kv.value} readOnly={ !editable || !kv.key || kv.key.length === 0 } onChange={(e) => this.onValueChange(kv.key, e.target.value)} /> 
      }.bind(this) },
    ]
    const rows = keys.map((k, i) => {
      return {
        ind: i,
        key: k,
        value: value[k]
      }
    })

    if (editable) {
      rows.push({ ind: rows.length, key: '', value: '' })
    }
    
    return (
      <Table 
        size="small"
        pagination={false} 
        dataSource={rows} 
        columns={columns} 
        rowKey={r => r.ind}
      />
    )
  }
}