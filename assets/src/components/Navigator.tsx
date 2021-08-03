import React from "react";
import clsx from "clsx";
import {
  createStyles,
  Theme,
  withStyles,
  WithStyles,
} from "@material-ui/core/styles";
import Divider from "@material-ui/core/Divider";
import Drawer, { DrawerProps } from "@material-ui/core/Drawer";
import List from "@material-ui/core/List";
import ListItem from "@material-ui/core/ListItem";
import ListItemIcon from "@material-ui/core/ListItemIcon";
import ListItemText from "@material-ui/core/ListItemText";
import AllInboxIcon from "@material-ui/icons/AllInbox";
import ReportProblemIcon from "@material-ui/icons/ReportProblem";
import SearchIcon from "@material-ui/icons/Search";
import { Omit } from "@material-ui/types";

import { Link } from "react-router-dom";

const categories = [
  {
    id: "Menu",
    children: [
      {
        id: "Repository",
        icon: <AllInboxIcon />,
        link: "/repository",
      },
      {
        id: "Package",
        icon: <SearchIcon />,
        link: "/package",
      },
      {
        id: "Vulnerability",
        icon: <ReportProblemIcon />,
        link: "/vuln",
      },
    ],
  },
];

const styles = (theme: Theme) =>
  createStyles({
    categoryHeader: {
      paddingTop: theme.spacing(2),
      paddingBottom: theme.spacing(2),
    },
    categoryHeaderPrimary: {
      color: theme.palette.common.white,
    },
    item: {
      paddingTop: 1,
      paddingBottom: 1,
      color: "rgba(255, 255, 255, 0.7)",
      "&:hover,&:focus": {
        backgroundColor: "rgba(255, 255, 255, 0.08)",
      },
    },
    itemCategory: {
      backgroundColor: "#232f3e",
      boxShadow: "0 -1px 0 #404854 inset",
      paddingTop: theme.spacing(2),
      paddingBottom: theme.spacing(2),
    },
    firebase: {
      fontSize: 24,
      color: theme.palette.common.white,
    },
    itemActiveItem: {
      color: "#4fc3f7",
    },
    itemPrimary: {
      fontSize: "inherit",
    },
    itemIcon: {
      minWidth: "auto",
      marginRight: theme.spacing(2),
    },
    divider: {
      marginTop: theme.spacing(2),
    },
  });

export interface NavigatorProps
  extends Omit<DrawerProps, "classes">,
    WithStyles<typeof styles> {}

function Navigator(props: NavigatorProps) {
  const { classes, ...other } = props;

  return (
    <Drawer variant="permanent" {...other}>
      <List disablePadding>
        <ListItem
          className={clsx(
            classes.firebase,
            classes.item,
            classes.itemCategory
          )}>
          Octovy
        </ListItem>
        {categories.map(({ id, children }) => (
          <React.Fragment key={id}>
            <ListItem className={classes.categoryHeader}>
              <ListItemText
                classes={{
                  primary: classes.categoryHeaderPrimary,
                }}>
                {id}
              </ListItemText>
            </ListItem>
            {children.map(({ id: childId, icon, link }) => (
              <ListItem
                key={childId}
                component={Link}
                to={link}
                className={clsx(classes.item)}>
                <ListItemIcon className={classes.itemIcon}>{icon}</ListItemIcon>
                <ListItemText
                  classes={{
                    primary: classes.itemPrimary,
                  }}>
                  {childId}
                </ListItemText>
              </ListItem>
            ))}
            <Divider className={classes.divider} />
          </React.Fragment>
        ))}
      </List>
    </Drawer>
  );
}

export default withStyles(styles)(Navigator);
