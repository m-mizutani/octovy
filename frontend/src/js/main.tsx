import * as React from "react";
import * as ReactDOM from "react-dom";
import Paperbase from "./components/Paperbase/Paperbase";

function App() {
  return <Paperbase />;
}

ReactDOM.render(<App />, document.querySelector("#app"));
