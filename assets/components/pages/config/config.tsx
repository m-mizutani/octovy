import React from "react";
import Grid from "@mui/material/Grid";
import Container from "@mui/material/Container";
import Typography from "@mui/material/Typography";

import * as app from "../../app";
import Severities from "./severity";

export default function Config() {
  return (
    <app.Main>
      <Container>
        <Grid container spacing={2}>
          <Grid item>
            <Typography variant="h4">Config</Typography>
          </Grid>
        </Grid>
      </Container>
      <Container style={{ marginTop: 48 }}>
        <Grid>
          <Typography variant="h5">Custom Severity</Typography>
        </Grid>
        <Grid style={{ margin: 15 }}>
          <Severities />
        </Grid>
      </Container>
    </app.Main>
  );
}
