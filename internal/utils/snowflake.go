package utils

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/bwmarrin/snowflake"
)

var (
	snowflakeNode *snowflake.Node
	once          sync.Once
)

// InitSnowflake initializes the global snowflake node
func InitSnowflake(nodeID int64) error {
	var err error
	once.Do(func() {
		snowflakeNode, err = snowflake.NewNode(nodeID)
	})
	return err
}

// GenerateID generates a new Snowflake ID
func GenerateID() int64 {
	if snowflakeNode == nil {
		panic("snowflake node not initialized, call InitSnowflake first")
	}
	return snowflakeNode.Generate().Int64()
}

// ParseID converts string ID to int64
func ParseID(idStr string) (int64, error) {
	if idStr == "" {
		return 0, fmt.Errorf("ID string is empty")
	}
	return strconv.ParseInt(idStr, 10, 64)
}

// FormatID converts int64 ID to string
func FormatID(id int64) string {
	return strconv.FormatInt(id, 10)
}

// ParseIDSlice converts string slice to int64 slice
func ParseIDSlice(idStrs []string) ([]int64, error) {
	if len(idStrs) == 0 {
		return nil, fmt.Errorf("ID string slice is empty")
	}

	ids := make([]int64, len(idStrs))
	for i, idStr := range idStrs {
		id, err := ParseID(idStr)
		if err != nil {
			return nil, fmt.Errorf("invalid ID at index %d: %w", i, err)
		}
		ids[i] = id
	}
	return ids, nil
}

// FormatIDSlice converts int64 slice to string slice
func FormatIDSlice(ids []int64) []string {
	if len(ids) == 0 {
		return []string{}
	}

	idStrs := make([]string, len(ids))
	for i, id := range ids {
		idStrs[i] = FormatID(id)
	}
	return idStrs
}
