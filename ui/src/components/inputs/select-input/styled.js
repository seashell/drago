import styled, { css } from 'styled-components'
import { layout, space } from 'styled-system'

import Select from 'react-styled-select'

export const Container = styled.div`
  width: 100%;
  ${layout}
  ${space}
`

export const StyledSelect = styled(Select)`
  ${({ theme: { colors, shadows } }) => css`
    --styled-select-arrow-zone__width: 32px;

    --styled-select-arrow__color: #9b9ba5;
    --styled-select-arrow__size: 10;

    --styled-select-clear-zone__width: 17px;

    --styled-select-clear__color: #999;
    --styled-select-clear__font-size: 14px;

    --styled-select-control__border-color: ${colors.neutralLighter};
    --styled-select-control__border-color--focused: #40a3f5;
    --styled-select-control__cursor--disabled: not-allowed;
    --styled-select-control__min-height: 36px;

    --styled-select-input__height: 23px;
    --styled-select-input__line-height: 23px;
    --styled-select-input__padding: 12px;

    --styled-select-menu-outer__background-color: #fff;
    --styled-select-menu-outer__border-color: ${colors.neutralLight};
    --styled-select-menu-outer__border-radius: 4px;
    --styled-select-menu-outer__border-style: solid;
    --styled-select-menu-outer__border-width: 1px;
    --styled-select-menu-outer__box-shadow: ${shadows.medium};
    --styled-select-menu-outer__margin: 5px 0 0 0;
    --styled-select-menu-outer__max-height: 200px;
    --styled-select-menu-outer__padding: 0;

    --styled-select-menu__border-radius: 2px;
    --styled-select-menu__max-height: 198px;

    --styled-select-multi-value-wrapper__padding: 3px 0 3px 5px;

    --styled-select-multi-value__background-color: #eee;
    --styled-select-multi-value__border: 1px solid #aaa;
    --styled-select-multi-value__border--hover: 1px solid #777;
    --styled-select-multi-value__border-radius: 3px;
    --styled-select-multi-value__box-shadow: rgba(0, 0, 0, 0.2) 0px 0px 3px;
    --styled-select-multi-value__font-size: 0.9em;
    --styled-select-multi-value__line-height: 1.4;
    --styled-select-multi-value__margin: 2px 5px 2px 0;

    --styled-select-no-results__color: #999;
    --styled-select-no-results__font-family: Lato, sans-serif;
    --styled-select-no-results__font-size: 14px;
    --styled-select-no-results__padding: 8px 10px;

    --styled-select-option__background-color: #fff;
    --styled-select-option__background-color--focused: ${colors.primaryLightest};
    --styled-select-option__background-color--selected: #fff;
    --styled-select-option__color: ${colors.neutralDarker};
    --styled-select-option__color--focused: ${colors.neutralDarker};
    --styled-select-option__color--selected: ${colors.neutralDarker};
    --styled-select-option__font-family: Lato, sans-serif;
    --styled-select-option__padding: 16px 14px;

    --styled-select-placeholder__color: ${colors.neutralDarker};
    --styled-select-placeholder__font-family: Roboto, sans-serif;
    --styled-select-placeholder__font-size: 14px;
    --styled-select-placeholder__line-height: 48px;
    --styled-select-placeholder__padding: 0 14px;

    --styled-select-value-icon__background-color: transparent;
    --styled-select-value-icon__background-color--hover: rgba(0, 0, 0, 0.1);
    --styled-select-value-icon__font-family: arial;
    --styled-select-value-icon__padding: 1px 5px;

    --styled-select-value-label__font-family: Tahoma, Helvetica, Arial, sans-serif;
    --styled-select-value-label__padding: 1px 6px;

    --styled-select-value-wrapper__align-content: space-around;
    --styled-select-value-wrapper__align-items: center;
    --styled-select-value-wrapper__box-sizing: border-box;
    --styled-select-value-wrapper__display: flex;
    --styled-select-value-wrapper__flex: 2 100%;
    --styled-select-value-wrapper__padding: 0 0 0 5px;

    --styled-select-value__color: Tahoma, Helvetica, Arial, sans-serif;
    --styled-select-value__font-size: 14px;
    --styled-select-value__line-height: 34px;
    --styled-select-value__max-width: 100%;
    --styled-select-value__overflow: hidden;
    --styled-select-value__padding: 0 5px;
    --styled-select-value__text-overflow: ellipsis;
    --styled-select-value__white-space: nowrap;

    --styled-select__background-color: ${colors.white};
    --styled-select__border-radius: 2px;
    --styled-select__border-style: solid;
    --styled-select__border-width: 1px;
    --styled-select__box-sizing: border-box;
    --styled-select__color: #777;
    --styled-select__cursor--disabled: not-allowed;
    --styled-select__opacity--disabled: 0.5;
    --styled-select__pointer-events--disabled: none;
    --styled-select__position: relative;
  `}
`
