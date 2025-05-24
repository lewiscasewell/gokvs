package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	storeConfig := StoreConfig{Persist: true, SnapshotFile: "dump.rdb", SnapshotInterval: 30 * time.Second}
	store := NewStore(storeConfig)

	ln, err := net.Listen("tcp", "127.0.0.1:6379")
	if err != nil {
		panic(err)
	}
	fmt.Println("mini redis running locally on 127.0.0.1:6379")

	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		redisConn := NewRedisConn(conn)
		go handleConnection(redisConn, store)
	}
}
