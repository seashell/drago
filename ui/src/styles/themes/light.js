const colors = {
  white: '#fff',
  black: '#131d2a',
  transparent: 'rgba(0,0,0,0)',
  primaryDark: '#131d29',
  primary: '#00a1ff',
  primaryLight: 'rgba(67, 209, 255, 0.5)',
  primaryLighter: 'rgba(67, 209, 255, 0.2)',
  primaryLightest: 'rgba(67, 209, 255, 0.1)',
  secondary: '#2fe98c',
  neutralDarkest: '#555c66',
  neutralDarker: '#666d77',
  neutralDark: '#8d949b',
  neutral: '#b1b6bd',
  neutralLight: '#c3c9ce',
  neutralLighter: '#ececf0',
  neutralLightest: '#fafbfc',
  green: '#00c582',
  yellow: '#ffc85e',
  violet: '#c66ce0',
  blue: '#43d1ff',
  purple: '#8b74ff',
  orange: '#ff824c',
  success: '#21ba45',
  warning: '#fbbd08',
  danger: '#ff5353',
}

const gradients = {
  primary: 'linear-gradient(45deg, #00a1ff 0%, #1effc3 100%)',
}

const shadows = {
  light: '1px 2px 4px 0 rgba(0, 0, 0, 0.05)',
  medium: '0px 2px 8px rgba(0, 0, 0, 0.1);',
  heavy: '1px 2px 4px 0 rgba(0, 0, 0, 0.2)',
  primary: '0 4px 9px rgba(0, 161, 255, 0.33)',
}

const borders = {
  thin: `1px solid`,
  medium: `2px solid`,
  thick: `3px solid`,
  weird: `4px solid ${colors.primary}`,
  discrete: `2px solid ${colors.neutralLighter}`,
}

const breakpoints = ['200px', '52em', '64em']

const mediaQueries = {
  small: `@media screen and (min-width: ${breakpoints[0]})`,
  medium: `@media screen and (min-width: ${breakpoints[1]})`,
  large: `@media screen and (min-width: ${breakpoints[2]})`,
}

export default {
  colors,
  shadows,
  borders,
  containers: {
    'grid-12': {
      display: 'grid',
      gridTemplate: 'auto / repeat(12, 76px)',
      gridGap: '24px',
    },
    'grid-6': {
      display: 'grid',
      gridTemplate: 'auto / repeat(6, 72px)',
      gridGap: '16px',
    },
    padded: {
      padding: '24px 32px',
    },
  },
  buttons: {
    glowing: {
      background: gradients.primary,
      boxShadow: shadows.primary,
      color: colors.white,
    },
    primary: { background: colors.primary, color: colors.white },
    secondary: { background: colors.secondary, color: colors.white },
    neutral: { background: colors.neutralLighter, color: colors.neutralDarker },
    warning: { background: colors.warning, color: colors.white },
    danger: { background: colors.danger, color: colors.white },
    primaryInverted: {
      background: colors.transparent,
      border: borders.thin,
      borderColor: colors.primary,
      color: colors.primary,
    },
    dangerInverted: {
      background: colors.transparent,
      color: colors.danger,
      border: borders.thin,
      borderColor: colors.danger,
    },
  },
  colorStyles: {
    darkHighContrast: {
      background: colors.primaryDark,
      color: colors.neutralLight,
      borderColor: colors.neutralLight,
      fill: colors.neutralLight,
    },
    lightHighContrast: {
      background: colors.white,
      color: colors.neutralDark,
      borderColor: colors.neutralDark,
      fill: colors.neutralDark,
    },
  },
  textStyles: {
    title: {
      fontSize: 36,
      fontFamily: 'Raleway',
      fontWeight: 'normal',
      lineHeight: '40px',
      color: colors.neutralDarkest,
    },
    subtitle: {
      fontSize: 24,
      fontFamily: 'Raleway',
      fontWeight: 'bold',
      lineHeight: '38px',
      color: colors.neutralDarker,
    },
    description: {
      fontSize: 16,
      fontFamily: 'sans-serif',
      fontWeight: '400',
      lineHeight: '24px',
      color: colors.neutralDark,
    },
    bodyText: {
      fontSize: 16,
      fontFamily: 'Lato',
      fontWeight: '400',
      lineHeight: '18px',
      color: colors.neutralDarker,
    },
    detail: {
      fontSize: 14,
      fontFamily: 'sans-serif',
      fontWeight: '400',
      lineHeight: '14px',
      color: colors.neutralDarker,
    },
    code: {
      fontSize: 14,
      fontFamily: 'open-sans',
      fontWeight: '400',
      lineHeight: '16px',
      color: colors.neutralDark,
    },
  },
  fontSizes: [10, 12, 14, 16, 20, 24, 32, 48, 64, 72],
  sizes: [0.5, 1, 2, 4, 8, 16, 32, 48, 64, 72, 96, 128],
  space: [0, 4, 8, 16, 32, 64, 128, 256, 512],
  breakpoints,
  mediaQueries,
}
