import reset from 'styled-reset'
import { variant } from 'styled-system'
import { createGlobalStyle } from 'styled-components'

import lightTheme from './themes/light'

const GlobalStyles = createGlobalStyle`
    ${reset}

    :root {
        height: 100%;
        font-family: Lato;
        text-rendering: optimizeLegibility;
        outline: none;
    }

    body {
        height: 100%;
    }

    body::-webkit-scrollbar {
        width: 8px;
        height: 8px;
        background-color: ${(props) => props.theme.colors.background3};
    }
    
    body::-webkit-scrollbar-track {
        border-radius: 4px;
        background-color: ${(props) => props.theme.colors.background3};
    }
    
    body::-webkit-scrollbar-thumb {
        border-radius: 4px;
        background-color: ${(props) => props.theme.colors.border1};
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
