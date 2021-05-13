import { createStyles, makeStyles, Theme } from "@material-ui/core/styles";

const useStyles = makeStyles((theme: Theme) =>
  createStyles({
    paper: {
      // maxWidth: 936,
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

    typographyTitle: {
      marginTop: theme.spacing(0),
      marginBottom: theme.spacing(1),
      fontWeight: "bold",
    },

    pkgChip: {
      marginLeft: theme.spacing(1),
    },

    progressIcon: {
      marginTop: "2px",
      marginLeft: "15px",
    },

    formControl: {
      margin: theme.spacing(1),
      minWidth: 120,
    },
    selectEmpty: {
      marginTop: theme.spacing(2),
    },
    packageList: {
      marginTop: theme.spacing(2),
    },
    packageTable: {},
    packageTableHeader: {
      background: "#ddd",
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

    pkgList: {
      height: "600px",
      width: "100%",
    },
  })
);
export default useStyles;
