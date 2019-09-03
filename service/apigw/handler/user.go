package handler

import (
	"context"
	"log"
	"net/http"

	// 加入k8s作为registry center
	_ "github.com/micro/go-plugins/registry/kubernetes"

	"github.com/gin-gonic/gin"
	ratelimit2 "github.com/juju/ratelimit"
	micro "github.com/micro/go-micro"

	"github.com/micro/go-plugins/wrapper/breaker/hystrix"
	"github.com/micro/go-plugins/wrapper/ratelimiter/ratelimit"

	cmn "filestore-server/common"
	cfg "filestore-server/config"
	userProto "filestore-server/service/account/proto"
	dlProto "filestore-server/service/download/proto"
	upProto "filestore-server/service/upload/proto"
	"filestore-server/util"
)

var (
	userCli userProto.UserService
	upCli   upProto.UploadService
	dlCli   dlProto.DownloadService
)

func init() {
	//配置请求容量及qps
	bRate := ratelimit2.NewBucketWithRate(100, 1000)
	service := micro.NewService(
		micro.Flags(cmn.CustomFlags...),
		micro.WrapClient(ratelimit.NewClientWrapper(bRate, false)), //加入限流功能, false为不等待(超限即返回请求失败)
		micro.WrapClient(hystrix.NewClientWrapper()),               // 加入熔断功能, 处理rpc调用失败的情况(cirucuit breaker)
	)
	// 初始化， 解析命令行参数等
	service.Init()

	cli := service.Client()
	// tracer, err := tracing.Init("apigw service", "<jaeger-agent-host>")
	// if err != nil {
	// 	log.Println(err.Error())
	// } else {
	// 	cli = client.NewClient(
	// 		client.Wrap(mopentracing.NewClientWrapper(tracer)),
	// 	)
	// }

	// 初始化一个account服务的客户端
	userCli = userProto.NewUserService("go.micro.service.user", cli)
	// 初始化一个upload服务的客户端
	upCli = upProto.NewUploadService("go.micro.service.upload", cli)
	// 初始化一个download服务的客户端
	dlCli = dlProto.NewDownloadService("go.micro.service.download", cli)
}

// SignupHandler : 响应注册页面
func SignupHandler(c *gin.Context) {
	c.Redirect(http.StatusFound, "/static/view/signup.html")
}

// DoSignupHandler : 处理注册post请求
func DoSignupHandler(c *gin.Context) {
	username := c.Request.FormValue("username")
	passwd := c.Request.FormValue("password")

	resp, err := userCli.Signup(context.TODO(), &userProto.ReqSignup{
		Username: username,
		Password: passwd,
	})

	if err != nil {
		log.Println(err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": resp.Code,
		"msg":  resp.Message,
	})
}

// SigninHandler : 响应登录页面
func SigninHandler(c *gin.Context) {
	c.Redirect(http.StatusFound, "/static/view/signin.html")
}

// DoSigninHandler : 处理登录post请求
func DoSigninHandler(c *gin.Context) {
	username := c.Request.FormValue("username")
	password := c.Request.FormValue("password")

	rpcResp, err := userCli.Signin(context.TODO(), &userProto.ReqSignin{
		Username: username,
		Password: password,
	})

	if err != nil {
		log.Println(err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}

	if rpcResp.Code != cmn.StatusOK {
		c.JSON(200, gin.H{
			"msg":  "登录失败",
			"code": rpcResp.Code,
		})
		return
	}

	// // 动态获取上传入口地址
	// upEntryResp, err := upCli.UploadEntry(context.TODO(), &upProto.ReqEntry{})
	// if err != nil {
	// 	log.Println(err.Error())
	// } else if upEntryResp.Code != cmn.StatusOK {
	// 	log.Println(upEntryResp.Message)
	// }

	// // 动态获取下载入口地址
	// dlEntryResp, err := dlCli.DownloadEntry(context.TODO(), &dlProto.ReqEntry{})
	// if err != nil {
	// 	log.Println(err.Error())
	// } else if dlEntryResp.Code != cmn.StatusOK {
	// 	log.Println(dlEntryResp.Message)
	// }

	// 登录成功，返回用户信息
	cliResp := util.RespMsg{
		Code: int(cmn.StatusOK),
		Msg:  "登录成功",
		Data: struct {
			Location      string
			Username      string
			Token         string
			UploadEntry   string
			DownloadEntry string
		}{
			Location: "/static/view/home.html",
			Username: username,
			Token:    rpcResp.Token,
			// UploadEntry:   upEntryResp.Entry,
			// DownloadEntry: dlEntryResp.Entry,
			UploadEntry:   cfg.UploadLBHost,
			DownloadEntry: cfg.DownloadLBHost,
		},
	}
	c.Data(http.StatusOK, "application/json", cliResp.JSONBytes())
}

// UserInfoHandler ： 查询用户信息
func UserInfoHandler(c *gin.Context) {
	// 1. 解析请求参数
	username := c.Request.FormValue("username")

	resp, err := userCli.UserInfo(context.TODO(), &userProto.ReqUserInfo{
		Username: username,
	})

	if err != nil {
		log.Println(err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}

	// 3. 组装并且响应用户数据
	cliResp := util.RespMsg{
		Code: 0,
		Msg:  "OK",
		Data: gin.H{
			"Username": username,
			"SignupAt": resp.SignupAt,
			// TODO: 完善其他字段信息
			"LastActive": resp.LastActiveAt,
		},
	}
	c.Data(http.StatusOK, "application/json", cliResp.JSONBytes())
}
