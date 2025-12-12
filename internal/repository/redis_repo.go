package repository

import (
	"collab-editor/internal/domain"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisRepo struct {
	client *redis.Client
}

func NewRedisRepo(connString string) (*RedisRepo, error) {
	opt, err := redis.ParseURL(connString)
	if err != nil {
		return nil, fmt.Errorf("invalid redis url: %v", err)
	}

	client := redis.NewClient(opt)

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %v", err)
	}

	return &RedisRepo{client: client}, nil
}

func (r *RedisRepo) GetByID(id string) (*domain.Document, error) {
	ctx := context.Background()

	val, err := r.client.Get(ctx, "doc:"+id).Result()

	if err == redis.Nil {
		return &domain.Document{ID: id, Operations: []domain.Operation{}}, nil
	} else if err != nil {
		return nil, err
	}

	var doc domain.Document
	err = json.Unmarshal([]byte(val), &doc)
	return &doc, err
}

func (r *RedisRepo) Save(doc *domain.Document) error {
	ctx := context.Background()

	data, err := json.Marshal(doc)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, "doc:"+doc.ID, data, 24*time.Hour).Err()
}
