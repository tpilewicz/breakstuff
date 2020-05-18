package common

import (
	"fmt"
	"github.com/stretchr/testify/mock"
)

type MockStoreModifier struct {
	mock.Mock
}

func (store *MockStoreModifier) Get(key string) (int, error) {
	rtns := store.Called(key)
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
func (store *MockStoreModifier) Set(key string, value int) error {
	rtns := store.Called(key, value)
	return rtns.Error(0)
}
