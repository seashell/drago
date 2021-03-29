import { ToastContainer as BaseToastContainer } from 'react-toastify'
import 'react-toastify/dist/ReactToastify.min.css'
import styled from 'styled-components'
import IconButton from '_components/icon-button'

export const ToastContainer = styled(BaseToastContainer).attrs({
  position: 'top-right',
  hideProgressBar: true,
  toastClassName: 'toast',
  bodyClassName: 'body',
})`
  margin: 0;
  padding: 0;
  width: max-content;
  height: max-content;
  margin-top: 100px;
  background: transparent;
  .toast {
    padding: 0;
    margin: 0;
    background: none;
    overflow: visible;
    margin-bottom: 16px;
  }
`
export const ToastBody = styled.div`
  display: flex;
  align-items: center;
  min-height: 96px;
  width: 400px;
  overflow: visible;
  position: relative;
  border-radius: 2px;
  min-height: 86px;
  width: 390px;
  border: 2px solid ${(props) => props.theme.colors[props.color]};
  background: ${(props) => props.theme.colors.white};
`

export const IconContainer = styled.div`
  position: absolute;
  left: -10px;
  top: -10px;
`

export const ToastContent = styled.p`
  padding: 10px;
  padding-left: 24px;
  padding-right: 32px;
  display: block;
  line-break: anywhere;
  color: ${(props) => props.theme.colors[props.color]};
`

export const DismissButton = styled(IconButton).attrs({
  size: '11px',
})`
  position: absolute;
  right: 4px;
  top: 4px;
`

export const CloseIcon = styled(IconButton)`
  position: absolute;
  right: 10px;
  top: 10px;
`
