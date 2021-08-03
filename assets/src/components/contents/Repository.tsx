import React from "react";

import Paper from "@material-ui/core/Paper";
import Toolbar from "@material-ui/core/Toolbar";
import AppBar from "@material-ui/core/AppBar";
import Grid from "@material-ui/core/Grid";
import TextField from "@material-ui/core/TextField";

import { useParams } from "react-router-dom";
import { Redirect, useLocation } from "react-router-dom";

import * as scan from "./Scan";
import useStyles from "./Style";

type repoStatus = {
  err?: any;
  reportID?: string;
};

type repoProps = {
  enablePackageLink?: boolean;
};

export function Content(props: repoProps) {
  const classes = useStyles();

  const { owner, repoName, branch } = useParams();
  const [branchInput, setBranchInput] = React.useState<string>("");
  const [redirect, setRedirect] = React.useState<string>("");
  const [status, setStatus] = React.useState<repoStatus>({});

  const getRepositoryInfo = () => {
    if (branch) {
      setBranchInput(branch);

      fetch(`api/v1/repo/${owner}/${repoName}/${branch}`)
        .then((res) => res.json())
        .then(
          (result) => {
            console.log("branch:", { result });
            setStatus({
              reportID: result.data.ReportSummary.ReportID,
            });
          },
          (error) => {
            setStatus({
              err: error,
            });
          }
        );
    } else {
      fetch(`api/v1/repo/${owner}/${repoName}`)
        .then((res) => res.json())
        .then(
          (result) => {
            console.log("repo:", { result });
            setRedirect(result.data.DefaultBranch);
          },
          (error) => {
            setStatus({
              err: error,
            });
          }
        );
    }
  };

  const onKeyUpBranch = (e: any) => {
    if (e.which === 13) {
      setRedirect(e.target.value);
    }
  };
  const onChangeBranch = (e: any) => {
    setBranchInput(e.target.value);
  };
  const doRedirect = () => {
    if (redirect) {
      return <Redirect to={`/repository/${owner}/${repoName}/${redirect}`} />;
    }
  };

  const location = useLocation();
  React.useEffect(getRepositoryInfo, [location]);

  return (
    <Paper className={classes.paper}>
      {doRedirect()}
      <AppBar position="static" color="default" elevation={1}>
        <Toolbar>
          <Grid container spacing={3} alignItems="center">
            <Grid component="h3">
              {owner}/{repoName}
            </Grid>

            <Grid item xs>
              <TextField
                value={branchInput}
                onChange={onChangeBranch}
                onKeyUp={onKeyUpBranch}
                InputProps={{
                  className: classes.searchInput,
                }}
              />
            </Grid>
          </Grid>
        </Toolbar>
      </AppBar>

      <scan.Report
        reportID={status.reportID}
        packageLink={props.enablePackageLink}
      />
    </Paper>
  );
}
