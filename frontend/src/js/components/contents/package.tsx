import React from "react";

import Toolbar from "@material-ui/core/Toolbar";
import Paper from "@material-ui/core/Paper";
import AppBar from "@material-ui/core/AppBar";
import Grid from "@material-ui/core/Grid";
import TextField from "@material-ui/core/TextField";
import Tooltip from "@material-ui/core/Tooltip";
import IconButton from "@material-ui/core/IconButton";
import SearchIcon from "@material-ui/icons/Search";
import RefreshIcon from "@material-ui/icons/Refresh";
import MenuItem from "@material-ui/core/MenuItem";
import FormControl from "@material-ui/core/FormControl";
import Select from "@material-ui/core/Select";

import Table from "@material-ui/core/Table";
import TableBody from "@material-ui/core/TableBody";
import TableCell from "@material-ui/core/TableCell";
import TableContainer from "@material-ui/core/TableContainer";
import TableHead from "@material-ui/core/TableHead";
import TableRow from "@material-ui/core/TableRow";

import { Redirect, useLocation } from "react-router-dom";

import useStyles from "./style";
import * as model from "./model";

interface errorResponse {
  Error: string;
}

interface packageState {
  error?: errorResponse;
  isLoaded?: boolean;
  items?: model.packageRecord[];
}

export default function Package() {
  const classes = useStyles();
  const [pkgType, setPkgType] = React.useState("bundler");
  const [pkgName, setPkgName] = React.useState("");

  const submitSearch = () => {
    const qs = new URLSearchParams({
      type: pkgType,
      name: pkgName,
    });
    setRedirect(qs);
  };

  const [redirect, setRedirect] = React.useState<URLSearchParams>(null);
  const doRedirect = () => {
    if (redirect !== null) {
      const dist = redirect.toString();
      return <Redirect to={`/package?${dist}`} />;
    }
  };

  const [pkgState, setPkgState] = React.useState<packageState>({ items: [] });
  const updatePackageList = () => {
    const hashParts = window.location.hash.split("?");
    if (hashParts.length !== 2) {
      setPkgState({ items: [] });
      setPkgType("");
      setPkgName("");
      return;
    }

    const qs = new URLSearchParams(hashParts[1]);
    setPkgType(qs.get("type"));
    setPkgName(qs.get("name"));
    const apiPath = `api/v1/package?` + qs.toString();
    fetch(apiPath)
      .then((res) => res.json())
      .then(
        (result) => {
          console.log(result);
          setPkgState({
            isLoaded: true,
            items: result.data,
          });
        },
        (error) => {
          setPkgState({
            isLoaded: true,
            error,
          });
        }
      );
  };
  const location = useLocation();
  React.useEffect(updatePackageList, [location]);

  const renderPackageList = () => {
    if (!pkgState.isLoaded) {
      return;
    } else if (pkgState.error !== undefined) {
      return <div>Error: {pkgState.error}</div>;
    } else {
      return (
        <TableContainer component={Paper}>
          <Table size="small" aria-label="a dense table">
            <TableHead>
              <TableRow>
                <TableCell align="left">Repository</TableCell>
                <TableCell align="right">Branch</TableCell>
                <TableCell align="right">Source</TableCell>
                <TableCell align="right">Version</TableCell>
                <TableCell align="right">Detected</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {pkgState.items.map((pkg, idx) => {
                return (
                  <TableRow key={idx}>
                    <TableCell component="th" scope="row">
                      {`${pkg.Detected.Owner}/${pkg.Detected.RepoName}`}
                    </TableCell>
                    <TableCell align="right">{pkg.Detected.Branch}</TableCell>
                    <TableCell align="right">{pkg.Source}</TableCell>
                    <TableCell align="right">{pkg.Version}</TableCell>
                    <TableCell align="right">{pkg.Detected.CommitID}</TableCell>
                  </TableRow>
                );
              })}
            </TableBody>
          </Table>
        </TableContainer>
      );
    }
  };

  return (
    <Paper className={classes.paper}>
      <AppBar position="static" color="default" elevation={0}>
        <Toolbar>
          <Grid container spacing={2} alignItems="center">
            <Grid item>
              <SearchIcon color="inherit" />
            </Grid>

            <Grid item>
              <FormControl className={classes.formControl}>
                <Select
                  id="pkg-type-select"
                  value={pkgType}
                  inputProps={{ "aria-label": "Without label" }}
                  onChange={(e) => {
                    setPkgType(e.target.value as string);
                  }}>
                  <MenuItem value="bundler">Bundler</MenuItem>
                  <MenuItem value="npm">NPM</MenuItem>
                  <MenuItem value="yarn">yarn</MenuItem>
                  <MenuItem value="gomod">Go Modules</MenuItem>
                  <MenuItem value="pipenv">Pipenv</MenuItem>
                </Select>
              </FormControl>
            </Grid>

            <Grid item xs>
              <TextField
                fullWidth
                placeholder="Package name"
                value={pkgName}
                onChange={(e) => {
                  setPkgName(e.target.value as string);
                }}
                InputProps={{
                  disableUnderline: true,
                }}
              />
            </Grid>
            <Grid item>
              <Tooltip title="Reload">
                <IconButton onClick={submitSearch}>
                  <RefreshIcon color="inherit" />
                </IconButton>
              </Tooltip>
            </Grid>
          </Grid>
        </Toolbar>

        <Toolbar>
          <Grid container spacing={2} alignItems="center">
            <Grid item className={classes.packageList}>
              {renderPackageList()}
            </Grid>
          </Grid>
        </Toolbar>
      </AppBar>
      {doRedirect()}
    </Paper>
  );
}
