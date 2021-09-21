import React, { useState } from "react";
import { useRouter } from "next/router";
import Link from "next/link";

import Grid from "@mui/material/Grid";
import Container from "@mui/material/Container";
import Typography from "@mui/material/Typography";
import Alert from "@mui/material/Alert";

import Paper from "@mui/material/Paper";
import Table from "@mui/material/Table";
import TableBody from "@mui/material/TableBody";
import TableCell from "@mui/material/TableCell";
import TableContainer from "@mui/material/TableContainer";
import TableHead from "@mui/material/TableHead";
import TableRow from "@mui/material/TableRow";

import * as model from "../../components/model";
import * as app from "../../components/app";
import { bgcolor, borderBottom } from "@mui/system";

type scanStatus = {
  isLoaded: boolean;
  data?: model.scan;
  err?: any;
};

function Scan() {
  const router = useRouter();
  const scanID = router.query.id;

  const [status, setStatus] = React.useState<scanStatus>({
    isLoaded: false,
  });

  const updatePackages = () => {
    if (!scanID) {
      return;
    }

    fetch(`/api/v1/scan/${scanID}`)
      .then((res) => res.json())
      .then(
        (result) => {
          console.log("result:", { result });
          setStatus({
            isLoaded: true,
            data: result.data,
          });
        },
        (error) => {
          console.log("error:", { error });
          setStatus({
            isLoaded: true,
            err: error,
          });
        }
      );
  };

  React.useEffect(updatePackages, [router.query.id]);

  if (!status.isLoaded) {
    return <app.Main>Loading...</app.Main>;
  } else if (status.err) {
    return (
      <app.Main>
        <Alert>{status.err}</Alert>
      </app.Main>
    );
  }

  const scan = status.data;
  const repo = scan.edges.repository[0];
  const vulnPkgMap = {};
  status.data.edges.packages.forEach((pkg) => {
    if (!vulnPkgMap[pkg.source]) {
      vulnPkgMap[pkg.source] = new Array<model.packageRecord>();
    }
    if (pkg.vuln_ids !== undefined) {
      vulnPkgMap[pkg.source].push(pkg);
    }
  });

  return (
    <app.Main>
      <Container>
        <Grid container spacing={2}>
          <Grid item>
            <Typography variant="h5">
              {repo.owner}/{repo.name}
            </Typography>
          </Grid>
          <Grid item>
            <Typography style={{ marginTop: 3 }}>
              <Link href={repo.url + "/commit/" + status.data.commit_id}>
                {status.data.commit_id.substr(0, 7)}
              </Link>
            </Typography>
          </Grid>
        </Grid>
      </Container>
      {Object.keys(vulnPkgMap).map((key) => {
        const url = `${repo.url}/blob/${scan.commit_id}/${key}`;
        return renderPackageSource(key, vulnPkgMap[key], url);
      })}
    </app.Main>
  );
}

export default Scan;

function renderPackageSource(
  source: string,
  pkgs: model.packageRecord[],
  url: string
) {
  return (
    <Container key={source} style={{ margin: "30px 0px" }}>
      <Grid style={{ marginBottom: 10 }}>
        <Typography variant="h6">
          <Link href={url}>{source}</Link>
        </Typography>
      </Grid>
      <Grid>
        {pkgs.length > 0 ? (
          renderPackageTable(pkgs)
        ) : (
          <Typography>✅ No vulnerability found</Typography>
        )}
      </Grid>
    </Container>
  );
}

function renderPackageTable(pkgs: model.packageRecord[]) {
  return (
    <TableContainer component={Paper}>
      <Table sx={{ minWidth: 650 }} size="small" aria-label="a dense table">
        <TableHead>
          <TableRow style={{ background: "#eee" }}>
            <TableCell style={{ minWidth: 200 }}>Package</TableCell>
            <TableCell>Version</TableCell>
            <TableCell>Vulnerability</TableCell>
            <TableCell>Title</TableCell>
            <TableCell>Status</TableCell>
          </TableRow>
        </TableHead>
        <TableBody>{pkgs.map(renderPackageRow)}</TableBody>
      </Table>
    </TableContainer>
  );
}

function renderPackageRow(pkg: model.packageRecord) {
  return pkg.edges.vulnerabilities.map((vuln, idx) => {
    const pkgStyle =
      idx < pkg.edges.vulnerabilities.length - 1
        ? { borderBottom: "none" }
        : {};
    return (
      <TableRow key={`${pkg.source}:${pkg.name}:${pkg.version}:${vuln.id}`}>
        <TableCell style={pkgStyle}>{idx === 0 ? pkg.name : ""}</TableCell>
        <TableCell style={pkgStyle}>{idx === 0 ? pkg.version : ""}</TableCell>
        <TableCell>{vuln.id}</TableCell>
        <TableCell>{vuln.title}</TableCell>
        <TableCell>status</TableCell>
      </TableRow>
    );
  });
}
