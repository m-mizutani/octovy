# DynamoDB schema

## Use case

- Find repositories
- Find package versions by a repository
- Find package versions of each repository by package type & name

## Index schema

### Repository

- PK: `list:repository`
- SK: `{Org}/{RepoName}`

### Package

- PK: `pkg:{Org}/{RepoName}@{Branch}`
- SK: `{Source}|{PkgType}:{PkgName}@{Version}`
- PK2: `pkg:{PkgType}:{PkgName}`
- SK2: `{Org}/{RepoName}@{Branch}|{Version}`

### ScanResult

- PK: `scan:{Org}/{RepoName}@{Branch}`
- SK: `{CommitTimestamp}/{CommitID}`
- PK2: `scan:{Org}/{RepoName}`
- SK2: `{CommitID}`

### Vulnerability

- PK: `vuln:{Org}/{RepoName}@{Branch}`
- SK: `{Source}|{PkgType}:{PkgName}@{Version}`
- PK2: `vuln:{PkgType}:{PkgName}`
- SK2: `{Org}/{RepoName}@{Branch}|{Version}`

### Advisory

- PK: `list:advisory`
- SK: `{VulnID}`
- PK2: `list:advisory`
- SK2: `{DetectedTimestamp}`
