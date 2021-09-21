import * as React from "react";
import {
  createTheme,
  createStyles,
  ThemeProvider,
  makeStyles,
  Theme,
} from "@mui/material/styles";
import Typography from "@mui/material/Typography";
import Box from "@mui/material/Box";
import AppBar from "@mui/material/AppBar";
import Toolbar from "@mui/material/Toolbar";
import Container from "@mui/material/Container";
import Paper from "@mui/material/Paper";

import { pink } from "@mui/material/colors";

type mainProp = {
  children?: React.ReactNode;
};

export function Main(props: mainProp) {
  return (
    <ThemeProvider theme={theme}>
      <div
        style={{ display: "flex", minHeight: "100vh", padding: 0, margin: 0 }}>
        <Box sx={{ flexGrow: 1 }}>
          <AppBar position="static">
            <Toolbar>
              <Typography variant="h5">Octovy</Typography>
            </Toolbar>
          </AppBar>
          <main
            style={{
              flex: 1,
              padding: theme.spacing(6, 4),
              background: "#eaeff1",
            }}>
            <Container>
              <Paper elevation={3} style={{ padding: 20, height: "100vh" }}>
                {props.children}
              </Paper>
            </Container>
          </main>
        </Box>
      </div>
    </ThemeProvider>
  );
}

let theme = createTheme({
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
    h5: {
      fontFamily: "Kanit",
      fontWeight: "bold",
      fontSize: 20,
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
