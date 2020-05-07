import React from 'react'
import { toast as baseToast } from 'react-toastify'

import { icons } from '_assets/'

import { ToastContainer, ToastBody, ToastContent } from './styled'

const displayToast = (content, icon, options = {}) => {
  baseToast(
    <ToastBody>
      {icon}
      <ToastContent>{content}</ToastContent>
    </ToastBody>,
    { ...options, closeButton: <icons.Times />, autoClose: 3000 }
  )
}

const custom = (text, { icon, ...options }) => displayToast(text, icon, options)
const error = (text, options) => displayToast(text, <icons.Error />, options)
const warning = (text, options) => displayToast(text, <icons.Warning />, options)
const success = (text, options) => displayToast(text, <icons.Success />, options)

export { ToastContainer }

export default { custom, error, warning, success }
