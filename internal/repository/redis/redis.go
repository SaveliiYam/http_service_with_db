package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/SaveliiYam/http_service_with_db/internal/models"
	"github.com/redis/go-redis/v9"
)

type RedisRepository struct {
	client *redis.Client
}

const (
	keyUser   = "users:%d"
	idCounter = "user_counter"
)

func NewRedisRepository() (*RedisRepository, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &RedisRepository{client: client}, nil
}

func (r *RedisRepository) GetByID(
	ctx context.Context,
	id int,
) (models.User, error) {
	userRaw, err := r.client.Get(ctx, fmt.Sprintf(keyUser, id)).Bytes()
	if err != nil {
		return models.User{}, err
	}

	var user models.User

	if err := json.Unmarshal(userRaw, &user); err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (r *RedisRepository) CreateUser(
	ctx context.Context,
	name string,
) (models.User, error) {
	id, err := r.client.Incr(ctx, idCounter).Result()
	if err != nil {
		return models.User{}, err
	}

	user := models.User{
		ID:   int(id),
		Name: name,
	}

	userRaw, err := json.Marshal(user)
	if err != nil {
		return models.User{}, err
	}

	if err := r.client.Set(ctx, fmt.Sprintf(keyUser, id), userRaw, 0).Err(); err != nil {
		return models.User{}, err
	}

	return user, nil
}
