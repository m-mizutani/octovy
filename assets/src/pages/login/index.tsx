import * as React from "react";
import * as app from "@/src/components/app";
import Alert from "@mui/material/Alert";
import Button from "@mui/material/Button";
import GitHubIcon from "@mui/icons-material/GitHub";
import Grid from "@mui/material/Grid";

import { useRouter } from "next/router";

export function Login() {
  const router = useRouter();
  const { login_error } = router.query;
  const callbackTo = router.query.callback;
  const githubLogin = `/auth/github${
    callbackTo ? `?callback=${callbackTo}` : ""
  }`;
  return (
    <app.Main>
      {login_error ? (
        <Grid container justifyContent="center" style={{ marginTop: 30 }}>
          <Grid item minWidth={480}>
            <Alert severity="error">{login_error}</Alert>
          </Grid>
        </Grid>
      ) : (
        ""
      )}

      <Grid container justifyContent="center">
        <Grid item style={{ marginTop: 50 }}>
          <Button
            variant="outlined"
            startIcon={<GitHubIcon />}
            href={githubLogin}>
            Login with GitHub
          </Button>
        </Grid>
      </Grid>
    </app.Main>
  );
}
