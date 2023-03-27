package server

import (
	"context"
	"encoding/json"

	"github.com/redis/go-redis/v9"
)

type IStore interface {
	get(string) (Student, error)
	put(string, Student) error
	delete(string) error
}

type RedisStore struct {
	redisClient *redis.Client
	ctx         context.Context
}

func (r *RedisStore) get(key string) (Student, error) {
	res, err := r.redisClient.Get(r.ctx, key).Result()
	if err != nil {
		return Student{}, err
	}
	var student Student
	json.Unmarshal([]byte(res), &student)
	return student, nil
}

func (r *RedisStore) put(key string, student Student) error {
	val, _ := json.Marshal(student)
	err := r.redisClient.Set(r.ctx, key, val, 0).Err()
	return err
}

func (r *RedisStore) delete(key string) error {
	err := r.redisClient.Del(r.ctx, key).Err()
	return err
}

func NewRedisStore() *RedisStore {
	ctx := context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // default DB
	})

	return &RedisStore{
		redisClient: rdb,
		ctx:         ctx,
	}
}
