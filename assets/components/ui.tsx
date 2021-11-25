import Chip from "@mui/material/Chip";
import * as model from "@/components/model";

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
