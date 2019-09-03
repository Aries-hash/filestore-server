package route

import (
	"filestore-server/assets"
	"filestore-server/service/apigw/handler"
	"net/http"
	"strings"

	"github.com/gin-gonic/contrib/static"
	assetfs "github.com/moxiaomomo/go-bindata-assetfs"

	"github.com/gin-gonic/gin"
)

type binaryFileSystem struct {
	fs http.FileSystem
}

func (b *binaryFileSystem) Open(name string) (http.File, error) {
	return b.fs.Open(name)
}

func (b *binaryFileSystem) Exists(prefix string, filepath string) bool {

	if p := strings.TrimPrefix(filepath, prefix); len(p) < len(filepath) {
		if _, err := b.fs.Open(p); err != nil {
			return false
		}
		return true
	}
	return false
}

func BinaryFileSystem(root string) *binaryFileSystem {
	fs := &assetfs.AssetFS{
		Asset:     assets.Asset,
		AssetDir:  assets.AssetDir,
		AssetInfo: assets.AssetInfo,
		Prefix:    root,
	}
	return &binaryFileSystem{
		fs,
	}
}

// Router : 网关api路由
func Router() *gin.Engine {
	router := gin.Default()

	//	router.Static("/static/", "./static")
	// 将静态文件打包到bin文件
	router.Use(static.Serve("/static/", BinaryFileSystem("static")))

	// 注册
	router.GET("/user/signup", handler.SignupHandler)
	router.POST("/user/signup", handler.DoSignupHandler)
	// 登录
	router.GET("/user/signin", handler.SigninHandler)
	router.POST("/user/signin", handler.DoSigninHandler)

	router.Use(handler.Authorize())

	// 用户查询
	router.POST("/user/info", handler.UserInfoHandler)

	// 用户文件查询
	router.POST("/file/query", handler.FileQueryHandler)
	// 用户文件修改(重命名)
	router.POST("/file/update", handler.FileMetaUpdateHandler)

	return router
}
