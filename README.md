# Octovy

Octovy is a GitHub App that scans your repository's code for potentially vulnerable dependencies. It utilizes [trivy](https://github.com/aquasecurity/trivy) to detect software vulnerabilities. When triggered by events like `push` and `pull_request` from GitHub, Octovy scans the repository for dependency vulnerabilities and performs the following actions:

- Adds a comment to the pull request, summarizing the vulnerabilities found
- Inserts the scan results into BigQuery

![architecture](https://github.com/m-mizutani/octovy/assets/605953/4366161f-a4ff-4abb-9766-0fb4df818cb1)

Octovy adds a comment to the pull request when it detects new vulnerabilities between the head of the PR and the merge destination.

<img width="755" alt="comment example" src="https://github.com/m-mizutani/octovy/assets/605953/052a6362-c284-4857-921c-5c3c2f32065b">

## Setup

### 1. Creating a GitHub App

Start by creating a GitHub App [here](https://github.com/settings/apps). You can use any name and description you like. However, ensure you set the following configurations:

- **General**
  - **Webhook URL**: `https://<your domain>/webhook/github`
  - **Webhook secret**: A string of your choosing (e.g. `mysecret_XOIJPOIFEA`)

- **Permissions & events**
  - Repository Permissions
    - **Checks**: Set to Read & Write
    - **Contents**: Set to Read-only
    - **Metadata**: Set to Read-only
    - **Pull Requests**: Set to Read & Write
  - Subscribe to events
    - **Pull request**
    - **Push**

Once you have completed the setup, make sure to take note of the following information from the **General** section for future reference:

- **App ID** (e.g. `123456`)
- **Private Key**: Click `Generate a private key` and download the key file (e.g. `your-app-name.2023-08-14.private-key.pem`)

### 2. Setting Up Cloud Resources

- **Cloud Storage**: Create a Cloud Storage bucket dedicated to storing the scan results exclusively for Octovy's use.
- **BigQuery** (Optional): Create a BigQuery dataset and table for storing the scan results. Octovy will automatically update the schema. The default table name should be `scans`.

### 3. Deploying Octovy

The recommended method of deploying Octovy is via a container image, available at `ghcr.io/m-mizutani/octovy`. This image is built using GitHub Actions and published to the GitHub Container Registry.

To run Octovy, set the following environment variables:

#### Required Environment Variables
- `OCTOVY_ADDR`: The address to bind the server to (e.g. `:8080`)
- `OCTOVY_GITHUB_APP_ID`: The GitHub App ID
- `OCTOVY_GITHUB_APP_PRIVATE_KEY`: The path to the private key file
- `OCTOVY_GITHUB_APP_SECRET`: The secret string used to verify the webhook request from GitHub
- `OCTOVY_CLOUD_STORAGE_BUCKET`: The name of the Cloud Storage bucket

#### Optional Environment Variables
- `OCTOVY_TRIVY_PATH`: The path to the trivy binary. If you uses the our container image, you don't need to set this variable.
- `OCTOVY_CLOUD_STORAGE_PREFIX`: The prefix for the Cloud Storage object
- `OCTOVY_BIGQUERY_PROJECT_ID`: The name of the BigQuery dataset
- `OCTOVY_BIGQUERY_DATASET_ID`: The name of the BigQuery table
- `OCTOVY_BIGQUERY_TABLE_ID`: The name of the BigQuery table
- `OCTOVY_BIGQUERY_IMPERSONATE_SERVICE_ACCOUNT`: The service account to impersonate when accessing BigQuery
- `OCTOVY_SENTRY_DSN`: The DSN for Sentry
- `OCTOVY_SENTRY_ENV`: The environment for Sentry

## Configuration

### Ignore list

The developer can ignore specific vulnerabilities by adding them to the ignore list. The config file is written in CUE. See CUE definition in [pkg/domain/model/schema/ignore.cue](pkg/domain/model/schema/ignore.cue).

The config file should be placed in `.octovy` directory at the root of the repository. Octovy checks all files in the `.octovy` directory recursively and loads them. (e.g. `.octovy/ignore.cue`)

The following is an example of the ignore list configuration:

```cue
package octovy

IgnoreList: [
  {
    Target: "Gemfile.lock"
    Vulns: [
      {
        ID:        "CVE-2020-8130"
        ExpiresAt: "2024-08-01T00:00:00Z"
        Comment:   "This is not used"
      },
    ]
  },
]
```

`package` name should be `octovy`. `IgnoreList` is a list of `Ignore` struct.

- `Target` is the file path to ignore. That should be matched `Target` of trivy
- `Vulns` is a list of `IgnoreVuln` struct.
  - `ID` (required):  the vulnerability ID to ignore. (e.g. `CVE-2022-2202`)
  - `ExpiresAt` (required): The expiration date of the ignore. It should be in RFC3339 format. (e.g. `2023-08-01T00:00:00`). The date must be in 90 days and if it's over 90 days, Octovy will ignore it.
  - `Comment` (optional): The developer's comment


## License

Octovy is licensed under the Apache License 2.0. Copyright 2023 Masayoshi Mizutani <mizutani@hey.com>