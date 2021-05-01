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
import Table from "@material-ui/core/Table";
import TableBody from "@material-ui/core/TableBody";
import TableCell from "@material-ui/core/TableCell";
import TableContainer from "@material-ui/core/TableContainer";
import TableHead from "@material-ui/core/TableHead";
import TableRow from "@material-ui/core/TableRow";

import { Link as RouterLink } from "react-router-dom";

import useStyles from "./style";

interface errorResponse {
  Error: string;
}

interface repoState {
  error?: errorResponse;
  isLoaded?: boolean;
  items?: repoInfo[];
  allItems?: repoInfo[];
}

interface repoInfo {
  Owner: string;
  RepoName: string;
  Branches?: string[];
  URL: string;
}

export default function RepositoryList() {
  const classes = useStyles();

  const [repoState, setRepoState] = React.useState<repoState>({});
  const reloadRepoState = () => {
    fetch("api/v1/repo")
      .then((res) => res.json())
      .then(
        (result) => {
          console.log(result);
          setRepoState({
            isLoaded: true,
            items: result.data,
            allItems: result.data,
          });
        },
        (error) => {
          setRepoState({
            isLoaded: true,
            error,
          });
        }
      );
  };

  React.useEffect(reloadRepoState, []);

  if (repoState.error) {
    return <div>Error: {repoState.error.Error}</div>;
  } else if (!repoState.isLoaded) {
    return <div>Loading...</div>;
  } else {
    const onChange = (e: any) => {
      setRepoState({
        isLoaded: true,
        items: repoState.allItems.filter((item) => {
          const key = item.Owner + "/" + item.RepoName;
          return key.includes(e.target.value);
        }),
        allItems: repoState.allItems,
      });
    };

    return (
      <Paper className={classes.paper}>
        <AppBar position="static" color="default" elevation={0}>
          <Toolbar>
            <Grid container spacing={2} alignItems="center">
              <Grid item>
                <SearchIcon color="inherit" />
              </Grid>
              <Grid item xs>
                <TextField
                  fullWidth
                  placeholder="filter"
                  onChange={onChange}
                  InputProps={{
                    disableUnderline: true,
                  }}
                />
              </Grid>
              <Grid item>
                <Tooltip title="Reload">
                  <IconButton>
                    <RefreshIcon color="inherit" />
                  </IconButton>
                </Tooltip>
              </Grid>
            </Grid>
          </Toolbar>
        </AppBar>

        <TableContainer component={Paper}>
          <Table aria-label="repo table" size="small">
            <TableHead>
              <TableRow>
                <TableCell>Repository</TableCell>
                <TableCell align="right">Branches</TableCell>
                <TableCell align="right">Link</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {repoState.items.map((item) => (
                <TableRow key={item.Owner + "/" + item.RepoName}>
                  <TableCell component="th" scope="row">
                    <RouterLink
                      to={`/repository/${item.Owner}/${item.RepoName}`}
                      style={{ textDecoration: "none" }}>
                      {item.Owner + "/" + item.RepoName}
                    </RouterLink>
                  </TableCell>
                  <TableCell align="right">{item.Branches}</TableCell>
                  <TableCell align="right">
                    <RouterLink to={item.URL}>github</RouterLink>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </TableContainer>
      </Paper>
    );
  }
}
