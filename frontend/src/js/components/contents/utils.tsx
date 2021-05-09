import React from "react";
import MuiAlert, { AlertProps } from "@material-ui/lab/Alert";
import useStyles from "./style";

function Alert(props: AlertProps) {
  return <MuiAlert elevation={6} variant="filled" {...props} />;
}

type messageProps = {
  msg: string;
};

export function errorMessage(props: messageProps) {
  const classes = useStyles();

  return <Alert severity="error">{props.msg}</Alert>;
}

export function warnMessage(props: messageProps) {
  const classes = useStyles();

  return <Alert severity="warning">{props.msg}</Alert>;
}
