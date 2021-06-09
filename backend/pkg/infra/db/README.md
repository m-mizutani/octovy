# DynamoDB schema

## Use case

- Find repositories
- Find package versions by a repository
- Find package versions of each repository by package type & name

## Index schema

### Owner

- PK: `list:owner`
- SK: `{Owner}`

### Repository

- PK: `list:repository`
- SK: `{Owner}/{RepoName}`

### Branch

- PK: `branch:{Owner}/{RepoName}`
- SK: `{Branch}`
- PK2: `branch:{Owner}/{RepoName}`
- SK2: `{LastScannedAt}/{Branch}`

### PackageRecord

- PK: `pkg:{Owner}/{RepoName}@{Branch}`
- SK: `{Source}|{PkgType}|{PkgName}|{Version}`
- PK2: `pkg:{PkgType}|{PkgName}`
- SK2: `{Owner}/{RepoName}@{Branch}|{Version}`

### ScanLog

- PK: `scan_log:{Owner}/{RepoName}@{Branch}`
- SK: `{ScannedAt}/{CommitID}`
- PK2: `scan_log:{Owner}/{RepoName}`
- SK2: `{CommitID}/{ScannedAt}`

### ScanReport

- PK: `scan_report:{ReportID}`
- SK: `*`

### Vulnerability Package Map

- PK: `vuln_pkg:{VulnID}`
- SK: `{Source}|{PkgType}:{PkgName}@{Version}`

### Vulnerability Info

- PK: `list:vulnerability`
- SK: `{VulnID}`
- PK2: `list:vulnerability`
- SK2: `{DetectedTimestamp}/{VulnID}`

### VulnStatus

- PK: `vuln_status:{Owner}/{RepoName}`
- SK: `{Source}|{PkgName}|{VulnID}`

### VulnStatusLog

- PK: `vuln_status_log:{Owner}/{RepoName}`
- SK: `{Source}|{PkgName}|{VulnID}|{CreatedAt}`
