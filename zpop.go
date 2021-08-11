package main

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
)

func main() {
	conn, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// Add test data using a pipeline.
	key := "zset1"
	for i, member := range []string{"red", "blue", "green"} {
		conn.Send("ZADD", key, i, member)
	}
	res, _ := conn.Do("")
	fmt.Printf("1: %v\n", res)

	// Pop using WATCH/MULTI/EXEC
	v, err := zpop(conn, key)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(v)

}

func zpop(conn2 redis.Conn, key string) (result string, err error) {
	defer func() {
		if err != nil {
			conn2.Do("DISCARD")
		}
	}()

	for {
		if _, err := conn2.Do("WATCH", key); err != nil {
			return "", err
		}

		members, err := redis.Strings(conn2.Do("ZRANGE", key, 0, 0)) // 0 0 即只取第一个
		if err != nil {
			return "", err
		}

		fmt.Printf("2: %v\n", members)
		if len(members) != 1 {
			return "", redis.ErrNil
		}

		conn2.Send("MULTI")
		conn2.Send("ZREM", key, members[0])
		queued, err := conn2.Do("EXEC")
		if err != nil {
			return "", err
		}

		fmt.Printf("3: %v\n", queued)
		if queued != nil {
			result = members[0]
			break
		}
	}
	return result, nil
}
