import React from "react";
import { Grid, Alert, Typography } from "@mui/material";
import Link from "next/link";
import Avatar from "@mui/material/Avatar";
import Stack from "@mui/material/Stack";

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

import * as ui from "@/components/ui";
import ReactTimeAgo from "react-time-ago";
import TimeAgo from "javascript-time-ago";
import en from "javascript-time-ago/locale/en.json";
TimeAgo.addDefaultLocale(en);

export default VulnStatuses;

function VulnStatuses(props: { repo: model.repository }) {
  const statusSet = props.repo.edges.status
    .map((idx) => {
      return idx.edges.latest;
    })
    .filter((status) => {
      return status.status !== "none";
    });
  if (statusSet.length === 0) {
    return <Typography>No status</Typography>;
  }

  return (
    <Grid marginTop={3}>
      <TableContainer>
        <Table sx={{ minWidth: 650 }} aria-label="simple table" size="small">
          <TableHead>
            <TableRow>
              <TableCell width={24}></TableCell>
              <TableCell width={150}>
                <Typography variant="h6">Vulnerability</Typography>
              </TableCell>
              <TableCell width={128}>
                <Typography variant="h6">Status</Typography>
              </TableCell>
              <TableCell width={128}>
                <Typography variant="h6">Expires at</Typography>
              </TableCell>
              <TableCell width={128}>
                <Typography variant="h6">By</Typography>
              </TableCell>
              <TableCell>
                <Typography variant="h6">Comment</Typography>
              </TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {statusSet.map((vulnStatus) => {
              return (
                <TableRow key={`repo-label-${vulnStatus.created_at}`}>
                  <TableCell>
                    <ui.StatusIcon
                      status={vulnStatus.status}
                      expiresAt={vulnStatus.expires_at}
                    />
                  </TableCell>
                  <TableCell>
                    <Typography>
                      <Link href={`/vulnerability/${vulnStatus.vuln_id}`}>
                        {vulnStatus.vuln_id}
                      </Link>
                    </Typography>
                  </TableCell>
                  <TableCell>
                    <Typography>{vulnStatus.status}</Typography>
                  </TableCell>
                  <TableCell>
                    {vulnStatus.expires_at ? (
                      <ReactTimeAgo
                        date={new Date(vulnStatus.expires_at * 1000)}
                        locale="en-US"
                      />
                    ) : (
                      <></>
                    )}
                  </TableCell>
                  <TableCell>
                    <Stack direction="row" spacing={1}>
                      <Avatar
                        alt={vulnStatus.edges.author.name}
                        src={vulnStatus.edges.author.avatar_url}
                        sx={{ width: 24, height: 24 }}
                      />
                      <Typography style={{ fontSize: 14 }}>
                        <Link href={vulnStatus.edges.author.url}>
                          {vulnStatus.edges.author.login}
                        </Link>
                      </Typography>
                    </Stack>
                  </TableCell>
                  <TableCell>
                    <Typography>{vulnStatus.comment}</Typography>
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
