# Octovy [![Go Report Card](https://goreportcard.com/badge/github.com/m-mizutani/octovy)](https://goreportcard.com/report/github.com/m-mizutani/octovy) [![Build Status](https://travis-ci.com/m-mizutani/octovy.svg?branch=master)](https://travis-ci.com/m-mizutani/octovy)

![SampleView](https://user-images.githubusercontent.com/605953/120887167-48f7eb80-c62c-11eb-877d-79f081367c81.png)
https://octovy.io

`Octovy` is a GitHub App to scan vulnerability of package system (such as RubyGems, NPM, etc.) for GitHub repository. It detects a package lock file such as `Gemfile.lock` and checks if the package includes vulnerability based on package version. After that, Octovy stores scan report to database that can be accessed via Web UI and sends a result to [GitHub Check](https://docs.github.com/en/rest/reference/checks) as CI. A conclusion of GitHub Check is only `success` (No vulnerable packages) or `neutral` (Vulnerable package found) for now.

![GitHub Check](https://user-images.githubusercontent.com/605953/120887551-82c9f180-c62e-11eb-8049-1f5e448b4dc5.png)

Basic idea of Octovy is based on [Trivy](https://github.com/aquasecurity/trivy).

## How to use

Octovy provides 2 modes: Public or Private.

1. Public mode is available at https://octovy.io and you can install GitHub App from https://github.com/apps/octovy
2. Private mode can be deployed as your own AWS CDK stack. See [Deployment](#Deployment) section for installation step.

Public mode feature is limited because of scalability and access control perspective. I recommend to deploy your own Octovy as Private mode if you want to control access to vulnerability information of your repository.

## Deployment

Octovy consists of GitHub App and AWS CDK stack. Therefore deployment steps are slightly complex.

### Prerequisite

- npm >= 7.10.0
- AWS CDK >= 1.90.0

### 1) Create GitHub App

- Move https://github.com/settings/apps and click `New GitHub App`
- Fill `GitHub App name` and `Homepage URL`
- **Disable** Webhook -> `Active` for now
- Change `Repository permissions`
  - `Checks` to `Read & Write`
  - `Contents` to `Read-only`
  - `Pull requests` to `Read & Write`
- Choose `Any account` in `Where can this GitHub App be installed?` if you wan to use the App in other's repository

Then click `Create GitHub App` and save following information.

1. Get `App ID` in `About`
2. Create a private key by `Generate a private key` button and it will be downloaded to your PC automatically.

### 2) Create AWS resources

- Create a S3 bucket
- Create a secret of Secrets Manager. The secret must have following secret values:
  - `github_app_id`: Put `App ID` of your GitHub App
  - `github_app_private_key`: Put a private key that is **encoded to base64**

### 3) Deploy to CDK

Create your CDK configuration and clone octovy code.

```
$ mkdir your-octovy-deploy
$ cd your-octovy-deploy
$ cdk init --language=typescript
$ npm i @aws-cdk/aws-apigateway@1.90.0
$ git clone https://github.com/m-mizutani/octovy.git
$ cd octovy && npm i && cd ..
```

Edit a deployment configuration in `bin` directory (e.g. `bin/your-octovy-deploy.ts`)

```ts
#!/usr/bin/env node
import "source-map-support/register";
import * as cdk from "@aws-cdk/core";
import * as apigateway from "@aws-cdk/aws-apigateway";

import { OctovyStack } from "../octovy/lib/octovy-stack";

const app = new cdk.App();
new OctovyStack(app, "your-octovy-stack", {
  stage: "public",
  secretsARN:
    "arn:aws:secretsmanager:ap-northeast-1:11111111111:secret:octovy-xxxxxx",
  s3Region: "ap-northeast-1",
  s3Bucket: "your-octovy-bucket",
  s3Prefix: "production/",
  webhookEndpointTypes: [apigateway.EndpointType.REGIONAL],
  apiEndpointTypes: [apigateway.EndpointType.REGIONAL],
});
```

Then run `cdk deploy`. After deployment, you should see API

```
 ✅  your-octovy-stack

Outputs:
your-octovy-stack.octovyapiEndpointXXXXXX = https://xxxxxxxxx.execute-api.ap-northeast-1.amazonaws.com/prod/
your-octovy-stack.octovywebhookEndpointYYYYYY = https://yyyyyyy.execute-api.ap-northeast-1.amazonaws.com/prod/
```

### 4) Configure GitHub App again

Back to GitHub app configuration page like https://github.com/settings/apps/my-octovy and do additional configurations.

- Enable Webhook (check `Active`) and set Webhook URL that is webhookEndpoint + `webhook/github` (e.g. `https://yyyyyyy.execute-api.ap-northeast-1.amazonaws.com/prod/webhook/github`). Then click `Save changes`
- Move to `Permissions & events` -> `Subscribe to events` and check following events and `Save changes`
  - `Pull request`
  - `Push`

### Deployment options

- `lambdaRoleARN`: You can use pre-configured IAM role for Lambda
- `githubEndpoint`: API endpoint for GitHub Enterprise
- `vpcConfig`: VPC information if GitHub Enterprise is in VPC network
- `domainConfig`: You can assign your own domain name to Web UI
- `dynamoPITR`: You can enable Point in Time Recovery of DynamoDB
- `frontendURL`: Web UI if you use own domain name
- `webhookEndpointTypes`: API endpoint type for GitHub App webhook
- `apiEndpointTypes`: API endpoint type for Web UI
- `sentryDSN`: DSN URL of https://sentry.io
- `sentryEnv`: Environment name of sentry

## Development

### Invoke local server

You need 2 consoles: 1) webpack dev server and 2) API server. API server requires actual DynamoDB table. You can use dynamodb-local also however please notice dynamodb-local can not be stored scan results for now. Please note that AWS credential is required to access DynamoDB if you use DynamoDB on AWS.

- webpack dev server
  - Move `./frontend/`
  - Run `npm run dev-public` or `npm run dev-private`
- API server
  - Move root of repository
  - Run `go run ./cmd/octovy/ api -r [your-aws-region] -t [dynamodb-table-name]`

After invoking webpack dev server and API server, access to http://localhost:8080

### Test

1. Run dynamodb-local such as `docker run -d -p 127.0.0.1:8000:8000 amazon/dynamodb-local`
2. Run go test `go test ./backend/...`

## Acknowledge

`Octovy` is massively inspired by [Trivy](https://github.com/aquasecurity/trivy) and has a similar mechanism with trivy to detect vulnerability. Additionally Octovy leverages [trivy-db](https://github.com/aquasecurity/trivy-db) as vulnerability/advisory database. I appreciate trivy authors for publishing great OSS.
