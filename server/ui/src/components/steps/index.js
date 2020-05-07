import styled from 'styled-components'
import { space, layout, grid, flex } from 'styled-system'

export const StepsContainer = styled.div`
    display: flex;
    counter-reset: step-counter;

    > * {
      opacity: 0.3;
      :nth-child(-n+${props => props.currentStep}){
        opacity: 1;
      }
    }

    ${space}
    ${layout}
    ${grid}
    ${flex}
`

export const Step = styled.p`
  color: ${props => props.theme.colors.primary};
  display: flex;
  align-items: center;
  :before {
    content: counter(step-counter);
    display: flex;
    align-items: center;
    justify-content: center;
    margin-right: 8px;
    border-radius: 50%;
    counter-increment: step-counter;
    height: 28px;
    width: 28px;
    border: 1px solid ${props => props.theme.colors.primary};
  }
  z-index: -1;
`
