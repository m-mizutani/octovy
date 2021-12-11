import React from "react";
import { useRouter } from "next/router";
import Link from "next/link";

import Grid from "@mui/material/Grid";
import Container from "@mui/material/Container";
import Typography from "@mui/material/Typography";
import Alert from "@mui/material/Alert";
import Button from "@mui/material/Button";
import Modal from "@mui/material/Modal";

import Paper from "@mui/material/Paper";
import Table from "@mui/material/Table";
import TableBody from "@mui/material/TableBody";
import TableCell from "@mui/material/TableCell";
import TableContainer from "@mui/material/TableContainer";
import TableHead from "@mui/material/TableHead";
import TableRow from "@mui/material/TableRow";

import * as model from "@/components/model";
import * as app from "@/components/app";

import Package from "./package";
import CopyClipboard from "./raw";

type scanStatus = {
  isLoaded: boolean;
  data?: model.scan;
  db?: model.vulnStatusDB;
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
          if (result.error) {
            setStatus({ isLoaded: true, err: result.error });
          } else {
            const scan: model.scan = result.data;
            setStatus({
              isLoaded: true,
              data: scan,
              db: new model.vulnStatusDB(scan.edges.repository[0].edges.status),
            });
          }
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
    return (
      <app.Main>
        <Typography variant="h5">Loading...</Typography>
      </app.Main>
    );
  } else if (status.err) {
    return (
      <app.Main>
        <Alert severity="error">{status.err}</Alert>
      </app.Main>
    );
  }

  const scan = status.data;
  const repo = scan.edges.repository[0];
  const vulnPkgMap = {};

  if (status.data.edges.packages) {
    status.data.edges.packages.forEach((pkg) => {
      if (!vulnPkgMap[pkg.source]) {
        vulnPkgMap[pkg.source] = new Array<model.packageRecord>();
      }
      if (pkg.vuln_ids !== undefined) {
        vulnPkgMap[pkg.source].push(pkg);
      }
    });
  }

  return (
    <app.Main>
      <Container>
        <Grid container spacing={2}>
          <Grid item>
            <Typography variant="h4">
              <Link href={`/repository/${repo.owner}/${repo.name}`}>
                <a style={{ color: "inherit", textDecoration: "none" }}>
                  {repo.owner}/{repo.name}
                </a>
              </Link>
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
      {vulnPkgMap ? (
        Object.keys(vulnPkgMap).map((key) => {
          const url = `${repo.url}/blob/${scan.commit_id}/${key}`;
          return renderPackageSource(
            repo,
            key,
            vulnPkgMap[key],
            url,
            status.db
          );
        })
      ) : (
        <Container style={{ margin: "30px 0px" }}>
          <Typography>✅ No vulnerable package</Typography>
        </Container>
      )}
      <Container style={{ margin: "30px 0px" }}>
        <Grid>
          <CopyClipboard scanID={router.query.id as string} />
        </Grid>
      </Container>
    </app.Main>
  );
}

export default Scan;

function renderPackageSource(
  repo: model.repository,
  source: string,
  pkgs: model.packageRecord[],
  url: string,
  db: model.vulnStatusDB
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
          <TableContainer component={Paper}>
            <Table
              sx={{ minWidth: 650 }}
              size="small"
              aria-label="a dense table">
              <TableHead>
                <TableRow style={{ background: "#eee" }}>
                  <TableCell style={{ minWidth: 160 }}>Package</TableCell>
                  <TableCell>Version</TableCell>
                  <TableCell>Vulnerability</TableCell>
                  <TableCell>Title</TableCell>
                  <TableCell style={{ minWidth: 160 }}>Status</TableCell>
                  <TableCell style={{ minWidth: 120 }}>Comment</TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {pkgs.map((pkg) => {
                  return pkg.edges.vulnerabilities.map((vuln, idx) => {
                    const k = `${pkg.source}:${pkg.name}:${pkg.version}:${vuln.id}`;
                    return (
                      <Package
                        key={k}
                        repo={repo}
                        pkg={pkg}
                        vuln={vuln}
                        idx={idx}
                        vulnDB={db}
                      />
                    );
                  });
                })}
              </TableBody>
            </Table>
          </TableContainer>
        ) : (
          <Typography>✅ No vulnerability found</Typography>
        )}
      </Grid>
    </Container>
  );
}
