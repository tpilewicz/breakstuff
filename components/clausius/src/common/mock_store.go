package common

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/stretchr/testify/mock"
)

type MockStoreModifier struct {
	mock.Mock
	table string
}

func (modifier *MockStoreModifier) Table() string {
	return modifier.table
}

// uh so *MockStoreModifier implements StoreModifier, but the dynamo one is not
// a pointer?
func (modifier *MockStoreModifier) Get(key string) (int, error) {
	rtns := modifier.Called(key)
	i := rtns.Int(0)
	switch errValue := rtns.Get(1).(type) {
	case error:
		return i, errValue
	case nil:
		return i, nil
	default:
		return i, fmt.Errorf("Wait, this is not an error: %v. I should be returning an error.", errValue)
	}
}
func (modifier *MockStoreModifier) Set(key string, value int) error {
	rtns := modifier.Called(key, value)
	return rtns.Error(0)
}

func (modifier *MockStoreModifier) BatchGetItem(input *dynamodb.BatchGetItemInput) (*dynamodb.BatchGetItemOutput, error) {
	rtns := modifier.Called(input)
	switch outputValue := rtns.Get(0).(type) {
	case *dynamodb.BatchGetItemOutput:
		return outputValue, rtns.Error(1)
	default:
		return nil, rtns.Error(1)
	}
}
