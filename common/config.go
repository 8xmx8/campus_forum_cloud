package common

import (
	"time"
)

type RedisClientOpt struct {
	Type string `json:",default=node,options=node|cluster|sentinel"` // node:单节点模式、cluster:普通集群模式、sentinel:哨兵集群模式
	// Network type to use, either tcp or unix.
	// Default is tcp.
	Network string `json:",optional"`

	// Redis server address in "host:port" format.
	Addr string `json:",optional"`

	// Username to authenticate the current connection when Redis ACLs are used.
	// See: https://redis.io/commands/auth.
	Username string `json:",optional"`

	// Password to authenticate the current connection.
	// See: https://redis.io/commands/auth.
	Password string `json:",optional"`

	// Redis DB to select after connecting to a server.
	// See: https://redis.io/commands/select.
	DB int `json:",optional"`

	// Dial timeout for establishing new connections.
	// Default is 5 seconds.
	DialTimeout time.Duration `json:",optional"`

	// Timeout for socket reads.
	// If timeout is reached, read commands will fail with a timeout error
	// instead of blocking.
	//
	// Use value -1 for no timeout and 0 for default.
	// Default is 3 seconds.
	ReadTimeout time.Duration `json:",optional"`

	// Timeout for socket writes.
	// If timeout is reached, write commands will fail with a timeout error
	// instead of blocking.
	//
	// Use value -1 for no timeout and 0 for default.
	// Default is ReadTimout.
	WriteTimeout time.Duration `json:",optional"`

	// Maximum number of socket connections.
	// Default is 10 connections per every CPU as reported by runtime.NumCPU.
	PoolSize int `json:",optional"`
}
