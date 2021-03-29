import PropTypes from 'prop-types'
import React, { useContext } from 'react'
import { toast } from 'react-toastify'

import { icons } from '_assets/'

import {
  ToastContainer,
  ToastBody,
  IconContainer,
  DismissButton,
  ToastContent,
} from '_components/toast'

const displayToast = (content, icon, color, options = {}) => {
  toast(
    ({ closeToast }) => (
      <ToastBody color={color}>
        <IconContainer>{icon}</IconContainer>
        <ToastContent color={color}>{content}</ToastContent>
        <DismissButton icon={<icons.Times />} color={color} onClick={() => closeToast()} />
      </ToastBody>
    ),
    { ...options, autoClose: 3000, closeButton: false }
  )
}

const ToastContext = React.createContext()

export const useToast = () => useContext(ToastContext)

const ToastProvider = ({ children }) => (
  <>
    <ToastContext.Provider
      value={{
        custom: (text, { icon, color, ...options }) => displayToast(text, icon, color, options),
        error: (text, options) => displayToast(text, <icons.Error />, 'danger', options),
        warning: (text, options) => displayToast(text, <icons.Warning />, 'warning', options),
        success: (text, options) => displayToast(text, <icons.Success />, 'success', options),
      }}
    >
      {children}
      <ToastContainer />
    </ToastContext.Provider>
  </>
)

ToastProvider.propTypes = {
  children: PropTypes.node.isRequired,
}

export default ToastProvider
