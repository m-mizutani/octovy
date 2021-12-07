import React from "react";
import { Grid, Alert, Typography } from "@mui/material";
import Stack from "@mui/material/Stack";

import * as app from "@/components/app";
import * as model from "@/components/model";

import Labels from "./label";
import Scans from "./scan";
import VulnStatuses from "./vulnStatus";

export default function Repository(props: { owner: string; repo: string }) {
  type repoStatus = {
    isLoaded: boolean;
    repository?: model.repository;
    err?: any;
  };
  const [status, setStatus] = React.useState<repoStatus>({
    isLoaded: false,
  });

  const getRepo = () => {
    if (props.owner === undefined || props.repo === undefined) {
      return;
    }

    fetch(`/api/v1/repository/${props.owner}/${props.repo}`)
      .then((res) => res.json())
      .then(
        (result) => {
          console.log("repo:", { result });
          if (result.error) {
            setStatus({ isLoaded: true, err: result.error });
          } else {
            setStatus({
              isLoaded: true,
              repository: result.data,
            });
          }
        },
        (error) => {
          console.log("got error:", { error });
          setStatus({
            isLoaded: true,
            err: error.message,
          });
        }
      );
  };
  React.useEffect(getRepo, [props.owner, props.repo]);

  const render = () => {
    if (!status.isLoaded) {
      return <Alert severity="info">Loading...</Alert>;
    } else if (status.err) {
      return <Alert severity="error">{status.err}</Alert>;
    }

    return (
      <>
        <Grid marginTop={5}>
          <Typography variant="h5">Recent scan reports</Typography>
          <Scans repo={status.repository} />
        </Grid>
        <Grid marginTop={5}>
          <Typography variant="h5">Status</Typography>
          <VulnStatuses repo={status.repository} />
        </Grid>
        <Grid marginTop={5}>
          <Typography variant="h5">Labels</Typography>
          <Labels repo={status.repository} />
        </Grid>
      </>
    );
  };

  return (
    <app.Main>
      <Grid>
        <Stack direction="row" spacing={2}>
          <Typography variant="h4">
            {props.owner}/{props.repo}
          </Typography>
        </Stack>
      </Grid>
      {render()}
    </app.Main>
  );
}
