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
import Menu from "@material-ui/core/Menu";
import MenuItem from "@material-ui/core/MenuItem";
import Link from "@material-ui/core/Link";
import { Redirect, useLocation } from "react-router-dom";
import ImageSearchIcon from "@material-ui/icons/ImageSearch";
import Tooltip from "@material-ui/core/Tooltip";

import * as model from "./contents/Model";

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
    title: {
      fontWeight: "bold",
      fontSize: 32,
      letterSpacing: 0.5,
      fontFamily: ["Kanit"].join(","),
      color: "#fff",
      textDecoration: "none",
    },
  });

type headerProps = {
  enablePackageSearch?: boolean;
};

function Header(props: headerProps) {
  const classes = makeStyles(styles)();
  const [errMsg, setErrMsg] = useState<string>();
  const [user, setUser] = useState<model.user>();
  const [menuAnchorEl, setMenuAnchorEl] =
    React.useState<null | HTMLElement>(null);
  const [callback, setCallback] = useState<string>("");

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
  const routerLocation = useLocation();
  useEffect(() => {
    setCallback(location.hash.substr(2));
  }, [routerLocation]);

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

  const handleMenuClose = () => {
    setMenuAnchorEl(null);
  };
  const handleMenuClick = (event: React.MouseEvent<HTMLButtonElement>) => {
    setMenuAnchorEl(event.currentTarget);
  };

  const renderSearchIcon = () => {
    if (!props.enablePackageSearch) {
      return;
    }

    return (
      <Grid item>
        <Tooltip title="Search packages">
          <IconButton color="inherit" component={RouterLink} to="/package">
            <ImageSearchIcon />
          </IconButton>
        </Tooltip>
      </Grid>
    );
  };
  const renderLoginStatus = () => {
    if (user) {
      return (
        <Grid item>
          <IconButton
            color="inherit"
            onClick={handleMenuClick}
            className={classes.iconButtonAvatar}>
            <Avatar src={user.AvatarURL} alt={user.Name} />
          </IconButton>
          <Menu
            id="simple-menu"
            anchorEl={menuAnchorEl}
            keepMounted
            open={Boolean(menuAnchorEl)}
            onClose={handleMenuClose}>
            <MenuItem onClick={handleMenuClose}>
              <Link href="auth/logout">Logout</Link>
            </MenuItem>
          </Menu>
        </Grid>
      );
    } else {
      return (
        <Button
          size="small"
          variant="contained"
          href={`auth/github?callback=${callback}`}
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
              <Typography color="inherit" variant="h4">
                <RouterLink to="/" className={classes.title}>
                  Octovy
                </RouterLink>
              </Typography>
            </Grid>

            {renderSearchIcon()}
            {renderLoginStatus()}
          </Grid>
        </Toolbar>
      </AppBar>

      {renderLoginErrorMessage()}
    </React.Fragment>
  );
}

export default withStyles(styles)(Header);
