import styled from 'styled-components'
import { ToastContainer as BaseToastContainer } from 'react-toastify'

import 'react-toastify/dist/ReactToastify.min.css'
import IconButton from '_components/icon-button'

export const ToastContainer = styled(BaseToastContainer).attrs({
  position: 'top-right',
  hideProgressBar: true,
  toastClassName: 'toast',
  bodyClassName: 'body',
})`
  .toast {
    background: white;
    border-radius: 4px;
  }
  margin-top: 100px;
`

export const ToastBody = styled.div`
  display: grid;
  grid-template-columns: auto 1fr;
  grid-column-gap: 16px;
  position: relative;
  align-items: center;
  padding: 10px;
`

export const ToastContent = styled.p`
  color: ${props => props.theme.neutralLight};
`

export const CloseIcon = styled(IconButton)`
  position: absolute;
  right: 10px;
  top: 10px;
`
