import React from "react";
import { Grid, Alert, Typography } from "@mui/material";
import Link from "next/link";

import PlagiarismIcon from "@mui/icons-material/Plagiarism";
import {
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
} from "@mui/material";

import * as model from "@/components/model";

import ReactTimeAgo from "react-time-ago";
import TimeAgo from "javascript-time-ago";
import en from "javascript-time-ago/locale/en.json";
TimeAgo.addDefaultLocale(en);

export default Scans;

function Scans(props: { repo: model.repository }) {
  type status = {
    isLoaded: boolean;
    scans?: model.scan[];
    err?: any;
  };
  const [status, setStatus] = React.useState<status>({
    isLoaded: false,
  });

  const get = () => {
    fetch(`/api/v1/repository/${props.repo.owner}/${props.repo.name}/scan`)
      .then((res) => res.json())
      .then(
        (result) => {
          console.log("scans:", { result });
          if (result.error) {
            setStatus({ isLoaded: true, err: result.error });
          } else {
            setStatus({
              isLoaded: true,
              scans: result.data,
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
  React.useEffect(get, []);

  if (!status.isLoaded) {
    return <Alert severity="info">Loading...</Alert>;
  } else if (status.err) {
    return <Alert severity="error">{status.err}</Alert>;
  } else if (!status.scans) {
    return (
      <Grid marginTop={3}>
        <Typography>No scan results</Typography>
      </Grid>
    );
  }

  return (
    <Grid marginTop={3}>
      <TableContainer>
        <Table sx={{ minWidth: 650 }} aria-label="simple table" size="small">
          <TableHead>
            <TableRow>
              <TableCell width={24}></TableCell>
              <TableCell width={128}>
                <Typography variant="h6">Scanned at</Typography>
              </TableCell>
              <TableCell>
                <Typography variant="h6">Target</Typography>
              </TableCell>
              <TableCell width={128}>
                <Typography variant="h6">Packages</Typography>
              </TableCell>
              <TableCell width={128}>
                <Typography variant="h6">Vulnerables</Typography>
              </TableCell>
            </TableRow>
          </TableHead>

          <TableBody>
            {status.scans.map((scan) => {
              return (
                <TableRow key={`repo-label-${scan.id}`}>
                  <TableCell>
                    <PlagiarismIcon />
                  </TableCell>
                  <TableCell>
                    <Typography>
                      <Link href={`/scan/${scan.id}`}>
                        <a>
                          <ReactTimeAgo
                            date={new Date(scan.scanned_at * 1000)}
                            locale="en-US"
                          />
                        </a>
                      </Link>
                    </Typography>
                  </TableCell>
                  <TableCell>
                    <Typography>
                      <Link href={`${props.repo.url}/tree/${scan.branch}`}>
                        <a>{scan.branch}</a>
                      </Link>{" "}
                      (
                      <Link href={`${props.repo.url}/commit/${scan.commit_id}`}>
                        <a>{scan.commit_id.substr(0, 7)}</a>
                      </Link>
                      )
                    </Typography>
                  </TableCell>
                  <TableCell>
                    <Typography>
                      {scan.edges.packages ? scan.edges.packages.length : 0}
                    </Typography>
                  </TableCell>
                  <TableCell>
                    <Typography>
                      {scan.edges.packages
                        ? scan.edges.packages.filter((p) => {
                            return p.vuln_ids !== undefined;
                          }).length
                        : 0}
                    </Typography>
                  </TableCell>
                </TableRow>
              );
            })}
          </TableBody>
        </Table>
      </TableContainer>
    </Grid>
  );
}
