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

  readonly webhookEndpointTypes?: apigateway.EndpointType[];
  readonly apiEndpointTypes?: apigateway.EndpointType[];

  readonly sentryDSN?: string;
  readonly sentryEnv?: string;
}

const defaultEndpointTypes = [apigateway.EndpointType.PRIVATE];

export class OctovyStack extends cdk.Stack {
  readonly metaTable: dynamodb.Table;
  readonly scanRequestQueue: sqs.Queue;

  constructor(scope: cdk.Construct, id: string, props: OctovyProps) {
    super(scope, id, props);

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

    const envVars: { [key: string]: string } = {
      TABLE_NAME: this.metaTable.tableName,
      SECRETS_ARN: props.secretsARN,
      SCAN_REQUEST_QUEUE: this.scanRequestQueue.queueUrl,
      GITHUB_ENDPOINT: props.githubEndpoint || "",

      S3_REGION: props.s3Region,
      S3_BUCKET: props.s3Bucket,
      S3_PREFIX: props.s3Prefix || "",

      SENTRY_DSN: props.sentryDSN || "",
      SENTRY_ENV: props.sentryEnv || "",
    };

    const apiHandler = new lambda.Function(this, "apiHandler", {
      runtime: lambda.Runtime.GO_1_X,
      handler: "handler",
      role: lambdaRole,
      code: asset,
      timeout: cdk.Duration.seconds(30),
      memorySize: 128,
      environment: { ...envVars, ...{ LAMBDA_FUNC_ID: "apiHandler" } },

      vpc,
      securityGroups,
    });

    const scanRepo = new lambda.Function(this, "scanRepo", {
      runtime: lambda.Runtime.GO_1_X,
      handler: "handler",
      role: lambdaRole,
      code: asset,
      timeout: cdk.Duration.seconds(300),
      memorySize: 1024,
      environment: { ...envVars, ...{ LAMBDA_FUNC_ID: "scanRepo" } },
      events: [new SqsEventSource(this.scanRequestQueue)],

      vpc,
      securityGroups,
    });

    const updateDB = new lambda.Function(this, "updateDB", {
      runtime: lambda.Runtime.GO_1_X,
      handler: "handler",
      role: lambdaRole,
      code: asset,
      timeout: cdk.Duration.seconds(300),
      memorySize: 1024,
      environment: { ...envVars, ...{ LAMBDA_FUNC_ID: "updateDB" } },
    });
    const rule = new events.Rule(this, "PeriodicUpdateDB", {
      schedule: events.Schedule.rate(cdk.Duration.hours(1)),
    });
    rule.addTarget(new targets.LambdaFunction(updateDB));

    // API gateway
    /// Webhook endpoint
    const webhookGW = new apigateway.LambdaRestApi(this, "octovy-webhook", {
      handler: apiHandler,
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
      handler: apiHandler,
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

    // Scan
    const apiScan = apiRoot.addResource("scan");
    const apiScanReport = apiScan
      .addResource("report")
      .addResource("{report_id}");
    apiScanReport.addMethod("GET");

    // Configure lambda permission if lambdaRole is not set
    if (props.lambdaRoleARN === undefined) {
      this.metaTable.grantFullAccess(apiHandler);
      this.metaTable.grantFullAccess(scanRepo);

      this.scanRequestQueue.grantSendMessages(apiHandler);

      const secret = secretsmanager.Secret.fromSecretCompleteArn(
        this,
        "secret",
        props.secretsARN
      );
      secret.grantRead(scanRepo);

      const bucket = s3.Bucket.fromBucketName(
        this,
        "data-bucket",
        props.s3Bucket
      );
      bucket.grantReadWrite(updateDB);
      bucket.grantRead(scanRepo);
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
