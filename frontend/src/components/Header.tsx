import React, { useEffect, useState } from "react";
import AppBar from "@material-ui/core/AppBar";
import Grid from "@material-ui/core/Grid";
import Toolbar from "@material-ui/core/Toolbar";
import Typography from "@material-ui/core/Typography";
import {
  createStyles,
  Theme,
  withStyles,
  makeStyles,
} from "@material-ui/core/styles";
import Button from "@material-ui/core/Button";
import { Route, Switch, Link as RouterLink } from "react-router-dom";
import GitHubIcon from "@material-ui/icons/GitHub";
import Alert from "@material-ui/lab/Alert";
import Avatar from "@material-ui/core/Avatar";
import IconButton from "@material-ui/core/IconButton";

import * as model from "./contents/Model";
import { error } from "console";

const lightColor = "rgba(255, 255, 255, 0.7)";

const styles = (theme: Theme) =>
  createStyles({
    secondaryBar: {
      zIndex: 0,
    },
    menuButton: {
      marginLeft: -theme.spacing(1),
    },
    iconButtonAvatar: {
      padding: 4,
    },
    link: {
      textDecoration: "none",
      color: lightColor,
      "&:hover": {
        color: theme.palette.common.white,
      },
    },
    button: {
      borderColor: lightColor,
    },
  });

function Header() {
  const classes = makeStyles(styles)();
  const [errMsg, setErrMsg] = useState<string>();
  const [user, setUser] = useState<model.user>();

  useEffect(() => {
    const hashParts = window.location.hash.split("?");
    if (hashParts.length === 2) {
      const qs = new URLSearchParams(hashParts[1]);
      const err = qs.get("login_error");
      if (err) {
        setErrMsg(err);
      }
    }
  }, []);

  useEffect(() => {
    fetch("api/v1/user")
      .then((res) => res.json())
      .then(
        (result) => {
          setUser(result.data);
        },
        (error) => {
          setErrMsg(error);
        }
      );
  }, []);

  const renderLoginErrorMessage = () => {
    if (errMsg) {
      return (
        <Alert
          severity="error"
          onClose={() => {
            setErrMsg(undefined);
          }}>
          {errMsg}
        </Alert>
      );
    } else {
      return;
    }
  };

  const renderLoginStatus = () => {
    if (user) {
      console.log({ user });
      return (
        <IconButton color="inherit" className={classes.iconButtonAvatar}>
          <Avatar src={user.AvatarURL} alt={user.Name} />
        </IconButton>
      );
    } else {
      return (
        <Button
          size="small"
          variant="contained"
          href="auth/github"
          startIcon={<GitHubIcon />}>
          Login with GitHub
        </Button>
      );
    }
  };

  return (
    <React.Fragment>
      <AppBar color="primary" position="sticky" elevation={0}>
        <Toolbar>
          <Grid container spacing={1} alignItems="center">
            <Grid item xs>
              <Typography color="inherit" variant="h5" component="h1">
                <Switch>
                  <Route path="/repository">Repository</Route>
                  <Route path="/package">Package</Route>
                  <Route path="/vuln">Vulnerability</Route>
                  <Route path="/scan/report/">Scan Report</Route>
                </Switch>
              </Typography>
            </Grid>
            <Grid item>{renderLoginStatus()}</Grid>
          </Grid>
        </Toolbar>
      </AppBar>
      {renderLoginErrorMessage()}
    </React.Fragment>
  );
}

export default withStyles(styles)(Header);
