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
    },
    addUser: {
      marginRight: theme.spacing(1),
    },
    contentWrapper: {
      margin: "40px 16px",
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
    branchTab: {
      flexGrow: 1,
      backgroundColor: theme.palette.background.paper,
    },
    pkgList: {
      height: "600px",
      width: "100%",
    },
  })
);
export default useStyles;
