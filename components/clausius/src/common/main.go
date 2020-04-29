package common

import (
	"fmt"
	"github.com/go-redis/redis/v7"
	"os"
	"strconv"
	"time"
)

const defaultCellValue = 1

func GetGridSize() (int, int, error) {
	nb_rows, err := strconv.Atoi(os.Getenv("NB_ROWS"))
	if err != nil {
		return 0, 0, err
	}
	nb_cols, err := strconv.Atoi(os.Getenv("NB_COLS"))
	if err != nil {
		return 0, 0, err
	}
	return nb_rows, nb_cols, nil
}

type Store struct {
	StoreModifier
}

type StoreModifier interface {
	Get(key string) (string, error)
	Set(key string, value interface{}, expiration time.Duration) error
}

type RedisStoreModifier struct {
	client redis.Client
}

func (modifier RedisStoreModifier) Get(key string) (string, error) {
	return modifier.client.Get(key).Result()
}

func (modifier RedisStoreModifier) Set(key string, value interface{}, expiration time.Duration) error {
	return modifier.client.Set(key, value, expiration).Err()
}

func ConnectToFunes() (Store, error) {
	opts, err := redis.ParseURL(os.Getenv("FUNES_URL"))
	if err != nil {
		return Store{}, err
	}
	return Store{
		RedisStoreModifier{*redis.NewClient(opts)},
	}, nil
}

func (store Store) GetGrid(nb_rows int, nb_cols int) (map[string]int, error) {
	m := make(map[string]int)
	var err error
	for y := 0; y < nb_rows; y++ {
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
	if getErr == redis.Nil {
		setErr := store.SetCell(x, y, defaultCellValue)
		if setErr != nil {
			return 0, setErr
		}
		return defaultCellValue, nil
	}
	return got, getErr
}

// TODO: write the test
func (store Store) revertState(x int, y int) error {
	state, err := store.GetCell(x, y)
	if err != nil {
		return err
	}
	otherState := 1 - state
	return store.SetCell(x, y, otherState)
}

func (store Store) GetCell(x int, y int) (int, error) {
	key := BuildKey(x, y)
	value, err := store.Get(key)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(value)
}

func (store Store) SetCell(x int, y int, v int) error {
	key := BuildKey(x, y)
	return store.Set(key, v, 0)
}

func BuildKey(x int, y int) string {
	return fmt.Sprintf("x:%v,y:%v", x, y)
}
