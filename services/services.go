package services

import (
	"github.com/aws/aws-sdk-go/service/apigateway"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/kraneware/kws/config"
	"sync"

	"github.com/aws/aws-sdk-go/service/glue"
	"github.com/aws/aws-sdk-go/service/glue/glueiface"

	"github.com/aws/aws-sdk-go/service/sagemaker"
	"github.com/aws/aws-sdk-go/service/ssm"

	"github.com/aws/aws-sdk-go/service/rds"

	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sqs"
)

var (
	lambdaClient     *lambda.Lambda // nolint:gochecknoglobals
	lambdaClientInit sync.Once      // nolint:gochecknoglobals

	apigwClient     *apigateway.APIGateway // nolint:gochecknoglobals
	apigwClientInit sync.Once              // nolint:gochecknoglobals

	s3Client     *s3.S3    // nolint:gochecknoglobals
	s3ClientInit sync.Once // nolint:gochecknoglobals

	stsClient     *sts.STS  // nolint:gochecknoglobals
	stsClientInit sync.Once // nolint:gochecknoglobals

	snsClient     *sns.SNS  // nolint:gochecknoglobals
	snsClientInit sync.Once // nolint:gochecknoglobals

	sqsClient     *sqs.SQS  // nolint:gochecknoglobals
	sqsClientInit sync.Once // nolint:gochecknoglobals

	cwClient     *cloudwatch.CloudWatch // nolint:gochecknoglobals
	cwClientInit sync.Once              // nolint:gochecknoglobals

	cwLogsClient     *cloudwatchlogs.CloudWatchLogs // nolint:gochecknoglobals
	cwLogsClientInit sync.Once                      // nolint:gochecknoglobals

	//xrayClient     *xray2.XRay // nolint:gochecknoglobals
	//xrayClientInit sync.Once   // nolint:gochecknoglobals

	rdsClient     *rds.RDS  // nolint:gochecknoglobals
	rdsClientInit sync.Once // nolint:gochecknoglobals

	sagemakerClient     *sagemaker.SageMaker // nolint:gochecknoglobals
	sagemakerClientInit sync.Once            // nolint:gochecknoglobals

	ssmClient     *ssm.SSM  // nolint:gochecknoglobals
	ssmClientInit sync.Once // nolint:gochecknoglobals

	glueClient     glueiface.GlueAPI // nolint:gochecknoglobals
	glueClientInit sync.Once         // nolint:gochecknoglobals

	secretClient     *secretsmanager.SecretsManager // nolint:gochecknoglobals
	secretClientInit sync.Once                      // nolint:gochecknoglobals
)

// LambdaClient returns an Lambda client singleton
func LambdaClient() *lambda.Lambda {
	lambdaClientInit.Do(func() {
		c := config.SessionConfig()
		if config.Endpoints.Lambda != "" {
			c = c.WithEndpoint(config.Endpoints.Lambda)
		}
		lambdaClient = lambda.New(config.NewSession(c))
		xray.AWS(lambdaClient.Client)
	})

	return lambdaClient
}

// SNSClient returns an SNS client singleton
func SNSClient() *sns.SNS {
	snsClientInit.Do(func() {
		c := config.SessionConfig()
		if config.Endpoints.SNS != "" {
			c = c.WithEndpoint(config.Endpoints.SNS)
		}
		snsClient = sns.New(config.NewSession(c))
		xray.AWS(snsClient.Client)
	})
	return snsClient
}

// SNSClientInRegion SNSClient returns an SNS client singleton
func SNSClientInRegion(region string) *sns.SNS {
	c := config.SessionConfig().WithRegion(region)
	if config.Endpoints.SNS != "" {
		c = c.WithEndpoint(config.Endpoints.SNS)
	}
	snsClient := sns.New(config.NewSession(c))
	xray.AWS(snsClient.Client)
	return snsClient
}

// SQSClient returns an SQS client singleton
func SQSClient() *sqs.SQS {
	sqsClientInit.Do(func() {
		c := config.SessionConfig()
		if config.Endpoints.SQS != "" {
			c = c.WithEndpoint(config.Endpoints.SQS)
		}
		sqsClient = sqs.New(config.NewSession(c))

		xray.AWS(sqsClient.Client)
	})

	return sqsClient
}

// S3Client returns an S3 client singleton
func S3Client() *s3.S3 {
	s3ClientInit.Do(func() {
		c := config.SessionConfig()
		if config.Endpoints.S3 != "" {
			c = config.LocalS3Config(c, config.Endpoints.S3)
		}
		s3Client = s3.New(config.NewSession(c))
		xray.AWS(s3Client.Client)
	})

	return s3Client
}

// S3Downloader returns a new S3 downloader
func S3Downloader() *s3manager.Downloader {
	c := config.SessionConfig()
	if config.Endpoints.S3 != "" {
		c = config.LocalS3Config(c, config.Endpoints.S3)
	}
	return s3manager.NewDownloader(config.NewSession(c))
}

// S3Uploader return a new S3 uploader
func S3Uploader() *s3manager.Uploader {
	c := config.SessionConfig()
	if config.Endpoints.S3 != "" {
		c = config.LocalS3Config(c, config.Endpoints.S3)
	}
	return s3manager.NewUploader(config.NewSession(c))
}

// CWLogsClient returns a new CloudWatch Logs client
func CWLogsClient() *cloudwatchlogs.CloudWatchLogs {
	cwLogsClientInit.Do(func() {
		c := config.SessionConfig()
		if config.Endpoints.CloudWatchLogs != "" {
			c = c.WithEndpoint(config.Endpoints.CloudWatchLogs)
		}
		cwLogsClient = cloudwatchlogs.New(config.NewSession(c))
	})

	return cwLogsClient
}

// CWClient returns a new CLoudWatch client
func CWClient() *cloudwatch.CloudWatch {
	cwClientInit.Do(func() {
		c := config.SessionConfig()
		if config.Endpoints.CloudWatch != "" {
			c = c.WithEndpoint(config.Endpoints.CloudWatch)
		}
		cwClient = cloudwatch.New(config.NewSession(c))
	})

	return cwClient
}

// RDSClient returns a new RDS client
func RDSClient() *rds.RDS {
	rdsClientInit.Do(func() {
		c := config.SessionConfig()
		if config.Endpoints.RDS != "" {
			c = c.WithEndpoint(config.Endpoints.RDS)
		}
		rdsClient = rds.New(config.NewSession(c))
	})

	return rdsClient
}

// SagemakerClient returns a new Sagemaker client
func SagemakerClient() (svc *sagemaker.SageMaker) { // nolint:interfacer
	sagemakerClientInit.Do(func() {
		c := config.SessionConfig()
		if config.Endpoints.Sagemaker != "" {
			c = c.WithEndpoint(config.Endpoints.Sagemaker)
		}
		sagemakerClient = sagemaker.New(config.NewSession(c))
		xray.AWS(sagemakerClient.Client)
	})

	return sagemakerClient
}

// SSMClient returns a new client for AWS Systems Manager Agent
func SSMClient() (svc *ssm.SSM) { // nolint:interfacer
	ssmClientInit.Do(func() {
		c := config.SessionConfig()
		if config.Endpoints.SSM != "" {
			c = c.WithEndpoint(config.Endpoints.SSM)
		}
		ssmClient = ssm.New(config.NewSession(c))
		xray.AWS(ssmClient.Client)
	})

	return ssmClient
}

func GlueClient() (svc glueiface.GlueAPI) { // nolint:interfacer
	glueClientInit.Do(func() {
		c := config.SessionConfig()
		glueClient = glue.New(config.NewSession(c))
	})

	return glueClient
}

func STSClient() (svc *sts.STS) { // nolint:interfacer
	stsClientInit.Do(func() {
		c := config.SessionConfig()
		if config.Endpoints.STS != "" {
			c = c.WithEndpoint(config.Endpoints.STS)
		}
		stsClient = sts.New(config.NewSession(c))
	})

	return stsClient
}

// APIGWClient returns an apigw client singleton
func APIGWClient() *apigateway.APIGateway {
	apigwClientInit.Do(func() {
		c := config.SessionConfig()
		if config.Endpoints.APIGateway != "" {
			c = c.WithEndpoint(config.Endpoints.APIGateway)
		}
		apigwClient = apigateway.New(config.NewSession(c))
		xray.AWS(apigwClient.Client)
	})

	return apigwClient
}

func SecretClient() *secretsmanager.SecretsManager {
	secretClientInit.Do(func() {
		c := config.SessionConfig()
		if config.Endpoints.SecretsManager != "" {
			c = c.WithEndpoint(config.Endpoints.SecretsManager)
		}
		secretClient = secretsmanager.New(config.NewSession(c))
	})

	return secretClient
}
