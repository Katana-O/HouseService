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
	"github.com/tedcy/fdfs_client"
	"image/png"
	"main/bj38web/web/model"
	getCaptcha "main/bj38web/web/proto/getCaptcha"
	go_micro_srv_house "main/bj38web/web/proto/house"
	go_micro_srv_user "main/bj38web/web/proto/user"
	go_micro_srv_userOrder "main/bj38web/web/proto/userOrder"
	"main/bj38web/web/utils"
	"net/http"
	"os"
	"path"

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
	fileHeader, err := ctx.FormFile("avatar")
	if err != nil {
		fmt.Println("PostAvatar 1 err:", err)
		return
	}
	fdfsClient, err := fdfs_client.NewClientWithConfig("../service/user/conf/fdfs.conf")
	if err != nil {
		fmt.Println("PostAvatar 2 err:", err)
		return
	}
	f, err := fileHeader.Open()
	if err != nil {
		fmt.Println("PostAvatar 3 err:", err)
		return
	}
	fileBuf := make([]byte, fileHeader.Size)
	f.Read(fileBuf)
	fileExtName := path.Ext(fileHeader.Filename)[1:]
	dir, _ := os.Getwd()
	fmt.Println("fileHeader.Size:", fileHeader.Size, " extName:", fileExtName, " dir:", dir)
	fdfsFile, err := os.Open("../service/user/conf/fdfs.conf")
	if err != nil {
		fmt.Println("PostAvatar 3.5 err:", err)
		return
	}
	fmt.Println("fdfsFile:", fdfsFile)
	remoteID, err := fdfsClient.UploadByBuffer(fileBuf, fileExtName)

	userName := sessions.Default(ctx).Get("userName")

	model.UpdateAvatar(userName.(string), remoteID)

	if err != nil {
		fmt.Println("PostAvatar 4 err:", err)
		return
	}
	fmt.Println("remoteID:", remoteID)

	resp := make(map[string]interface{})
	resp["errno"] = utils.RECODE_OK
	resp["errmsg"] = utils.RecodeText(utils.RECODE_OK)

	temp := make(map[string]interface{})
	temp["avatar_url"] = "http://192.168.1.161:8888/" + remoteID
	resp["data"] = temp
	ctx.JSON(http.StatusOK, resp)
}

type AuthStu struct {
	IdCard   string `json:"id_card"`
	RealName string `json:"real_name"`
}

func PostUserAuth(ctx * gin.Context) {
	//????????????
	var auth AuthStu
	err := ctx.Bind(&auth)
	//????????????
	if err != nil {
		fmt.Println("??????????????????", err)
		return
	}

	session := sessions.Default(ctx)
	userName := session.Get("userName")

	//???????????? ?????????
	microClient := go_micro_srv_user.NewUserService("go.micro.srv.user", utils.GetMicroClient())

	//??????????????????
	fmt.Println("controller PostUserAuth 1")
	resp, err := microClient.AuthUpdate(context.TODO(), &go_micro_srv_user.AuthReq{
		UserName: userName.(string),
		RealName: auth.RealName,
		IdCard:   auth.IdCard,
	})
	if err != nil {
		fmt.Println("Call authUpdate fail, err:", err)
		return
	}

	//????????????
	ctx.JSON(http.StatusOK, resp)
	fmt.Println("controller PostUserAuth 2")
}

func GetUserHouses(ctx * gin.Context) {
	fmt.Println("GetUserHouses 1")
	userName := sessions.Default(ctx).Get("userName")
	microClient := go_micro_srv_house.NewHouseService("go.micro.srv.house", utils.GetMicroClient())
	resp, _ := microClient.GetHouseInfo(
		context.TODO(),
		&go_micro_srv_house.GetReq{
			UserName: userName.(string),
		},
	)
	//fmt.Println("GetUserHouses resp:", resp)
	ctx.JSON(http.StatusOK, resp)
	fmt.Println("GetUserHouses 2")
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
	temp["avatar_url"] = "http://192.168.1.161:8888/"+ user.Avatar_url
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
		//fmt.Println("get Area len == 0")
		model.GlobalConn.Find(&areas)
		areaBuf, _ := json.Marshal(areas)
		conn.Do("set", "areaData", areaBuf)
	} else {
		//fmt.Println("get Area len != 0")
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

type HouseStu struct {
	Acreage   string   `json:"acreage"`
	Address   string   `json:"address"`
	AreaId    string   `json:"area_id"`
	Beds      string   `json:"beds"`
	Capacity  string   `json:"capacity"`
	Deposit   string   `json:"deposit"`
	Facility  []string `json:"facility"`
	MaxDays   string   `json:"max_days"`
	MinDays   string   `json:"min_days"`
	Price     string   `json:"price"`
	RoomCount string   `json:"room_count"`
	Title     string   `json:"title"`
	Unit      string   `json:"unit"`
}

func PostHousesImage(ctx * gin.Context) {
	fmt.Println("PostHousesImage 1")
	//????????????
	houseId := ctx.Param("id")

	fileHeader, err := ctx.FormFile("house_image")

	//????????????
	if houseId == "" || err != nil {
		fmt.Println("?????????????????????", err, " houseId:", houseId)
		return
	}

	//???????????? ??????,??????,????????????  fastdfs
	if fileHeader.Size > 50000000 {
		fmt.Println("????????????,???????????????")
		return
	}

	fileExt := path.Ext(fileHeader.Filename)
	if fileExt != ".png" && fileExt != ".jpg" {
		fmt.Println("??????????????????,???????????????")
		return
	}

	//????????????????????????
	file, _ := fileHeader.Open()
	buf := make([]byte, fileHeader.Size)
	file.Read(buf)

	//????????????  ???????????????
	microClient := go_micro_srv_house.NewHouseService("go.micro.srv.house", utils.GetMicroClient())

	//??????????????????
	resp, _ := microClient.UploadHouseImg(context.TODO(), &go_micro_srv_house.ImgReq{
		HouseId: houseId,
		ImgData: buf,
		FileExt: fileExt,
	})
	fmt.Println("PostHousesImage 2")
	//????????????
	ctx.JSON(http.StatusOK, resp)
}

func GetHouseInfo(ctx *gin.Context) {
	//????????????
	houseId := ctx.Param("id")
	//????????????
	if houseId == "" {
		fmt.Println("??????????????????")
		return
	}
	userName := sessions.Default(ctx).Get("userName")
	//????????????
	microClient := go_micro_srv_house.NewHouseService("go.micro.srv.house", utils.GetMicroClient())
	//??????????????????
	resp, _ := microClient.GetHouseDetail(context.TODO(), &go_micro_srv_house.DetailReq{
		HouseId:  houseId,
		UserName: userName.(string),
	})

	//????????????
	ctx.JSON(http.StatusOK, resp)
}

func GetIndex(ctx *gin.Context) {
	//????????????
	microClient := go_micro_srv_house.NewHouseService("go.micro.srv.house", utils.GetMicroClient())

	//??????????????????
	resp, _ := microClient.GetIndexHouse(context.TODO(), &go_micro_srv_house.IndexReq{})

	ctx.JSON(http.StatusOK, resp)
}

//????????????
func GetHouses(ctx *gin.Context) {
	//????????????
	//areaId
	aid := ctx.Query("aid")
	//start day
	sd := ctx.Query("sd")
	//end day
	ed := ctx.Query("ed")
	//????????????
	sk := ctx.Query("sk")
	//page  ?????????
	//ctx.Query("p")
	//????????????
	if aid == "" || sd == "" || ed == "" || sk == "" {
		fmt.Println("?????????????????????")
		return
	}

	//????????????   ?????????  ?????????????????????????????????,????????????time.Parse()  ??????????????????????????????,???????????????????????????,??????????????????
	/*sdTime ,_:=time.Parse("2006-01-02 15:04:05",sd+" 00:00:00")
	edTime,_ := time.Parse("2006-01-02",ed)*/

	/*sdTime,_ :=time.Parse("2006-01-02",sd)
	edTime,_ := time.Parse("2006-01-02",ed)
	d := edTime.Sub(sdTime)
	fmt.Println(d.Hours())*/

	microClient := go_micro_srv_house.NewHouseService("go.micro.srv.house", utils.GetMicroClient())
	//??????????????????
	resp, _ := microClient.SearchHouse(context.TODO(), &go_micro_srv_house.SearchReq{
		Aid: aid,
		Sd:  sd,
		Ed:  ed,
		Sk:  sk,
	})

	//????????????
	ctx.JSON(http.StatusOK, resp)

}

func PostHouses(ctx * gin.Context) {
	fmt.Println("PostHouses 1")
	//????????????   	bind?????????????????????????????????
	var house HouseStu
	err := ctx.Bind(&house)

	//????????????
	if err != nil {
		fmt.Println("??????????????????", err)
		return
	}

	//???????????????
	userName := sessions.Default(ctx).Get("userName")

	//????????????  ???????????????
	microClient := go_micro_srv_house.NewHouseService("go.micro.srv.house", utils.GetMicroClient())

	//??????????????????
	resp, _ := microClient.PubHouse(context.TODO(), &go_micro_srv_house.Request{
		Acreage:   house.Acreage,
		Address:   house.Address,
		AreaId:    house.AreaId,
		Beds:      house.Beds,
		Capacity:  house.Capacity,
		Deposit:   house.Deposit,
		Facility:  house.Facility,
		MaxDays:   house.MaxDays,
		MinDays:   house.MinDays,
		Price:     house.Price,
		RoomCount: house.RoomCount,
		Title:     house.Title,
		Unit:      house.Unit,
		UserName:  userName.(string),
	})

	//????????????
	ctx.JSON(http.StatusOK, resp)
	fmt.Println("PostHouses 2")
}

type OrderStu struct {
	EndDate   string `json:"end_date"`
	HouseId   string `json:"house_id"`
	StartDate string `json:"start_date"`
}


//?????????
func PostOrders(ctx *gin.Context) {
	//????????????
	var order OrderStu
	err := ctx.Bind(&order)

	//????????????
	if err != nil {
		fmt.Println("??????????????????", err)
		return
	}
	//???????????????
	userName := sessions.Default(ctx).Get("userName")

	//????????????  ?????????????????????
	microClient := go_micro_srv_userOrder.NewUserOrderService("go.micro.srv.userOrder", utils.GetMicroClient())
	//????????????
	resp, _ := microClient.CreateOrder(context.TODO(), &go_micro_srv_userOrder.Request{
		StartDate: order.StartDate,
		EndDate:   order.EndDate,
		HouseId:   order.HouseId,
		UserName:  userName.(string),
	})

	//????????????
	ctx.JSON(http.StatusOK, resp)
}

//??????????????????
func GetUserOrder(ctx *gin.Context) {
	//??????get????????????
	role := ctx.Query("role")
	//????????????
	if role == "" {
		fmt.Println("??????????????????")
		return
	}

	//????????????  ?????????
	microClient := go_micro_srv_userOrder.NewUserOrderService("go.micro.srv.userOrder",utils.GetMicroClient())
	//??????????????????
	resp,_ :=microClient.GetOrderInfo(context.TODO(),&go_micro_srv_userOrder.GetReq{
		Role:role,
		UserName:sessions.Default(ctx).Get("userName").(string),
	})

	//????????????
	ctx.JSON(http.StatusOK,resp)
}

//??????????????????
func PutOrders(ctx*gin.Context){
	//????????????
	id := ctx.Param("id")
	var statusStu StatusStu
	err := ctx.Bind(&statusStu)

	//????????????
	if err != nil || id == "" {
		fmt.Println("??????????????????",err)
		return
	}

	//????????????   ??????????????????
	microClient := go_micro_srv_userOrder.NewUserOrderService("go.micro.srv.userOrder",utils.GetMicroClient())
	//????????????????????????
	resp,_ := microClient.UpdateStatus(context.TODO(),&go_micro_srv_userOrder.UpdateReq{
		Action:statusStu.Action,
		Reason:statusStu.Reason,
		Id:id,
	})

	//????????????
	ctx.JSON(http.StatusOK,resp)
}

type StatusStu struct {
	Action string `json:"action"`
	Reason string `json:"reason"`
}
