import React from "react";

import Paper from "@material-ui/core/Paper";
import Toolbar from "@material-ui/core/Toolbar";
import Tab from "@material-ui/core/Tab";
import Tabs from "@material-ui/core/Tabs";
import AppBar from "@material-ui/core/AppBar";
import Grid from "@material-ui/core/Grid";
import { DataGrid } from "@material-ui/data-grid";

import { useParams } from "react-router-dom";
import useStyles from "./style";

interface packageRecord {
  id?: string;
  Owner: string;
  RepoName: string;
  Branch: string;
  Source: string;
  Name: string;
  Version: string;
}

const defColumns = [
  { field: "Source", headerName: "Source", width: 250 },
  { field: "Name", headerName: "Package name", width: 500 },
  { field: "Version", headerName: "Version", width: 130 },
];

export default function Repository() {
  const classes = useStyles();

  const { owner, repoName } = useParams();
  const [branches, setBranches] = React.useState<string[]>([]);
  const [target, setTarget] = React.useState<string>();
  const [packages, setPackages] = React.useState<packageRecord[]>([]);
  const [err, setErr] = React.useState("");

  React.useEffect(() => {
    fetch(`api/v1/repo/${owner}/${repoName}`)
      .then((res) => res.json())
      .then(
        (result) => {
          console.log(result);
          setBranches(result.data.Branches);
          if (result.data.DefaultBranch) {
            setTarget(result.data.DefaultBranch);
          } else if (result.data.Branches.length > 0) {
            setTarget(result.data.Branches[0]);
          }
        },
        (error) => {
          setErr(error);
        }
      );
  }, []);

  React.useEffect(() => {
    if (target === undefined) {
      return;
    }

    fetch(`api/v1/repo/${owner}/${repoName}/${target}/package`)
      .then((res) => res.json())
      .then(
        (result) => {
          console.log(result);
          result.data.forEach((pkg: packageRecord) => {
            pkg.id = pkg.Source + "|" + pkg.Name + "|" + pkg.Version;
          });
          setPackages(result.data);
        },
        (error) => {
          setErr(error);
        }
      );
  }, [target]);

  const [tab, setTab] = React.useState(0);
  const handleTabChange = (event: React.ChangeEvent<{}>, newValue: number) => {
    setTab(newValue);
  };

  return (
    <Paper className={classes.paper}>
      <AppBar position="static" color="default" elevation={1}>
        <Toolbar>
          <Grid component="h3">
            {owner}/{repoName}
          </Grid>
        </Toolbar>
      </AppBar>

      <AppBar
        position="static"
        color="default"
        elevation={1}
        className={classes.branchTab}>
        <Tabs
          value={tab}
          onChange={handleTabChange}
          aria-label="simple tabs example">
          {branches.map((branch) => {
            return <Tab label={branch} key={branch} />;
          })}
        </Tabs>
      </AppBar>

      <AppBar
        position="static"
        color="default"
        elevation={1}
        className={classes.pkgList}>
        <DataGrid
          rows={packages}
          columns={defColumns}
          pageSize={50}
          checkboxSelection
        />
      </AppBar>
    </Paper>
  );
}
