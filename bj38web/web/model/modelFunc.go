package model

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/gomodule/redigo/redis"
)

var RedisPool redis.Pool
// 创建函数, 初始化Redis连接池
func InitRedis()  {
	RedisPool = redis.Pool{
		MaxIdle:20,
		MaxActive:50,
		MaxConnLifetime:60 * 5,
		IdleTimeout:60,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", "127.0.0.1:8083")
		},
	}
}

// 校验图片验证码
func CheckImgCode(uuid, imgCode string) bool {
	conn := RedisPool.Get()
	defer conn.Close()

	// 查询 redis 数据
	code, err := redis.String(conn.Do("get", uuid))
	if err != nil {
		fmt.Println("查询错误 err:", err)
		return false
	}

	// 返回校验结果
	return code == imgCode
}

// 处理登录业务,根据手机号/密码 获取用户名
func Login(mobile string, pwd string) (string, error) {
	fmt.Println("Login 111")
	var user User
	// 对参数 pwd 做md5 hash
	m5 := md5.New()
	m5.Write([]byte(pwd))
	pwd_hash := hex.EncodeToString(m5.Sum(nil))
	// GlobalConn > Gorm > MySql
	err := GlobalConn.Where("mobile = ?", mobile).Select("name").
		Where("password_hash = ?", pwd_hash).Find(&user).Error
	if err != nil {
		fmt.Println("Login err:", err)
	}

	return user.Name, err
}

func GetUserInfo(userName string) (User, error) {
	// SQL: select * from user where name = userName;
	var user User
	err := GlobalConn.Where("name = ?", userName).First(&user).Error
	return user, err
}

func UpdateUserName(newName string, oldName string) error {
	// update user set name = <newname> where name = <oldname>
	return GlobalConn.Model(new(User)).Where("name = ?", oldName).Update("name", newName).Error
}








