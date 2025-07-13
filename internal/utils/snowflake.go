package utils

import (
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
