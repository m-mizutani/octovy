import * as React from "react";
import * as app from "../components/app";
import * as model from "../components/model";

import Alert from "@mui/material/Alert";
import Grid from "@mui/material/Grid";
import Paper from "@mui/material/Paper";
import Table from "@mui/material/Table";
import TableBody from "@mui/material/TableBody";
import TableCell from "@mui/material/TableCell";
import TableContainer from "@mui/material/TableContainer";
import TableHead from "@mui/material/TableHead";
import TableRow from "@mui/material/TableRow";
import GitHubIcon from "@mui/icons-material/GitHub";

import Link from "next/link";
import Stack from "@mui/material/Stack";

import strftime from "strftime";
import TimeAgo from "javascript-time-ago";
import en from "javascript-time-ago/locale/en.json";

import { Typography } from "@mui/material";

type repoStatus = {
  isLoaded: boolean;
  repositories?: model.repository[];
  err?: any;
};

export default function Index() {
  const [status, setStatus] = React.useState<repoStatus>({
    isLoaded: false,
  });

  const updateRepositories = () => {
    fetch(`/api/v1/repository`)
      .then((res) => res.json())
      .then(
        (result) => {
          console.log("result:", { result });
          if (result.error) {
            setStatus({ isLoaded: true, err: result.error });
          } else {
            setStatus({
              isLoaded: true,
              repositories: result.data,
            });
          }
        },
        (error) => {
          console.log("got error:", { error });
          setStatus({
            isLoaded: true,
            err: error.message,
          });
        }
      );
  };

  React.useEffect(() => {
    TimeAgo.addDefaultLocale(en);
    updateRepositories();
  }, []);

  return (
    <app.Main>
      <Grid container justifyContent="center">
        <Grid item style={{ marginTop: 100 }}>
          <Typography variant="h1">Octovy</Typography>
        </Grid>
      </Grid>
      <Grid
        container
        justifyContent="center"
        style={{ marginTop: 30, width: "100%" }}>
        <Body status={status} />
      </Grid>
    </app.Main>
  );
}

type bodyProps = {
  status: repoStatus;
};

function Body(props: bodyProps) {
  console.log("status=>", props.status);
  if (!props.status.isLoaded) {
    return <Alert severity="info">Loading...</Alert>;
  } else if (props.status.err) {
    return <Alert severity="error">Error: {props.status.err}</Alert>;
  }

  return (
    <TableContainer component={Paper}>
      <Table
        sx={{ minWidth: 650 }}
        style={{ width: "100%" }}
        size="small"
        aria-label="a dense table">
        <TableHead>
          <TableRow style={{ background: "#eee" }}>
            <TableCell>Repository</TableCell>
            <TableCell>Recent scan of default branch</TableCell>
          </TableRow>
        </TableHead>
        <TableBody>{props.status.repositories.map(Repository)}</TableBody>
      </Table>
    </TableContainer>
  );
}

function Repository(repo: model.repository) {
  return (
    <TableRow key={repo.owner + "/" + repo.name}>
      <TableCell>
        <Stack direction="row" spacing={2}>
          <Link href={repo.url}>
            <a style={{ color: "inherit" }}>
              <GitHubIcon />
            </a>
          </Link>
          <Link href={`/repository/${repo.owner}/${repo.name}`}>
            {repo.owner + "/" + repo.name}
          </Link>
        </Stack>
      </TableCell>
      <TableCell>
        <Scan repo={repo} scan={repo.edges.latest} />
      </TableCell>
    </TableRow>
  );
}

function Scan(props: { repo: model.repository; scan: model.scan }) {
  if (!props.scan) {
    return <Typography style={{ fontSize: 14 }}>No scan</Typography>;
  }

  const timeAgo = new TimeAgo("en-US");
  const scan = props.scan;
  const ts = new Date(scan.scanned_at * 1000);
  const vulnPkgs: model.packageRecord[] = scan.edges.packages
    ? scan.edges.packages.filter((pkg) => {
        return pkg.vuln_ids;
      })
    : [];
  return (
    <Typography style={{ fontSize: 14 }}>
      <Link href={"/scan/" + scan.id}>
        {vulnPkgs.length > 0
          ? `⚠️ ${vulnPkgs.length} vulnerable packages`
          : "✅ No vulnerabilities"}
      </Link>
      {" / "}
      <Link href={props.repo.url + "/tree/" + scan.commit_id}>
        {scan.commit_id.substr(0, 7)}
      </Link>
      <span> {timeAgo.format(ts)}</span>
    </Typography>
  );
}
