import React from 'react'
import DeviceCard from '_components/device-card'
import List from '_components/list'

const DevicesView = () => (
  <List>
    <DeviceCard id="123" name="device-1" />
    <DeviceCard id="123" name="device-2" />
    <DeviceCard id="123" name="device-3" />
  </List>
)

export default DevicesView
