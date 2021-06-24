import React, { useEffect } from "react";

import Paper from "@material-ui/core/Paper";
import Toolbar from "@material-ui/core/Toolbar";
import AppBar from "@material-ui/core/AppBar";
import Grid from "@material-ui/core/Grid";
import GitHubIcon from "@material-ui/icons/GitHub";
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

import FormControl from "@material-ui/core/FormControl";
import Select from "@material-ui/core/Select";
import InputLabel from "@material-ui/core/InputLabel";
import MenuItem from "@material-ui/core/MenuItem";
import Avatar from "@material-ui/core/Avatar";

import Typography from "@material-ui/core/Typography";
import { useParams } from "react-router-dom";
import { Link as RouterLink } from "react-router-dom";
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
    ownerTitle: {
      fontSize: "32px",
      fontWeight: "bold",
      letterSpacing: 0.5,
    },
    ownerGrid: {
      marginTop: theme.spacing(4),
      marginBottom: theme.spacing(0),
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
    formControl: {
      margin: theme.spacing(1),
      minWidth: 256,
    },
    largeAvatar: {
      width: theme.spacing(5),
      height: theme.spacing(5),
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

type OwnersProps = {
  selected?: string;
};

function Owners(props: OwnersProps) {
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
  const [redirectTo, setRedirectTo] = React.useState<string>();

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

  React.useEffect(fetchOwners, []);

  const doRedirect = () => {
    if (redirectTo) {
      return <Redirect to={`/repository/${redirectTo}`} />;
    }
  };

  const handleChange = (event: React.ChangeEvent<{ value: unknown }>) => {
    setRedirectTo(event.target.value as string);
  };

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
      <FormControl className={classes.formControl}>
        {doRedirect()}
        <InputLabel id="owner-select-label">Owner</InputLabel>
        <Select
          labelId="owner-select-label"
          id="owner-select"
          value={props.selected || ""}
          onChange={handleChange}>
          {status.owners.map((owner, idx) => (
            <MenuItem key={idx} value={owner.Name}>
              {owner.Name}
            </MenuItem>
          ))}
        </Select>
      </FormControl>
    );
  }
}

function Repositories() {
  const classes = repoStyles();
  const repoClasses = repoStyles();

  const { owner } = useParams();
  const [repoState, setRepoState] = React.useState<repoState>({});
  const [filterScan, setFilterScan] = React.useState<boolean>(true);

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

  const renderRepositories = () => {
    const renderRepoIcon = (repo: model.repository) => {
      if (repo.Branch.ReportSummary.PkgCount === 0) {
        return "☑️";
      } else if (repo.Branch.ReportSummary.VulnCount === 0) {
        return "✅";
      } else {
        return "⚠️";
      }
    };

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
        <React.Fragment>
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
        </React.Fragment>
      );
    }
  };

  return (
    <React.Fragment>
      <Grid>
        <Owners selected={owner} />
      </Grid>
      {owner ? (
        <Grid
          container
          spacing={2}
          alignItems="center"
          className={repoClasses.ownerGrid}>
          <Grid item>
            {owner.AvatarURL ? (
              <Avatar
                src={owner.AvatarURL}
                alt={owner.Name}
                className={classes.largeAvatar}
              />
            ) : (
              <GitHubIcon color="inherit" fontSize="large" />
            )}
          </Grid>
          <Grid item xs>
            <Typography className={repoClasses.ownerTitle}>{owner}</Typography>
          </Grid>
        </Grid>
      ) : undefined}

      {renderRepositories()}
    </React.Fragment>
  );
}

type ContentProps = {
  ownerList?: boolean;
};

export function Content(props: ContentProps) {
  const common = useStyles();
  const classes = repoStyles();

  return (
    <Grid item xs={12}>
      <Grid container spacing={4} alignItems="center" justify="center">
        <Grid item xs className={classes.repositoryListGrid}>
          <Paper elevation={3} square className={common.paper}>
            <Grid className={common.contentWrapper}>
              <Repositories />
            </Grid>
          </Paper>
        </Grid>
      </Grid>
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
