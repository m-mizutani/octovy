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

interface vulnDetail {
  Title: string;
  Description: string;
  Severity: string;
  CweIDs: string[];
  //VendorSeverity   :VendorSeverity ;
  //CVSS             :VendorCVSS     ;
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
