package util

import (
	"github.com/bwmarrin/snowflake"
	"github.com/google/uuid"
)

var node *snowflake.Node

func init() {
	node, _ = snowflake.NewNode(1)
}

func SnowflakeId() int64 {
	return node.Generate().Int64()
}

func UUID() string {
	id, _ := uuid.NewV7()
	return id.String()
}
