export interface GitHubRepo {
  Owner: string;
  RepoName: string;
  Branch: string;
  CommitID: string;
  URL: string;
}

interface pkg {
  Type: string;
  Name: string;
  Version: string;
  Vulnerabilities: string[];
}

export interface packageSource {
  Source: string;
  Packages: pkg[];
}

interface scanTarget {
  Owner: string;
  RepoName: string;
  Branch: string;
  CommitID: string;
  UpdatedAt: number;
  RequestedAt: number;
}

export interface scanResult {
  Target: scanTarget;
  ScannedAt: number;
  Sources: packageSource[];
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

export interface cvss {
  V2Vector: string;
  V3Vector: string;
  V2Score: number;
  V3Score: number;
}

interface vulnDetail {
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
