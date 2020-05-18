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

const defaultCellValue = 1

func GetGridSize() (int, int, error) {
	nbRows, err := strconv.Atoi(os.Getenv("NB_ROWS"))
	if err != nil {
		return 0, 0, err
	}
	nb_cols, err := strconv.Atoi(os.Getenv("NB_COLS"))
	if err != nil {
		return 0, 0, err
	}
	return nbRows, nb_cols, nil
}

type Store struct {
	StoreModifier
}

type StoreModifier interface {
	Get(key string) (int, error)
	Set(key string, value int) error
}

type DynamoStoreModifier struct {
	client dynamodb.DynamoDB
	table  string
}

type Cell struct {
	X int
	Y int
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

func (c Cell) IsValid(nbRows int, nbCols int) bool {
	validX := 0 <= c.X && c.X < nbCols
	validY := 0 <= c.Y && c.Y < nbRows
	return validX && validY
}

func (modifier DynamoStoreModifier) Get(key string) (int, error) {
	result, err := modifier.client.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(modifier.table),
		Key: map[string]*dynamodb.AttributeValue{
			"Key": {
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
	vStr, err := strconv.Itoa(value)
	if err != nil {
		return "", err
	}

	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":v": {
				N: aws.String(vStr),
			},
		},
		TableName: aws.String(modifier.table),
		Key: map[string]*dynamodb.AttributeValue{
			"Key": {
				S: aws.String(key),
			},
		},
		ReturnValues:     aws.String("UPDATED_NEW"),
		UpdateExpression: aws.String("set V = :v"),
	}

	_, err = modifier.client.UpdateItem(input)
	return err
}

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
	return Store{modifier}, nil
}

func (store Store) GetGrid(nbRows int, nb_cols int) (map[string]int, error) {
	m := make(map[string]int)
	var err error
	for y := 0; y < nbRows; y++ {
		for x := 0; x < nb_cols; x++ {
			m[BuildKey(x, y)], err = store.GetOrSetCell(x, y)
			if err != nil {
				return m, err
			}
		}
	}
	return m, nil
}

func (store Store) GetOrSetCell(x int, y int) (int, error) {
	got, getErr := store.GetCell(x, y)
	switch errV := getErr.(type) {
	case *ItemNotFound:
		setErr := store.SetCell(x, y, defaultCellValue)
		if setErr != nil {
			return 0, setErr
		}
		return defaultCellValue, nil
	default:
		return got, errV
	}
}

// TODO: use dynamodb atomic operation
func (store Store) RevertState(x int, y int) error {
	state, err := store.GetCell(x, y)
	if err != nil {
		return err
	}
	otherState := 1 - state
	return store.SetCell(x, y, otherState)
}

func (store Store) GetCell(x int, y int) (int, error) {
	key := BuildKey(x, y)
	return store.Get(key)
}

func (store Store) SetCell(x int, y int, v int) error {
	key := BuildKey(x, y)
	return store.Set(key, v)
}

func BuildKey(x int, y int) string {
	return fmt.Sprintf("x:%v,y:%v", x, y)
}
