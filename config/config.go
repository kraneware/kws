package config

import (
	"github.com/aws/aws-sdk-go/aws/credentials"
)

// Credentials defines any custom credentials for AWS
var Credentials *credentials.Credentials // nolint:gochecknoglobals

// Region defines a custom region for AWS
var Region string // nolint:gochecknoglobals

// Endpoints definees the global variables for definition of endpoints
var Endpoints AwsEndpointSet // nolint:gochecknoglobals

// AwsEndpointSet tracks any custom endpoints especially used when using localstack
type AwsEndpointSet struct {
	DynamoDB       string
	S3             string
	Lambda         string
	SNS            string
	SQS            string
	CloudWatch     string
	CloudWatchLogs string
	XRay           string
	RDS            string
	Sagemaker      string
	SSM            string
	APIGateway     string
	EC2            string
	SecretsManager string
}
