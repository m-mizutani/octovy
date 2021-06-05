import React from "react";

import Paper from "@material-ui/core/Paper";
import Toolbar from "@material-ui/core/Toolbar";
import AppBar from "@material-ui/core/AppBar";
import Grid from "@material-ui/core/Grid";
import TextField from "@material-ui/core/TextField";

import { useParams } from "react-router-dom";

import * as scan from "./Scan";
import useStyles from "./Style";

type repoProps = {
  enablePackageLink?: boolean;
};

export function Content(props: repoProps) {
  const classes = useStyles();
  const { reportID } = useParams();
  return (
    <Paper className={classes.paper}>
      <scan.Report reportID={reportID} packageLink={props.enablePackageLink} />
    </Paper>
  );
}
