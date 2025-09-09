package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/tkoleo84119/nail-salon-backend/internal/infra/redis"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
)

const (
	authCacheExpiration = time.Hour

	staffCacheKeyPrefix    = "auth:staff:"
	customerCacheKeyPrefix = "auth:customer:"
)

type AuthCache struct {
	redis *redis.Client
}

func NewAuthCache(redisClient *redis.Client) AuthCacheInterface {
	return &AuthCache{
		redis: redisClient,
	}
}

// GetStaffContext from Redis
func (c *AuthCache) GetStaffContext(ctx context.Context, userID int64) (*common.StaffContext, error) {
	key := fmt.Sprintf("%s%d", staffCacheKeyPrefix, userID)

	val, err := c.redis.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var staffContext common.StaffContext
	if err := json.Unmarshal([]byte(val), &staffContext); err != nil {
		// JSON unmarshal failed, delete the wrong cache data
		_ = c.redis.Del(ctx, key)
		return nil, fmt.Errorf("failed to unmarshal staff context: %w", err)
	}

	return &staffContext, nil
}

// SetStaffContext to Redis
func (c *AuthCache) SetStaffContext(ctx context.Context, userID int64, staffContext *common.StaffContext) error {
	key := fmt.Sprintf("%s%d", staffCacheKeyPrefix, userID)

	data, err := json.Marshal(staffContext)
	if err != nil {
		return fmt.Errorf("failed to marshal staff context: %w", err)
	}

	err = c.redis.Set(ctx, key, data, authCacheExpiration).Err()
	if err != nil {
		return fmt.Errorf("failed to set staff context to redis: %w", err)
	}

	return nil
}

// DeleteStaffContext from Redis
func (c *AuthCache) DeleteStaffContext(ctx context.Context, userID int64) error {
	key := fmt.Sprintf("%s%d", staffCacheKeyPrefix, userID)

	err := c.redis.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("failed to delete staff context from redis: %w", err)
	}

	return nil
}

// GetCustomerContext from Redis
func (c *AuthCache) GetCustomerContext(ctx context.Context, customerID int64) (*common.CustomerContext, error) {
	key := fmt.Sprintf("%s%d", customerCacheKeyPrefix, customerID)

	val, err := c.redis.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var customerContext common.CustomerContext
	if err := json.Unmarshal([]byte(val), &customerContext); err != nil {
		// JSON unmarshal failed, delete the wrong cache data
		_ = c.redis.Del(ctx, key)
		return nil, fmt.Errorf("failed to unmarshal customer context: %w", err)
	}

	return &customerContext, nil
}

// SetCustomerContext to Redis
func (c *AuthCache) SetCustomerContext(ctx context.Context, customerID int64, customerContext *common.CustomerContext) error {
	key := fmt.Sprintf("%s%d", customerCacheKeyPrefix, customerID)

	data, err := json.Marshal(customerContext)
	if err != nil {
		return fmt.Errorf("failed to marshal customer context: %w", err)
	}

	err = c.redis.Set(ctx, key, data, authCacheExpiration).Err()
	if err != nil {
		return fmt.Errorf("failed to set customer context to redis: %w", err)
	}

	return nil
}

// DeleteCustomerContext from Redis
func (c *AuthCache) DeleteCustomerContext(ctx context.Context, customerID int64) error {
	key := fmt.Sprintf("%s%d", customerCacheKeyPrefix, customerID)

	err := c.redis.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("failed to delete customer context from redis: %w", err)
	}

	return nil
}