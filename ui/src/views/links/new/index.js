/* eslint-disable react/prop-types */
import React, { useEffect, useState } from 'react'
import styled from 'styled-components'
import validator from 'validator'
import * as _ from 'lodash'
import * as Yup from 'yup'

import { useFormik } from 'formik'
import { useLocation, useNavigate } from '@reach/router'

import { useMutation, useLazyQuery } from 'react-apollo'
import { CREATE_LINK } from '_graphql/actions/links'
import { GET_NETWORKS } from '_graphql/actions/networks'
import { GET_INTERFACES } from '_graphql/actions/interfaces'
import { GET_HOSTS } from '_graphql/actions/hosts'

import { icons } from '_assets/'

import Flex from '_components/flex'
import Box from '_components/box'
import Link from '_components/link'
import Text from '_components/text'
import Button from '_components/button'
import IconButton from '_components/icon-button'
import SelectInput from '_components/inputs/select-input'
import NumberInput from '_components/inputs/number-input'
import TagsInput from '_components/inputs/tags-input'
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
  height: '40px',
  width: '40px',
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
    <IconContainer mr="12px">
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

const HostOption = ({ innerRef, innerProps, ...props }) => (
  <SearchResult innerRef={innerRef} {...innerProps} {...props}>
    <IconContainer mr="12px">
      <IconButton ml="auto" icon={<icons.Host />} />
    </IconContainer>
    {props.data.name}
  </SearchResult>
)

const SelectedHostValue = ({ innerRef, innerProps, ...props }) => (
  <div innerRef={innerRef} {...innerProps} {...props}>
    {props.data.name}
  </div>
)

const InterfaceOption = ({ innerRef, innerProps, ...props }) => (
  <SearchResult innerRef={innerRef} {...innerProps} {...props}>
    <IconContainer mr="12px">
      <IconButton ml="auto" icon={<icons.Interface />} />
    </IconContainer>
    {props.data.name} {props.data.ipAddress ? <>({props.data.ipAddress})</> : null}
  </SearchResult>
)

const SelectedInterfaceValue = ({ innerRef, innerProps, ...props }) => (
  <div innerRef={innerRef} {...innerProps} {...props}>
    {props.data.name} {props.data.ipAddress ? <>({props.data.ipAddress})</> : null}
  </div>
)

const filterNetworks = (option, searchText) => {
  if (option.data.name.toLowerCase().includes(searchText.toLowerCase())) {
    return true
  }
  return false
}

const filterHosts = (option, searchText) => {
  if (option.data.name.toLowerCase().includes(searchText.toLowerCase())) {
    return true
  }
  return false
}

const filterInterfaces = (option, searchText) => {
  if (option.data.name.toLowerCase().includes(searchText.toLowerCase())) {
    return true
  }
  return false
}

const NewLinkView = () => {
  const navigate = useNavigate()
  const location = useLocation()

  const { hostId } = location.state || {}

  const [selectedNetworkId, setSelectedNetworkId] = useState('')
  const [selectedHostId, setSelectedHostId] = useState('')

  const formik = useFormik({
    initialValues: {
      fromInterfaceId: null,
      toInterfaceId: null,
      allowedIps: [],
      persistentKeepalive: null,
    },
    validationSchema: Yup.object().shape({
      fromInterfaceId: Yup.string()
        .required()
        .nullable(),
      toInterfaceId: Yup.string()
        .required()
        .nullable(),
      allowedIps: Yup.array()
        .of(
          Yup.string().test(
            'allowedIps',
            'Allowed IPs must contain only valid IP ranges in CIDR notation',
            value => validator.isIPRange(value)
          )
        )
        .nullable(),
      persistentKeepalive: Yup.number()
        .positive()
        .nullable(),
    }),
  })

  const [getSourceInterfaces, getSourceInterfacesQuery] = useLazyQuery(GET_INTERFACES, {
    variables: { hostId },
  })

  const [getTargetNetworks, getTargetNetworksQuery] = useLazyQuery(GET_NETWORKS)

  const [getTargetHosts, getTargetHostsQuery] = useLazyQuery(GET_HOSTS)

  const [getTargetInterfaces, getTargetInterfacesQuery] = useLazyQuery(GET_INTERFACES, {
    variables: { hostId },
  })

  const [createLink, createLinkMutation] = useMutation(CREATE_LINK, {
    onCompleted: () => {
      toast.success('Link created')
      navigate(`/ui/hosts/${hostId}/links`)
    },
    onError: () => {
      toast.error('Error creating link')
    },
  })

  useEffect(() => {
    window.scrollTo(0, 0)

    if (hostId === undefined) {
      navigate(-1)
    }
    setSelectedNetworkId(null)
    setSelectedHostId(null)

    formik.resetForm()

    getSourceInterfaces()
    getTargetNetworks()
  }, [location])

  const handleCancelButtonClick = () => {
    navigate(`/ui/hosts/${hostId}/links`)
  }

  const handleSelectedNetworkChanged = network => {
    if (network !== null) {
      setSelectedNetworkId(network.id)
      getTargetHosts({ variables: { networkId: network.id } })
    } else {
      setSelectedNetworkId(null)
      setSelectedHostId(null)
      formik.setFieldValue('toInterface', null)
    }
  }

  const handleSelectedHostChanged = host => {
    if (host !== null) {
      setSelectedHostId(host.id)
      getTargetInterfaces({ variables: { hostId: host.id } })
    } else {
      setSelectedHostId(null)
      formik.setFieldValue('toInterface', null)
    }
  }

  const handleCreateButtonClick = () => {
    formik.validateForm().then(errors => {
      if (_.isEmpty(errors)) {
        createLink({ variables: { ...formik.values } })
      } else {
        toast.error('Form has errors')
      }
    })
  }

  const isLoading =
    getSourceInterfacesQuery.loading || getTargetNetworksQuery.loading || createLinkMutation.loading

  if (isLoading) {
    return <Spinner />
  }

  const sourceInterfaceOptions =
    !isLoading && getSourceInterfacesQuery.data ? getSourceInterfacesQuery.data.result.items : []

  const targetNetworkOptions =
    !isLoading && getTargetNetworksQuery.data ? getTargetNetworksQuery.data.result.items : []

  const targetHostOptions =
    !isLoading && getTargetHostsQuery.data ? getTargetHostsQuery.data.result.items : []

  const targetInterfaceOptions =
    !isLoading && getTargetInterfacesQuery.data
      ? getTargetInterfacesQuery.data.result.items.filter(
          el => el.id !== formik.values.fromInterfaceId
        )
      : []

  return (
    <Container>
      <Text textStyle="title" mb={4}>
        New link
      </Text>
      <Box flexDirection="column">
        <form>
          <Text mt={3} mb={2}>
            Source *
          </Text>
          {withValidityIndicator(
            <SelectInput
              id="fromInterfaceId"
              name="fromInterfaceId"
              value={
                sourceInterfaceOptions.find(el => el.id === formik.values.fromInterfaceId) || null
              }
              onChange={value => {
                setSelectedNetworkId(null)
                setSelectedHostId(null)
                formik.setFieldValue('fromInterfaceId', value !== null ? value.id : value)
              }}
              placeholder="Source interface"
              options={sourceInterfaceOptions}
              optionComponent={InterfaceOption}
              singleValueComponent={SelectedInterfaceValue}
              filterOption={filterInterfaces}
              isClearable
            />,
            formik.errors.fromInterfaceId
          )}

          <Text mt={3} mb={2}>
            Target *
          </Text>

          <SelectInput
            value={targetNetworkOptions.find(el => el.id === selectedNetworkId) || null}
            placeholder="Target network"
            options={targetNetworkOptions}
            onChange={e => handleSelectedNetworkChanged(e)}
            singleValueComponent={SelectedNetworkValue}
            optionComponent={NetworkOption}
            filterOption={filterNetworks}
            isClearable
          />

          <SelectInput
            value={targetHostOptions.find(el => el.id === selectedHostId) || null}
            options={targetHostOptions}
            isDisabled={selectedNetworkId === null}
            onChange={e => handleSelectedHostChanged(e)}
            singleValueComponent={SelectedHostValue}
            optionComponent={HostOption}
            filterOption={filterHosts}
            placeholder="Target host"
            isClearable
            my={2}
          />

          {withValidityIndicator(
            <SelectInput
              id="toInterfaceId"
              name="toInterfaceId"
              value={
                targetInterfaceOptions.find(el => el.id === formik.values.toInterfaceId) || null
              }
              onChange={value => formik.setFieldValue('toInterfaceId', value.id)}
              isDisabled={selectedNetworkId === null || selectedHostId === null}
              placeholder="Target interface"
              options={targetInterfaceOptions}
              optionComponent={InterfaceOption}
              singleValueComponent={SelectedInterfaceValue}
              filterOption={filterInterfaces}
              isClearable
            />,
            formik.errors.toInterfaceId
          )}

          <Text my={3}>Allowed IPs</Text>
          {withValidityIndicator(
            <TagsInput
              id="allowedIps"
              name="allowedIps"
              value={formik.values.allowedIps.map(el => ({ value: el, label: el }))}
              onChange={values => {
                formik.setFieldValue(
                  'allowedIps',
                  values !== null ? values.map(el => el.value) : []
                )
              }}
              placeholder="Add CIDR address and hit Enter"
            />,
            formik.errors.allowedIps
          )}

          <Text my={3}>Persistent keepalive</Text>
          {withValidityIndicator(
            <NumberInput
              name="persistentKeepalive"
              {...formik.getFieldProps('persistentKeepalive')}
              onChange={e =>
                formik.setFieldValue(
                  'persistentKeepalive',
                  Number.isNaN(parseInt(e.target.value, 10)) ? null : parseInt(e.target.value, 10)
                )
              }
              placeholder={20}
            />,
            formik.errors.persistentKeepalive
          )}
          <Button width="100%" my={3} mx="auto" onClick={handleCreateButtonClick}>
            Create
          </Button>
        </form>
        <Link mx="auto" to="" onClick={handleCancelButtonClick}>
          Cancel
        </Link>
        <FormikState {...formik} />
      </Box>
    </Container>
  )
}

NewLinkView.propTypes = {}

export default NewLinkView
