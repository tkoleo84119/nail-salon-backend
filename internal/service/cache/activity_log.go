package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/tkoleo84119/nail-salon-backend/internal/infra/redis"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

const (
	activityLogKey        = "activity_logs"
	activityLogMaxEntries = 50
	activityLogExpiration = 24 * time.Hour
)

type ActivityLogCache struct {
	redis *redis.Client
}

func NewActivityLogCache(redisClient *redis.Client) ActivityLogCacheInterface {
	return &ActivityLogCache{
		redis: redisClient,
	}
}

// LogActivity to Redis List
func (c *ActivityLogCache) LogActivity(ctx context.Context, activityType common.ActivityLogType, message string) error {
	entry := common.ActivityLogEntry{
		ID:        utils.FormatID(utils.GenerateID()),
		Type:      activityType,
		Message:   message,
		Timestamp: time.Now(),
	}

	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("failed to marshal activity log entry: %w", err)
	}

	// push new record to list
	if err := c.redis.LPush(ctx, activityLogKey, data).Err(); err != nil {
		return fmt.Errorf("failed to push activity log to redis: %w", err)
	}

	// keep list with max 50 elements
	if err := c.redis.LTrim(ctx, activityLogKey, 0, activityLogMaxEntries-1).Err(); err != nil {
		return fmt.Errorf("failed to trim activity log list: %w", err)
	}

	// set expiration time
	if err := c.redis.Expire(ctx, activityLogKey, activityLogExpiration).Err(); err != nil {
		return fmt.Errorf("failed to set activity log expiration: %w", err)
	}

	return nil
}

// GetRecentActivities from Redis List
func (c *ActivityLogCache) GetRecentActivities(ctx context.Context, limit int) (*common.ActivityLogResponse, error) {
	if limit <= 0 {
		limit = activityLogMaxEntries
	}

	if limit > activityLogMaxEntries {
		limit = activityLogMaxEntries
	}

	// get specified range records
	results, err := c.redis.LRange(ctx, activityLogKey, 0, int64(limit-1)).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get activity logs from redis: %w", err)
	}

	activities := make([]common.ActivityLogEntry, 0, len(results))
	for _, result := range results {
		var entry common.ActivityLogEntry
		if err := json.Unmarshal([]byte(result), &entry); err != nil {
			// if some record unmarshal failed, skip it instead of returning error
			continue
		}
		activities = append(activities, entry)
	}

	return &common.ActivityLogResponse{
		Activities: activities,
		Total:      len(activities),
	}, nil
}

// LogCustomerRegister to Redis List
func (c *ActivityLogCache) LogCustomerRegister(ctx context.Context, customerName string, lineName string) error {
	message := ""
	if lineName == "" {
		message = fmt.Sprintf("顧客 %s 完成註冊", customerName)
	} else {
		message = fmt.Sprintf("顧客 %s (LINE：%s) 完成註冊", customerName, lineName)
	}
	return c.LogActivity(ctx, common.ActivityCustomerRegister, message)
}

// LogCustomerLogin to Redis List
func (c *ActivityLogCache) LogCustomerBrowse(ctx context.Context, customerName string, lineName string) error {
	message := ""
	if lineName == "" {
		message = fmt.Sprintf("顧客 %s 進入網頁瀏覽", customerName)
	} else {
		message = fmt.Sprintf("顧客 %s (LINE：%s) 進入網頁瀏覽", customerName, lineName)
	}
	return c.LogActivity(ctx, common.ActivityCustomerBrowse, message)
}

// LogCustomerBooking to Redis List
func (c *ActivityLogCache) LogCustomerBooking(ctx context.Context, customerName string, lineName string, storeName string) error {
	message := ""
	if lineName == "" {
		message = fmt.Sprintf("顧客 %s 建立預約 (門市：%s)", customerName, storeName)
	} else {
		message = fmt.Sprintf("顧客 %s (LINE：%s) 建立預約 (門市：%s)", customerName, lineName, storeName)
	}
	return c.LogActivity(ctx, common.ActivityCustomerBooking, message)
}

// LogCustomerBookingUpdate to Redis List
func (c *ActivityLogCache) LogCustomerBookingUpdate(ctx context.Context, customerName string, lineName string, storeName string) error {
	message := ""
	if lineName == "" {
		message = fmt.Sprintf("顧客 %s 修改預約 (門市：%s)", customerName, storeName)
	} else {
		message = fmt.Sprintf("顧客 %s (LINE：%s) 修改預約 (門市：%s)", customerName, lineName, storeName)
	}
	return c.LogActivity(ctx, common.ActivityCustomerBookingUpdate, message)
}

// LogCustomerBookingCancel to Redis List
func (c *ActivityLogCache) LogCustomerBookingCancel(ctx context.Context, customerName string, lineName string, storeName string) error {
	message := ""
	if lineName == "" {
		message = fmt.Sprintf("顧客 %s 取消預約 (門市：%s)", customerName, storeName)
	} else {
		message = fmt.Sprintf("顧客 %s (LINE：%s) 取消預約 (門市：%s)", customerName, lineName, storeName)
	}
	return c.LogActivity(ctx, common.ActivityCustomerBookingCancel, message)
}

// LogAdminBookingCreate to Redis List
func (c *ActivityLogCache) LogAdminBookingCreate(ctx context.Context, staffName string, customerName string, lineName string, storeName string) error {
	message := ""
	if lineName == "" {
		message = fmt.Sprintf("員工 %s 為顧客 %s 建立預約 (門市：%s)", staffName, customerName, storeName)
	} else {
		message = fmt.Sprintf("員工 %s 為顧客 %s (LINE：%s) 建立預約 (門市：%s)", staffName, customerName, lineName, storeName)
	}
	return c.LogActivity(ctx, common.ActivityAdminBookingCreate, message)
}

// LogAdminBookingUpdate to Redis List
func (c *ActivityLogCache) LogAdminBookingUpdate(ctx context.Context, staffName string, customerName string, lineName string, storeName string) error {
	message := ""
	if lineName == "" {
		message = fmt.Sprintf("員工 %s 為顧客 %s 修改預約 (門市：%s)", staffName, customerName, storeName)
	} else {
		message = fmt.Sprintf("員工 %s 為顧客 %s (LINE：%s) 修改預約 (門市：%s)", staffName, customerName, lineName, storeName)
	}
	return c.LogActivity(ctx, common.ActivityAdminBookingUpdate, message)
}

// LogAdminBookingCancel to Redis List
func (c *ActivityLogCache) LogAdminBookingCancel(ctx context.Context, staffName string, customerName string, lineName string, storeName string) error {
	message := ""
	if lineName == "" {
		message = fmt.Sprintf("員工 %s 為顧客 %s 取消預約 (門市：%s)", staffName, customerName, storeName)
	} else {
		message = fmt.Sprintf("員工 %s 為顧客 %s (LINE：%s) 取消預約 (門市：%s)", staffName, customerName, lineName, storeName)
	}
	return c.LogActivity(ctx, common.ActivityAdminBookingCancel, message)
}

// LogAdminBookingCompleted to Redis List
func (c *ActivityLogCache) LogAdminBookingCompleted(ctx context.Context, staffName string, customerName string, lineName string, checkoutCount int, storeName string) error {
	message := ""
	if lineName == "" {
		message = fmt.Sprintf("員工 %s 為顧客 %s 的預約完成結帳 %d 筆 (門市：%s)", staffName, customerName, checkoutCount, storeName)
	} else {
		message = fmt.Sprintf("員工 %s 為顧客 %s (LINE：%s) 的預約完成結帳 %d 筆 (門市：%s)", staffName, customerName, lineName, checkoutCount, storeName)
	}
	return c.LogActivity(ctx, common.ActivityAdminBookingCompleted, message)
}
