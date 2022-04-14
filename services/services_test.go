package services_test

import (
	"github.com/kraneware/kws/services"
	"testing"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
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

})
