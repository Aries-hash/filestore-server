package main
import(
	"filestore-server/handler"
	"fmt"
	"net/http"
)

func main()  {
	//处理静态资源映射
	http.Handle("/static/",http.StripPrefix("/static",http.FileServer(http.Dir("./static"))))
	//url映射
	//文件接口相关
	http.HandleFunc("/file/upload",handler.HTTPInterceptor(handler.UploadHandler))
	http.HandleFunc("/file/upload/suc",handler.HTTPInterceptor(handler.UploadSucHandler))
	http.HandleFunc("/file/meta",handler.HTTPInterceptor(handler.GetFileMetaHandler))
	http.HandleFunc("/file/download",handler.HTTPInterceptor(handler.DownloadHandler))
	http.HandleFunc("/file/update",handler.HTTPInterceptor(handler.FileMetaUpdateHandler))
	http.HandleFunc("/file/delete",handler.HTTPInterceptor(handler.FileDeleteHandler))
	http.HandleFunc("/file/query",handler.HTTPInterceptor(handler.FileQueryHandler))
	http.HandleFunc("/file/fastupload",handler.HTTPInterceptor(handler.TryFastUploadHandler))
	
	//分块上传接口
	// 初始化分块信息
    http.HandleFunc("/file/mpupload/init",handler.HTTPInterceptor(handler.InitialMultipartUploadHandler))
	// 上传分块
    http.HandleFunc("/file/mpupload/uppart",handler.HTTPInterceptor(handler.UploadPartHandler))
	// 通知分块上传完成
	http.HandleFunc("/file/mpupload/complete",handler.HTTPInterceptor(handler.CompleteUploadHandler))
	// 取消上传分块

	// 查看分块上传的整体状态

    //用户相关接口 
	http.HandleFunc("/user/signup",handler.SignupHandler)
    http.HandleFunc("/user/signin",handler.SignInHandler)
    http.HandleFunc("/user/info",handler.HTTPInterceptor(handler.UserInfoHandler))
    //监听端口
	err:=http.ListenAndServe(":8080",nil)
	if err != nil {
	   fmt.Printf("Failed to start server,err:%s",err.Error)
	}
}
