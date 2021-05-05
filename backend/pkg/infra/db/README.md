# DynamoDB schema

## Use case

- Find repositories
- Find package versions by a repository
- Find package versions of each repository by package type & name

## Index schema

### Repository

- PK: `list:repository`
- SK: `{Owner}/{RepoName}`

### PackageRecord

- PK: `pkg:{Owner}/{RepoName}@{Branch}`
- SK: `{Source}|{PkgType}|{PkgName}|{Version}`
- PK2: `pkg:{PkgType}|{PkgName}`
- SK2: `{Owner}/{RepoName}@{Branch}|{Version}`
- PK3: `pkg:{Owner}/{RepoName}@{Branch}`
- SK3: `{Source}|{PkgType}|{PkgName}|{Version}`

PK3 and SK3 are available while the package exists in the branch.

### ScanResult

- PK: `scan:{Owner}/{RepoName}@{Branch}`
- SK: `{ScannedAt}/{CommitID}`
- PK2: `scan:{Owner}/{RepoName}`
- SK2: `{CommitID}/{ScannedAt}`

### Vulnerability Status

- PK: `vuln:{Owner}/{RepoName}@{Branch}`
- SK: `{Source}|{PkgType}:{PkgName}@{Version}`
- PK2: `vuln:{PkgType}:{PkgName}`
- SK2: `{Owner}/{RepoName}@{Branch}|{Version}`

### Vulnerability Info

- PK: `list:vulnerability`
- SK: `{VulnID}`
- PK: `list:vulnerability`
- SK: `{DetectedTimestamp}/{VulnID}`
