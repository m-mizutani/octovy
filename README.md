# Octovy [![Go Report Card](https://goreportcard.com/badge/github.com/m-mizutani/octovy)](https://goreportcard.com/report/github.com/m-mizutani/octovy) [![Unit test](https://github.com/m-mizutani/octovy/actions/workflows/test.yml/badge.svg)](https://github.com/m-mizutani/octovy/actions/workflows/test.yml) [![Vulnerability scan](https://github.com/m-mizutani/octovy/actions/workflows/trivy.yml/badge.svg)](https://github.com/m-mizutani/octovy/actions/workflows/trivy.yml) [![Security scan](https://github.com/m-mizutani/octovy/actions/workflows/gosec.yml/badge.svg)](https://github.com/m-mizutani/octovy/actions/workflows/gosec.yml)

![SampleView](https://user-images.githubusercontent.com/605953/137612896-ce9bc9b7-9af5-4963-bd02-6372a81f0108.png)
Demo site: https://octovy.dev

## Overview

`Octovy` is a vulnerability management tool for 3rd party OSS packages based on [Trivy](https://github.com/aquasecurity/trivy).

![Comment to PR](https://user-images.githubusercontent.com/605953/137613080-ba866f19-cfa6-40b8-ab41-d7e2269356f2.png)

## Features

- **Package vulnerability detection in organization-wide**: Vulnerability detection and handling needs an organization-wide effort. As the law of the "weakest link", the weakest product, service or system determines the level of security in an organization. Octovy stores this data and presents the necessary information to security administrator.
- **Vulnerability management**:

## Architecture

![architecture](https://user-images.githubusercontent.com/605953/137614140-f5005f39-0ead-49bf-a097-fc6507697305.jpg)

Octovy runs as individual container with [Trivy](https://github.com/aquasecurity/trivy).



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

`OCTOVY_GITHUB_APP_PRIVATE_KEY`, `OCTOVY_GITHUB_SECRET`, `OCTOVY_GITHUB_WEBHOOK_SECRET` and `OCTOVY_DB_CONFIG` may contain secret values. I highly recommend to use secret variable management service (e.g. [Secret Manager](https://cloud.google.com/secret-manager) of Google Cloud and [AWS Secrets Manager](https://aws.amazon.com/jp/secrets-manager/)).

## License

The MIT License, Copyright 2021 Masayoshi Mizutani <mizutani@hey.com>