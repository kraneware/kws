package kc2_test

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/kraneware/kws/kc2"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

func TestAWSServices(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "AWS EC2 Test Suite")
}

var _ = Describe("EC2 Test Suite", func() {

	Context("EC2 Test", func() {
		It("should create ec2 instance", func() {
			svc := kc2.EC2Client()
			vols := kc2.LoadAllVolumes(
				svc,
				[]*ec2.Filter{
					&ec2.Filter{
						Name: aws.String("volume-type"),
						Values: []*string{
							aws.String("gp2"),
						},
					},
				},
				[]string{"us-east-1", "us-east-2", "us-west-1", "us-west-2"},
			)

			Expect(vols).ToNot(BeNil())
		})
	})
})
