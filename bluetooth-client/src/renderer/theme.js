import { createMuiTheme } from '@material-ui/core/styles';

const theme = createMuiTheme({
  spacing: 4,
  palette: {
    primary: {
      main: '#00838e',
      dark: '#005661',
      light: '#4fb3be'
    },
    secondary: {
      main: '#00838e',
      dark: '#005661',
      light: '#4fb3be'
    },
    error: {
      main: '#A40000',
    },
    background: {
      default: '#303030',
      paper: '#424242'
    },
  },
});

export default theme;