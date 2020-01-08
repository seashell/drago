import React from 'react'
import { useFormState } from 'react-use-form-state'
import Image from 'react-graceful-image'

import Box from '_components/box'
import Link from '_components/link'
import Text from '_components/text'
import Button from '_components/button'
import TextInput from '_components/inputs/text-input'

const NewDevice = () => {
  const [formState, { text }] = useFormState()

  return (
    <Box type="grid-12">
      <Box justifyContent="center">
        <Text textStyle="title">New device</Text>
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
          Please provide some details about your new device.
        </Text>
        <TextInput {...text('projectName')} placeholder="Device name" mb={2} />
        <TextInput {...text('allowedIPs')} placeholder="10.0.8.0/24" mb={2} />

        <Button width="100%" borderRadius={3}>
          Continue
        </Button>
      </Box>

      <Box justifyContent="center" gridColumn="4 / span 6">
        <Link color="primary" to="/">
          Cancel
        </Link>
      </Box>
    </Box>
  )
}

export default NewDevice
