import React, { useEffect } from 'react'
import PropTypes from 'prop-types'
import * as _ from 'lodash'
import * as Yup from 'yup'

import { useFormik } from 'formik'
import { useNavigate, useLocation, useParams } from '@reach/router'

import { useQuery, useMutation } from 'react-apollo'
import { GET_HOST, UPDATE_HOST } from '_graphql/actions/hosts'

import toast from '_components/toast'
import { Dragon as Spinner } from '_components/spinner'
import TextInput from '_components/inputs/text-input'
import Button from '_components/button'
import Text from '_components/text'
import Box from '_components/box'

import { withValidityIndicator } from '_hocs'
import FormikState from '_utils/formik-utils'

const HostAttributesTab = () => {
  const navigate = useNavigate()
  const location = useLocation()
  const urlParams = useParams()

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

  const getHostQuery = useQuery(GET_HOST, {
    variables: { id: urlParams.hostId },
    fetchPolicy: 'cache-and-network',
    onCompleted: data => {
      formik.setValues(
        {
          name: data.result.name,
          advertiseAddress: data.result.advertiseAddress,
        },
        true
      )
    },
    onError: () => {
      toast.error('Error fetching host details')
      navigate('/ui/hosts/')
    },
  })

  const [updateHost, updateHostMutation] = useMutation(UPDATE_HOST, {
    variables: { id: urlParams.hostId, ...formik.values },
    onCompleted: () => {
      toast.success('Host updated')
      getHostQuery.refetch()
    },
    onError: () => {
      toast.error('Error updating host')
      navigate('/ui/hosts/')
    },
  })

  useEffect(() => {
    getHostQuery.refetch()
  }, [location])

  const handleSaveButtonClick = () => {
    formik.validateForm().then(errors => {
      if (_.isEmpty(errors)) {
        updateHost({ id: urlParams.hostId, ...formik.values })
        getHostQuery.refetch()
      } else {
        toast.error('Form has errors')
      }
    })
  }

  const isLoading = getHostQuery.loading || updateHostMutation.loading
  if (isLoading) {
    return <Spinner />
  }

  return (
    <Box flexDirection="column">
      <form>
        <Text my={3}>Name</Text>
        {withValidityIndicator(
          <TextInput name="name" {...formik.getFieldProps('name')} placeholder="wg0" />,
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
        <Button width="100%" borderRadius={3} mt={3} mb={4} onClick={handleSaveButtonClick}>
          Save
        </Button>
      </form>
      <FormikState {...formik} />
    </Box>
  )
}

HostAttributesTab.propTypes = {
  hostId: PropTypes.string.isRequired,
}

export default HostAttributesTab
