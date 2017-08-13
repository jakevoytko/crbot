package main

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

// Creates a new Redis client, connects to it, and returns a RedisStringMap.
func NewRedisStringMap(bucket string) (*RedisStringMap, error) {
	// Initialize redis.
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	if _, err := redisClient.Ping().Result(); err != nil {
		return nil, err
	}
	return &RedisStringMap{
		redisClient: redisClient,
		bucket:      bucket,
	}, nil
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
