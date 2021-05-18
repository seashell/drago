import { useQuery } from '@apollo/client'
import { useFormik } from 'formik'
import PropTypes from 'prop-types'
import React, { useEffect } from 'react'
import styled from 'styled-components'
import * as Yup from 'yup'
import { icons } from '_assets/'
import Box from '_components/box'
import Button from '_components/button'
import Icon from '_components/icon'
import IconButton from '_components/icon-button'
import NumberInput from '_components/inputs/number-input'
import TagsInput from '_components/inputs/tags-input'
import Link from '_components/link'
import Separator from '_components/separator'
import { Dragon as Spinner } from '_components/spinner'
import Text from '_components/text'
import { GET_CONNECTION, GET_PEER } from '_graphql/queries'

const Container = styled(Box).attrs({
  border: 'discrete',
})`
  flex-direction: column;
  :not(:last-child) {
    border-bottom: none;
  }
`

const HeaderContainer = styled(Box).attrs({
  px: 3,
})`
  display: grid;
  grid-template-columns: 330px 30px 120px 120px auto;

  height: 72px;
  cursor: pointer;
  align-items: center;
`

const HiddenContentContainer = styled(Box).attrs({
  p: 3,
})`
  flex-direction: column;
`

const StyledIcon = styled(Icon)`
  width: 36px;
  height: 36px;
  padding: 4px;
  border-radius: 4px;
  background: ${(props) => props.theme.colors.neutralLighter};
  align-items: center;
  justify-content: center;
`

const Badge = styled(Box)`
  color: ${(props) => props.theme.colors.neutralDark};
  background: ${(props) => props.theme.colors.neutralLighter};
  padding: 4px 8px;
  border-radius: 2px;
  justify-content: center;
`

const ConfigurationGrid = styled(Box)`
  display: grid;
  grid-template: 1fr / repeat(12, 1fr);
  grid-column-gap: 8px;
  grid-row-gap: 16px;
`

const StyledLink = styled(Link)`
  color: ${(props) => props.theme.colors.neutralDarker};
  font-size: 12px;
  line-height: inherit;
  font-weight: bold;
  display: inline-flex;
  align-items: center;
  :hover {
    text-decoration: underline;
  }
`

const StyledSpinner = styled(Spinner).attrs({
  size: 40,
})`
  position: relative;
  display: flex;
  align-items: center;
`

const ConnectionCard = ({
  id,
  networkName,
  fromInterfaceId,
  toInterfaceId,
  onClick,
  onChange,
  onDelete,
  isExpanded,
  showSpinner,
}) => {
  const formik = useFormik({
    initialValues: {
      id,
      allowedIPs: [],
      persistentKeepalive: null,
      interfaceId: fromInterfaceId,
    },
    enableReinitialize: true,
    validationSchema: Yup.object().shape({
      persistentKeepalive: Yup.number().positive().integer().required().nullable(),
    }),
    validateOnBlur: true,
    validateOnMount: true,
    validateOnChange: true,
  })

  const handleGetConnectionsQueryData = (data) => {
    const fromInterfaceSettings = data.result.PeerSettings[fromInterfaceId]
    formik.setFieldValue('allowedIPs', fromInterfaceSettings.RoutingRules.AllowedIPs)
    formik.setFieldValue('persistentKeepalive', data.result.PersistentKeepalive)
  }

  const getConnectionQuery = useQuery(GET_CONNECTION, {
    variables: { connectionId: id },
    onCompleted: handleGetConnectionsQueryData,
  })

  const getSelfQuery = useQuery(GET_PEER, {
    variables: { interfaceId: fromInterfaceId },
  })

  const getPeerQuery = useQuery(GET_PEER, {
    variables: { interfaceId: toInterfaceId },
  })

  useEffect(() => {
    if (isExpanded) {
      getConnectionQuery.refetch()
      getSelfQuery.refetch()
      getPeerQuery.refetch()
    }
  }, [isExpanded])

  const handleClick = () => {
    onClick()
  }

  const handleDeleteButtonClick = (e) => {
    e.preventDefault()
    e.stopPropagation()
    onDelete()
  }

  const handleSaveButtonClick = (e) => {
    e.preventDefault()
    e.stopPropagation()
    onChange(id, formik.values)
  }

  const self = getSelfQuery.data ? getSelfQuery.data.result : { ID: '', Node: { ID: '' } }
  const peer = getPeerQuery.data ? getPeerQuery.data.result : { ID: '', Node: { ID: '' } }

  return (
    <Container>
      <HeaderContainer onClick={handleClick}>
        <Box>
          <StyledIcon mr="12px" icon={<icons.Connection />} color="neutralDarker" />
          <Box flexDirection="column">
            <Text textStyle="subtitle" fontSize="14px">
              {id.split('-')[0]}
            </Text>
            <Text textStyle="detail" fontSize="12px">
              Connected to{' '}
              <StyledLink to={`/ui/clients/${peer.Node.ID}/`} state={{ connectionId: id }}>
                {toInterfaceId !== null && toInterfaceId.split('-')[0]} ({peer.Node.Name})
                <Icon icon={<icons.ExternalLink />} size="10px" color="neutralDarker" ml="1" />
              </StyledLink>
            </Text>
          </Box>
        </Box>

        <Box width="48px" justifyContent="center" />

        <Badge alignItems="center" justifySelf="center">
          <Icon mr={1} icon={<icons.Interface />} color="neutralDarker" size="16px" />
          <Text textStyle="detail" fontSize="12px">
            {fromInterfaceId !== null && fromInterfaceId.split('-')[0]}
          </Text>
        </Badge>

        <Badge alignItems="center" justifySelf="center">
          <Icon mr={1} icon={<icons.Network />} color="neutralDarker" size="16px" />
          <Text textStyle="detail" fontSize="12px">
            {networkName}
          </Text>
        </Badge>

        <IconButton
          icon={<icons.Trash />}
          color="neutralDark"
          size="16px"
          ml="auto"
          onClick={handleDeleteButtonClick}
          style={{ visibility: isExpanded ? 'hidden' : 'visible' }}
        />
      </HeaderContainer>
      {isExpanded && (
        <>
          <Separator mx={'10px'} width="auto" />
          <HiddenContentContainer>
            <Box alignItems="center" mb={2}>
              <Text textStyle="subtitle" fontSize="16px">
                Configuration
              </Text>
            </Box>

            <ConfigurationGrid>
              <Box flexDirection="column" gridColumn="1 / span 10">
                <Text textStyle="detail">Allowed IPs</Text>
                <TagsInput
                  value={formik.values.allowedIPs.map((el) => ({ label: el, value: el }))}
                  onChange={(values) =>
                    formik.setFieldValue(
                      'allowedIPs',
                      values === null ? [] : values.map((el) => el.value)
                    )
                  }
                  placeholder="192.168.0.1/24"
                />
              </Box>

              <Box flexDirection="column" gridColumn="11 / span 2">
                <Text textStyle="detail" mb={'14px'}>
                  Persistent Keepalive
                </Text>
                <NumberInput
                  name="persistentKeepalive"
                  {...formik.getFieldProps('persistentKeepalive')}
                  placeholder={1000}
                />
              </Box>
            </ConfigurationGrid>
            <Box alignItems="center" mt={3}>
              <div>
                <Text textStyle="detail">
                  * This configuration applies to interface <strong>{self.ID.split('-')[0]}</strong>{' '}
                  {self.Address ? `(${self.Address})` : ''} of the node identified by{' '}
                  <strong>{self.Node.ID.split('-')[0]}</strong>.
                </Text>
                <Text textStyle="detail">
                  ** Interface <strong>{self.ID.split('-')[0]}</strong> is connected to interface{' '}
                  <strong>{peer.ID.split('-')[0]}</strong> {peer.Address ? `(${peer.Address})` : ''}{' '}
                  on node <strong>{peer.Node.ID.split('-')[0]}</strong>
                </Text>
              </div>
              <Box ml="auto" alignItems="center">
                {showSpinner && <StyledSpinner />}
                <Button ml={3} variant="primary" height="40px" onClick={handleSaveButtonClick}>
                  Save
                </Button>
              </Box>
            </Box>
          </HiddenContentContainer>
        </>
      )}
    </Container>
  )
}

ConnectionCard.propTypes = {
  id: PropTypes.string.isRequired,
  networkName: PropTypes.string.isRequired,
  fromInterfaceId: PropTypes.string.isRequired,
  toInterfaceId: PropTypes.string.isRequired,
  updatedAt: PropTypes.string,
  createdAt: PropTypes.string,
  showSpinner: PropTypes.bool,
  onClick: PropTypes.func,
  onChange: PropTypes.func,
  onDelete: PropTypes.func,
  isExpanded: PropTypes.bool,
}

ConnectionCard.defaultProps = {
  createdAt: undefined,
  updatedAt: undefined,
  showSpinner: false,
  onClick: () => {},
  onChange: () => {},
  onDelete: () => {},
  isExpanded: false,
}

export default ConnectionCard
