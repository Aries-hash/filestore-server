package client

import (
	"context"
	"encoding/json"

	"github.com/mitchellh/mapstructure"

	"github.com/micro/go-micro"

	"filestore-server/service/dbproxy/orm"
	dbProto "filestore-server/service/dbproxy/proto"
)

// FileMeta : 文件元信息结构
type FileMeta struct {
	FileSha1 string
	FileName string
	FileSize int64
	Location string
	UploadAt string
}

var (
	dbCli dbProto.DBProxyService
)

func Init(service micro.Service) {
	// 初始化一个dbproxy服务的客户端
	dbCli = dbProto.NewDBProxyService("go.micro.service.dbproxy", service.Client())
}

func TableFileToFileMeta(tfile orm.TableFile) FileMeta {
	return FileMeta{
		FileSha1: tfile.FileHash,
		FileName: tfile.FileName.String,
		FileSize: tfile.FileSize.Int64,
		Location: tfile.FileAddr.String,
	}
}

// execAction : 向dbproxy请求执行action
func execAction(funcName string, paramJson []byte) (*dbProto.RespExec, error) {
	return dbCli.ExecuteAction(context.TODO(), &dbProto.ReqExec{
		Action: []*dbProto.SingleAction{
			&dbProto.SingleAction{
				Name:   funcName,
				Params: paramJson,
			},
		},
	})
}

// parseBody : 转换rpc返回的结果
func parseBody(resp *dbProto.RespExec) *orm.ExecResult {
	if resp == nil || resp.Data == nil {
		return nil
	}
	resList := []orm.ExecResult{}
	_ = json.Unmarshal(resp.Data, &resList)
	// TODO:
	if len(resList) > 0 {
		return &resList[0]
	}
	return nil
}

func ToTableUser(src interface{}) orm.TableUser {
	user := orm.TableUser{}
	mapstructure.Decode(src, &user)
	return user
}

func ToTableFile(src interface{}) orm.TableFile {
	file := orm.TableFile{}
	mapstructure.Decode(src, &file)
	return file
}

func ToTableFiles(src interface{}) []orm.TableFile {
	file := []orm.TableFile{}
	mapstructure.Decode(src, &file)
	return file
}

func ToTableUserFile(src interface{}) orm.TableUserFile {
	ufile := orm.TableUserFile{}
	mapstructure.Decode(src, &ufile)
	return ufile
}

func ToTableUserFiles(src interface{}) []orm.TableUserFile {
	ufile := []orm.TableUserFile{}
	mapstructure.Decode(src, &ufile)
	return ufile
}

func GetFileMeta(filehash string) (*orm.ExecResult, error) {
	uInfo, _ := json.Marshal([]interface{}{filehash})
	res, err := execAction("/file/GetFileMeta", uInfo)
	return parseBody(res), err
}

func GetFileMetaList(limitCnt int) (*orm.ExecResult, error) {
	uInfo, _ := json.Marshal([]interface{}{limitCnt})
	res, err := execAction("/file/GetFileMetaList", uInfo)
	return parseBody(res), err
}

// OnFileUploadFinished : 新增/更新文件元信息到mysql中
func OnFileUploadFinished(fmeta FileMeta) (*orm.ExecResult, error) {
	uInfo, _ := json.Marshal([]interface{}{fmeta.FileSha1, fmeta.FileName, fmeta.FileSize, fmeta.Location})
	res, err := execAction("/file/OnFileUploadFinished", uInfo)
	return parseBody(res), err
}

func UpdateFileLocation(filehash, location string) (*orm.ExecResult, error) {
	uInfo, _ := json.Marshal([]interface{}{filehash, location})
	res, err := execAction("/file/UpdateFileLocation", uInfo)
	return parseBody(res), err
}

func UserSignup(username, encPasswd string) (*orm.ExecResult, error) {
	uInfo, _ := json.Marshal([]interface{}{username, encPasswd})
	res, err := execAction("/user/UserSignup", uInfo)
	return parseBody(res), err
}

func UserSignin(username, encPasswd string) (*orm.ExecResult, error) {
	uInfo, _ := json.Marshal([]interface{}{username, encPasswd})
	res, err := execAction("/user/UserSignin", uInfo)
	return parseBody(res), err
}

func GetUserInfo(username string) (*orm.ExecResult, error) {
	uInfo, _ := json.Marshal([]interface{}{username})
	res, err := execAction("/user/GetUserInfo", uInfo)
	return parseBody(res), err
}

func UserExist(username string) (*orm.ExecResult, error) {
	uInfo, _ := json.Marshal([]interface{}{username})
	res, err := execAction("/user/UserExist", uInfo)
	return parseBody(res), err
}

func UpdateToken(username, token string) (*orm.ExecResult, error) {
	uInfo, _ := json.Marshal([]interface{}{username, token})
	res, err := execAction("/user/UpdateToken", uInfo)
	return parseBody(res), err
}

func QueryUserFileMeta(username, filehash string) (*orm.ExecResult, error) {
	uInfo, _ := json.Marshal([]interface{}{username, filehash})
	res, err := execAction("/ufile/QueryUserFileMeta", uInfo)
	return parseBody(res), err
}

func QueryUserFileMetas(username string, limit int) (*orm.ExecResult, error) {
	uInfo, _ := json.Marshal([]interface{}{username, limit})
	res, err := execAction("/ufile/QueryUserFileMetas", uInfo)
	return parseBody(res), err
}

// OnUserFileUploadFinished : 新增/更新文件元信息到mysql中
func OnUserFileUploadFinished(username string, fmeta FileMeta) (*orm.ExecResult, error) {
	uInfo, _ := json.Marshal([]interface{}{username, fmeta.FileSha1,
		fmeta.FileName, fmeta.FileSize})
	res, err := execAction("/ufile/OnUserFileUploadFinished", uInfo)
	return parseBody(res), err
}

func RenameFileName(username, filehash, filename string) (*orm.ExecResult, error) {
	uInfo, _ := json.Marshal([]interface{}{username, filehash, filename})
	res, err := execAction("/ufile/RenameFileName", uInfo)
	return parseBody(res), err
}
