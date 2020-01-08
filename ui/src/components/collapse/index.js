import React, { useState } from 'react'
import PropTypes from 'prop-types'
import styled from 'styled-components'
import { space } from 'styled-system'
import { icons } from '_assets'

const Header = styled.div.attrs({
  role: 'button',
})`
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: space-between;
  ${space}
`

Header.defaultProps = {
  paddingTop: 2,
  paddingBottom: 2,
}

const Content = styled.div`
  ${space}
`

const Collapse = ({ title, children, ...props }) => {
  const [isCollapseOpen, setCollapseOpen] = useState(props.isOpen)

  const handleHeaderClick = () => {
    setCollapseOpen(!isCollapseOpen)
  }

  return (
    <>
      <Header className="header" onClick={handleHeaderClick} {...props}>
        {title}
        {isCollapseOpen ? <icons.ArrowUp className="indicator" /> : <icons.ArrowDown />}
      </Header>
      {isCollapseOpen && <Content>{children}</Content>}
    </>
  )
}

Collapse.propTypes = {
  header: PropTypes.node,
  title: PropTypes.node.isRequired,
  isOpen: PropTypes.bool,
  children: PropTypes.node,
}

Collapse.defaultProps = {
  header: null,
  children: undefined,
  isOpen: false,
}

export default Collapse
