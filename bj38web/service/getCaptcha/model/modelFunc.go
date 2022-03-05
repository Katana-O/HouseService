package model

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
)

// Save image id to redis DB.
func SaveImgCode(imgCode string, uuid string) error {
	// 1. connect to DB
	conn, err := redis.Dial("tcp", "127.0.0.1:8083")
	if err != nil {
		fmt.Println("err:", err)
		return err
	}
	defer conn.Close()
	// 2. manipulate Redis
	_, err = conn.Do("setex", uuid, 60 * 5, imgCode)

	// 3. get Result
	// r, e := redis.String(reply, err)
	// fmt.Println(r, e)
	return err
}

