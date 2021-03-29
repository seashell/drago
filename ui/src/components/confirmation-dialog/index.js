import PropTypes from 'prop-types'
import React, { useContext, useState } from 'react'
import { Portal } from 'react-portal'
import Modal from 'styled-react-modal'
import Dialog from './dialog'

const ConfirmationDialogContext = React.createContext()

export const useConfirmationDialog = () => useContext(ConfirmationDialogContext)

const defaultOptions = {
  title: 'Are you sure?',
  details: 'This cannot be undone',
  isDestructive: false,
  onConfirm: () => {},
  onCancel: () => {},
}

const ConfirmationDialogProvider = ({ children }) => {
  const [isOpen, setOpen] = useState(false)
  const [options, setOptions] = useState(defaultOptions)

  const onConfirmButtonClick = () => {
    options.onConfirm()
    closeConfirmationDialog()
  }

  const onCancelButtonClick = () => {
    options.onCancel()
    closeConfirmationDialog()
  }

  const openConfirmationDialog = () => {
    setOpen(true)
  }

  const closeConfirmationDialog = () => {
    setOpen(false)
  }

  return (
    <>
      <ConfirmationDialogContext.Provider
        value={{
          confirm: (opts) => {
            setOptions(defaultOptions)
            if (typeof opts === 'object' && opts !== null) {
              const merged = Object.assign(options, opts)
              setOptions(merged)
            }
            openConfirmationDialog()
          },
        }}
      >
        {children}
      </ConfirmationDialogContext.Provider>
      <Portal>
        <Modal isOpen={isOpen} onBackgroundClick={closeConfirmationDialog}>
          <Dialog
            title={options.title}
            details={options.details}
            isDestructive={options.isDestructive}
            onConfirm={onConfirmButtonClick}
            onCancel={onCancelButtonClick}
          />
        </Modal>
      </Portal>
    </>
  )
}

ConfirmationDialogProvider.propTypes = {
  children: PropTypes.node.isRequired,
}

export default ConfirmationDialogProvider
