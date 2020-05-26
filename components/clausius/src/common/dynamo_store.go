package common

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"os"
	"strconv"
)

func ConnectToFunes() (Store, error) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	svc := dynamodb.New(sess)
	table := os.Getenv("FUNES_TABLE")
	if table == "" {
		return Store{}, fmt.Errorf("Need to set env FUNES_TABLE!")
	}
	modifier := DynamoStoreModifier{client: *svc, table: table}
	return Store{StoreModifier: modifier}, nil
}

type dynamoItem struct {
	Key string
	V   int
}

type ItemNotFound struct {
	key string
}

func (e *ItemNotFound) Error() string {
	return fmt.Sprintf("No dynamodb item found for key: %v", e.key)
}

type DynamoStoreModifier struct {
	client dynamodb.DynamoDB
	table  string
}

func (modifier DynamoStoreModifier) Table() string {
	return modifier.table
}

func (modifier DynamoStoreModifier) Get(key string) (int, error) {
	result, err := modifier.client.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(modifier.Table()),
		Key: map[string]*dynamodb.AttributeValue{
			"K": {
				S: aws.String(key),
			},
		},
	})
	if err != nil {
		return 0, err
	}
	if len(result.Item) == 0 {
		return 0, &ItemNotFound{key}
	}

	item := dynamoItem{}
	err = dynamodbattribute.UnmarshalMap(result.Item, &item)
	return item.V, err
}

func (modifier DynamoStoreModifier) Set(key string, value int) error {
	vStr := strconv.Itoa(value)

	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":v": {
				N: aws.String(vStr),
			},
		},
		TableName: aws.String(modifier.Table()),
		Key: map[string]*dynamodb.AttributeValue{
			"K": {
				S: aws.String(key),
			},
		},
		ReturnValues:     aws.String("UPDATED_NEW"),
		UpdateExpression: aws.String("set V = :v"),
	}

	_, err := modifier.client.UpdateItem(input)
	return err
}

func (modifier DynamoStoreModifier) BatchGetItem(input *dynamodb.BatchGetItemInput) (*dynamodb.BatchGetItemOutput, error) {
	return modifier.client.BatchGetItem(input)
}

func BuildAllAttrValues(keys []string) []map[string]*dynamodb.AttributeValue {
	attrValues := []map[string]*dynamodb.AttributeValue{}
	for _, key := range keys {
		attrValue := BuildAttrValue(key)
		attrValues = append(attrValues, attrValue)
	}
	return attrValues
}

func BuildAttrValue(key string) map[string]*dynamodb.AttributeValue {
	return map[string]*dynamodb.AttributeValue{
		"K": &dynamodb.AttributeValue{
			S: aws.String(key),
		},
	}
}

func BuildRequestItems(attrValues []map[string]*dynamodb.AttributeValue, tableName string) map[string]*dynamodb.KeysAndAttributes {
	return map[string]*dynamodb.KeysAndAttributes{
		tableName: {
			Keys:                 attrValues,
			ProjectionExpression: aws.String("V,K"),
		},
	}
}

func BuildBGIInput(requestItems map[string]*dynamodb.KeysAndAttributes) *dynamodb.BatchGetItemInput {
	return &dynamodb.BatchGetItemInput{
		RequestItems: requestItems,
	}
}

func FillMap(output *dynamodb.BatchGetItemOutput, tableName string) (map[string]int, error) {
	result := map[string]int{}
	for _, item := range output.Responses[tableName] {
		key := *item["K"].S
		v, err := strconv.Atoi(*item["V"].N)
		if err != nil {
			return result, err
		}
		result[key] = v
	}
	return result, nil
}
