export interface scan {
  id: string;
  branch: string;
  commit_id: string;
  requested_at: number;
  scanned_at: number;
  check_id: string;
  pull_request_target: string;

  edges: {
    packages: packageRecord[];
    repository: repository[];
  };
}

export type repository = {
  owner: string;
  name: string;
  url: string;
  default_branch: string;
  edges: {
    scan: scan[];
    status: vulnStatus[];
  };
};

export type packageRecord = {
  type: string;
  source: string;
  name: string;
  version: string;
  vuln_ids: string[];
  edges: {
    vulnerabilities: vulnerability[];
  };
};

export type vulnStatusType =
  | "none"
  | "snoozed"
  | "mitigated"
  | "unaffected"
  | "fixed";

export type vulnStatus = {
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
};

export type vulnerability = {
  id: string;
  first_seen_at: number;
  last_modified_at: number;
  title: string;
  description: string;
  cwe_id: string[];
  severity: string;
  cvss: string[];
  references: string[];
};

export interface user {
  UserID: string;
  Login: string;
  Name: string;
  AvatarURL: string;
  URL: string;
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
  }

  getStatus(src: string, pkgName: string, vulnID: string): vulnStatus {
    const key = vulnStatusDB.toKey(src, pkgName, vulnID);
    return this.vulnMap[key];
  }
}
