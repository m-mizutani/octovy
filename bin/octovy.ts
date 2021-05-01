#!/usr/bin/env node
import "source-map-support/register";
import * as cdk from "@aws-cdk/core";
import { OctovyStack } from "../lib/octovy-stack";

const app = new cdk.App();
new OctovyStack(app, process.env.STACK_NAME || "octovy", {
  secretsARN: process.env.SECRETS_ARN!,

  lambdaRoleARN: process.env.LAMBDA_ROLE_ARN,
  githubEndpoint: process.env.GITHUB_ENDPOINT,

  sentryDSN: process.env.SENTRY_DSN,
  sentryEnv: process.env.SENTRY_ENV,
});
