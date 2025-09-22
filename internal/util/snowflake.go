package util

import "github.com/bwmarrin/snowflake"

var node *snowflake.Node

func init() {
	node, _ = snowflake.NewNode(1)
}

func SnowflakeId() int64 {
	return node.Generate().Int64()
}
