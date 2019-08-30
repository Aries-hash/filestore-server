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
	//file
	http.HandleFunc("/file/upload",handler.HTTPInterceptor(handler.UploadHandler))
	http.HandleFunc("/file/upload/suc",handler.HTTPInterceptor(handler.UploadSucHandler))
	http.HandleFunc("/file/meta",handler.HTTPInterceptor(handler.GetFileMetaHandler))
	http.HandleFunc("/file/download",handler.HTTPInterceptor(handler.DownloadHandler))
	http.HandleFunc("/file/update",handler.HTTPInterceptor(handler.FileMetaUpdateHandler))
	http.HandleFunc("/file/delete",handler.HTTPInterceptor(handler.FileDeleteHandler))
	http.HandleFunc("/file/query",handler.HTTPInterceptor(handler.FileQueryHandler))
	http.HandleFunc("/file/fastupload",handler.HTTPInterceptor(handler.TryFastUploadHandler))
    //user 
	http.HandleFunc("/user/signup",handler.SignupHandler)
    http.HandleFunc("/user/signin",handler.SignInHandler)
    http.HandleFunc("/user/info",handler.HTTPInterceptor(handler.UserInfoHandler))
    //监听端口
	err:=http.ListenAndServe(":8080",nil)
	if err != nil {
	   fmt.Printf("Failed to start server,err:%s",err.Error)
	}
}
