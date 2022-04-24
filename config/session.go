package config

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

// NewSession creates a new AWS session to interact with
func NewSession(config *aws.Config) *session.Session {
	return session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Config:            *config,
	}))
}

func SessionConfig() *aws.Config {
	c := aws.NewConfig()

	if Region != "" {
		c = c.WithRegion(Region)
	}

	if Credentials != nil {
		c = c.WithCredentials(Credentials)
	}

	return c
}

func LocalS3Config(c *aws.Config, endpoint string) *aws.Config {
	return c.WithEndpoint(endpoint).WithS3ForcePathStyle(true).WithDisableSSL(true)
}
