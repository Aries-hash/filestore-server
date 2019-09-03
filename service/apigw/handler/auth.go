package handler

import (
	"filestore-server/common"
	"filestore-server/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

// IsTokenValid : token是否有效
func IsTokenValid(token string) bool {
	if len(token) != 40 {
		return false
	}
	// TODO: 判断token的时效性，是否过期
	// TODO: 从数据库表tbl_user_token查询username对应的token信息
	// TODO: 对比两个token是否一致
	return true
}

// Authorize : http请求拦截器
func Authorize() gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.Request.FormValue("username")
		token := c.Request.FormValue("token")

		//验证登录token是否有效
		if len(username) < 3 || !IsTokenValid(token) {
			// w.WriteHeader(http.StatusForbidden)
			// token校验失败则跳转到登录页面
			c.Abort()
			resp := util.NewRespMsg(
				int(common.StatusTokenInvalid),
				"token无效",
				nil,
			)
			c.JSON(http.StatusOK, resp)
			return
		}
		c.Next()
	}
}
