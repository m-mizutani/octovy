import * as React from "react";
import * as ReactDOM from "react-dom";
import * as octovy from "../components/octovy";

import {
  HashRouter as Router,
  Route,
  Switch,
  Redirect,
} from "react-router-dom";

import * as repositoryList from "../components/contents/RepositoryList";
import * as repository from "../components/contents/Repository";
import * as packages from "../components/contents/Packages";
import * as vulnerability from "../components/contents/Vulnerability";

function App() {
  return (
    <octovy.Frame hasNavigator={true}>
      <Switch>
        <Route path="/repository/:owner/:repoName/:branch">
          <repository.Content enablePackageLink={true} />
        </Route>
        <Route path="/repository/:owner/:repoName">
          <repository.Content enablePackageLink={true} />
        </Route>
        <Route path="/repository/:owner">
          <repositoryList.Content ownerList={true} />
        </Route>
        <Route path="/repository">
          <repositoryList.Content ownerList={true} />
        </Route>
        <Route path="/package">
          <packages.Content />
        </Route>
        <Route path="/vuln" exact>
          <vulnerability.Content />
        </Route>
        <Route path="/vuln/:vulnID">
          <vulnerability.Content />
        </Route>

        <Route path="/" exact>
          <Redirect to="/repository" />
        </Route>
      </Switch>
    </octovy.Frame>
  );
}

ReactDOM.render(<App />, document.querySelector("#app"));
