import Chip from "@mui/material/Chip";
import * as model from "@/components/model";

import ReportProblemIcon from "@mui/icons-material/ReportProblem";
import AccessAlarmIcon from "@mui/icons-material/AccessAlarm";
import BuildIcon from "@mui/icons-material/Build";
import BeenhereIcon from "@mui/icons-material/Beenhere";
import Tooltip from "@mui/material/Tooltip";

import { makeStyles } from "@mui/styles";

const useStyles = makeStyles((theme) => ({
  vulnStatusIcon: {
    marginTop: 4,
    marginRight: 1,
    marginLeft: 0,
    marginBottom: 0,
  },
}));

function labelColor(hex: string) {
  var r = parseInt(hex.substr(1, 2), 16);
  var g = parseInt(hex.substr(3, 2), 16);
  var b = parseInt(hex.substr(5, 2), 16);

  return (r * 299 + g * 587 + b * 114) / 1000 < 128 ? "white" : "black";
}

export function RepoLabel(props: {
  label: model.repoLabel;
  size?: "small" | "medium";
}) {
  return (
    <Chip
      label={props.label.name}
      size={props.size}
      style={{
        marginTop: 3,
        backgroundColor: props.label.color,
        color: labelColor(props.label.color),
      }}
    />
  );
}

export function StatusIcon(props: {
  status: model.vulnStatusType;
  expiresAt?: number;
}) {
  const classes = useStyles();
  switch (props.status) {
    case "none":
      return <ReportProblemIcon className={classes.vulnStatusIcon} />;
    case "mitigated":
      return <BuildIcon className={classes.vulnStatusIcon} />;
    case "unaffected":
      return <BeenhereIcon className={classes.vulnStatusIcon} />;
    case "snoozed":
      const now = new Date();
      if (props.expiresAt) {
        const diff = props.expiresAt - now.getTime() / 1000;
        const expiresIn =
          diff > 86400
            ? Math.floor(diff / 86000) + " days left"
            : Math.floor(diff / 3600) + " hours left";

        return (
          <Tooltip title={expiresIn}>
            <AccessAlarmIcon className={classes.vulnStatusIcon} />
          </Tooltip>
        );
      } else {
        return <AccessAlarmIcon className={classes.vulnStatusIcon} />;
      }
  }
  return;
}
