package main

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
)

func main () {
	// 1. connect to DB
	conn, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	defer conn.Close()
	// 2. manipulate DB
	reply, err := conn.Do("set", "hello", "world")

	// 3. get Result
	r, e := redis.String(reply, err)
	fmt.Println(r, e)
}
