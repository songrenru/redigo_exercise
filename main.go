package main

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
)

var conn redis.Conn

func init() {
	c, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		panic(err)
	}
	conn = c
}

func main() {
	defer conn.Close()

	//Do()
	//Pipelining()
	//Transcation()
	PubSub()
}

func Do() {
	rep, err := conn.Do("SET", "K1", "v1")
	if err != nil {
		panic(err)
	}

	fmt.Println(rep)
}

func Pipelining() {
	conn.Send("SET", "K2", "V2")
	conn.Send("GET", "K2")
	conn.Flush()
	res, _ := conn.Receive() // reply from SET
	fmt.Printf("%t\n", res)
	res2, _ := conn.Receive() // reply from GET
	fmt.Printf("%t, %s\n", res2, res2.([]byte))
}

func Transcation() {
	conn.Send("MULTI")

	conn.Send("SET", "K3", "V3")
	conn.Send("SET", "K4", "V4")
	r, err := conn.Do("EXEC")
	fmt.Println(r, err) // output: [OK OK] <nil>

	//conn.Send("GET", "K3")
	//conn.Send("GET", "K4")
	//r, err := conn.Do("")
	//fmt.Println(r, err) // output: [OK QUEUED QUEUED] <nil>
}

func PubSub() {
	// c.Send()方式订阅的，需要解析返回结果的结构，故不推荐

	psc := redis.PubSubConn{Conn: conn}
	psc.Subscribe("example")
	for {
		switch v := psc.Receive().(type) {
		case redis.Message:
			fmt.Printf("%s: message: %s\n", v.Channel, v.Data)
		case redis.Subscription:
			fmt.Printf("%s: %s %d\n", v.Channel, v.Kind, v.Count)
		case error:
			panic(v)
		}
	}
}
