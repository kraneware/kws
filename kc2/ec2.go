package kc2

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/kraneware/kws/config"
	"os"
	"sync"
)

var (
	ec2Client     *ec2.EC2  // nolint:gochecknoglobals
	ec2ClientInit sync.Once // nolint:gochecknoglobals
)

func EC2Client() *ec2.EC2 { // nolint:gochecknoglobals
	c := config.SessionConfig()

	ec2ClientInit.Do(func() {
		if config.Endpoints.EC2 != "" {
			c = c.WithEndpoint(config.Endpoints.EC2)
		}

		ec2Client = ec2.New(config.NewSession(c))
	})

	return ec2Client
}

func LoadAllVolumes(svc *ec2.EC2, filters []*ec2.Filter, regions []string) (volumes []*ec2.Volume) {
	volumes = make([]*ec2.Volume, 10)

	input := &ec2.DescribeVolumesInput{
		MaxResults: aws.Int64(100),
		Filters:    filters,
	}

	for _, region := range regions {
		config.Region = region
		svc = EC2Client()

		dvo, err := svc.DescribeVolumes(input)
		if err != nil {
			panic(err)
		}

		for i, v := range dvo.Volumes {
			fmt.Fprintf(os.Stdout, "%d: %+v", i, v)
			volumes = append(volumes, v)
		}

		nextToken := dvo.NextToken
		for nextToken != nil {
			input.NextToken = nextToken
			dvo, err := svc.DescribeVolumes(input)

			if err != nil {
				panic(err)
			}

			for i, v := range dvo.Volumes {
				fmt.Fprintf(os.Stdout, "%d: %+v", i, v)
				volumes = append(volumes, v)
			}

			nextToken = dvo.NextToken
		}
	}

	return volumes
}
