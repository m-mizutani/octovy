import * as React from "react";
import * as app from "../components/app";
import Grid from "@mui/material/Grid";

import { Typography } from "@mui/material";

export default function Index() {
  return (
    <app.Main>
      <Grid container justifyContent="center">
        <Grid item style={{ marginTop: 100 }}>
          <Typography variant="h1">Octovy</Typography>
        </Grid>
      </Grid>
    </app.Main>
  );
}
