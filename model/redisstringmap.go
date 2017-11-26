package model

import (
	"errors"

	"github.com/go-redis/redis"
)

// RedisStringMap is a StringMap implementation that uses Redis as a backing store.
// TODO(jake): Parameterize the Redis connection.
type RedisStringMap struct {
	redisClient *redis.Client
	bucket      string
}

// NewRedisStringMap works as advertised.
func NewRedisStringMap(redisClient *redis.Client, bucket string) *RedisStringMap {
	return &RedisStringMap{
		redisClient: redisClient,
		bucket:      bucket,
	}
}

// Has tests whether the map contains the key.
func (m *RedisStringMap) Has(key string) (bool, error) {
	result := m.redisClient.HExists(m.bucket, key)
	if result.Err() != nil {
		return false, result.Err()
	}
	return result.Val(), nil
}

// Get retrieves the given key.
func (m *RedisStringMap) Get(key string) (string, error) {
	result := m.redisClient.HGet(m.bucket, key)
	if result.Err() != nil {
		return "", result.Err()
	}
	return result.Val(), nil
}

// Set sets the given key.
func (m *RedisStringMap) Set(key, value string) error {
	// Don't care about overwriting vs new.
	if result := m.redisClient.HSet(m.bucket, key, value); result.Err() != nil {
		return result.Err()
	}
	return nil
}

// Delete removes the key.
func (m *RedisStringMap) Delete(key string) error {
	result := m.redisClient.HDel(m.bucket, key)
	if result.Err() != nil {
		return result.Err()
	} else if result.Val() != 1 {
		return errors.New("Did not successfully delete a key")
	}
	return nil
}

// GetAll returns all keys and values.
func (m *RedisStringMap) GetAll() (map[string]string, error) {
	result := m.redisClient.HGetAll(m.bucket)
	if result.Err() != nil {
		return nil, result.Err()
	}
	return result.Val(), nil
}

// ScanKeys returns all of the keys that match the given pattern.
func (m *RedisStringMap) ScanKeys(pattern string) ([]string, error) {
	// HScan does not guarantee that keys won't be returned multiple times, so use
	// a map to try to produce a consistent snapshot.
	keyMap := map[string]bool{}

	scanCmd := m.redisClient.HScan(m.bucket, 0 /* cursor */, pattern, 100 /* count */)

	iter := scanCmd.Iterator()
	for iter.Next() {
		keyMap[iter.Val()] = true
	}

	if scanCmd.Err() != nil {
		return nil, scanCmd.Err()
	}

	results := []string{}
	for key := range keyMap {
		results = append(results, key)
	}
	return results, nil
}
