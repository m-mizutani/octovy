import * as React from "react";

import Typography from "@mui/material/Typography";
import Box from "@mui/material/Box";
import AppBar from "@mui/material/AppBar";
import Toolbar from "@mui/material/Toolbar";
import Container from "@mui/material/Container";
import Paper from "@mui/material/Paper";

import Avatar from "@mui/material/Avatar";
import LoginIcon from "@mui/icons-material/Login";
import IconButton from "@mui/material/IconButton";
import MenuItem from "@mui/material/MenuItem";
import Menu from "@mui/material/Menu";

import { ThemeProvider } from "@mui/material/styles";

import Link from "next/link";
import { useRouter } from "next/router";

import theme from "./theme";
import Head from "./head";
import * as model from "./model";

type mainProp = {
  children?: React.ReactNode;
};

export function Main(props: mainProp) {
  const [user, setUser] = React.useState<model.user>();
  const [anchorEl, setAnchorEl] = React.useState(null);
  const router = useRouter();

  const getUser = () => {
    fetch(`/api/v1/user`)
      .then((res) => res.json())
      .then(
        (resp) => {
          console.log("get user resp:", { resp });
          setUser(resp.data);
        },
        (error) => {
          console.log("ignore error:", { error });
        }
      );
  };

  React.useEffect(getUser, []);

  const renderAvatar = () => {
    return user ? (
      <IconButton
        onClick={(e) => {
          setAnchorEl(e.currentTarget);
        }}>
        <Avatar
          alt={user.name}
          src={user.avatar_url}
          sx={{ width: 32, height: 32 }}
        />
      </IconButton>
    ) : (
      <Avatar sx={{ width: 32, height: 32 }}>
        <Link href="/login">
          <LoginIcon />
        </Link>
      </Avatar>
    );
  };

  const menuId = "user-menu-id";

  const handleLogout = () => {
    window.location.href = "/auth/logout";
    handleMenuClose();
  };
  const handleMenuClose = () => {
    setAnchorEl(null);
  };
  const renderUserMenu = (
    <Menu
      anchorEl={anchorEl}
      anchorOrigin={{
        vertical: "top",
        horizontal: "right",
      }}
      id={menuId}
      keepMounted
      transformOrigin={{
        vertical: "top",
        horizontal: "right",
      }}
      open={anchorEl !== null}
      onClose={handleMenuClose}>
      <MenuItem onClick={handleLogout}>Logout</MenuItem>
    </Menu>
  );

  return (
    <div>
      <Head />
      <ThemeProvider theme={theme}>
        <div
          style={{
            display: "flex",
            minHeight: "100vh",
            padding: 0,
            margin: 0,
          }}>
          <Box sx={{ flexGrow: 1 }}>
            <AppBar position="static">
              <Toolbar>
                <Typography variant="h5">
                  <a
                    href="/"
                    style={{ color: "inherit", textDecoration: "inherit" }}>
                    Octovy
                  </a>
                </Typography>

                <Box sx={{ flexGrow: 1 }} />
                <Box sx={{ display: { xs: "none", md: "flex" } }}>
                  {renderAvatar()}
                </Box>
              </Toolbar>
            </AppBar>
            <main
              style={{
                flex: 1,
                padding: theme.spacing(6, 4),
                background: "#eaeff1",
              }}>
              <Container>
                <Paper
                  elevation={3}
                  style={{ padding: 20, minHeight: "100vh" }}>
                  {props.children}
                </Paper>
              </Container>
            </main>
          </Box>
        </div>
      </ThemeProvider>
      {renderUserMenu}
    </div>
  );
}
