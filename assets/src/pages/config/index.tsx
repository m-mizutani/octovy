import React from "react";
import Grid from "@mui/material/Grid";
import Container from "@mui/material/Container";
import Typography from "@mui/material/Typography";

import * as app from "../../components/app";
import Severities from "./severity";
import RepoLabels from "./repoLabel";

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
      <Container style={{ marginTop: 48 }}>
        <Grid>
          <Typography variant="h5">Repository Label</Typography>
        </Grid>
        <Grid style={{ margin: 15 }}>
          <RepoLabels />
        </Grid>
      </Container>
    </app.Main>
  );
}
