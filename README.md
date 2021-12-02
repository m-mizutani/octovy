# Octovy [![Go Report Card](https://goreportcard.com/badge/github.com/m-mizutani/octovy)](https://goreportcard.com/report/github.com/m-mizutani/octovy) [![Unit test](https://github.com/m-mizutani/octovy/actions/workflows/test.yml/badge.svg)](https://github.com/m-mizutani/octovy/actions/workflows/test.yml) [![Vulnerability scan](https://github.com/m-mizutani/octovy/actions/workflows/trivy.yml/badge.svg)](https://github.com/m-mizutani/octovy/actions/workflows/trivy.yml) [![Security scan](https://github.com/m-mizutani/octovy/actions/workflows/gosec.yml/badge.svg)](https://github.com/m-mizutani/octovy/actions/workflows/gosec.yml)

![SampleView](https://user-images.githubusercontent.com/605953/137612896-ce9bc9b7-9af5-4963-bd02-6372a81f0108.png)
Demo site: https://octovy.dev

## Overview

`Octovy` is a vulnerability management tool for 3rd party OSS packages based on [Trivy](https://github.com/aquasecurity/trivy). It works as GitHub App and scan source code of a repository that is installed the GitHub App by Trivy. The scan result is stored into database and developer and security administrator can see and manage vulnerability via Web console.

## Features

- **Organization-wide vulnerability detection**: Vulnerability detection and handling needs an organization-wide effort. Octovy scans all repositories that are installed GitHub App. It prepends misconfiguration of each repository. Also Octovy stores all scanned vulnerability package list and presents the necessary information to security administrator.
    - List newly detected vulnerabilities in your organization
    - List all repositories that have specified vulnerability
- **Vulnerability management**: Octovy provides Web user interface to manage vulnerability status. A user can change status and put a comment to share vulnerability handling decision with a team. Status can be selected from below:
  - `To be fixed`: Vulnerability should be fixed later
  - `Snoozed`: Waiting vulnerability fix. E.g.) a package author have not update vulnerable code.
  - `Unaffected`: The vulnerability is not used in your product.
  - `Mitigated`: Developer have changed settings to disable the vulnerability.

Also, Octovy notifies changes of vulnerability in Pull Request of GitHub. Developer can see new/fixed package vulnerabilities by own commit in a comment of the PR.

![Comment to PR](https://user-images.githubusercontent.com/605953/137613080-ba866f19-cfa6-40b8-ab41-d7e2269356f2.png)


## Architecture

![architecture](https://user-images.githubusercontent.com/605953/137614140-f5005f39-0ead-49bf-a097-fc6507697305.jpg)

## Usage

### Prerequisite

- Prepare your own domain name. (e.g. `octovy.dev`)
- PostgreSQL 13 database

### Setup GitHub App

Replace `{your-domain}` to your own domain name.

1. Create your own GitHub app at https://github.com/settings/apps/
2. Configure `General` tab
    - Set `Callback URL` to `https://{your-domain}/auth/github/callback`
    - Set `Webhook URL` to `https://{your-domain}/webhook/github`
    - (Optional) Set `Webhook secret` if you need. The secret value should be provided as environment variable `OCTOVY_GITHUB_WEBHOOK_SECRET` to octovy runtime.
    - Generate `Client secrets`
    - Generate `Private keys`
3. Configure `Permissions & events` tab
    - In `Repository permissions`
        - Change `Contents` to `Read-only`
        - Change `Pull requests` to `Read & Write`
    - In `Subscribe to events`
        - Enable `Pull request`
        - Enable `Push`

If you want to use auto generated URL (e.g. provided by API gateway of AWS or Cloud Run of Google Cloud), `Callback URL` and `Webhook URL` can be configured later.

Please note to remember to push `Save changes` button.

### Deploy container image

Octovy container image is published into both of GitHub Container Registry `ghcr.io/m-mizutani/octovy` and Google Container Registry `gcr.io/octovy/octovy`.

| Registry                  | Commit | Release | Latest |
|:--------------------------|:------:|:-------:|:------:|
| GitHub Container Registry |   x    |    x    |   x    |
| Google Container Registry |        |    x    |   x    |

- Commit: Images built by all push event on `main` branch. Tag is commit ID (e.g. `ghcr.io/m-mizutani/octovy:2e96dedacb63c7c8ddf51fccac7780822081057a`)
- Release: Image built by release. Tag is version number (e.g. `ghcr.io/m-mizutani/octovy:v0.1.0`)
- Latest: Image built by latest release. Tag is `latest`.

Run container image with following environment variables.

- General
    - `OCTOVY_FRONTEND_URL`: Set `https://{your-domain}`
    - `OCTOVY_ADDR`: Recommend to use `0.0.0.0`
    - `OCTOVY_PORT`: (Optional) Can change port number of octovy if you needed
    - `OCTOVY_LOG_LEVEL`: (Optional) Choose log level from `trace`, `debug`, `
    - `OCTOVY_LOG_FORMAT`: (Optional) Recommend to use `json` in cloud environment.
    - `GIN_MODE`: (Optional) Set `release` if you want to avoid debug log of gin-gonic.
- GitHub App
    - `OCTOVY_GITHUB_APP_ID`: Set App ID of your GitHub App
    - `OCTOVY_GITHUB_CLIENT_ID`: Set Client ID of your GitHub App
    - `OCTOVY_GITHUB_APP_PRIVATE_KEY`: Set private key value (content of key file) of your GitHub App
    - `OCTOVY_GITHUB_SECRET`: Set Client secret of your GitHub App
    - `OCTOVY_GITHUB_WEBHOOK_SECRET`: (Optional) Set webhook secret that you set
- Database
    - `OCTOVY_DB_TYPE`: Database type. Recommend to use `postgres`
    - `OCTOVY_DB_CONFIG`: DSN of your database. Example: `host=x.x.x.x port=5432 user=octovy_app dbname=octovy_db password=xxxxxx`
- Custom GitHub check rule
    - `OCTOVY_CHECK_POLICY_DATA`: Check result policy in Rego (plain text)
    - `OCTOVY_CHECK_POLICY_FILE`: Check result policy in Rego (file path)

`OCTOVY_GITHUB_APP_PRIVATE_KEY`, `OCTOVY_GITHUB_SECRET`, `OCTOVY_GITHUB_WEBHOOK_SECRET` and `OCTOVY_DB_CONFIG` may contain secret values. I highly recommend to use secret variable management service (e.g. [Secret Manager](https://cloud.google.com/secret-manager) of Google Cloud and [AWS Secrets Manager](https://aws.amazon.com/jp/secrets-manager/)).

An example of deploy script to Cloud Run is available in [tools/deploy_cloud_run.sh](tools/deploy_cloud_run.sh).

### Custom GitHub check policy

You can define custom policy for result of GitHub check run by [Rego](https://www.openpolicyagent.org/docs/latest/).

#### Example

A following example is a policy to make CI fail if the commit has a package that has `CVE-2021-0000` vulnerability.

```rego
package octovy.check

default result = "success"

result = "failure" {
    vulnID := input.sources[_].packages[_].vuln_ids[_]
    vulnID == "CVE-2021-0000"
}
```

#### Policy specification

- Package
    - `package octovy.check` is required at head line of policy
- Input
    - `model.ScanReport` of scan result
- Output:
    - `result` as string type (required): It must be either one of `conclusion` in [GitHub check parameters](https://docs.github.com/en/rest/reference/checks#update-a-check-run--parameters).
    - `msg` as string type (optional): The message will be appeared in title of check result if given.

## License

The MIT License, Copyright 2021 Masayoshi Mizutani <mizutani@hey.com>