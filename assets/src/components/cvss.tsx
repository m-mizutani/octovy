export class V3 {
  src: string;
  params: { [key: string]: string };
  constructor(s: string) {
    const parts = s.split("/");
    console.log("cvss:", parts);
    this.params = {};
    parts.forEach((p) => {
      const v = p.split(":");
      this.params[v[0]] = v[1];
    });
  }

  AccessVector() {
    const value = this.params["AV"];
    return value
      ? {
          N: "Network",
          A: "Adjacent",
          L: "Local",
          P: "Physical",
        }[value]
      : "N/A";
  }

  AttackComplexity() {
    const value = this.params["AC"];
    return value
      ? {
          L: "Low",
          H: "High",
        }[value]
      : "N/A";
  }

  PrivilegesRequired() {
    const value = this.params["PR"];
    return value
      ? {
          N: "None",
          L: "Low",
          H: "High",
        }[value]
      : "N/A";
  }

  UserInteraction() {
    const value = this.params["UI"];
    return value
      ? {
          N: "None",
          R: "Required",
        }[value]
      : "N/A";
  }

  Scope() {
    const value = this.params["S"];
    return value
      ? {
          U: "Unchanged",
          C: "Changed",
        }[value]
      : "N/A";
  }

  Confidentiality() {
    const value = this.params["C"];
    return value
      ? {
          N: "None",
          L: "Low",
          H: "High",
        }[value]
      : "N/A";
  }

  Integrity() {
    const value = this.params["C"];
    return value
      ? {
          N: "None",
          L: "Low",
          H: "High",
        }[value]
      : "N/A";
  }

  Availability() {
    const value = this.params["C"];
    return value
      ? {
          N: "None",
          L: "Low",
          H: "High",
        }[value]
      : "N/A";
  }

  ExploitCodeMaturity() {
    const value = this.params["E"];
    return value
      ? {
          X: "Not Defined",
          H: "High",
          F: "Functional",
          P: "Proof-of-Concept",
          U: "Unproven",
        }[value]
      : "N/A";
  }

  RemediationLevel() {
    const value = this.params["RL"];
    return value
      ? {
          X: "Not Defined",
          U: "Unavailable",
          W: "Workaround",
          T: "Temporary Fix",
          O: "Official Fix",
        }[value]
      : "N/A";
  }

  ReportConfidence() {
    const value = this.params["RC"];
    return value
      ? {
          X: "Not Defined",
          C: "Confirmed",
          R: "Reasonable",
          U: "Unknown",
        }[value]
      : "N/A";
  }
}
