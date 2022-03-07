package handler

import (
	"context"
	_ "context"
	"fmt"
	"user/model"
	user "user/proto/user"
	"user/utils"

	_ "github.com/micro/go-micro/util/log"

	//user "user/proto/user"
)

type User struct{}

func (e *User) Register(ctx context.Context, req *user.RegReq, rsp *user.Response) error {
	// model.CheckSmsCode()
	fmt.Println("handler.user.go Register, mobile:", req.Mobile, " pw:", req.Password)
	err := model.RegisterUser(req.Mobile, req.Password)
	if err != nil {
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(utils.RECODE_DBERR)
	} else {
		rsp.Errno = utils.RECODE_OK
		rsp.Errmsg = utils.RecodeText(utils.RECODE_OK)
	}
	return nil
}

func (e *User) AuthUpdate(ctx context.Context, req *user.AuthReq, rsp *user.AuthResp) error {
	fmt.Println("AuthUpdate@handler 1")
	err := model.SaveRealName(req.UserName, req.RealName, req.IdCard)
	if err != nil {
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(utils.RECODE_DBERR)
		return nil
	}
	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(utils.RECODE_OK)
	fmt.Println("AuthUpdate@handler 2")
	return nil
}
