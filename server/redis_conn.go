package main

import (
	"fmt"
	"net"
)

type RedisConn struct {
	net.Conn
}

func NewRedisConn(c net.Conn) *RedisConn {
	return &RedisConn{Conn: c}
}

func (rc *RedisConn) WriteError(msg string) {
	fmt.Fprintf(rc, "-ERR %s\r\n", msg)
}

func (rc *RedisConn) WriteErrorWrongArgCount() {
	rc.WriteError("wrong number of arguments")
}

func (rc *RedisConn) WriteSimpleString(msg string) {
	fmt.Fprintf(rc, "+%s\r\n", msg)
}

func (rc *RedisConn) WriteOk() {
	rc.WriteSimpleString("OK")
}

func (rc *RedisConn) WriteBulkString(msg string) {
	fmt.Fprintf(rc, "$%d\r\n%s\r\n", len(msg), msg)
}

func (rc *RedisConn) WriteNullBulkString() {
	fmt.Fprint(rc, "$-1\r\n")
}

func (rc *RedisConn) WriteArray(n int) {
	fmt.Fprintf(rc, "*%d\r\n", n)
}
