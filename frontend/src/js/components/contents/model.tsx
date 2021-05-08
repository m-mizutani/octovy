interface GitHubRepo {
  Owner: string;
  RepoName: string;
  Branch: string;
  CommitID: string;
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
