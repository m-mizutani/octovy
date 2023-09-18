# Octovy [![Go Report Card](https://goreportcard.com/badge/github.com/m-mizutani/octovy)](https://goreportcard.com/report/github.com/m-mizutani/octovy) [![Test](https://github.com/m-mizutani/octovy/actions/workflows/test.yml/badge.svg)](https://github.com/m-mizutani/octovy/actions/workflows/test.yml) [![trivy](https://github.com/m-mizutani/octovy/actions/workflows/trivy.yml/badge.svg)](https://github.com/m-mizutani/octovy/actions/workflows/trivy.yml) [![gosec](https://github.com/m-mizutani/octovy/actions/workflows/gosec.yml/badge.svg)](https://github.com/m-mizutani/octovy/actions/workflows/gosec.yml)

Octovy is a GitHub App to detect vulnerable dependencies in your repository by [trivy](https://github.com/aquasecurity/trivy) and save the result to database.

![architecture](https://github.com/m-mizutani/octovy/assets/605953/81eeb92d-a4e9-4baf-aae0-ace6b9dc447f)

## Setup

### 1. Create GitHub App

Create GitHub App from [here](https://github.com/settings/apps). You can use any name and description, but you need to set following configurations.

- **General**
  - **Webhook URL**: `https://<your domain>/webhook/github`
  - **Webhook secret**: Any string (e.g. `mysecret_XOIJPOIFEA`)
- **Permissions & events**
  - Repository Permissions
    - **Contents**: Read-only
    - **Metadata**: Read-only
    - **Pull Requests**: Read & Write
  - Subscribe to events
    - **Pull request**
    - **Push**

Additionally, save following information from **General** section for later use.

- **App ID** (e.g. `123456`)
- **Private Key**: Click `Generate a private key` and download the key file (e.g. `your-app-name.2023-08-14.private-key.pem`)

### 2. Setup Database

Octovy requires PostgreSQL database. You can use any PostgreSQL instance, but we recommend to use Cloud based database services like [Google Cloud SQL](https://cloud.google.com/sql) and [Amazon RDS](https://aws.amazon.com/rds/).

For database migration, [sqldef] is recommended. Installation steps should be reffered to [sqldef document](https://github.com/k0kubun/sqldef). Then you can migrate database schema by following command.

```bash
# NOTICE: Be careful not to save password to shell history
$ export PGPASSWORD=[db_password]
$ psqldef -U [db_user] -p [db_port] -h [db_host] -f database/schema.sql [db_name]
```

### 3. Deploy Octovy

A recommended way to deploy Octovy is using container image. You can use `ghcr.io/m-mizutani/octovy`. The image is built by GitHub Actions and published to GitHub Container Registry.

Following environment variables are required to run Octovy.

- GitHub App
  - `OCTOVY_GITHUB_APP_ID`: App ID of your GitHub App
  - `OCTOVY_GITHUB_APP_PRIVATE_KEY`: Private key of your GitHub App
  - `OCTOVY_GITHUB_SECRET`: Webhook secret of your GitHub App
- Network
  - `OCTOVY_ADDR`: Listening address (e.g. `0.0.0.0:8080`)
- Database
  - `OCTOVY_DB_HOST`: Hostname of your PostgreSQL database
  - `OCTOVY_DB_PORT`: Port number of your PostgreSQL database
  - `OCTOVY_DB_USER`: Username of your PostgreSQL database
  - `OCTOVY_DB_PASSWORD`: Password of your PostgreSQL database
  - `OCTOVY_DB_NAME`: Database name of your PostgreSQL database
- Logging
  - `OCTOVY_LOG_LEVEL`: Log level (e.g. `debug`, `info`, `warn`, `error`)
  - `OCTOVY_LOG_FORMAT`: Log format, recommend to use `json`

## License

Apache License 2.0. Copyright 2021 Masayoshi Mizutani <mizutani@hey.com>
