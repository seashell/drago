import React from 'react'
import styled from 'styled-components'
import * as _ from 'lodash'
import * as Yup from 'yup'

import { useFormik } from 'formik'
import { useNavigate } from '@reach/router'

import { useMutation } from 'react-apollo'
import { CREATE_NETWORK } from '_graphql/actions/networks'

import toast from '_components/toast'
import Flex from '_components/flex'
import Box from '_components/box'
import Link from '_components/link'
import Text from '_components/text'
import Button from '_components/button'
import TextInput from '_components/inputs/text-input'
import { Dragon as Spinner } from '_components/spinner'

import { withValidityIndicator } from '_hocs/'

import FormikState from '_utils/formik-utils'

const Container = styled(Flex)`
  flex-direction: column;
`

const NewNetwork = () => {
  const navigate = useNavigate()

  const formik = useFormik({
    initialValues: {
      name: null,
      ipAddressRange: null,
    },
    validationSchema: Yup.object().shape({
      name: Yup.string()
        .required()
        .nullable(),
      ipAddressRange: Yup.string()
        .required()
        .nullable(),
    }),
  })

  const [createNetwork, createNetworkMutation] = useMutation(CREATE_NETWORK, {
    variables: {},
    onCompleted: () => {
      toast.success('Network created')
      navigate('/ui/networks/', { replace: true })
    },
    onError: () => {
      toast.error('Error creating network')
      navigate('/ui/networks', { replace: true })
    },
  })

  const handleSaveButtonClick = () => {
    formik.validateForm().then(errors => {
      if (_.isEmpty(errors)) {
        createNetwork({ variables: { ...formik.values } })
      } else {
        toast.error('Form has errors')
      }
    })
  }

  const isLoading = createNetworkMutation.loading

  if (isLoading) {
    return <Spinner />
  }

  return (
    <Container>
      <Text textStyle="title" mb={4}>
        New network
      </Text>

      <Box flexDirection="column">
        <form>
          <Text my={3}>Name *</Text>
          {withValidityIndicator(
            <TextInput name="name" {...formik.getFieldProps('name')} placeholder="a-network" />,
            formik.errors.name
          )}
          <Text my={3}>Address range *</Text>
          {withValidityIndicator(
            <TextInput
              name="ipAddressRange"
              {...formik.getFieldProps('ipAddressRange')}
              placeholder="10.0.8.0/24"
            />,
            formik.errors.ipAddressRange
          )}
          <Button width="100%" borderRadius={3} mt={3} mb={3} onClick={handleSaveButtonClick}>
            Save
          </Button>
        </form>
        <FormikState {...formik} />
      </Box>

      <Box justifyContent="center" gridColumn="4 / span 6">
        <Link color="primary" to="../">
          Cancel
        </Link>
      </Box>
    </Container>
  )
}

export default NewNetwork
