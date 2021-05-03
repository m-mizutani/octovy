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
- SK: `{Source}|{PkgType}|{PkgName}|{Version}`
- PK2: `pkg:{PkgType}:{PkgName}`
- SK2: `{Org}/{RepoName}@{Branch}|{Version}`
- PK3: `pkg:{Org}/{RepoName}@{Branch}`
- SK3: `{Source}|{PkgType}|{PkgName}|{Version}`

PK3 and SK3 are available while the package exists in the branch.

### ScanResult

- PK: `scan:{Org}/{RepoName}`
- SK: `{CommitTimestamp}/{CommitID}`
- PK2: `scan:{Org}/{RepoName}`
- SK2: `{CommitID}`

### Vulnerability Status

- PK: `vuln:{Org}/{RepoName}@{Branch}`
- SK: `{Source}|{PkgType}:{PkgName}@{Version}`
- PK2: `vuln:{PkgType}:{PkgName}`
- SK2: `{Org}/{RepoName}@{Branch}|{Version}`

### Vulnerability Info

- PK: `list:vulninfo`
- SK: `{VulnID}`
- PK: `list:vulninfo`
- SK: `{DetectedTimestamp}`
