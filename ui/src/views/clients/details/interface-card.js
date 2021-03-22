import { useFormik } from 'formik'
import PropTypes from 'prop-types'
import React from 'react'
import styled from 'styled-components'
import * as Yup from 'yup'
import { icons } from '_assets/'
import Box from '_components/box'
import Button from '_components/button'
import Icon from '_components/icon'
import IconButton from '_components/icon-button'
import NumberInput from '_components/inputs/number-input'
import TextInput from '_components/inputs/text-input'
import Separator from '_components/separator'
import { Dragon as Spinner } from '_components/spinner'
import Text from '_components/text'

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
  grid-template-columns: 200px auto 160px 120px 120px auto;

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

const StyledSpinner = styled(Spinner).attrs({
  size: 40,
})`
  position: relative;
  display: flex;
  align-items: center;
`

const InterfaceCard = ({
  id,
  showSpinner,
  name,
  connectionsCount,
  address,
  listenPort,
  dns,
  mtu,
  table,
  preUp,
  postUp,
  preDown,
  postDown,
  networkName,
  onClick,
  onChange,
  onDelete,
  isExpanded,
  hasPublicKey,
  publicKey,
}) => {
  const formik = useFormik({
    initialValues: {
      id,
      address,
      listenPort,
      dns,
      table,
      mtu,
      preUp,
      preDown,
      postUp,
      postDown,
    },
    enableReinitialize: true,
    validationSchema: Yup.object().shape({
      address: Yup.string().required().nullable(),
      listenPort: Yup.number().positive().integer().required().nullable(),
    }),
    validateOnBlur: true,
    validateOnMount: true,
    validateOnChange: true,
  })

  const handleClick = () => {
    onClick()
  }

  const handleDeleteButtonClick = (e) => {
    e.preventDefault()
    e.stopPropagation()
    onDelete()
  }

  const handleSaveButtonClick = () => {
    const values = { ...formik.values }
    onChange(id, values)
  }

  return (
    <Container>
      <HeaderContainer onClick={handleClick}>
        <Box>
          <StyledIcon mr="12px" icon={<icons.Interface />} color="neutralDarker" />
          <Box flexDirection="column" justifyContent="center">
            <Text textStyle="subtitle" fontSize="14px">
              {id.split('-')[0]}
            </Text>
            {
              <Text textStyle="detail" fontSize="12px">
                {name}
              </Text>
            }
          </Box>
        </Box>

        <Text textStyle="detail" fontSize="12px" title="Managed address" justifySelf="center">
          {address}
        </Text>

        {hasPublicKey ? (
          <Badge alignItems="center" justifySelf="center">
            <Icon icon={<icons.Key />} color="neutralDarker" size="14px" mr={'4px'} />
          </Badge>
        ) : (
          <span />
        )}

        <Badge alignItems="center" justifySelf="center">
          <Icon icon={<icons.ConnectionSmall />} color="neutralDarker" size="14px" mr={'4px'} />
          <Text textStyle="detail" fontSize="10px" title="Managed address">
            {connectionsCount}
          </Text>
        </Badge>

        <Badge alignItems="center" justifySelf="center">
          <Icon mr={1} icon={<icons.Network />} color="neutralDarker" size="16px" />
          <Text textStyle="detail" fontSize="12px">
            {networkName}
          </Text>
        </Badge>

        <IconButton
          onClick={handleDeleteButtonClick}
          icon={<icons.Trash />}
          color="neutralDark"
          size="16px"
          ml="auto"
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
              {publicKey && (
                <Text textStyle="detail" ml={2}>
                  (PK: {publicKey})
                </Text>
              )}
            </Box>

            <ConfigurationGrid mt={3}>
              <Box flexDirection="column" gridColumn="span 6">
                <Text textStyle="detail" mb={1}>
                  Address
                </Text>
                <TextInput
                  name="address"
                  {...formik.getFieldProps('address')}
                  placeholder="192.168.0.1/24"
                />
              </Box>

              <Box flexDirection="column" gridColumn="span 2">
                <Text textStyle="detail" mb={1}>
                  Listen Port
                </Text>
                <NumberInput
                  name="listenPort"
                  {...formik.getFieldProps('listenPort')}
                  placeholder="51820"
                />
              </Box>

              <Box flexDirection="column" gridColumn="span 2">
                <Text textStyle="detail" mb={1}>
                  DNS
                </Text>
                <TextInput name="dns" {...formik.getFieldProps('dns')} placeholder="8.8.8.8" />
              </Box>

              <Box flexDirection="column" gridColumn="span 1">
                <Text textStyle="detail" mb={1}>
                  Table
                </Text>
                <TextInput name="table" {...formik.getFieldProps('table')} placeholder="12345" />
              </Box>

              <Box flexDirection="column" gridColumn="span 1">
                <Text textStyle="detail" mb={1}>
                  MTU
                </Text>
                <TextInput name="mtu" {...formik.getFieldProps('mtu')} placeholder="1500" />
              </Box>

              <Box flexDirection="column" gridColumn="1 / span 6">
                <Text textStyle="detail" mb={1}>
                  Pre up
                </Text>
                <TextInput
                  name="preUp"
                  {...formik.getFieldProps('preUp')}
                  placeholder="/bin/echo 'Hello world!'"
                />
              </Box>

              <Box flexDirection="column" gridColumn="span 6">
                <Text textStyle="detail" mb={1}>
                  Pre down
                </Text>
                <TextInput
                  name="preDown"
                  {...formik.getFieldProps('preDown')}
                  placeholder="/bin/echo 'Goodbye world!'"
                />
              </Box>

              <Box flexDirection="column" gridColumn="span 6">
                <Text textStyle="detail" mb={1}>
                  Post up
                </Text>
                <TextInput
                  name="postUp"
                  {...formik.getFieldProps('postUp')}
                  placeholder="/bin/echo 'World?'"
                />
              </Box>

              <Box flexDirection="column" gridColumn="span 6">
                <Text textStyle="detail" mb={1}>
                  Post down
                </Text>
                <TextInput
                  name="postDown"
                  {...formik.getFieldProps('postDown')}
                  placeholder="/bin/echo ''"
                />
              </Box>
            </ConfigurationGrid>

            <Box ml="auto" alignItems="center">
              {showSpinner && <StyledSpinner />}
              <Button variant="primary" mt={3} height="40px" onClick={handleSaveButtonClick}>
                Save
              </Button>
            </Box>
          </HiddenContentContainer>
        </>
      )}
    </Container>
  )
}

InterfaceCard.propTypes = {
  id: PropTypes.string,
  name: PropTypes.string,
  showSpinner: PropTypes.bool,
  connectionsCount: PropTypes.number,
  listenPort: PropTypes.number,
  dns: PropTypes.string,
  mtu: PropTypes.number,
  table: PropTypes.string,
  preUp: PropTypes.string,
  postUp: PropTypes.string,
  preDown: PropTypes.string,
  postDown: PropTypes.string,
  address: PropTypes.string,
  networkName: PropTypes.string,
  isExpanded: PropTypes.bool,
  onClick: PropTypes.func,
  onChange: PropTypes.func,
  onDelete: PropTypes.func,
  hasPublicKey: PropTypes.bool,
  publicKey: PropTypes.string,
}

InterfaceCard.defaultProps = {
  id: '',
  name: '',
  showSpinner: false,
  connectionsCount: 0,
  listenPort: 0,
  dns: '',
  mtu: 0,
  table: '',
  preUp: '',
  postUp: '',
  preDown: '',
  postDown: '',
  address: '',
  networkName: '',
  isExpanded: false,
  hasPublicKey: false,
  publicKey: '',
  onClick: () => {},
  onChange: () => {},
  onDelete: () => {},
}

export default InterfaceCard
