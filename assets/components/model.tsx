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
    latest: scan;
    scan: scan[];
    main: scan[];
    status: vulnStatusIndex[];
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

type vulnStatusIndex = {
  edges: {
    latest?: vulnStatus;
    status: vulnStatus[];
  };
};

export type vulnStatus = {
  comment: string;
  created_at: number;
  expires_at: number;
  pkg_name: string;
  pkg_type: string;
  source: string;
  status: vulnStatusType;
  vuln_id: string;
  edges: {
    author: user;
  };
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
  edges: {
    custom_severity?: severity;
  };
};

export interface user {
  id: string;
  login: string;
  name: string;
  avatar_url: string;
  url: string;
}

export type vulnStatusAttrs = {
  comment: string;
  expires_at: number;
  created_at: number;
  status: vulnStatusType;
  author_name: string;
  author_url: string;
  author_avatar: string;
};

// TODO: Migrate vulnStatusDB to backend pkg/domain/model
export class vulnStatusDB {
  readonly vulnMap: { [key: string]: vulnStatusAttrs };
  static toKey(src: string, pkgName: string, vulnID: string): string {
    return `${src}|${pkgName}|${vulnID}`;
  }
  constructor(statusIndex: vulnStatusIndex[]) {
    this.vulnMap = {};
    (statusIndex || []).forEach((idx) => {
      if (!idx.edges.latest) {
        return;
      }

      const status = idx.edges.latest;
      const key = vulnStatusDB.toKey(
        status.source,
        status.pkg_name,
        status.vuln_id
      );
      const attrs: vulnStatusAttrs = {
        comment: status.comment,
        expires_at: status.expires_at,
        created_at: status.created_at,
        status: status.status,
        author_name: status.edges.author.login,
        author_url: status.edges.author.url,
        author_avatar: status.edges.author.avatar_url,
      };
      console.log("insert", { key }, { attrs });
      this.vulnMap[key] = attrs;
    });
  }

  get(pkg: packageRecord, vulnID: string): vulnStatusAttrs {
    const key = vulnStatusDB.toKey(pkg.source, pkg.name, vulnID);
    return (
      this.vulnMap[key] || {
        comment: "",
        expires_at: 0,
        created_at: 0,
        status: "none",
        author_name: "",
        author_url: "",
        author_avatar: "",
      }
    );
  }
}

export type severity = {
  id: number;
  label: string;
};
