import { useFormik } from 'formik'
import PropTypes from 'prop-types'
import React, { useEffect } from 'react'
import styled from 'styled-components'
import * as Yup from 'yup'
import { icons } from '_assets/'
import Box from '_components/box'
import Button from '_components/button'
import IconButton from '_components/icon-button'
import TextInput from '_components/inputs/text-input'
import List from '_components/list'
import Text from '_components/text'

const Container = styled(Box)`
  flex-direction: column;
`

const RowsContainer = styled(List)`
  > * {
    :not(:last-child) {
      border-bottom: 1px solid ${(props) => props.theme.colors.neutralLighter};
    }
  }
`

const KeyValue = ({ isEditable, rows, onDelete, onAdd, ...props }) => {
  const formik = useFormik({
    enableReinitialize: true,
    initialValues: {
      key: '',
      value: '',
    },
    validationSchema: Yup.object().shape({
      key: Yup.string().required(),
      value: Yup.string(),
    }),
    validateOnChange: true,
    validateOnMount: true,
    validateOnBlur: true,
  })

  useEffect(() => {
    formik.validateForm()
  }, [])

  const handleAddRow = (k, v) => {
    onAdd(k, v)
    formik.setValues({ key: '', value: '' })
    formik.validateForm()
  }

  const handleDeleteRow = (k) => {
    onDelete(k)
  }

  return (
    <Container {...props}>
      {isEditable && (
        <Box alignItems="center" width="100%" my={1}>
          <TextInput flex="1" name="key" placeholder="Key" {...formik.getFieldProps('key')} />
          <TextInput
            ml={2}
            name="value"
            placeholder="Value"
            {...formik.getFieldProps('value')}
            flex="1"
          />
          <Button
            width="130px"
            variant="primary"
            disabled={!formik.isValid}
            ml={2}
            onClick={() => handleAddRow(formik.values.key, formik.values.value)}
          >
            Add
          </Button>
        </Box>
      )}
      <RowsContainer>
        {rows.map((el) => (
          <Box height="48px" alignItems="center" width="100%" my={1}>
            <Text textStyle="body" width="100%">
              {el.key}
            </Text>
            <Text textStyle="body" width="100%" ml={2}>
              {el.value}
            </Text>
            <IconButton
              style={{ visibility: isEditable ? 'visible' : 'hidden' }}
              mr="20px"
              icon={<icons.Trash />}
              size="20px"
              hoverEffect
              ml={2}
              onClick={() => handleDeleteRow(el.key)}
            />
          </Box>
        ))}
      </RowsContainer>
    </Container>
  )
}

KeyValue.propTypes = {
  rows: PropTypes.arrayOf(
    PropTypes.shape({
      key: PropTypes.string,
      value: PropTypes.string,
    })
  ),
  isEditable: PropTypes.bool,
  onAdd: PropTypes.func,
  onDelete: PropTypes.func,
}

KeyValue.defaultProps = {
  rows: [],
  isEditable: true,
  onAdd: () => {},
  onDelete: () => {},
}

export default KeyValue
