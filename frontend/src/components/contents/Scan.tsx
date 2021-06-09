import React, { useState } from "react";

import Paper from "@material-ui/core/Paper";
import Grid from "@material-ui/core/Grid";
import Link from "@material-ui/core/Link";
import { createStyles, makeStyles, Theme } from "@material-ui/core/styles";

import Table from "@material-ui/core/Table";
import TableBody from "@material-ui/core/TableBody";
import TableCell from "@material-ui/core/TableCell";
import TableContainer from "@material-ui/core/TableContainer";
import TableHead from "@material-ui/core/TableHead";
import TableRow from "@material-ui/core/TableRow";

import Accordion from "@material-ui/core/Accordion";
import AccordionSummary from "@material-ui/core/AccordionSummary";
import AccordionDetails from "@material-ui/core/AccordionDetails";
import ExpandMoreIcon from "@material-ui/icons/ExpandMore";

import Typography from "@material-ui/core/Typography";
import Chip from "@material-ui/core/Chip";
import Tooltip from "@material-ui/core/Tooltip";

import Avatar from "@material-ui/core/Avatar";
import Button from "@material-ui/core/Button";
import Select from "@material-ui/core/Select";
import MenuItem from "@material-ui/core/MenuItem";

import { Link as RouterLink } from "react-router-dom";

import strftime from "strftime";

import useStyles from "./Style";
import * as model from "./Model";
import { red, orange, pink } from "@material-ui/core/colors";

const scanStyles = makeStyles((theme: Theme) =>
  createStyles({
    sectionTitle: {
      marginTop: theme.spacing(4),
      marginBottom: theme.spacing(2),
      fontWeight: "bold",
      fontSize: "18px",
    },
    packageTableHeader: {
      background: "#ddd",
    },
    packageTableNameRow: {
      width: "20%",
      position: "static",
    },
    packageTableVersionRow: {
      width: "20%",
      position: "static",
    },
    packageTableVulnRow: {
      width: "60%",
    },
    packageTableVulnCell: {
      "& > *": {
        margin: theme.spacing(0.5),
      },
    },
    vulnSourceTitle: {
      marginTop: theme.spacing(3),
      marginBottom: theme.spacing(1.5),
    },
    vulnImpactCell: {
      display: "flex",
      "& > *": {
        margin: theme.spacing(0.3),
      },
    },
  })
);

type reportProps = {
  reportID?: string;
  packageLink?: boolean;
};

type reportStatus = {
  isLoaded: boolean;
  err?: any;
  report?: model.scanReport;
  displayed: model.packageSource[];
  vulnSources: model.packageSource[];
  statusDB?: model.vulnStatusDB;
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

function filterVulnerability(
  sources: model.packageSource[]
): model.packageSource[] {
  if (!sources) {
    return [];
  }

  return sources
    .map((src): model.packageSource => {
      return {
        Source: src.Source,
        Packages: src.Packages.filter((pkg) => pkg.Vulnerabilities.length > 0),
      };
    })
    .filter((src) => {
      return src.Packages.length > 0;
    });
}

export function Report(props: reportProps) {
  const classes = useStyles();
  const scanClasses = scanStyles();

  const [status, setStatus] = React.useState<reportStatus>({
    isLoaded: false,
    displayed: [],
    vulnSources: [],
  });

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
            displayed: result.data.Sources,
            vulnSources: filterVulnerability(result.data.Sources),
            statusDB: new model.vulnStatusDB(result.data.VulnStatuses),
          });
        },
        (error) => {
          setStatus({
            isLoaded: true,
            err: error,
            displayed: [],
            vulnSources: [],
          });
        }
      );
  };

  const renderVulnerabilities = () => {
    const sources = status.vulnSources;
    if (sources.length === 0) {
      return (
        <Typography className={classes.typographyText}>
          ✅ No vulnerability found
        </Typography>
      );
    }

    const renderSource = (src: model.packageSource) => {
      return src.Packages.map((pkg) => {
        return pkg.Vulnerabilities.map(
          (vulnID, idx): JSX.Element => (
            <PackageRow
              key={`pkg-row-${pkg.Name}-${vulnID}-${idx}`}
              idx={idx}
              owner={status.report.Target.Owner}
              repoName={status.report.Target.RepoName}
              vulnID={vulnID}
              pkg={pkg}
              src={src.Source}
              vuln={status.report.Vulnerabilities[vulnID].Detail}
              status={status.statusDB.getStatus(src.Source, pkg.Name, vulnID)}
            />
          )
        );
      }).reduce((p, c) => [...p, ...c]);
    };

    return (
      <div>
        {sources.map((src, idx) => (
          <Grid key={idx}>
            <Typography className={scanClasses.vulnSourceTitle}>
              <Link
                href={
                  status.report.Target.URL +
                  "/blob/" +
                  status.report.Target.CommitID +
                  "/" +
                  src.Source
                }>
                {src.Source}
              </Link>
            </Typography>
            <TableContainer component={Paper}>
              <Table size="small" aria-label="simple table">
                <TableHead className={scanClasses.packageTableHeader}>
                  <TableRow>
                    <TableCell style={{ width: "20%" }}>Package</TableCell>
                    <TableCell style={{ minWidth: "160px" }}>VulnID</TableCell>
                    <TableCell>Title</TableCell>
                    <TableCell style={{ minWidth: "120px" }}>Impact</TableCell>
                    <TableCell style={{ minWidth: "120px" }}>Status</TableCell>
                    <TableCell style={{ minWidth: "100px" }}></TableCell>
                  </TableRow>
                </TableHead>
                <TableBody>{renderSource(src)}</TableBody>
              </Table>
            </TableContainer>
          </Grid>
        ))}
      </div>
    );
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
                  <Grid
                    item
                    xs={2}
                    key={"report-meta-" + idx}
                    className={classes.reportMetaGrid}>
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
            <Typography className={scanClasses.sectionTitle}>
              Vulnerabilities
            </Typography>
          </Grid>

          <Grid item>{renderVulnerabilities()}</Grid>

          <Grid item>
            <Typography className={scanClasses.sectionTitle}>
              All Detected Packages
            </Typography>
          </Grid>

          {status.displayed.map((src, idx) => {
            return (
              <Accordion key={idx}>
                <AccordionSummary
                  expandIcon={<ExpandMoreIcon />}
                  aria-controls="panel1a-content">
                  <Typography>
                    {src.Source} ({src.Packages.length})
                  </Typography>
                </AccordionSummary>
                <AccordionDetails>
                  <TableContainer component={Paper}>
                    <Table size="small" aria-label="simple table">
                      <TableHead className={scanClasses.packageTableHeader}>
                        <TableRow>
                          <TableCell
                            className={scanClasses.packageTableNameRow}>
                            Name
                          </TableCell>
                          <TableCell
                            className={scanClasses.packageTableVersionRow}>
                            Version
                          </TableCell>
                          <TableCell
                            className={scanClasses.packageTableVulnRow}>
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
                            <TableCell
                              className={scanClasses.packageTableVulnCell}>
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
                </AccordionDetails>
              </Accordion>
            );
          })}
        </div>
      );
    }
  };

  React.useEffect(updatePackages, [props.reportID]);

  return <div className={classes.contentWrapper}>{packageView()}</div>;
}

type PackageRowProps = {
  idx: number;
  pkg: model.pkg;
  vulnID: string;
  vuln: model.vulnDetail;
  src: string;
  owner: string;
  repoName: string;
  status?: model.vulnStatus;
};

function PackageRow(props: PackageRowProps) {
  const scanClasses = scanStyles();
  const [vulnStatus, setVulnStatus] = useState<model.vulnStatus>(
    props.status || {
      RepoName: props.repoName,
      Owner: props.owner,
      Comment: "",
      CreatedAt: 0,
      ExpiresAt: 0,
      PkgName: props.pkg.Name,
      PkgType: props.pkg.Type,
      Source: props.src,
      Status: "none",
      VulnID: props.vulnID,
    }
  );

  const renderCVSS = (cvss?: { [key: string]: model.cvss }) => {
    const naMsg = "No CVSS";
    if (!cvss) {
      return naMsg;
    }
    const providers = ["nvd", "redhat"];
    const results = providers
      .map((provider) => cvss[provider])
      .filter((v) => v !== undefined);
    if (results.length === 0 || !results[0].V3Vector) {
      return naMsg;
    }

    const vectors = {};
    results[0].V3Vector.split("/").forEach((c) => {
      const v = c.split(":");
      vectors[v[0]] = v[1];
    });
    const metrics = {
      C: "Confidentiality",
      I: "Integrity",
      A: "Availability",
    };
    const styles = {
      C: { backgroundColor: red[600] },
      I: { backgroundColor: pink[300] },
      A: { backgroundColor: orange[300] },
    };
    return (
      <div className={scanClasses.vulnImpactCell}>
        {Object.keys(metrics).map((m, idx) => {
          if (vectors[m] === "L" || vectors[m] === "H") {
            return (
              <Tooltip title={`${metrics[m]} (${vectors[m]})`} key={idx}>
                <Avatar style={styles[m]}>{m}</Avatar>
              </Tooltip>
            );
          }
        })}
      </div>
    );
  };

  type vulnStatusRequest = {
    Status: string;
    Source: string;
    PkgType: string;
    PkgName: string;
    VulnID: string;
    ExpiresAt: number;
  };

  const onChangeStatus = (event: React.ChangeEvent<{ value: unknown }>) => {
    const now = new Date();
    const newStatus = event.target.value as model.vulnStatusType;
    const req: vulnStatusRequest = {
      Status: newStatus,
      ExpiresAt:
        newStatus !== "snoozed"
          ? 0
          : Math.floor(now.getTime() / 1000) + 86400 * 14,
      PkgName: props.pkg.Name,
      PkgType: props.pkg.Type,
      VulnID: props.vulnID,
      Source: props.src,
    };

    fetch(`api/v1/status/${props.owner}/${props.repoName}`, {
      method: "POST",
      body: JSON.stringify(req),
    })
      .then((res) => res.json())
      .then(
        (result) => {
          console.log("status:", { result });
          setVulnStatus(result.data);
        },
        (error) => {
          console.log("Error:", error);
        }
      );
  };

  const renderStatus = (status?: model.vulnStatus) => {
    if (vulnStatus.Status === "snoozed") {
      const now = new Date();
      const diff = vulnStatus.ExpiresAt - now.getTime() / 1000;
      if (diff > 86400) {
        return Math.floor(diff / 86000) + " days";
      } else {
        return Math.floor(diff / 3600) + " hours";
      }
    }

    return;
  };

  return (
    <TableRow key={props.idx}>
      <TableCell
        component="th"
        scope="row"
        style={
          props.idx < props.pkg.Vulnerabilities.length - 1
            ? { borderBottom: "none" }
            : {}
        }>
        {props.idx === 0 ? `${props.pkg.Name} (${props.pkg.Version})` : ""}
      </TableCell>
      <TableCell>
        {" "}
        <Chip
          component={RouterLink}
          to={"/vuln/" + props.vulnID}
          size="small"
          label={props.vulnID}
          color={vulnStatus.Status === "none" ? "secondary" : "default"}
          clickable
        />
      </TableCell>
      <TableCell>{props.vuln.Title}</TableCell>
      <TableCell>{renderCVSS(props.vuln.CVSS)}</TableCell>
      <TableCell>
        <Select
          value={vulnStatus.Status}
          onChange={onChangeStatus}
          style={{ fontSize: "12px" }}>
          <MenuItem value={"none"}>InProgress</MenuItem>
          <MenuItem value={"snoozed"}>Snoozed</MenuItem>
          <MenuItem value={"mitigated"}>Mitigated</MenuItem>
        </Select>
      </TableCell>
      <TableCell>{renderStatus()}</TableCell>
    </TableRow>
  );
}
