/* eslint-disable react/prop-types */
import React, { useEffect } from 'react'
import styled from 'styled-components'
import * as _ from 'lodash'
import * as Yup from 'yup'

import { useLocation, useNavigate, useParams } from '@reach/router'
import { useFormik } from 'formik'

import { useQuery, useMutation } from 'react-apollo'
import { GET_INTERFACE, UPDATE_INTERFACE } from '_graphql/actions/interfaces'
import { GET_NETWORKS } from '_graphql/actions/networks'

import { icons } from '_assets/'
import Flex from '_components/flex'
import Box from '_components/box'
import Link from '_components/link'
import Text from '_components/text'
import Button from '_components/button'
import TextInput from '_components/inputs/text-input'
import SelectInput from '_components/inputs/select-input'
import IconButton from '_components/icon-button'
import { Dragon as Spinner } from '_components/spinner'
import toast from '_components/toast'

import { withValidityIndicator } from '_hocs/'
import FormikState from '_utils/formik-utils'

const Container = styled(Flex)`
  flex-direction: column;
`

const SearchResult = styled(Box)`
  cursor: pointer;
  height: 48px;
  padding: 8px;
  align-items: center;
  :hover {
    background: #eee;
  }
`

const IconContainer = styled(Box).attrs({
  display: 'flex',
  bg: 'neutralLighter',
  borderRadius: '4px',
})`
  position: relative;
  button {
    margin-right: auto;
  }
  align-items: center;
  justify-content: center;
`

const NetworkOption = ({ innerRef, innerProps, ...props }) => (
  <SearchResult innerRef={innerRef} {...innerProps} {...props}>
    <IconContainer height="32px" width="32px" mr="12px">
      <IconButton ml="auto" icon={<icons.Network />} />
    </IconContainer>
    {props.data.name} ({props.data.ipAddressRange})
  </SearchResult>
)

const SelectedNetworkValue = ({ innerRef, innerProps, ...props }) => (
  <div innerRef={innerRef} {...innerProps} {...props}>
    {props.data.name} ({props.data.ipAddressRange})
  </div>
)

const filterNetworks = (option, searchText) => {
  if (option.data.name.toLowerCase().includes(searchText.toLowerCase())) {
    return true
  }
  return false
}

const InterfaceDetailsView = () => {
  const navigate = useNavigate()
  const location = useLocation()
  const urlParams = useParams()

  const { hostId } = location.state || {}

  if (hostId === undefined) {
    navigate(`/ui/hosts`)
  }

  const formik = useFormik({
    initialValues: {
      name: null,
      networkId: null,
      ipAddress: null,
      listenPort: null,
    },
    validationSchema: Yup.object().shape({
      name: Yup.string()
        .required()
        .nullable(),
      networkId: Yup.string().nullable(),
      ipAddress: Yup.string().nullable(),
      listenPort: Yup.number()
        .positive()
        .nullable(),
    }),
  })

  const getNetworksQuery = useQuery(GET_NETWORKS)

  const getInterfaceQuery = useQuery(GET_INTERFACE, {
    variables: { id: urlParams.interfaceId },
    fetchPolicy: 'cache-and-network',
    onCompleted: data => {
      formik.setValues(
        {
          name: data.result.name,
          networkId: data.result.networkId,
          ipAddress: data.result.ipAddress,
          listenPort: data.result.listenPort,
        },
        true
      )
    },
    onError: () => {
      toast.error('Error fetching interface details')
      navigate(-1)
    },
  })

  const [updateInterface, updateInterfaceMutation] = useMutation(UPDATE_INTERFACE, {
    variables: {},
    onCompleted: () => {
      toast.success('Interface updated')
      navigate(`/ui/hosts/${hostId}/interfaces`)
    },
    onError: () => {
      toast.error('Error updating interface')
    },
  })

  useEffect(() => {
    window.scrollTo(0, 0)
    formik.resetForm()
    getNetworksQuery.refetch()
    getInterfaceQuery.refetch()
  }, [location])

  const handleSaveButtonClick = () => {
    formik.validateForm().then(errors => {
      if (_.isEmpty(errors)) {
        updateInterface({ variables: { id: urlParams.interfaceId, ...formik.values } })
        navigate(`/ui/hosts/${hostId}/interfaces`)
      } else {
        toast.error('Form has errors')
      }
    })
  }

  const isLoading =
    getNetworksQuery.loading || getInterfaceQuery.loading || updateInterfaceMutation.loading

  if (isLoading) {
    return <Spinner />
  }

  const networkOptions = !isLoading ? getNetworksQuery.data.result.items : []

  return (
    <Container>
      <Box alignItems="center" mb={4}>
        <IconContainer mr="12px">
          <IconButton ml="auto" size={32} icon={<icons.Interface />} />
        </IconContainer>
        <Text textStyle="title">{getInterfaceQuery.data.result.name}</Text>
      </Box>
      {isLoading ? (
        <Spinner />
      ) : (
        <Box flexDirection="column">
          <form>
            <Text my={3}>Network</Text>
            {withValidityIndicator(
              <SelectInput
                id="networkId"
                name="networkId"
                value={networkOptions.find(el => el.id === formik.values.networkId)}
                onChange={value =>
                  formik.setFieldValue('networkId', value !== null ? value.id : null)
                }
                placeholder="Select network"
                options={networkOptions}
                optionComponent={NetworkOption}
                singleValueComponent={SelectedNetworkValue}
                filterOption={filterNetworks}
                isClearable
              />,
              formik.errors.networkId && formik.errors.networkId
            )}

            <Text my={3}>Name *</Text>
            {withValidityIndicator(
              <TextInput name="name" {...formik.getFieldProps('name')} placeholder="wg0" />,
              formik.errors.name && formik.errors.name
            )}

            <Text my={3}>IP Address</Text>
            {withValidityIndicator(
              <TextInput
                name="ipAddress"
                {...formik.getFieldProps('ipAddress')}
                onChange={e =>
                  formik.setFieldValue('ipAddress', e.target.value === '' ? null : e.target.value)
                }
                placeholder="0.0.0.0/0"
              />,
              formik.errors.ipAddress && formik.errors.ipAddress
            )}

            <Text my={3}>Listen port</Text>
            {withValidityIndicator(
              <TextInput
                name="listenPort"
                {...formik.getFieldProps('listenPort')}
                onChange={e =>
                  formik.setFieldValue('listenPort', e.target.value === '' ? null : e.target.value)
                }
                placeholder="51820"
              />,
              formik.errors.listenPort && formik.errors.listenPort
            )}
            <Button width="100%" mt={3} mb={3} mx="auto" onClick={handleSaveButtonClick}>
              Save
            </Button>
          </form>
          <Link mx="auto" to={`/ui/hosts/${hostId}/interfaces`}>
            Cancel
          </Link>
          <FormikState {...formik} />
        </Box>
      )}
    </Container>
  )
}

InterfaceDetailsView.propTypes = {}

export default InterfaceDetailsView
