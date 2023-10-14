# Octovy

Octovy is a GitHub application designed to identify and alert you to any dependencies in your repository that could be potentially vulnerable. It uses [trivy](https://github.com/aquasecurity/trivy) for detection and then stores the results in a database for your reference.

![architecture](https://github.com/m-mizutani/octovy/assets/605953/a58c93e1-cfbf-4ff7-9427-1fc385cf7b9c)

## Setup

### 1. Creating a GitHub App

Start by creating a GitHub App [here](https://github.com/settings/apps). You can use any name and description you like. However, ensure you set the following configurations:

- **General**
  - **Webhook URL**: `https://<your domain>/webhook/github`
  - **Webhook secret**: A string of your choosing (e.g. `mysecret_XOIJPOIFEA`)

- **Permissions & events**
  - Repository Permissions
    - **Contents**: Set to Read-only
    - **Metadata**: Set to Read-only
    - **Pull Requests**: Set to Read & Write
  - Subscribe to events
    - **Pull request**
    - **Push**

Once complete, note down the following information from the **General** section for later:

- **App ID** (e.g. `123456`)
- **Private Key**: Click `Generate a private key` and download the key file (e.g. `your-app-name.2023-08-14.private-key.pem`)

### 2. Setting Up the Database

Octovy requires a PostgreSQL database. You can use any PostgreSQL instance you like, but we recommend cloud-based database services such as [Google Cloud SQL](https://cloud.google.com/sql) or [Amazon RDS](https://aws.amazon.com/rds/).

For database migration, [sqldef](https://github.com/k0kubun/sqldef) is recommended. After installing sqldef, you can migrate your database schema using the command below. Be sure to replace the placeholders with your actual database information.

```bash
# NOTICE: Be careful not to save the password to shell history
$ export PGPASSWORD=[db_password]
$ psqldef -U [db_user] -p [db_port] -h [db_host] -f database/schema.sql [db_name]
```

### 3. Deploying Octovy

The recommended method of deploying Octovy is via a container image, available at `ghcr.io/m-mizutani/octovy`. This image is built using GitHub Actions and published to the GitHub Container Registry.

To run Octovy, set the following environment variables:

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

Octovy is licensed under the Apache License 2.0. Copyright 2023 Masayoshi Mizutani <mizutani@hey.com>