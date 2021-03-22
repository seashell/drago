import { useMutation } from '@apollo/client'
import { useNavigate } from '@reach/router'
import { useFormik } from 'formik'
import * as _ from 'lodash'
import React from 'react'
import styled from 'styled-components'
import * as Yup from 'yup'
import Box from '_components/box'
import Button from '_components/button'
import Flex from '_components/flex'
import TextInput from '_components/inputs/text-input'
import Link from '_components/link'
import { Dragon as Spinner } from '_components/spinner'
import Text from '_components/text'
import { CREATE_NETWORK } from '_graphql/mutations'
import FormikState from '_utils/formik-utils'
import { withValidityIndicator } from '_utils/hocs'
import { useToast } from '_utils/toast-provider'

const Container = styled(Flex)`
  flex-direction: column;
`

const NewNetwork = () => {
  const navigate = useNavigate()
  const { success, error } = useToast()

  const formik = useFormik({
    initialValues: {
      name: '',
      addressRange: '',
    },
    validationSchema: Yup.object().shape({
      name: Yup.string().required().nullable(),
      addressRange: Yup.string().required().nullable(),
    }),
  })

  const [createNetwork, createNetworkMutation] = useMutation(CREATE_NETWORK, {
    variables: {},
    onCompleted: () => {
      success('Network created')
      navigate('/ui/networks/')
    },
  })

  const handleSaveButtonClick = () => {
    formik.validateForm().then((errors) => {
      if (_.isEmpty(errors)) {
        createNetwork({ variables: { ...formik.values } })
          .then(() => {
            navigate('/ui/networks')
          })
          .catch(() => {})
      } else {
        error('Form has errors')
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
        <Text textStyle="detail">Name *</Text>
        {withValidityIndicator(
          <TextInput
            name="name"
            {...formik.getFieldProps('name')}
            placeholder="a-network"
            height="40px"
          />,
          formik.errors.name
        )}
        <Text textStyle="detail" mt={4}>
          Address Range *
        </Text>
        {withValidityIndicator(
          <TextInput
            name="addressRange"
            {...formik.getFieldProps('addressRange')}
            placeholder="10.0.8.0/24"
            height="40px"
          />,
          formik.errors.addressRange
        )}
        <Button width="100%" borderRadius={3} mt={4} mb={3} onClick={handleSaveButtonClick}>
          Save
        </Button>
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
