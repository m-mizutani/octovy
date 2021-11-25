import React from "react";
import { Grid, Alert, Typography } from "@mui/material";

import {
  Dialog,
  DialogActions,
  DialogContent,
  DialogContentText,
  DialogTitle,
} from "@mui/material";

import {
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
} from "@mui/material";

import { TwitterPicker } from "react-color";

import Button from "@mui/material/Button";
import TextField from "@mui/material/TextField";

import * as model from "@/components/model";
import * as ui from "@/components/ui";

export default function RepoLabels() {
  type status = {
    isLoaded: boolean;
    resp?: model.repoLabel[];
    err?: any;
  };
  const [status, setStatus] = React.useState<status>({
    isLoaded: false,
  });
  const [displayDialog, setDisplayDialog] = React.useState<boolean>(false);

  const getList = () => {
    fetch(`/api/v1/repo-label`)
      .then((res) => res.json())
      .then(
        (result) => {
          console.log("labels:", { result });
          setStatus({ isLoaded: true, err: result.error, resp: result.data });
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

  const close = (refresh: boolean) => {
    setDisplayDialog(false);
    if (refresh) {
      getList();
    }
  };

  React.useEffect(getList, []);

  const NewButton = () => {
    return (
      <Button
        variant="contained"
        style={{ marginTop: 30, height: 32 }}
        onClick={() => {
          setDisplayDialog(true);
        }}>
        New
      </Button>
    );
  };

  if (!status.isLoaded) {
    return <Alert severity="info">Loading...</Alert>;
  } else if (status.err) {
    return <Alert severity="error">{status.err}</Alert>;
  } else if (!status.resp || status.resp.length === 0) {
    return (
      <>
        <Typography>No label</Typography>
        <EditDialog display={displayDialog} close={close} />
        <NewButton />
      </>
    );
  }

  return (
    <>
      <TableContainer>
        <Table sx={{ minWidth: 650 }} aria-label="simple table" size="small">
          <TableBody>
            {status.resp.map((label) => {
              return (
                <RepoLabel
                  key={`repo-label-${label.id}`}
                  label={label}
                  refresh={() => {
                    getList();
                  }}
                />
              );
            })}
          </TableBody>
        </Table>
      </TableContainer>
      <EditDialog display={displayDialog} close={close} />
      <NewButton />
    </>
  );
}

function RepoLabel(props: { label: model.repoLabel; refresh: () => void }) {
  const [displayEditDialog, setDisplayEditDialog] =
    React.useState<boolean>(false);
  const closeEditDialog = (refresh: boolean) => {
    setDisplayEditDialog(false);
    if (refresh) {
      props.refresh();
    }
  };

  const [displayDeleteDialog, setDisplayDeleteDialog] =
    React.useState<boolean>(false);
  const closeDeleteDialog = (refresh: boolean) => {
    setDisplayDeleteDialog(false);
    if (refresh) {
      props.refresh();
    }
  };

  return (
    <TableRow>
      <EditDialog
        repoLabel={props.label}
        display={displayEditDialog}
        close={closeEditDialog}
      />
      <DeleteDialog
        repoLabel={props.label}
        display={displayDeleteDialog}
        close={closeDeleteDialog}
      />

      <TableCell width={48}>
        <ui.RepoLabel label={props.label} />
      </TableCell>
      <TableCell align="left">{props.label.description}</TableCell>
      <TableCell width={48}>
        <Button
          variant="outlined"
          color="info"
          size="small"
          onClick={() => {
            setDisplayEditDialog(true);
          }}>
          Edit
        </Button>
      </TableCell>
      <TableCell width={48}>
        <Button
          variant="outlined"
          color="error"
          size="small"
          onClick={() => {
            setDisplayDeleteDialog(true);
          }}>
          Delete
        </Button>
      </TableCell>
    </TableRow>
  );
}

function EditDialog(props: {
  repoLabel?: model.repoLabel;
  display: boolean;
  close: (refresh: boolean) => void;
}) {
  type status = { err?: any };
  const [status, setStatus] = React.useState<status>({});

  const [name, setName] = React.useState<string>(
    props.repoLabel ? props.repoLabel.name : ""
  );
  const [description, setDescription] = React.useState<string>(
    props.repoLabel ? props.repoLabel.description : ""
  );
  const [color, setColor] = React.useState({
    hex: props.repoLabel ? props.repoLabel.color : "",
  });

  const sendRequest = () => {
    const body = JSON.stringify({
      name: name,
      description: description,
      color: color.hex,
    } as model.repoLabel);
    const baseURL = `/api/v1/repo-label`;

    const method = props.repoLabel ? "PUT" : "POST";
    const url = props.repoLabel ? `${baseURL}/${props.repoLabel.id}` : baseURL;
    console.log({ body });
    fetch(url, { method, body })
      .then((res) => res.json())
      .then(
        (result) => {
          console.log({ result });
          setName("");
          setDescription("");
          props.close(true);
        },
        (error) => {
          console.log("error:", { error });
          setStatus({ err: error });
        }
      );
  };

  return (
    <Dialog
      open={props.display}
      onClose={() => {
        props.close(false);
      }}
      aria-labelledby="alert-dialog-title"
      aria-describedby="alert-dialog-description">
      {status.err ? <Alert severity="error">{status.err}</Alert> : <></>}
      <DialogTitle id="alert-dialog-title">Edit repository label</DialogTitle>
      <DialogContent>
        <TextField
          id="label-name"
          label="Label name"
          variant="standard"
          size="small"
          value={name}
          onChange={(e) => {
            setName(e.target.value);
          }}
        />
      </DialogContent>
      <DialogContent>
        <TextField
          id="label-description"
          label="Description"
          variant="standard"
          size="small"
          value={description}
          onChange={(e) => {
            setDescription(e.target.value);
          }}
        />
      </DialogContent>
      <DialogContent>
        <TwitterPicker
          color={color}
          onChangeComplete={(c) => {
            console.log({ c });
            setColor(c);
          }}
        />
      </DialogContent>
      <DialogActions>
        <Button
          onClick={() => {
            props.close(false);
          }}
          autoFocus>
          Cancel
        </Button>
        <Button variant="contained" onClick={sendRequest}>
          {props.repoLabel ? "Update" : "Create"}
        </Button>
      </DialogActions>
    </Dialog>
  );
}

function DeleteDialog(props: {
  repoLabel: model.repoLabel;
  display: boolean;
  close: (refresh: boolean) => void;
}) {
  type status = { err?: any };
  const [status, setStatus] = React.useState<status>({});

  const sendRequest = () => {
    fetch(`/api/v1/repo-label/${props.repoLabel.id}`, { method: "DELETE" })
      .then((res) => res.json())
      .then(
        (result) => {
          console.log({ result });
          props.close(true);
        },
        (error) => {
          console.log("error:", { error });
          setStatus({ err: error });
        }
      );
  };

  const closeDialog = () => {
    props.close(false);
  };

  return (
    <Dialog
      open={props.display}
      onClose={() => {
        props.close(false);
      }}
      aria-labelledby="alert-dialog-title"
      aria-describedby="alert-dialog-description">
      {status.err ? <Alert severity="error">{status.err}</Alert> : <></>}
      <DialogTitle id="alert-dialog-title">
        {`Are you sure to delete "${props.repoLabel.name}" ?`}
      </DialogTitle>
      <DialogContent>
        <DialogContentText id="alert-dialog-description">
          {`All repository label of "${props.repoLabel.name}" will be deleted.`}
        </DialogContentText>
      </DialogContent>
      <DialogActions>
        <Button onClick={closeDialog} autoFocus>
          Cancel
        </Button>
        <Button variant="contained" color="error" onClick={sendRequest}>
          Delete
        </Button>
      </DialogActions>
    </Dialog>
  );
}
