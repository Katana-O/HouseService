package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/afocus/captcha"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"github.com/micro/go-micro"
	"image/png"
	getCaptcha "main/bj38web/web/proto/getCaptcha"
	go_micro_srv_user "main/bj38web/web/proto/user"
	"main/bj38web/web/utils"
	"net/http"
	"main/bj38web/web/model"
)

func GetImageCd(ctx * gin.Context) {
	uuid := ctx.Param("uuid")
	//consulReg := consul.NewRegistry()
	fmt.Println("uuid:", uuid)
	service := micro.NewService()
	microClient := getCaptcha.NewGetCaptchaService("getCaptcha", service.Client())
	fmt.Println("microClient:", microClient)
	resp, err := microClient.Call(context.TODO(), &getCaptcha.Request{Uuid: uuid})
	if err != nil {
		fmt.Println("cannot find method, ", err)
		return
	}
	var img captcha.Image
	json.Unmarshal(resp.Img, &img)
	png.Encode(ctx.Writer, img) // output image.
}

func PostAvatar(ctx * gin.Context) {
	fileHeader, _ := ctx.FormFile("avatar")
	err := ctx.SaveUploadedFile(fileHeader, "test/" + fileHeader.Filename)
	fmt.Println(err)
	resp := make(map[string]interface{})
	resp["errno"] = utils.RECODE_OK
	resp["errmsg"] = utils.RecodeText(utils.RECODE_OK)
	temp := make(map[string]interface{})
	///home/jjj/Desktop/testGO114/src/bj38web/web/test
	fmt.Println("fileheader.filename:", fileHeader.Filename)
	temp["avatar_url"] = "http://127.0.0.1:8080/test/" + fileHeader.Filename
	resp["data"] = temp
	ctx.JSON(http.StatusOK, resp)
}

func GetUserInfo(ctx * gin.Context) {
	resp := make(map[string]interface{})
	defer ctx.JSON(http.StatusOK, resp)
	s := sessions.Default(ctx)
	userName := s.Get("userName")
	if userName == nil {
		resp["errno"] = utils.RECODE_SESSIONERR
		resp["errmsg"] = utils.RecodeText(utils.RECODE_SESSIONERR)
		ctx.JSON(http.StatusOK, resp)
		return
	}
	user, err := model.GetUserInfo(userName.(string))
	if err != nil {
		resp["errno"] = utils.RECODE_DBERR
		resp["errmsg"] = utils.RecodeText(utils.RECODE_DBERR)
		return
	}

	resp["errno"] = utils.RECODE_OK
	resp["errmsg"] = utils.RecodeText(utils.RECODE_OK)
	temp := make(map[string]interface{})
	temp["user_id"] = user.ID
	temp["name"] = user.Name
	temp["mobile"] = user.Mobile
	temp["real_name"] = user.Real_name
	temp["id_card"] = user.Id_card
	temp["avatar_url"] = user.Avatar_url
	resp["data"] = temp
}

// Update user name
func PutUserInfo(ctx * gin.Context) {
	// Get current Username
	s := sessions.Default(ctx)
	userName := s.Get("userName")

	// Get New Username
	var nameData struct {
		Name string `json:"name"`
	}
	ctx.Bind(&nameData)

	// update user name
	resp := make(map[string]interface{})
	defer ctx.JSON(http.StatusOK, resp)

	// Update userName in DB.
	err := model.UpdateUserName(nameData.Name, userName.(string))
	if err != nil {
		// fmt.Println()
		resp["errno"] = utils.RECODE_DBERR
		resp["errmsg"] = utils.RecodeText(utils.RECODE_DBERR)
		return
	}

	// Update session data
	s.Set("userName", nameData.Name)
	err = s.Save()
	if err != nil {
		resp["errno"] = utils.RECODE_SESSIONERR
		resp["errmsg"] = utils.RecodeText(utils.RECODE_SESSIONERR)
		return
	}
	resp["errno"] = utils.RECODE_OK
	resp["errmsg"] = utils.RecodeText(utils.RECODE_OK)
	resp["data"] = nameData
}

func DeleteSession(ctx * gin.Context) {
	resp := make(map[string]interface{})
	s := sessions.Default(ctx)
	s.Delete("userName")
	err := s.Save()
	if err != nil {
		fmt.Println("DeleteSession err:", err)
		resp["errno"] = utils.RECODE_IOERR
		resp["errmsg"] = utils.RecodeText(utils.RECODE_IOERR)
	} else {
		resp["errno"] = utils.RECODE_OK
		resp["errmsg"] = utils.RecodeText(utils.RECODE_OK)
	}
	ctx.JSON(http.StatusOK, resp)
}

func GetSession(ctx * gin.Context) {
	resp := make(map[string]interface{})
	s := sessions.Default(ctx)
	userName := s.Get("userName")
	if userName == nil {
		resp["errno"] = utils.RECODE_SESSIONERR
		resp["errmsg"] = utils.RecodeText(utils.RECODE_SESSIONERR)
	} else {
		resp["errno"] = utils.RECODE_OK
		resp["errmsg"] = utils.RecodeText(utils.RECODE_OK)
		var nameData struct {
			Name string `json:"name"`
		}
		nameData.Name = userName.(string)
		resp["data"] = nameData
	}
	ctx.JSON(http.StatusOK, resp)
}

func PostLogin(ctx * gin.Context) {
	fmt.Println("POSTLOGIN 1")
	var LoginData struct {
		Mobile string `json:"mobile"`
		PassWord string `json:"password"`
	}
	ctx.Bind(&LoginData)
	fmt.Println("logindata:", LoginData)
	resp := make(map[string]interface{})
	userName, err := model.Login(LoginData.Mobile, LoginData.PassWord)
	if err == nil {
		fmt.Println("PostLogin 2222")
		resp["errno"] = utils.RECODE_OK
		resp["errmsg"] = utils.RecodeText(utils.RECODE_OK) + " yes yes yes "
		s := sessions.Default(ctx)
		s.Set("userName", userName)
		s.Save()
	} else {
		fmt.Println("PostLogin 3333")
		resp["errno"] = utils.RECODE_LOGINERR
		resp["errmsg"] = utils.RecodeText(utils.RECODE_LOGINERR) + " no no no "
	}
	ctx.JSON(http.StatusOK, resp)

	// Get Data from Redis/MySQL.

}

func GetArea(ctx * gin.Context) {


	// 1.
	// try get data from redis first.
	var areas []model.Area
	conn := model.RedisPool.Get()
	areaData, _ := redis.Bytes(conn.Do("get", "areaData"))
	if len(areaData) == 0 {
		// 2.
		// get data from mysql first
		fmt.Println("get Area len == 0")
		model.GlobalConn.Find(&areas)
		areaBuf, _ := json.Marshal(areas)
		conn.Do("set", "areaData", areaBuf)
	} else {
		fmt.Println("get Area len != 0")
		json.Unmarshal(areaData, &areas)
	}


	// write data to redis.
	resp := make(map[string]interface{})
	resp["errno"] = "0"
	resp["errmsg"] = utils.RecodeText(utils.RECODE_OK)
	resp["data"] = areas
	ctx.JSON(http.StatusOK, resp)
}

// Send register info
func PostRet(ctx * gin.Context) {
	fmt.Println("Post Ret!")
	var regData struct {
		Mobile   string `json:"mobile"`
		PassWord string `json:"password"`
		SmsCode  string `json:"sms_code"`
	}
	ctx.Bind(&regData)
	microService := utils.InitMicro()
	microClient := go_micro_srv_user.NewUserService("go.micro.srv.user", microService.Client())
	// rpc
	resp, err := microClient.Register(context.TODO(), &go_micro_srv_user.RegReq{
		Mobile:   regData.Mobile,
		Password: regData.PassWord,
		SmsCode:  regData.SmsCode,
	})
	if err != nil {
		fmt.Println("err:", err)
	}
	ctx.JSON(http.StatusOK, resp)
	fmt.Println("PostRet regData:", regData)
}
