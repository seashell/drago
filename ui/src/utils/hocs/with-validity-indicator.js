import PropTypes from 'prop-types'
import React from 'react'
import styled from 'styled-components'
import { icons } from '_assets/'
import Box from '_components/box'
import Popover from '_components/popover'
import Text from '_components/text'

const ErrorIconContainer = styled.div`
  position: absolute;
  right: -28px;
  top: 50%;
  transform: translateY(-50%);
`

const ErrorTooltip = styled.div`
  padding: 8px;
  max-width: 180px;
  background: #fff;
  border-radius: 4px;
`

const ErrorIndicator = ({ error }) => {
  if (error === undefined) return null
  return (
    <ErrorIconContainer>
      <Popover
        trigger="mouseenter"
        content={
          <ErrorTooltip>
            <Text fontSize="12px" textStyle="detail">
              {error}
            </Text>
          </ErrorTooltip>
        }
      >
        <icons.Error />
      </Popover>
    </ErrorIconContainer>
  )
}

ErrorIndicator.propTypes = {
  error: PropTypes.string,
}

ErrorIndicator.defaultProps = {
  error: undefined,
}

const withValidityIndicator = (input, error) => (
  <Box style={{ position: 'relative' }}>
    {input}
    <ErrorIndicator error={error} />
  </Box>
)

export default withValidityIndicator
