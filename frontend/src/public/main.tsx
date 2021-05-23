import * as React from "react";
import * as ReactDOM from "react-dom";
import * as octovy from "../components/octovy";
import Grid from "@material-ui/core/Grid";

import {
  HashRouter as Router,
  Route,
  Switch,
  Redirect,
} from "react-router-dom";

import * as repositoryList from "../components/contents/RepositoryList";
import * as repository from "../components/contents/Repository";
import * as packages from "../components/contents/Packages";
import * as vulnerability from "../components/contents/Vulnerability";
import { Link, Typography } from "@material-ui/core";
import { createStyles, makeStyles, Theme } from "@material-ui/core/styles";
import { Link as RouterLink } from "react-router-dom";

const useStyles = makeStyles((theme: Theme) =>
  createStyles({
    logoGrid: {
      marginBottom: "40px",
    },
    logoTitle: {
      fontWeight: "bold",
      fontSize: 48,
      letterSpacing: 0.5,
      fontFamily: ["Kanit"].join(","),
    },
  })
);

function App() {
  const classes = useStyles();
  return (
    <octovy.Frame hasNavigator={false}>
      <Grid
        container
        alignItems="center"
        justify="center"
        direction="column"
        className={classes.logoGrid}>
        <Grid>
          <Typography className={classes.logoTitle}>Octovy</Typography>
        </Grid>
        <Grid>
          <Typography variant="h6">
            Simple vulnerability scanner for GitHub repository
          </Typography>
        </Grid>
      </Grid>
      <Switch>
        <Route path="/repository/:owner/:repoName/:branch">
          <repository.Content />
        </Route>
        <Route path="/repository/:owner/:repoName">
          <repository.Content />
        </Route>
        <Route path="/repository/:owner">
          <repositoryList.Content />
        </Route>
        <Route path="/repository">
          <repositoryList.Content />
        </Route>
        <Route path="/vuln" exact>
          <vulnerability.Content />
        </Route>
        <Route path="/vuln/:vulnID">
          <vulnerability.Content />
        </Route>

        <Route path="/" exact>
          <repositoryList.Content />
        </Route>
      </Switch>
    </octovy.Frame>
  );
}

ReactDOM.render(<App />, document.querySelector("#app"));
