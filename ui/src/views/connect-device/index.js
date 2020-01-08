import React from 'react'
import Image from 'react-graceful-image'

import Box from '_components/box'
import Link from '_components/link'
import Text from '_components/text'

const ConnectDevice = () => (
  <Box type="grid-12">
    <Box justifyContent="center">
      <Text textStyle="title">Connect device</Text>
    </Box>

    <Box
      padding={5}
      gridColumn="4 / span 6"
      flexDirection="column"
      alignItems="center"
      border="discrete"
    >
      <Image src="path_to_image" width="100" height="100" />

      <Text textStyle="subtitle" py={3}>
        Device details
      </Text>
      <Text textStyle="description" textAlign="center" width="60%" mb={3}>
        In order to connect your device, download the Wireguard configuration file below, or scan
        the QR code.
      </Text>
    </Box>

    <Box justifyContent="center" gridColumn="4 / span 6">
      <Link color="primary" to="/">
        Back
      </Link>
    </Box>
  </Box>
)

export default ConnectDevice
