package order

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/dreamsofcode-io/orders-api/model"
	"github.com/redis/go-redis/v9"
)

type RedisRepo struct {
	Client *redis.Client
}

func orderIDKey(id uint64) string {
	return fmt.Sprintf("order:%d", id)
}

func(r *RedisRepo) Insert(ctx context.Context, order model.Order) error {
	data, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("failed to encode order: %w", err)
	}

	key := orderIDKey(order.OrderID)

	res := r.Client.SetNX(ctx, key, string(data), 0)	
	if err := res.Err(); err != nil {
		return fmt.Errorf("failed to set: %w", err)
	}

	if err := r.Client.SAdd(ctx, "order", key).Err(); err != nil {
		return fmt.Errorf("failed to add to orders set: %w", err)
	}

	return nil
}

var ErrNotExist = errors.New("order does not exist")

func(r *RedisRepo) FindByID(ctx context.Context, id uint64) (model.Order, error) {
	key := orderIDKey(id)

	value, err := r.Client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return model.Order{}, ErrNotExist
	} else if err != nil {
		return model.Order{}, fmt.Errorf("get order: %w", err)
	}

	var order model.Order
	err = json.Unmarshal([]byte(value), &order)
	if err != nil {
		return model.Order{}, fmt.Errorf("failed to decode order json : %w", err)
	}

	return order, nil
}


func(r *RedisRepo) DeleteByID(ctx context.Context, id uint64) error {
	key := orderIDKey(id)

	err := r.Client.Del(ctx, key).Err()
	if errors.Is(err, redis.Nil) {
		return ErrNotExist
	} else if err != nil {
		return fmt.Errorf("delete order: %w", err)
	}

	return nil
}


func(r *RedisRepo) Update(ctx context.Context, order model.Order) error {
	data, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("failed to encode order: %w", err)
	}

	key := orderIDKey(order.OrderID)

	err = r.Client.SetXX(ctx, key, string(data), 0).Err()
	if errors.Is(err, redis.Nil) {
		return ErrNotExist
	} else if err != nil {
		return fmt.Errorf("set order: %w", err)
	}

	return nil
}

type FindAllPage struct {
	Size	uint
	Offset	uint
}

func(r *RedisRepo) FindAll(ctx context.Context) ([]model.Order, error) {



	return nil, nil
}