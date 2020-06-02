/* eslint-disable react/prop-types */
import React, { useEffect } from 'react'
import styled from 'styled-components'
import validator from 'validator'
import * as _ from 'lodash'
import * as Yup from 'yup'

import { useFormik } from 'formik'
import { useLocation, useNavigate, useParams } from '@reach/router'

import { useMutation, useLazyQuery } from 'react-apollo'
import { GET_LINK, UPDATE_LINK } from '_graphql/actions/links'

import { icons } from '_assets/'
import toast from '_components/toast'
import Flex from '_components/flex'
import Box from '_components/box'
import Link from '_components/link'
import Text from '_components/text'
import Button from '_components/button'
import SelectInput from '_components/inputs/select-input'
import IconButton from '_components/icon-button'
import { Dragon as Spinner } from '_components/spinner'
import NumberInput from '_components/inputs/number-input'
import TagsInput from '_components/inputs/tags-input'
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
    {props.data.name} {props.data.ipAddress ? <>({props.data.ipAddress})</> : null} on{' '}
    {props.data.host.name}
  </div>
)

const LinkDetailsView = () => {
  const navigate = useNavigate()
  const location = useLocation()
  const urlParams = useParams()

  const { hostId } = location.state || {}

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

  const [getLink, getLinkQuery] = useLazyQuery(GET_LINK, {
    variables: { id: urlParams.linkId },
    fetchPolicy: 'cache-and-network',
    onCompleted: data => {
      formik.setValues(
        {
          fromInterfaceId: data.result.fromInterfaceId,
          toInterfaceId: data.result.toInterfaceId,
          allowedIps: data.result.allowedIps,
          persistentKeepalive: data.result.persistentKeepalive,
        },
        true
      )
    },
    onError: () => {
      toast.error('Error fetching link details')
      navigate(-1)
    },
  })

  const [updateLink, updateLinkMutation] = useMutation(UPDATE_LINK, {
    variables: { id: urlParams.linkId },
    onCompleted: () => {
      toast.success('Link updated')
      navigate(`/hosts/${hostId}/links`)
      getLink()
    },
    onError: () => {
      toast.error('Error updating link')
    },
  })

  useEffect(() => {
    window.scrollTo(0, 0)
    if (hostId === undefined) {
      navigate(-1)
    }
    formik.resetForm()
    getLink()
  }, [location])

  const handleCancelButtonClick = () => {
    navigate(`/hosts/${hostId}/links`)
  }

  const handleSaveButtonClick = () => {
    formik.validateForm().then(errors => {
      if (_.isEmpty(errors)) {
        updateLink({ variables: { id: urlParams.linkId, ...formik.values } })
      } else {
        toast.error('Form has errors')
      }
    })
  }

  const isLoading = getLinkQuery.loading || updateLinkMutation.loading

  if (isLoading) {
    return <Spinner />
  }

  const link =
    !isLoading && getLinkQuery.data
      ? getLinkQuery.data.result
      : { allowedIps: [], fromInterface: { host: {} }, toInterface: { host: {} } }

  return (
    <Container>
      <Box alignItems="center" mb={4}>
        <IconContainer mr="12px">
          <IconButton ml="auto" size={32} icon={<icons.Link />} />
        </IconContainer>
        <div>
          <Text textStyle="subtitle" fontSize="20px" mb={2}>
            From {link.fromInterface.name} on {link.fromInterface.host.name}
          </Text>
          <Text textStyle="subtitle" fontSize="20px">
            To {link.toInterface.name} on {link.toInterface.host.name}
          </Text>
        </div>
      </Box>
      <Box flexDirection="column">
        <form>
          <Text mt={3} mb={2}>
            Source *
          </Text>
          <SelectInput
            id="fromInterfaceId"
            name="fromInterfaceId"
            value={link.fromInterface}
            optionComponent={InterfaceOption}
            singleValueComponent={SelectedInterfaceValue}
            isClearable
            isDisabled
          />

          <Text mt={3} mb={2}>
            Target *
          </Text>
          <SelectInput
            id="toInterfaceId"
            name="toInterfaceId"
            value={link.toInterface}
            optionComponent={InterfaceOption}
            singleValueComponent={SelectedInterfaceValue}
            isClearable
            isDisabled
          />

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
          <Button width="100%" my={3} mx="auto" onClick={handleSaveButtonClick}>
            Save
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

export default LinkDetailsView
