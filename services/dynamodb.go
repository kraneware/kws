package services

import (
	"encoding/json"
	"github.com/kraneware/kws/config"
	"sync"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

var (
	dynamoDbClient     *dynamodb.DynamoDB // nolint:gochecknoglobals
	dynamoDbClientInit sync.Once          // nolint:gochecknoglobals
)

// DynamoDbClient returns an DynamoDB client singleton
func DynamoDbClient() *dynamodb.DynamoDB {
	dynamoDbClientInit.Do(func() {
		c := config.SessionConfig()
		if config.Endpoints.DynamoDB != "" {
			c = c.WithEndpoint(config.Endpoints.DynamoDB)
		}
		dynamoDbClient = dynamodb.New(config.NewSession(c))

		//xray.AWS(dynamoDbClient.Client)
	})

	return dynamoDbClient
}

// UnmarshalStreamImage coverts images incoming from DynamoDB streams to given struct
func UnmarshalStreamImage(attribute map[string]events.DynamoDBAttributeValue, out interface{}) (err error) {
	dbAttrMap := make(map[string]*dynamodb.AttributeValue)

	for k, v := range attribute {
		if err == nil {
			var dbAttr dynamodb.AttributeValue

			var bytes []byte
			bytes, err = json.Marshal(v)

			if err == nil {
				err = json.Unmarshal(bytes, &dbAttr)
				if err == nil {
					dbAttrMap[k] = &dbAttr
				}
			}
		}
	}

	if err == nil {
		err = dynamodbattribute.UnmarshalMap(dbAttrMap, out)
	}

	return err
}

// removeNullValues removes keys with null values.  This is needed due to issue in
// marshaling and unmashaling events.DynamoDBAttributeValue
func removeNullValues(in map[string]interface{}) {
	for k, v := range in {
		if v == nil {
			delete(in, k)
		} else if m, ok := v.(map[string]interface{}); ok {
			removeNullValues(m)
		} else if l, ok := v.([]interface{}); ok { // handle the case of a list of value
			for _, i := range l {
				if m, ok := i.(map[string]interface{}); ok {
					removeNullValues(m)
				}
			}
		}
	}
}

// MarshalStreamImage takes a struct and converts it into stream image for streaming DynamoDB events
func MarshalStreamImage(in interface{}) (result map[string]events.DynamoDBAttributeValue, err error) {
	result = make(map[string]events.DynamoDBAttributeValue)

	var avMap map[string]*dynamodb.AttributeValue
	avMap, err = dynamodbattribute.MarshalMap(in)
	if err == nil {
		for k, v := range avMap {
			var jsonBytes []byte
			jsonBytes, err = json.Marshal(v)
			if err == nil {
				var tempMap map[string]interface{}
				err = json.Unmarshal(jsonBytes, &tempMap)
				if err == nil {
					removeNullValues(tempMap)

					var bytes []byte
					bytes, err = json.Marshal(tempMap)

					if err == nil {
						var item events.DynamoDBAttributeValue
						err = json.Unmarshal(bytes, &item)
						if err == nil {
							result[k] = item
						}
					}
				}
			}
		}
	}

	return result, err
}
