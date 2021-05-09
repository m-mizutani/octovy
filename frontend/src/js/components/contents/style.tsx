import { createStyles, makeStyles, Theme } from "@material-ui/core/styles";

const useStyles = makeStyles((theme: Theme) =>
  createStyles({
    paper: {
      maxWidth: 936,
      margin: "auto",
      overflow: "hidden",
    },
    searchBar: {
      borderBottom: "1px solid rgba(0, 0, 0, 0.12)",
    },
    searchInput: {
      fontSize: theme.typography.fontSize,
    },
    block: {
      display: "block",
      margin: "10px",
    },
    contentWrapper: {
      margin: "40px 30px",
    },

    formControl: {
      margin: theme.spacing(1),
      minWidth: 120,
    },
    selectEmpty: {
      marginTop: theme.spacing(2),
    },
    packageList: {
      margin: theme.spacing(2),
    },
    packageTable: {},
    packageTableHeader: {
      background: "#eee",
    },
    packageTableNameRow: {
      width: "50%",
    },
    packageTableVersionRow: {
      width: "30%",
    },
    packageTableVulnRow: {
      width: "20%",
    },
    packageTableVulnCell: {
      "& > *": {
        margin: theme.spacing(0.5),
      },
    },

    vulnParagraph: {
      marginBottom: theme.spacing(5),
    },
    vulnDetailItem: {
      padding: theme.spacing(2),
    },
    vulnDetailItemTitle: {
      fontWeight: "bold",
      marginBottom: theme.spacing(1),
    },

    pageDivider: {
      marginTop: theme.spacing(2),
      marginBottom: theme.spacing(1),
      borderColor: "#aaa",
    },

    pkgList: {
      height: "600px",
      width: "100%",
    },
  })
);
export default useStyles;
