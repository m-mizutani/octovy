import * as React from "react";
import Typography from "@mui/material/Typography";
import Box from "@mui/material/Box";
import AppBar from "@mui/material/AppBar";
import Toolbar from "@mui/material/Toolbar";
import Container from "@mui/material/Container";
import Paper from "@mui/material/Paper";

import { ThemeProvider } from "@mui/material/styles";

import theme from "./theme";
import Head from "./head";

type mainProp = {
  children?: React.ReactNode;
};

export function Main(props: mainProp) {
  return (
    <div>
      <Head />
      <ThemeProvider theme={theme}>
        <div
          style={{
            display: "flex",
            minHeight: "100vh",
            padding: 0,
            margin: 0,
          }}>
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
    </div>
  );
}
