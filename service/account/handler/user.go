package handler

import (
	"context"
	"fmt"
	"time"

	"filestore-server/common"
	"filestore-server/config"
	cfg "filestore-server/config"
	proto "filestore-server/service/account/proto"
	dbcli "filestore-server/service/dbproxy/client"
	"filestore-server/util"
)

// User : 用于实现UserServiceHandler接口的对象
type User struct{}

// GenToken : 生成token
func GenToken(username string) string {
	// 40位字符:md5(username+timestamp+token_salt)+timestamp[:8]
	ts := fmt.Sprintf("%x", time.Now().Unix())
	tokenPrefix := util.MD5([]byte(username + ts + "_tokensalt"))
	return tokenPrefix + ts[:8]
}

// Signup : 处理用户注册请求
func (u *User) Signup(ctx context.Context, req *proto.ReqSignup, res *proto.RespSignup) error {
	username := req.Username
	passwd := req.Password

	// 参数简单校验
	if len(username) < 3 || len(passwd) < 5 {
		res.Code = common.StatusParamInvalid
		res.Message = "注册参数无效"
		return nil
	}

	// 对密码进行加盐及取Sha1值加密
	encPasswd := util.Sha1([]byte(passwd + cfg.PasswordSalt))
	// 将用户信息注册到用户表中
	dbResp, err := dbcli.UserSignup(username, encPasswd)
	if err == nil && dbResp.Suc {
		res.Code = common.StatusOK
		res.Message = "注册成功"
	} else {
		res.Code = common.StatusRegisterFailed
		res.Message = "注册失败"
	}
	return nil
}

// Signin : 处理用户登录请求
func (u *User) Signin(ctx context.Context, req *proto.ReqSignin, res *proto.RespSignin) error {
	username := req.Username
	password := req.Password

	encPasswd := util.Sha1([]byte(password + config.PasswordSalt))

	// 1. 校验用户名及密码
	dbResp, err := dbcli.UserSignin(username, encPasswd)
	if err != nil || !dbResp.Suc {
		res.Code = common.StatusLoginFailed
		return nil
	}

	// 2. 生成访问凭证(token)
	token := GenToken(username)
	upRes, err := dbcli.UpdateToken(username, token)
	if err != nil || !upRes.Suc {
		res.Code = common.StatusServerError
		return nil
	}

	// 3. 登录成功, 返回token
	res.Code = common.StatusOK
	res.Token = token
	return nil
}

// UserInfo ： 查询用户信息
func (u *User) UserInfo(ctx context.Context, req *proto.ReqUserInfo, res *proto.RespUserInfo) error {
	// 查询用户信息
	dbResp, err := dbcli.GetUserInfo(req.Username)
	if err != nil {
		res.Code = common.StatusServerError
		res.Message = "服务错误"
		return nil
	}
	// 查不到对应的用户信息
	if !dbResp.Suc {
		res.Code = common.StatusUserNotExists
		res.Message = "用户不存在"
		return nil
	}

	user := dbcli.ToTableUser(dbResp.Data)

	// 3. 组装并且响应用户数据
	res.Code = common.StatusOK
	res.Username = user.Username
	res.SignupAt = user.SignupAt
	res.LastActiveAt = user.LastActiveAt
	res.Status = int32(user.Status)
	// TODO: 需增加接口支持完善用户信息(email/phone等)
	res.Email = user.Email
	res.Phone = user.Phone
	return nil
}
