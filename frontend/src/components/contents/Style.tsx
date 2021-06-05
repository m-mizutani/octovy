import { createStyles, makeStyles, Theme } from "@material-ui/core/styles";

const useStyles = makeStyles((theme: Theme) =>
  createStyles({
    paper: {
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

    reportMetaParagraph: {
      marginBottom: theme.spacing(4),
    },
    reportMetaGrid: {
      wordBreak: "break-all",
      wordWrap: "break-word",
    },
    typographyTitle: {
      marginTop: theme.spacing(0),
      marginBottom: theme.spacing(1),
      fontSize: "18px",
      fontWeight: "bold",
    },
    typographyText: {
      fontSize: "14px",
    },

    formControl: {
      margin: theme.spacing(1),
      minWidth: 120,
    },
    packageTableHeader: {
      background: "#ddd",
    },

    progressIcon: {
      marginTop: "2px",
      marginLeft: "15px",
    },

    pkgList: {
      height: "600px",
      width: "100%",
    },

    vulnParagraph: {
      marginBottom: "30px",
    },
    vulnDetailItem: {
      padding: "6px",
      wordBreak: "break-all",
      wordWrap: "break-word",
    },
    vulnDetailItemTitle: {
      fontWeight: "bold",
      marginBottom: "3px",
    },
  })
);
export default useStyles;
