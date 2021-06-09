import { ModalManager } from "@material-ui/core";

export interface GitHubRepo {
  Owner: string;
  RepoName: string;
  Branch: string;
  CommitID: string;
  URL: string;
}

export interface pkg {
  Type: string;
  Name: string;
  Version: string;
  Vulnerabilities: string[];
}

export type vulnStatusType = "none" | "snoozed" | "mitigated" | "fixed";

export interface vulnStatus {
  Comment: string;
  CreatedAt: number;
  ExpiresAt: number;
  Owner: string;
  PkgName: string;
  PkgType: string;
  RepoName: string;
  Source: string;
  Status: vulnStatusType;
  VulnID: string;
}

export class vulnStatusDB {
  readonly vulnMap: { [key: string]: vulnStatus };
  static toKey(src: string, pkgName: string, vulnID: string): string {
    return `${src}|${pkgName}|${vulnID}`;
  }
  constructor(status: vulnStatus[]) {
    console.log({ status });
    this.vulnMap = {};
    status.forEach((status) => {
      const key = vulnStatusDB.toKey(
        status.Source,
        status.PkgName,
        status.VulnID
      );
      this.vulnMap[key] = status;
    });
    console.log({ map: this.vulnMap });
  }

  getStatus(src: string, pkgName: string, vulnID: string): vulnStatusType {
    const key = vulnStatusDB.toKey(src, pkgName, vulnID);
    console.log("lookup", { key });
    const status = this.vulnMap[key];
    return status ? status.Status : undefined;
  }
}

export interface packageSource {
  Source: string;
  Packages: pkg[];
}

export interface scanTarget {
  Owner: string;
  RepoName: string;
  Branch: string;
  CommitID: string;
  UpdatedAt: number;
  RequestedAt: number;
  URL: string;
}

export interface scanReport {
  ReportID: string;
  Target: scanTarget;
  ScannedAt: number;
  Sources: packageSource[];
  Vulnerabilities: { [key: string]: vulnerability };
  VulnStatues: vulnStatus[];
}

export interface packageRecord {
  key?: number;

  Detected: GitHubRepo;
  Source: string;
  Name: string;
  Type: string;
  Version: string;
  Vulnerabilities: string[];
}

export interface repository {
  Owner: string;
  RepoName: string;
  URL: string;
  DefaultBranch: string;
  Branch: branch;
}

interface branch {
  Branch: string;
  LastScannedAt: number;
  ReportSummary: scanReportSummary;
}

interface scanReportSummary {
  ReportID: string;
  PkgTypes: string[];
  PkgCount: number;
  VulnCount: number;
  VulnPkgCount: number;
}

export interface cvss {
  V2Vector: string;
  V3Vector: string;
  V2Score: number;
  V3Score: number;
}

export interface vulnDetail {
  Title: string;
  Description: string;
  Severity: string;
  CweIDs: string[];
  CVSS: { [key: string]: cvss };
  References: string[];
  PublishedDate: string;
  LastModifiedDate: string;
}

export interface vulnerability {
  VulnID: string;
  Detail: vulnDetail;
  FirstSeenAt: number;
  LastModifiedAt: number;
}
