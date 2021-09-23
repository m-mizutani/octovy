import React from "react";

import Select from "@mui/material/Select";
import MenuItem from "@mui/material/MenuItem";
import TableRow from "@mui/material/TableRow";
import TableCell from "@mui/material/TableCell";
import Grid from "@mui/material/Grid";

import ReportProblemIcon from "@mui/icons-material/ReportProblem";
import AccessAlarmIcon from "@mui/icons-material/AccessAlarm";
import BuildIcon from "@mui/icons-material/Build";
import BeenhereIcon from "@mui/icons-material/Beenhere";
import Tooltip from "@mui/material/Tooltip";
import Chip from "@mui/material/Chip";

import Dialog from "@mui/material/Dialog";
import DialogTitle from "@mui/material/DialogTitle";
import DialogActions from "@mui/material/DialogActions";
import DialogContent from "@mui/material/DialogContent";
import DialogContentText from "@mui/material/DialogContentText";
import TextField from "@mui/material/TextField";
import Button from "@mui/material/Button";
import Slider from "@mui/material/Slider";
import Alert from "@mui/material/Alert";

import { makeStyles } from "@mui/styles";
import { createTheme, ThemeProvider } from "@mui/material/styles";

import theme from "./theme";

const useStyles = makeStyles((theme) => ({
  vulnStatusIcon: {
    marginTop: 4,
    marginRight: 1,
    marginLeft: 0,
    marginBottom: 0,
  },
}));

import * as model from "./model";
import { Typography } from "@mui/material";

type vulnStatusRequest = {
  Status: string;
  Source: string;
  PkgType: string;
  PkgName: string;
  VulnID: string;
  ExpiresAt: number;
  Comment: string;
};

type packageProps = {
  idx: number;
  repo: model.repository;
  pkg: model.packageRecord;
  vuln: model.vulnerability;
  vulnDB: model.vulnStatusDB;
};

export default function Package(props: packageProps) {
  const [inputDialog, setInputDialog] = React.useState<string>();
  const [statusComment, setStatusComment] = React.useState<string>();
  const [statusDuration, setStatusDuration] = React.useState<number>(0);
  const [statusError, setStatusError] = React.useState<string>();
  const [vulnStatus, setVulnStatus] = React.useState<model.vulnStatusAttrs>(
    props.vulnDB.get(props.pkg, props.vuln.id)
  );

  const onChangeStatus = (event) => {
    const newStatus = event.target.value as model.vulnStatusType;
    console.log(newStatus);

    if (newStatus === "none") {
      updateVulnStatus(newStatus);
    } else {
      if (newStatus === "snoozed") {
        setStatusDuration(7);
      }
      setInputDialog(newStatus);
    }
  };

  const clearStatusDialog = () => {
    setStatusError(undefined);
    setStatusComment(undefined);
    setStatusDuration(0);
    setInputDialog(undefined);
  };

  const submitUpdate = () => {
    console.log("comment:", statusComment);
    if (!statusComment) {
      setStatusError("Comment required");
      return;
    }
    updateVulnStatus(inputDialog as model.vulnStatusType);
  };

  const dialogMessage = {
    snoozed: "Describe a reason for pending to update version",
    mitigated: "Describe how did you do to mitigate risk",
    unaffected: "Describe why unaffected",
  };
  const renderDialog = () => (
    <Dialog
      open={inputDialog !== undefined}
      maxWidth={"sm"}
      fullWidth
      onClose={() => {
        setInputDialog(undefined);
      }}>
      <DialogTitle id="vuln-status-update-dialog-title">
        Change status to {inputDialog}
      </DialogTitle>
      {statusError ? (
        <Alert severity="error" onClose={() => setStatusError(undefined)}>
          {statusError}
        </Alert>
      ) : undefined}
      <DialogContent>
        {inputDialog === "snoozed" ? (
          <div>
            <DialogContentText>
              Snooze duration: {statusDuration} days
            </DialogContentText>
            <Slider
              defaultValue={7}
              valueLabelDisplay="auto"
              onChange={(event: any, newValue: number) => {
                setStatusDuration(newValue);
              }}
              step={1}
              marks
              min={1}
              max={30}
            />
          </div>
        ) : (
          ""
        )}

        <DialogContentText>{dialogMessage[inputDialog]}</DialogContentText>
        <TextField
          autoFocus
          margin="dense"
          id="vuln-status-comment"
          label="Comment"
          onChange={(e) => {
            setStatusComment(e.target.value as string);
          }}
          onKeyPress={(e) => {
            if (e.code === "Enter") {
              submitUpdate();
            }
          }}
          fullWidth
        />
      </DialogContent>
      <DialogActions>
        <Button
          onClick={() => {
            setInputDialog(undefined);
          }}
          color="primary">
          Cancel
        </Button>
        <Button onClick={submitUpdate} color="primary">
          Update
        </Button>
      </DialogActions>
    </Dialog>
  );

  const updateVulnStatus = (newStatus: model.vulnStatusType) => {
    const now = new Date();
    const expiresAt =
      statusDuration > 0
        ? Math.floor(now.getTime() / 1000) + statusDuration * 86400
        : 0;
    const req: vulnStatusRequest = {
      Status: newStatus,
      ExpiresAt: expiresAt,
      PkgName: props.pkg.name,
      PkgType: props.pkg.type,
      VulnID: props.vuln.id,
      Source: props.pkg.source,
      Comment: statusComment,
    };

    const setErr = (errMsg) => {
      if (inputDialog) {
        setStatusError(errMsg);
      }
    };
    fetch(`/api/v1/status/${props.repo.owner}/${props.repo.name}`, {
      method: "POST",
      body: JSON.stringify(req),
    })
      .then((res) => res.json())
      .then(
        (result) => {
          console.log("status:", { result });
          if (result.error) {
            setErr(result.error);
          } else {
            setVulnStatus({
              status: result.data.status,
              expires_at: result.data.expires_at,
              comment: result.data.comment,
            });
            clearStatusDialog();
          }
        },
        (error) => {
          console.log("error:", error);
          setErr(error);
        }
      );
  };

  const pkgStyle =
    props.idx < props.pkg.edges.vulnerabilities.length - 1
      ? { borderBottom: "none" }
      : {};

  return (
    <TableRow
      key={`${props.pkg.source}:${props.pkg.name}:${props.pkg.version}:${props.vuln.id}`}>
      {renderDialog()}
      <TableCell style={pkgStyle}>
        {props.idx === 0 ? props.pkg.name : ""}
      </TableCell>
      <TableCell style={pkgStyle}>
        {props.idx === 0 ? props.pkg.version : ""}
      </TableCell>
      <TableCell>
        <Chip
          size="small"
          label={props.vuln.id}
          color={vulnStatus.status === "none" ? "secondary" : "default"}
          clickable
        />
      </TableCell>
      <TableCell>{props.vuln.title}</TableCell>
      <TableCell>
        <Grid container>
          <Grid item>
            <StatusIcon status={vulnStatus} />
          </Grid>
          <Grid item>
            <Select
              value={vulnStatus.status}
              onChange={onChangeStatus}
              style={{
                fontSize: "12px",
                height: 28,
                marginBottom: 5,
                marginLeft: 10,
              }}>
              <MenuItem value={"none"}>To be fixed</MenuItem>
              <MenuItem value={"snoozed"}>Snoozed</MenuItem>
              <MenuItem value={"mitigated"}>Mitigated</MenuItem>
              <MenuItem value={"unaffected"}>Unaffected</MenuItem>
            </Select>
          </Grid>
        </Grid>
      </TableCell>
      <TableCell>
        <Typography>{vulnStatus.comment}</Typography>
      </TableCell>
    </TableRow>
  );
}

type StatusIconProps = {
  status: model.vulnStatusAttrs;
};

function StatusIcon(props: StatusIconProps) {
  const classes = useStyles();
  switch (props.status.status) {
    case "none":
      return <ReportProblemIcon className={classes.vulnStatusIcon} />;
    case "mitigated":
      return <BuildIcon className={classes.vulnStatusIcon} />;
    case "unaffected":
      return <BeenhereIcon className={classes.vulnStatusIcon} />;
    case "snoozed":
      const now = new Date();
      const diff = props.status.expires_at - now.getTime() / 1000;
      const expiresIn =
        diff > 86400
          ? Math.floor(diff / 86000) + " days left"
          : Math.floor(diff / 3600) + " hours left";

      return (
        <Tooltip title={expiresIn}>
          <AccessAlarmIcon className={classes.vulnStatusIcon} />
        </Tooltip>
      );
  }
  return;
}
