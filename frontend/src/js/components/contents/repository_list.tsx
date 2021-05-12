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
import Alert from "@material-ui/lab/Alert";
import CircularProgress from "@material-ui/core/CircularProgress";
import { useLocation } from "react-router-dom";

import { useParams } from "react-router-dom";
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

function Owners() {
  const classes = useStyles();

  interface owner {
    Name: string;
  }

  type ownerStatus = {
    isLoaded: boolean;
    owners?: owner[];
    error?: any;
  };

  const [status, setStatus] = React.useState<ownerStatus>({ isLoaded: false });

  const fetchOwners = () => {
    fetch(`api/v1/repo`)
      .then((res) => res.json())
      .then(
        (result) => {
          console.log(result);
          setStatus({
            isLoaded: true,
            owners: result.data,
          });
        },
        (error) => {
          setStatus({
            isLoaded: true,
            error,
          });
        }
      );
  };

  const renderOwners = () => {
    if (!status.isLoaded) {
      return (
        <Alert severity="info">
          Loading...
          <CircularProgress size={16} className={classes.progressIcon} />
        </Alert>
      );
    } else if (status.error) {
      <Alert severity="error">Error: {status.error}</Alert>;
    } else if (status.owners) {
      return (
        <div>
          {status.owners.map((owner, idx) => {
            return (
              <div key={`owner-${idx}`}>
                <RouterLink to={`/repository/${owner.Name}`}>
                  {owner.Name}
                </RouterLink>
              </div>
            );
          })}
        </div>
      );
    }
  };

  React.useEffect(fetchOwners, []);

  return (
    <div>
      <Grid component="h2">Owners</Grid>
      {renderOwners()}
    </div>
  );
}

type RepositoriesProps = {
  owner?: string;
};

function Repositories(props: RepositoriesProps) {
  const classes = useStyles();

  const [repoState, setRepoState] = React.useState<repoState>({});
  const reloadRepoState = () => {
    if (!props.owner) {
      return;
    }

    fetch(`api/v1/repo/${props.owner}`)
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

  React.useEffect(reloadRepoState, [props.owner]);

  const renderRepositories = () => {
    if (!props.owner) {
      return <Alert severity="info">Choose owner</Alert>;
    } else if (repoState.error) {
      return <div>Error: {repoState.error.Error}</div>;
    } else if (!repoState.isLoaded) {
      return <Alert severity="info">Loading...</Alert>;
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
        <div>
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
        </div>
      );
    }
  };

  return (
    <div>
      <Grid component="h2">Repositories</Grid>
      {renderRepositories()}
    </div>
  );
}

export default function RepositoryList() {
  const classes = useStyles();
  const { owner } = useParams();

  return (
    <Grid item xs={12}>
      <Grid container spacing={4}>
        <Grid item xs={3}>
          <Paper elevation={3} square className={classes.paper}>
            <Grid className={classes.contentWrapper}>
              <Owners />
            </Grid>
          </Paper>
        </Grid>
        <Grid item xs={9}>
          <Paper elevation={3} square className={classes.paper}>
            <Grid className={classes.contentWrapper}>
              <Repositories owner={owner} />
            </Grid>
          </Paper>
        </Grid>
      </Grid>
    </Grid>
  );
}
