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
