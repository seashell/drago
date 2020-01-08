import styled from 'styled-components'
import { layout, space, flexbox } from 'styled-system'

const List = styled.ul`
  overflow: hidden;
  overflow-y: auto;

  ::-webkit-scrollbar {
    background: transparent;
    width: 4px;
    height: 4px;
  }

  ::-webkit-scrollbar-button {
    display: none;
  }

  ::-webkit-scrollbar-track {
    background-color: transparent;
  }

  ::-webkit-scrollbar-track-piece {
    background-color: transparent;
  }

  ::-webkit-scrollbar-thumb {
    background-color: #ececf0;
    border-radius: 2px;
  }

  ::-webkit-scrollbar-corner {
    display: none;
  }

  ::-webkit-resizer {
    display: none;
  }

  ${flexbox}
  ${layout}
  ${space}
`

export default List
