package main

import (
	"bufio"
	"strings"
	"time"
)

func handleConnection(rc *RedisConn, s *Store) {
	defer rc.Close()
	reader := bufio.NewReader(rc)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return
		}

		parts := strings.Fields(strings.TrimSpace(line))
		if len(parts) == 0 {
			continue
		}

		cmd := strings.ToUpper(parts[0])

		switch cmd {
		case "PING":
			if len(parts) > 1 {
				rc.WriteErrorWrongArgCount()
			}

			rc.WriteSimpleString("PONG")
		case "SET":
			if len(parts) < 3 {
				rc.WriteErrorWrongArgCount()
				continue
			}

			key := parts[1]
			val := parts[2]
			var expiry *time.Time
			if len(parts) == 5 && strings.ToUpper(parts[3]) == "EX" {
				dur, err := time.ParseDuration(parts[4])
				if err != nil {
					rc.WriteError("invalid expiration duration")
					continue
				}
				exp := time.Now().Add(dur)
				expiry = &exp
			}

			s.Set(key, val, expiry)
			rc.WriteOk()
		case "GET":
			if len(parts) != 2 {
				rc.WriteErrorWrongArgCount()
				continue
			}
			if val, ok := s.Get(parts[1]); ok {
				rc.WriteBulkString(val)
			} else {
				rc.WriteNullBulkString()
			}
		case "DEL":
			if len(parts) != 2 {
				rc.WriteErrorWrongArgCount()
				continue
			}
			if val, ok := s.Del(parts[1]); ok {
				rc.WriteBulkString(val)
			} else {
				rc.WriteNullBulkString()
			}
		case "KEYS":
			if len(parts) > 1 {
				rc.WriteErrorWrongArgCount()
				continue
			}
			keys := s.GetAll()

			if len(keys) == 0 {
				rc.WriteArray(0)
				continue
			}

			rc.WriteArray(len(keys))
			for _, key := range keys {
				rc.WriteBulkString(key)
			}
		case "FLUSHALL":
			if len(parts) > 1 {
				rc.WriteErrorWrongArgCount()
				continue
			}

			s.DelAll()
		default:
			rc.WriteError("unknown command")
		}
	}
}
