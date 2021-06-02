import React from "react";

import Paper from "@material-ui/core/Paper";
import Grid from "@material-ui/core/Grid";
import Checkbox from "@material-ui/core/Checkbox";
import Link from "@material-ui/core/Link";

import Table from "@material-ui/core/Table";
import TableBody from "@material-ui/core/TableBody";
import TableCell from "@material-ui/core/TableCell";
import TableContainer from "@material-ui/core/TableContainer";
import TableHead from "@material-ui/core/TableHead";
import TableRow from "@material-ui/core/TableRow";

import Typography from "@material-ui/core/Typography";
import Chip from "@material-ui/core/Chip";

import { Link as RouterLink } from "react-router-dom";

import strftime from "strftime";

import useStyles from "./Style";
import * as model from "./Model";

type reportProps = {
  reportID?: string;
  packageLink?: boolean;
};

type reportStatus = {
  isLoaded: boolean;
  err?: any;
  report?: model.scanReport;
  displayed: model.packageSource[];
};

function toCommitLink(target: model.scanTarget) {
  if (target.URL) {
    return (
      <Link href={target.URL + "/commit/" + target.CommitID}>
        {target.CommitID.substr(0, 7)}
      </Link>
    );
  } else {
    return target.CommitID.substr(0, 7);
  }
}

export function Report(props: reportProps) {
  const classes = useStyles();

  const [vulnFilter, setVulnFilter] = React.useState<boolean>(true);
  const [status, setStatus] = React.useState<reportStatus>({
    isLoaded: false,
    displayed: [],
  });

  console.log({ props });

  const updatePackages = () => {
    if (!props.reportID) {
      return;
    }

    fetch(`api/v1/scan/report/${props.reportID}`)
      .then((res) => res.json())
      .then(
        (result) => {
          console.log("report:", { result });
          setStatus({
            isLoaded: true,
            report: result.data,
            displayed: filterPackages(result.data.Sources),
          });
        },
        (error) => {
          setStatus({
            isLoaded: true,
            err: error,
            displayed: [],
          });
        }
      );
  };

  const filterPackages = (
    sources: model.packageSource[]
  ): model.packageSource[] => {
    if (sources === undefined) {
      return [];
    }

    return sources.map((src) => {
      return {
        Source: src.Source,
        Packages: src.Packages.filter((pkg) => {
          return !vulnFilter || pkg.Vulnerabilities.length > 0;
        }),
      };
    });
  };

  const updateVulnFilter = () => {
    if (!status.report) {
      return;
    }

    setStatus({
      isLoaded: status.isLoaded,
      report: status.report,
      displayed: filterPackages(status.report.Sources),
    });
  };

  const onChangeVulnFilter = (event: React.ChangeEvent<HTMLInputElement>) => {
    setVulnFilter(event.target.checked);
  };

  const renderPackageName = (pkgType: string, pkgName: string) => {
    if (props.packageLink) {
      return (
        <RouterLink to={`/package?name=${pkgName}&type=${pkgType}`}>
          {pkgName}
        </RouterLink>
      );
    } else {
      return pkgName;
    }
  };

  const packageView = () => {
    if (!props.reportID) {
      return <div></div>;
    } else if (!status.isLoaded) {
      return <div className={classes.contentWrapper}>Loading...</div>;
    } else if (Object.keys(status.displayed).length === 0) {
      return <div>No data</div>;
    } else {
      console.log({ status });

      type metadata = {
        title: string;
        data: string;
      };
      const reportMeta: metadata[] = [
        {
          title: "Repository",
          data:
            status.report.Target.Owner + "/" + status.report.Target.RepoName,
        },
        {
          title: "Scanned At",
          data: strftime("%F %T", new Date(status.report.ScannedAt * 1000)),
        },
        {
          title: "Branch",
          data: status.report.Target.Branch,
        },
        {
          title: "Commit",
          data: toCommitLink(status.report.Target),
        },
      ];

      return (
        <div>
          <Grid item className={classes.reportMetaParagraph}>
            <Grid container spacing={2}>
              {reportMeta.map((meta, idx) => {
                return (
                  <Grid item xs={2} key={"report-meta-" + idx}>
                    <Typography className={classes.typographyTitle}>
                      {meta.title}
                    </Typography>
                    <Typography className={classes.typographyText}>
                      {meta.data}
                    </Typography>
                  </Grid>
                );
              })}
            </Grid>
          </Grid>

          <Grid item>
            <Checkbox
              checked={vulnFilter}
              onChange={onChangeVulnFilter}
              inputProps={{ "aria-label": "primary checkbox" }}
            />
            Only vulnerables
          </Grid>

          {status.displayed.map((src, idx) => {
            return (
              <div key={idx}>
                <Grid item xs={12}>
                  <Grid component="h4"> {src.Source} </Grid>
                </Grid>

                <TableContainer component={Paper}>
                  <Table
                    className={classes.packageTable}
                    size="small"
                    aria-label="simple table">
                    <TableHead className={classes.packageTableHeader}>
                      <TableRow>
                        <TableCell className={classes.packageTableNameRow}>
                          Name
                        </TableCell>
                        <TableCell className={classes.packageTableVersionRow}>
                          Version
                        </TableCell>
                        <TableCell className={classes.packageTableVulnRow}>
                          Vulnerabilities
                        </TableCell>
                      </TableRow>
                    </TableHead>
                    <TableBody>
                      {src.Packages.map((pkg, idx) => (
                        <TableRow key={idx}>
                          <TableCell component="th" scope="row">
                            {renderPackageName(pkg.Type, pkg.Name)}
                          </TableCell>
                          <TableCell>{pkg.Version}</TableCell>
                          <TableCell className={classes.packageTableVulnCell}>
                            {pkg.Vulnerabilities.map((vulnID, idx) => {
                              return (
                                <Chip
                                  component={RouterLink}
                                  to={"/vuln/" + vulnID}
                                  key={idx}
                                  size="small"
                                  label={vulnID}
                                  color="secondary"
                                  clickable
                                />
                              );
                            })}
                          </TableCell>
                        </TableRow>
                      ))}
                    </TableBody>
                  </Table>
                </TableContainer>
              </div>
            );
          })}
        </div>
      );
    }
  };

  React.useEffect(updatePackages, [props.reportID]);
  React.useEffect(updateVulnFilter, [vulnFilter]);

  return <div className={classes.contentWrapper}>{packageView()}</div>;
}
