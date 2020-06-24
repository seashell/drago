import reset from 'styled-reset'
import { variant } from 'styled-system'
import { createGlobalStyle } from 'styled-components'

import 'typeface-roboto' // eslint-disable-line
import 'typeface-raleway' // eslint-disable-line
import 'typeface-lato' // eslint-disable-line

import lightTheme from './themes/light'

const GlobalStyles = createGlobalStyle`
    ${reset}
        
    :root {
        height: 100%;
        font-family: Lato !important;
        text-rendering: optimizeLegibility;
        outline: none;
    }

    body {
        height: 100%;
    }

    body::-webkit-scrollbar {
        width: 8px;
        height: 8px;
    }
    
    body::-webkit-scrollbar-track {
        border-radius: 4px;
    }
    
    body::-webkit-scrollbar-thumb {
        border-radius: 4px;
        background-color: #ccc;
    }

    * {
        outline: none;
        -webkit-tap-highlight-color: transparent;
    }

    button {
        cursor: pointer;
    }
`

const containers = variant({
  scale: 'containers',
  prop: 'type',
})

const themes = { light: lightTheme }

export { GlobalStyles, containers, themes }
