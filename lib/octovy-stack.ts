import * as cdk from "@aws-cdk/core";
import * as lambda from "@aws-cdk/aws-lambda";
import * as iam from "@aws-cdk/aws-iam";
import * as dynamodb from "@aws-cdk/aws-dynamodb";
import * as sqs from "@aws-cdk/aws-sqs";
import * as ec2 from "@aws-cdk/aws-ec2";
import * as s3 from "@aws-cdk/aws-s3";
import * as secretsmanager from "@aws-cdk/aws-secretsmanager";
import * as events from "@aws-cdk/aws-events";
import * as targets from "@aws-cdk/aws-events-targets";

import * as apigateway from "@aws-cdk/aws-apigateway";
import * as acm from "@aws-cdk/aws-certificatemanager";
import * as route53 from "@aws-cdk/aws-route53";
import * as alias from "@aws-cdk/aws-route53-targets";

import { SqsEventSource } from "@aws-cdk/aws-lambda-event-sources";
import * as path from "path";

export interface vpcConfig {
  vpcId: string;
  securityGroupIds: string[];
  subnetIds: string[];
}

export interface domainConfig {
  domainName: string;
  hostedZoneID: string;
  certARN: string;
}

type PullReqEvent = "opened" | "synchronize" | "ready_for_review" | "reopened";

type rules = {
  PullReqCommentTriggers?: PullReqEvent[];
  FailCheckIfVuln?: boolean;
};

type Stage = "private" | "public";

interface OctovyProps extends cdk.StackProps {
  readonly stage: Stage;

  readonly secretsARN: string;

  readonly s3Region: string;
  readonly s3Bucket: string;
  readonly s3Prefix?: string;

  readonly lambdaRoleARN?: string;
  readonly githubEndpoint?: string;
  readonly vpcConfig?: vpcConfig;
  readonly domainConfig?: domainConfig;
  readonly dynamoPITR?: boolean;

  readonly frontendURL?: string;
  readonly githubAppURL?: string;
  readonly homepageURL?: string;

  readonly rules?: rules;

  readonly webhookEndpointTypes?: apigateway.EndpointType[];
  readonly apiEndpointTypes?: apigateway.EndpointType[];

  readonly sentryDSN?: string;
  readonly sentryEnv?: string;
}

const defaultEndpointTypes = [apigateway.EndpointType.PRIVATE];

export class OctovyStack extends cdk.Stack {
  readonly metaTable: dynamodb.Table;
  readonly scanRequestQueue: sqs.Queue;
  readonly feedbackRequestQueue: sqs.Queue;

  readonly apiHandler: lambda.Function;
  readonly scanRepo: lambda.Function;
  readonly feedback: lambda.Function;
  readonly updateDB: lambda.Function;

  constructor(scope: cdk.Construct, id: string, props: OctovyProps) {
    super(scope, id, props);
    const stack = this;

    // DynamoDB
    this.metaTable = new dynamodb.Table(this, "metaTable", {
      partitionKey: { name: "pk", type: dynamodb.AttributeType.STRING },
      sortKey: { name: "sk", type: dynamodb.AttributeType.STRING },
      timeToLiveAttribute: "expires_at",
      pointInTimeRecovery: props.dynamoPITR || false,
      billingMode: dynamodb.BillingMode.PAY_PER_REQUEST,
    });
    this.metaTable.addGlobalSecondaryIndex({
      indexName: "secondary",
      partitionKey: { name: "pk2", type: dynamodb.AttributeType.STRING },
      sortKey: { name: "sk2", type: dynamodb.AttributeType.STRING },
      projectionType: dynamodb.ProjectionType.ALL,
    });

    // SQS
    this.scanRequestQueue = new sqs.Queue(this, "scanRequest", {
      visibilityTimeout: cdk.Duration.seconds(300),
    });
    this.feedbackRequestQueue = new sqs.Queue(this, "feedbackRequest", {
      visibilityTimeout: cdk.Duration.seconds(120),
    });
    // VPC
    var securityGroups: ec2.ISecurityGroup[] | undefined = undefined;
    var vpc: ec2.IVpc | undefined = undefined;
    // var vpcSubnets: ec2.SelectedSubnets | undefined = undefined;
    if (props.vpcConfig) {
      vpc = ec2.Vpc.fromVpcAttributes(this, "Vpc", {
        vpcId: props.vpcConfig.vpcId,
        availabilityZones: cdk.Fn.getAzs(),
        privateSubnetIds: props.vpcConfig.subnetIds,
      });
      securityGroups = props.vpcConfig.securityGroupIds.map((sgID) => {
        return ec2.SecurityGroup.fromSecurityGroupId(this, sgID, sgID);
      });
    }

    // Lambda function
    const rootPath = path.resolve(__dirname, "..");
    const asset = lambda.Code.fromAsset(rootPath, {
      bundling: {
        image: lambda.Runtime.GO_1_X.bundlingDockerImage,
        user: "root",
        command: ["make", "asset"],
        environment: {
          GOARCH: "amd64",
          GOOS: "linux",
          STAGE: props.stage,
        },
      },
    });

    const lambdaRole =
      props.lambdaRoleARN !== undefined
        ? iam.Role.fromRoleArn(this, "LambdaRole", props.lambdaRoleARN, {
            mutable: false,
          })
        : undefined;

    const rules = props.rules || {};
    const envVars: { [key: string]: string } = {
      TABLE_NAME: this.metaTable.tableName,
      SECRETS_ARN: props.secretsARN,
      SCAN_REQUEST_QUEUE: this.scanRequestQueue.queueUrl,
      FEEDBACK_REQUEST_QUEUE: this.feedbackRequestQueue.queueUrl,
      GITHUB_ENDPOINT: props.githubEndpoint || "",
      GITHUB_APP_URL: props.githubAppURL || "",
      HOMEPAGE_URL: props.homepageURL || "",

      RULE_PR_COMMENT_TRIGGERS: rules.PullReqCommentTriggers
        ? rules.PullReqCommentTriggers.join("|")
        : "",
      RULE_FAIL_CHECK_IF_VULN: rules.FailCheckIfVuln ? "YES" : "",

      S3_REGION: props.s3Region,
      S3_BUCKET: props.s3Bucket,
      S3_PREFIX: props.s3Prefix || "",

      SENTRY_DSN: props.sentryDSN || "",
      SENTRY_ENV: props.sentryEnv || "",
    };

    type lambdaConfig = {
      id: string;
      timeout: cdk.Duration;
      memorySize: number;
      events?: lambda.IEventSource[];
    };
    const lambdaFunctions: { [key: string]: lambda.Function } = {};

    const newLambda = (cfg: lambdaConfig): lambda.Function => {
      return new lambda.Function(this, cfg.id, {
        runtime: lambda.Runtime.GO_1_X,
        handler: "handler",
        role: lambdaRole,
        code: asset,
        timeout: cfg.timeout,
        memorySize: cfg.memorySize,
        environment: { ...envVars, ...{ LAMBDA_FUNC_ID: cfg.id } },
        events: cfg.events,

        vpc,
        securityGroups,
      });
    };

    this.apiHandler = newLambda({
      id: "apiHandler",
      timeout: cdk.Duration.seconds(30),
      memorySize: 128,
    });

    // API gateway
    /// Webhook endpoint
    const webhookGW = new apigateway.LambdaRestApi(this, "octovy-webhook", {
      handler: this.apiHandler,
      proxy: false,
      cloudWatchRole: false,
      endpointTypes: props.webhookEndpointTypes || defaultEndpointTypes,
      policy: new iam.PolicyDocument({
        statements: [
          new iam.PolicyStatement({
            actions: ["execute-api:Invoke"],
            resources: ["execute-api:/*/*"],
            effect: iam.Effect.ALLOW,
            principals: [new iam.AnyPrincipal()],
          }),
        ],
      }),
    });
    webhookGW.root
      .addResource("webhook")
      .addResource("github")
      .addMethod("POST");

    /// API endpoint
    const apiGW = new apigateway.LambdaRestApi(this, "octovy-api", {
      handler: this.apiHandler,
      proxy: false,
      cloudWatchRole: false,
      endpointTypes: props.apiEndpointTypes || defaultEndpointTypes,
      policy: new iam.PolicyDocument({
        statements: [
          new iam.PolicyStatement({
            actions: ["execute-api:Invoke"],
            resources: ["execute-api:/*/*"],
            effect: iam.Effect.ALLOW,
            principals: [new iam.AnyPrincipal()],
          }),
        ],
      }),
    });

    apiGW.root.addMethod("GET");
    apiGW.root.addResource("bundle.js").addMethod("GET");

    const apiRoot = apiGW.root.addResource("api").addResource("v1");

    // Repo
    const apiRepo = apiRoot.addResource("repo");
    if (props.stage === "private") {
      apiRepo.addMethod("GET");
    }

    const apiRepoOwner = apiRepo.addResource("{owner}");
    apiRepoOwner.addMethod("GET");

    const apiRepoOwnerName = apiRepoOwner.addResource("{repoName}");
    apiRepoOwnerName.addMethod("GET");

    const apiRepoOwnerNameBranch = apiRepoOwnerName.addResource("{branch}");
    apiRepoOwnerNameBranch.addMethod("GET");

    // Package
    if (props.stage === "private") {
      const apiPackage = apiRoot.addResource("package");
      apiPackage.addMethod("GET");
    }

    // Vulnerability
    const apiVuln = apiRoot.addResource("vuln");
    apiVuln.addResource("{vulnID}").addMethod("GET");

    // Vulnerability response
    if (props.stage === "private") {
      const apiVulnResp = apiRoot.addResource("status");
      apiVulnResp
        .addResource("{owner}")
        .addResource("{repoName}")
        .addMethod("POST");
    }

    // Scan
    const apiScan = apiRoot.addResource("scan");
    const apiScanReport = apiScan
      .addResource("report")
      .addResource("{report_id}");
    apiScanReport.addMethod("GET");

    // Metadata
    const apiMeta = apiRoot.addResource("meta");
    apiMeta.addResource("octovy").addMethod("GET");

    // Lambda without API handler
    envVars.FRONTEND_URL = props.frontendURL || apiGW.url;

    this.scanRepo = newLambda({
      id: "scanRepo",
      timeout: cdk.Duration.seconds(300),
      memorySize: 1024,
      events: [new SqsEventSource(this.scanRequestQueue)],
    });
    this.feedback = newLambda({
      id: "feedback",
      timeout: cdk.Duration.seconds(120),
      memorySize: 1024,
      events: [new SqsEventSource(this.feedbackRequestQueue)],
    });
    this.updateDB = newLambda({
      id: "updateDB",
      timeout: cdk.Duration.seconds(300),
      memorySize: 1024,
    });

    const rule = new events.Rule(this, "PeriodicUpdateDB", {
      schedule: events.Schedule.rate(cdk.Duration.hours(1)),
    });
    rule.addTarget(new targets.LambdaFunction(this.updateDB));

    // Configure lambda permission if lambdaRole is not set
    if (props.lambdaRoleARN === undefined) {
      this.metaTable.grantReadWriteData(this.apiHandler);
      this.metaTable.grantReadWriteData(this.scanRepo);
      this.metaTable.grantReadWriteData(this.feedback);

      this.scanRequestQueue.grantSendMessages(this.apiHandler);
      this.feedbackRequestQueue.grantSendMessages(this.scanRepo);

      const secret = secretsmanager.Secret.fromSecretCompleteArn(
        this,
        "secret",
        props.secretsARN
      );
      secret.grantRead(this.apiHandler);
      secret.grantRead(this.scanRepo);
      secret.grantRead(this.feedback);

      const bucket = s3.Bucket.fromBucketName(
        this,
        "data-bucket",
        props.s3Bucket
      );
      bucket.grantReadWrite(this.updateDB);
      bucket.grantRead(this.scanRepo);
    }

    // Configure original domain name
    if (props.domainConfig) {
      configureCustomDomainName(this, props.domainConfig, apiGW);
    }
  }
}

function configureCustomDomainName(
  stack: cdk.Stack,
  domainCfg: domainConfig,
  apiGW: apigateway.LambdaRestApi
) {
  const customDomain = new apigateway.DomainName(stack, "CustomDomain", {
    certificate: acm.Certificate.fromCertificateArn(
      stack,
      "Certificate",
      domainCfg.certARN
    ),
    domainName: domainCfg.domainName,
    endpointType: apigateway.EndpointType.REGIONAL,
  });

  customDomain.addBasePathMapping(apiGW);

  const hostedZone = route53.HostedZone.fromHostedZoneAttributes(
    stack,
    "customDomainHostedZone",
    {
      hostedZoneId: domainCfg.hostedZoneID,
      zoneName: domainCfg.domainName,
    }
  );

  new route53.ARecord(stack, "EndpointAlias", {
    zone: hostedZone,
    recordName: domainCfg.domainName,
    target: route53.RecordTarget.fromAlias(
      new alias.ApiGatewayDomain(customDomain)
    ),
  });
}
