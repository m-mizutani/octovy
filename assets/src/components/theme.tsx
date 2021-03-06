import { createTheme } from "@mui/material/styles";

import { pink } from "@mui/material/colors";

const theme = createTheme({
  palette: {
    primary: {
      light: "#757ce8",
      main: "#3f50b5",
      dark: "#002884",
      contrastText: "#fff",
    },
    secondary: {
      light: pink[100],
      main: pink[500],
      dark: pink[800],
      contrastText: "#fff",
    },
  },
  typography: {
    fontFamily: "Helvetica",

    h1: {
      fontFamily: "Kanit",
      fontWeight: "bold",
      fontSize: 48,
      letterSpacing: 0.5,
    },
    h4: {
      fontFamily: "Kanit",
      fontWeight: "bold",
      fontSize: 32,
      letterSpacing: 0.1,
    },
    h5: {
      fontFamily: "Kanit",
      fontWeight: "bold",
      fontSize: 24,
      letterSpacing: 0.1,
    },
    h6: {
      fontFamily: "Kanit",
      fontWeight: "bold",
      fontSize: 16,
      letterSpacing: 0.1,
    },
  },
  shape: {
    borderRadius: 8,
  },
  mixins: {
    toolbar: {
      minHeight: 48,
    },
  },
});

export default theme;
