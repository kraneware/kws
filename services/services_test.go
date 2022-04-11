package services_test

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/kraneware/kws/services"
	"testing"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	_ "github.com/aws/aws-xray-sdk-go/xray"

	_ "github.com/aws/aws-sdk-go/aws"
	_ "github.com/aws/aws-sdk-go/service/s3/s3manager"
	_ "github.com/aws/aws-sdk-go/service/xray"
)

func TestAWSServices(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "AWS Services Test Suite")
}

type nestedStruct struct {
	Stuff       string   `json:"stuff"`
	SomeStrings []string `json:"someStrings"`
}

type testStruct struct {
	ID          string       `json:"id"`
	Description string       `json:"description"`
	Nested      nestedStruct `json:"nested"`
}

var _ = Describe("DynamoDB Util", func() {
	Context("Mapping Test", func() {
		It("should marshal and unmarshal correctly", func() {
			t := testStruct{
				ID:          uuid.New().String(),
				Description: uuid.New().String(),
				Nested: nestedStruct{
					Stuff: uuid.New().String(),
					SomeStrings: []string{
						"abc",
						"def",
					},
				},
			}

			eventMap, err := services.MarshalStreamImage(t)
			Expect(err).Should(BeNil())

			var t1 testStruct
			Expect(services.UnmarshalStreamImage(eventMap, &t1)).Should(BeNil())

			Expect(t1).Should(Equal(t))
		})
	})
	Context("EC2 Test", func() {
		It("should create ec2 instance", func() {
			svc := services.EC2Client()
			vols := services.LoadAllVolumes(svc, []*ec2.Filter{
				{
					Name: aws.String("volume-type"),
					Values: []*string{
						aws.String("gp2"),
					},
				},
			}, []string{"us-east-1", "us-east-2", "us-west-1", "us-west-2"})

			Expect(vols).ToNot(BeNil())
		})
	})
})
