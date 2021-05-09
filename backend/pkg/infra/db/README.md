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

PK3 and SK3 are available while the package exists in the branch.

### ScanResult

- PK: `scan:{Owner}/{RepoName}@{Branch}`
- SK: `{ScannedAt}/{CommitID}`
- PK2: `scan:{Owner}/{RepoName}`
- SK2: `{CommitID}/{ScannedAt}`

### Vulnerability Package Map

- PK: `vulnpkg:{VulnID}`
- SK: `{Source}|{PkgType}:{PkgName}@{Version}`

### Vulnerability Info

- PK: `list:vulnerability`
- SK: `{VulnID}`
- PK2: `list:vulnerability`
- SK2: `{DetectedTimestamp}/{VulnID}`

