import React from "react";
import { Grid, Alert, Typography } from "@mui/material";
import Stack from "@mui/material/Stack";

import List from "@mui/material/List";
import ListItem from "@mui/material/ListItem";
import ListItemButton from "@mui/material/ListItemButton";
import ListItemIcon from "@mui/material/ListItemIcon";
import ListItemText from "@mui/material/ListItemText";
import Checkbox from "@mui/material/Checkbox";

import * as model from "@/components/model";
import * as ui from "@/components/ui";

export default RepoLabels;

function RepoLabels(props: { repo: model.repository }) {
  type labelsStatus = {
    isLoaded: boolean;
    labels?: model.repoLabel[];
    err?: any;
  };

  const [status, setStatus] = React.useState<labelsStatus>({
    isLoaded: false,
  });

  const get = () => {
    fetch(`/api/v1/repo-label`)
      .then((res) => res.json())
      .then(
        (result) => {
          console.log("labels:", { result });
          if (result.error) {
            setStatus({ isLoaded: true, err: result.error });
          } else {
            setStatus({
              isLoaded: true,
              labels: result.data,
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
  }

  return (
    <List sx={{ width: "100%", maxWidth: 560, bgcolor: "background.paper" }}>
      {status.labels.map((label) => {
        return <RepoLabel repo={props.repo} label={label} key={label.name} />;
      })}
    </List>
  );
}

function RepoLabel(props: { repo: model.repository; label: model.repoLabel }) {
  const initChecked = props.repo.edges.labels
    ? props.repo.edges.labels.filter((label) => {
        return label.id === props.label.id;
      }).length > 0
    : false;

  const [status, setStatus] = React.useState<{
    done: boolean;
    err?: any;
  }>({ done: false });
  const [checked, setChecked] = React.useState<boolean>(initChecked);

  const update = (value: boolean) => {
    const method = value ? "POST" : "DELETE";
    const url = `/api/v1/repo-label/${props.label.id}/assign/${props.repo.id}`;

    fetch(url, { method })
      .then((res) => res.json())
      .then(
        (result) => {
          console.log({ result });
          setChecked(value);
          setStatus({ done: true });
        },
        (error) => {
          console.log("error:", { error });
          setStatus({ done: true, err: error });
        }
      );
  };

  return (
    <ListItem
      secondaryAction={
        status.done ? (
          status.err ? (
            <Alert security="error">{status.err}</Alert>
          ) : (
            <Alert security="success">Updated</Alert>
          )
        ) : (
          <></>
        )
      }
      disablePadding>
      <ListItemButton
        onClick={() => {
          setStatus({ done: false });
          update(!checked);
        }}
        dense>
        <ListItemIcon>
          <Checkbox
            edge="start"
            checked={checked}
            tabIndex={-1}
            disableRipple
          />
        </ListItemIcon>
        <ListItemText
          primary={
            <Stack direction="row" spacing={2}>
              <ui.RepoLabel label={props.label} />
              <Typography paddingTop={1}>{props.label.description}</Typography>
            </Stack>
          }
        />
      </ListItemButton>
    </ListItem>
  );
}
