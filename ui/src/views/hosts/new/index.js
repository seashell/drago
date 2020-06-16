import React from 'react'
import PropTypes from 'prop-types'
import styled from 'styled-components'
import * as _ from 'lodash'
import * as Yup from 'yup'

import { useFormik } from 'formik'
import { useNavigate } from '@reach/router'

import { useMutation } from 'react-apollo'
import { CREATE_HOST } from '_graphql/actions/hosts'

import Flex from '_components/flex'
import Box from '_components/box'
import Link from '_components/link'
import Text from '_components/text'
import Button from '_components/button'
import TextInput from '_components/inputs/text-input'
import { Dragon as Spinner } from '_components/spinner'
import toast from '_components/toast'

import { withValidityIndicator } from '_hocs/'
import FormikState from '_utils/formik-utils'

const Container = styled(Flex)`
  flex-direction: column;
`

const NewHostView = ({ networkId }) => {
  const navigate = useNavigate()

  const formik = useFormik({
    initialValues: {
      name: null,
      advertiseAddress: null,
    },
    validationSchema: Yup.object().shape({
      name: Yup.string()
        .required()
        .nullable(),
      advertiseAddress: Yup.string().nullable(),
    }),
  })

  const [createHost, createHostMutation] = useMutation(CREATE_HOST, {
    variables: { networkId },
    onCompleted: () => {
      toast.success('Host created')
      navigate('/ui/hosts', { replace: true })
    },
    onError: () => {
      toast.error('Error creating host')
      navigate('/ui/hosts', { replace: true })
    },
  })

  const handleSaveButtonClick = () => {
    formik.validateForm().then(errors => {
      if (_.isEmpty(errors)) {
        createHost({ variables: { ...formik.values } })
      } else {
        toast.error('Form has errors')
      }
    })
  }

  const isLoading = createHostMutation.loading

  if (isLoading) {
    return <Spinner />
  }

  return (
    <Container>
      <Text textStyle="title" mb={4}>
        New host
      </Text>
      {
        <Box flexDirection="column">
          <form>
            <Text my={3}>Name</Text>
            {withValidityIndicator(
              <TextInput name="name" {...formik.getFieldProps('name')} placeholder="host-1" />,
              formik.errors.name
            )}
            <Text my={3}>Advertise address</Text>
            {withValidityIndicator(
              <TextInput
                name="advertiseAddress"
                {...formik.getFieldProps('advertiseAddress')}
                placeholder="wireguard.domain.io"
              />,
              formik.errors.advertiseAddress
            )}
            <Button width="100%" borderRadius={3} mt={3} mb={3} onClick={handleSaveButtonClick}>
              Save
            </Button>
          </form>
          <Link mx="auto" to="/ui/hosts">
            Cancel
          </Link>
          <FormikState {...formik} />
        </Box>
      }
    </Container>
  )
}

NewHostView.propTypes = {
  networkId: PropTypes.string.isRequired,
}

export default NewHostView
