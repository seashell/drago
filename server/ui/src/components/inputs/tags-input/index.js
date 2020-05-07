import React from 'react'
import PropTypes from 'prop-types'
import styled from 'styled-components'

import ReactTags from 'react-tag-autocomplete'

const StyledReactTags = styled(ReactTags)`
  .is-focused {
  }
  .react-tags__selected {
  }
  .react-tags__selected-tag {
  }
  .react-tags__selected-tag-name {
  }
  .react-tags__search {
  }
  .react-tags__search-input {
  }
  .react-tags__suggestions {
  }
  .suggestionActive {
  }
  .is-disabled {
    opacity: 0.3;
  }
`

const TagsInput = ({ tags }) => <ReactTags isDisabled tags={tags} suggestions={[]} />

TagsInput.propTypes = {
  tags: PropTypes.arrayOf(
    PropTypes.shape({
      id: PropTypes.string,
      text: PropTypes.string,
    })
  ),
}

TagsInput.defaultProps = {
  tags: [],
}

export default TagsInput
