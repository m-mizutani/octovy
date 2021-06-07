import React, { useEffect } from "react";

import Toolbar from "@material-ui/core/Toolbar";
import Paper from "@material-ui/core/Paper";
import AppBar from "@material-ui/core/AppBar";
import Grid from "@material-ui/core/Grid";
import TextField from "@material-ui/core/TextField";
import Tooltip from "@material-ui/core/Tooltip";
import IconButton from "@material-ui/core/IconButton";
import RefreshIcon from "@material-ui/icons/Refresh";
import Table from "@material-ui/core/Table";
import TableBody from "@material-ui/core/TableBody";
import TableCell from "@material-ui/core/TableCell";
import TableContainer from "@material-ui/core/TableContainer";
import TableHead from "@material-ui/core/TableHead";
import TableRow from "@material-ui/core/TableRow";
import Alert from "@material-ui/lab/Alert";
import CircularProgress from "@material-ui/core/CircularProgress";
import Link from "@material-ui/core/Link";
import Chip from "@material-ui/core/Chip";
import { Redirect, useLocation } from "react-router-dom";

import Typography from "@material-ui/core/Typography";
import { useParams } from "react-router-dom";
import { Link as RouterLink } from "react-router-dom";

import List from "@material-ui/core/List";
import ListItem from "@material-ui/core/ListItem";
import ListItemIcon from "@material-ui/core/ListItemIcon";
import ListItemText from "@material-ui/core/ListItemText";
import FolderIcon from "@material-ui/icons/Folder";
import Divider from "@material-ui/core/Divider";
import Checkbox from "@material-ui/core/Checkbox";

import strftime from "strftime";

import * as model from "./Model";
import useStyles from "./Style";

import { ClassNameMap } from "@material-ui/styles";
import { createStyles, makeStyles, Theme } from "@material-ui/core/styles";

const repoStyles = makeStyles((theme: Theme) =>
  createStyles({
    ownerListGrid: {
      minWidth: "256px",
    },
    ownerItemIcon: {
      minWidth: "auto",
      marginRight: theme.spacing(2),
    },
    ownerSearchBoxToolBar: {
      minHeight: "94px",
    },
    repositoryListGrid: {
      maxWidth: "1280px",
    },
    repositoryListTable: {
      marginTop: "30px",
    },
    pkgChip: {
      marginLeft: theme.spacing(1),
      marginBottom: theme.spacing(1),
    },
    tgNoData: {
      margin: theme.spacing(10),
    },
    noDataDescr: {
      fontSize: "16px",
      margin: theme.spacing(5),
    },
  })
);

interface errorResponse {
  Error: string;
}

interface repoState {
  error?: errorResponse;
  isLoaded?: boolean;
  displayed?: model.repository[];
  allItems?: model.repository[];
}

function Owners() {
  const classes = repoStyles();
  const common = useStyles();

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
          <CircularProgress size={16} className={common.progressIcon} />
        </Alert>
      );
    } else if (status.error) {
      <Alert severity="error">Error: {status.error}</Alert>;
    } else if (status.owners) {
      return (
        <div>
          <Divider />
          <List dense={true}>
            {status.owners.map((owner, idx) => {
              return (
                <ListItem key={`owner-${idx}`}>
                  <ListItemIcon className={classes.ownerItemIcon}>
                    <FolderIcon />
                  </ListItemIcon>
                  <RouterLink to={`/repository/${owner.Name}`}>
                    <ListItemText primary={owner.Name} />
                  </RouterLink>
                </ListItem>
              );
            })}
          </List>
        </div>
      );
    }
  };

  React.useEffect(fetchOwners, []);

  return (
    <div>
      <Typography variant="h6" className={common.typographyTitle}>
        Owners
      </Typography>
      {renderOwners()}
    </div>
  );
}

function Repositories() {
  const classes = repoStyles();

  const { owner } = useParams();
  const [inputOwner, setInputOwner] = React.useState<string>("");
  const [redirectTo, setRedirectTo] = React.useState<string>();
  const [repoState, setRepoState] = React.useState<repoState>({});
  const [filterScan, setFilterScan] = React.useState<boolean>(true);
  console.log("owner:", owner);
  const doRedirect = () => {
    if (redirectTo) {
      return <Redirect to={`/repository/${redirectTo}`} />;
    }
  };

  const filterRepos = (repos: model.repository[]): model.repository[] => {
    if (repos === undefined) {
      return [];
    }

    return repos.filter((repo) => {
      return !filterScan || repo.Branch.LastScannedAt > 0;
    });
  };

  const reloadRepoState = () => {
    setRepoState({ isLoaded: false, allItems: undefined });
    setInputOwner(owner || "");
    if (!owner) {
      return;
    }

    fetch(`api/v1/repo/${owner}`)
      .then((res) => res.json())
      .then(
        (result) => {
          console.log(result);
          setRepoState({
            isLoaded: true,
            displayed: filterRepos(result.data),
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

  React.useEffect(reloadRepoState, [owner]);
  React.useEffect(() => {
    setRepoState({
      isLoaded: true,
      displayed: filterRepos(repoState.allItems),
      allItems: repoState.allItems,
    });
  }, [filterScan]);

  const renderSearchBox = () => {
    return (
      <AppBar position="static" elevation={0} color="inherit">
        <Toolbar className={classes.ownerSearchBoxToolBar}>
          <Grid container spacing={2} alignItems="center">
            <Grid item xs>
              <TextField
                label="User / Organization"
                fullWidth
                value={inputOwner}
                onChange={(e) => {
                  setInputOwner(e.target.value as string);
                }}
                onKeyPress={(e) => {
                  if (e.code === "Enter") {
                    setRedirectTo(inputOwner);
                  }
                }}
                variant="outlined"
              />
            </Grid>
            <Grid item>
              <Tooltip title="Reload">
                <IconButton onClick={reloadRepoState}>
                  <RefreshIcon color="inherit" />
                </IconButton>
              </Tooltip>
            </Grid>
          </Grid>
        </Toolbar>
      </AppBar>
    );
  };

  const renderRepoIcon = (repo: model.repository) => {
    if (repo.Branch.ReportSummary.PkgCount === 0) {
      return "☑️";
    } else if (repo.Branch.ReportSummary.VulnCount === 0) {
      return "✅";
    } else {
      return "⚠️";
    }
  };

  const renderRepositories = () => {
    if (!owner) {
      return;
    } else if (repoState.error) {
      return <div>Error: {repoState.error.Error}</div>;
    } else if (!repoState.isLoaded) {
      return <Alert severity="info">Loading...</Alert>;
    } else if (repoState.displayed.length === 0) {
      return <NoData owner={owner} />;
    } else {
      return (
        <div>
          <Grid>
            <Checkbox
              checked={filterScan}
              onChange={(e) => {
                setFilterScan(e.target.checked);
              }}
              color="primary"
            />
            Only scanned
          </Grid>
          <TableContainer
            component={Paper}
            className={classes.repositoryListTable}>
            <Table aria-label="repo table" size="small">
              <TableHead>
                <TableRow>
                  <TableCell>Repository</TableCell>
                  <TableCell align="right">Last scanned</TableCell>
                  <TableCell align="right">Package types</TableCell>
                  <TableCell align="right">Packages</TableCell>
                  <TableCell align="right">Vulnerabilities</TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {repoState.displayed.map((item) => (
                  <TableRow key={item.Owner + "/" + item.RepoName}>
                    <TableCell component="th" scope="row">
                      {renderRepoIcon(item)}{" "}
                      <RouterLink
                        to={`/repository/${item.Owner}/${item.RepoName}`}
                        style={{ textDecoration: "none" }}>
                        {item.Owner + "/" + item.RepoName}
                      </RouterLink>{" "}
                      (<Link href={item.URL}>github</Link>)
                    </TableCell>
                    <TableCell align="right">
                      {renderUnixTime(item.Branch.LastScannedAt)}
                    </TableCell>
                    <TableCell align="right">
                      {renderPackageTypes(
                        item.Branch.ReportSummary.PkgTypes,
                        classes
                      )}
                    </TableCell>
                    <TableCell align="right">
                      {item.Branch.LastScannedAt
                        ? item.Branch.ReportSummary.PkgCount
                        : undefined}
                    </TableCell>
                    <TableCell align="right">
                      {item.Branch.LastScannedAt
                        ? item.Branch.ReportSummary.VulnCount
                        : undefined}
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
      {doRedirect()}
      {renderSearchBox()}
      {renderRepositories()}
    </div>
  );
}

type ContentProps = {
  ownerList?: boolean;
};

export function Content(props: ContentProps) {
  const common = useStyles();
  const classes = repoStyles();

  const renderContent = () => {
    if (props.ownerList) {
      return (
        <Grid container spacing={4}>
          <Grid item className={classes.ownerListGrid}>
            <Paper elevation={3} square className={common.paper}>
              <Grid className={common.contentWrapper}>
                <Owners />
              </Grid>
            </Paper>
          </Grid>
          <Grid item xs className={classes.repositoryListGrid}>
            <Paper elevation={3} square className={common.paper}>
              <Grid className={common.contentWrapper}>
                <Repositories />
              </Grid>
            </Paper>
          </Grid>
        </Grid>
      );
    } else {
      return (
        <Grid container spacing={4} alignItems="center" justify="center">
          <Grid item xs className={classes.repositoryListGrid}>
            <Paper elevation={3} square className={common.paper}>
              <Grid className={common.contentWrapper}>
                <Repositories />
              </Grid>
            </Paper>
          </Grid>
        </Grid>
      );
    }
  };

  return (
    <Grid item xs={12}>
      {renderContent()}
    </Grid>
  );
}

function renderUnixTime(ts: number) {
  if (ts === 0) {
    return <div>N/A</div>;
  }

  const dt = new Date(ts * 1000);
  return <div>{strftime("%F %T", dt)}</div>;
}

function renderPackageTypes(pkgTypes: string[], classes: ClassNameMap) {
  if (pkgTypes.length === 0) {
    return;
  }

  return (
    <div>
      {pkgTypes.map((t, i) => {
        return (
          <Chip key={i} label={t} size="small" className={classes.pkgChip} />
        );
      })}
    </div>
  );
}

type NoDataProps = {
  owner: string;
};

function NoData(props: NoDataProps) {
  const classes = repoStyles();
  interface metaData {
    AppURL: string;
  }
  const [meta, setMeta] = React.useState<metaData>();
  const getMeta = () => {
    fetch(`api/v1/meta/octovy`)
      .then((res) => res.json())
      .then(
        (result) => {
          console.log(result);
          setMeta(result.data);
        },
        (error) => {
          console.log(error);
        }
      );
  };
  useEffect(getMeta, []);

  return (
    <div className={classes.tgNoData}>
      <Typography variant="h4" align="center">
        No Data of "{props.owner}"
      </Typography>
      <Typography align="center" className={classes.noDataDescr}>
        Try to install <Link href={meta ? meta.AppURL : ""}>octovy</Link> to
        your repository
      </Typography>
    </div>
  );
}
