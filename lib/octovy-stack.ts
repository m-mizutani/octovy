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
import { SqsEventSource } from "@aws-cdk/aws-lambda-event-sources";
import * as path from "path";

export interface vpcConfig {
  vpcId: string;
  securityGroupIds: string[];
  subnetIds: string[];
}

interface OctovyProps extends cdk.StackProps {
  readonly secretsARN: string;

  readonly s3Region: string;
  readonly s3Bucket: string;
  readonly s3Prefix?: string;

  readonly lambdaRoleARN?: string;
  readonly githubEndpoint?: string;
  readonly vpcConfig?: vpcConfig;
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
      handler: "apiHandler",
      role: lambdaRole,
      code: asset,
      timeout: cdk.Duration.seconds(30),
      memorySize: 128,
      environment: envVars,

      vpc,
      securityGroups,
    });

    const scanRepo = new lambda.Function(this, "scanRepo", {
      runtime: lambda.Runtime.GO_1_X,
      handler: "scanRepo",
      role: lambdaRole,
      code: asset,
      timeout: cdk.Duration.seconds(300),
      memorySize: 1024,
      environment: envVars,
      events: [new SqsEventSource(this.scanRequestQueue)],

      vpc,
      securityGroups,
    });

    const updateDB = new lambda.Function(this, "updateDB", {
      runtime: lambda.Runtime.GO_1_X,
      handler: "updateDB",
      role: lambdaRole,
      code: asset,
      timeout: cdk.Duration.seconds(300),
      memorySize: 1024,
      environment: envVars,
    });
    const rule = new events.Rule(this, "PeriodicUpdateDB", {
      schedule: events.Schedule.rate(cdk.Duration.hours(1)),
    });
    rule.addTarget(new targets.LambdaFunction(updateDB));

    // API gateway
    const gw = new apigateway.LambdaRestApi(this, "api", {
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

    gw.root.addMethod("GET");
    gw.root.addResource("bundle.js").addMethod("GET");

    const apiRoot = gw.root.addResource("api").addResource("v1");
    apiRoot.addResource("webhook").addResource("github").addMethod("POST");

    // Repo
    const apiRepo = apiRoot.addResource("repo");
    apiRepo.addMethod("GET");

    const apiRepoOwner = apiRepo.addResource("{owner}");
    apiRepoOwner.addMethod("GET");

    const apiRepoOwnerName = apiRepoOwner.addResource("{repoName}");
    apiRepoOwnerName.addMethod("GET");
    apiRepoOwnerName
      .addResource("{branch}")
      .addResource("package")
      .addMethod("GET");

    // Package
    const apiPackage = apiRoot.addResource("package");
    apiPackage.addMethod("GET");

    // Vulnerability
    const apiVuln = apiRoot.addResource("vuln");
    apiVuln.addMethod("GET");

    // Scan
    const apiScan = apiRoot.addResource("scan");
    const apiScanRef = apiScan
      .addResource("{owner}")
      .addResource("{name}")
      .addResource("{ref}")
      .addResource("result");

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
  }
}
