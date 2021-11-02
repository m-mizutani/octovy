import React from "react";
import Grid from "@mui/material/Grid";
import Alert from "@mui/material/Alert";

import List from "@mui/material/List";
import ListItem from "@mui/material/ListItem";
import ListItemIcon from "@mui/material/ListItemIcon";
import ListItemText from "@mui/material/ListItemText";
import IconButton from "@mui/material/IconButton";
import ErrorOutlineIcon from "@mui/icons-material/ErrorOutline";
import EditIcon from "@mui/icons-material/Edit";
import SaveAltIcon from "@mui/icons-material/SaveAlt";
import DeleteIcon from "@mui/icons-material/Delete";

import Dialog from "@mui/material/Dialog";
import DialogActions from "@mui/material/DialogActions";
import DialogContent from "@mui/material/DialogContent";
import DialogContentText from "@mui/material/DialogContentText";
import DialogTitle from "@mui/material/DialogTitle";

import Stack from "@mui/material/Stack";
import TextField from "@mui/material/TextField";
import Button from "@mui/material/Button";

import * as model from "../../model";
import { Refresh } from "@mui/icons-material";

type severityState = {
  msg?: string;
  isError?: boolean;
};

function Severity(props: {
  sev: model.severity;
  refresh: React.Dispatch<React.SetStateAction<severityState>>;
}) {
  const [label, setLabel] = React.useState<string>(props.sev.label);
  const [onEdit, setOnEdit] = React.useState<boolean>(false);
  const [showDelete, setShowDelete] = React.useState<boolean>(false);

  React.useEffect(() => {
    document.addEventListener(
      "keydown",
      (e) => {
        if (e.keyCode === 27) {
          setLabel(props.sev.label);
          setOnEdit(false);
        }
      },
      false
    );
  }, []);

  const update = () => {
    const body = JSON.stringify({ Label: label });
    fetch(`/api/v1/severity/${props.sev.id}`, { method: "PUT", body })
      .then((res) => res.json())
      .then(
        (result) => {
          setOnEdit(false);
          props.refresh({ msg: "updated" });
        },
        (error) => {
          console.log("error:", { error });
          setOnEdit(false);
          props.refresh({ msg: error, isError: true });
        }
      );
  };

  const deleteSev = () => {
    fetch(`/api/v1/severity/${props.sev.id}`, { method: "DELETE" })
      .then((res) => res.json())
      .then(
        (result) => {
          closeDialog();
          props.refresh({ msg: `Severity "${label}" deleted` });
        },
        (error) => {
          console.log("error:", { error });
          closeDialog();
          props.refresh({ msg: error, isError: true });
        }
      );
  };

  if (onEdit) {
    return (
      <ListItem
        secondaryAction={
          <Stack spacing={2} direction="row">
            <IconButton edge="end" aria-label="Save" onClick={update}>
              <SaveAltIcon />
            </IconButton>
          </Stack>
        }>
        <ListItemIcon>
          <ErrorOutlineIcon />
        </ListItemIcon>

        <TextField
          fullWidth
          value={label}
          onChange={(e) => {
            setLabel(e.target.value);
          }}
          size="small"></TextField>
      </ListItem>
    );
  }

  const closeDialog = () => {
    setShowDelete(false);
  };

  return (
    <>
      <Dialog
        open={showDelete}
        onClose={closeDialog}
        aria-labelledby="alert-dialog-title"
        aria-describedby="alert-dialog-description">
        <DialogTitle id="alert-dialog-title">
          {`Are you sure to delete "${label}" ?`}
        </DialogTitle>
        <DialogContent>
          <DialogContentText id="alert-dialog-description">
            {`All vulnerability statuses of "${label}" will be deleted.`}
          </DialogContentText>
        </DialogContent>
        <DialogActions>
          <Button onClick={closeDialog} autoFocus>
            Cancel
          </Button>
          <Button onClick={deleteSev}>Delete</Button>
        </DialogActions>
      </Dialog>

      <ListItem
        secondaryAction={
          <Stack spacing={2} direction="row">
            <IconButton
              edge="end"
              onClick={() => {
                setOnEdit(true);
              }}>
              <EditIcon />
            </IconButton>
            <IconButton
              edge="end"
              onClick={() => {
                setShowDelete(true);
              }}>
              <DeleteIcon />
            </IconButton>
          </Stack>
        }>
        <ListItemIcon>
          <ErrorOutlineIcon />
        </ListItemIcon>
        <ListItemText primary={label} />
      </ListItem>
    </>
  );
}

export default function Severities() {
  type status = {
    isLoaded: boolean;
    resp?: model.severity[];
    err?: any;
  };

  const [status, setStatus] = React.useState<status>({
    isLoaded: false,
  });
  const [sevState, setSevState] = React.useState<severityState>({});
  const [label, setLabel] = React.useState<string>("");

  const getList = () => {
    fetch(`/api/v1/severity`)
      .then((res) => res.json())
      .then(
        (result) => {
          console.log("result:", { result });
          if (result.error) {
            setStatus({ isLoaded: true, err: result.error });
          } else {
            setStatus({
              isLoaded: true,
              resp: result.data,
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

  const create = () => {
    const body = JSON.stringify({ Label: label });
    fetch(`/api/v1/severity`, { method: "POST", body })
      .then((res) => res.json())
      .then(
        (result) => {
          console.log("result:", { result });
          if (result.error) {
            setStatus({ isLoaded: true, err: result.error });
          } else {
            setLabel("");
            getList();
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

  React.useEffect(getList, [sevState]);

  if (!status.isLoaded) {
    return <Alert severity="info">Loading...</Alert>;
  } else if (status.err) {
    return <Alert severity="error">{status.err}</Alert>;
  }

  const renderSevStateMsg = () => {
    if (!sevState.msg) {
      return "";
    }

    return (
      <Alert severity={sevState.isError ? "error" : "info"}>
        {sevState.msg}
      </Alert>
    );
  };

  return (
    <>
      {renderSevStateMsg()}
      <Grid>
        <List dense={false} style={{ width: 380 }}>
          {status.resp.map((sev, idx) => {
            return (
              <Severity
                key={`sev-${sev}-${idx}`}
                sev={sev}
                refresh={setSevState}
              />
            );
          })}
        </List>
      </Grid>
      <Grid style={{ marginTop: 10 }}>
        <Stack spacing={2} direction="row">
          <TextField
            id="severity-name"
            label="Severity name"
            variant="standard"
            size="small"
            value={label}
            onChange={(e) => {
              setLabel(e.target.value);
            }}
          />
          <Button
            variant="contained"
            style={{ marginTop: 10, height: 32 }}
            onClick={create}>
            Add
          </Button>
        </Stack>
      </Grid>
    </>
  );
}
