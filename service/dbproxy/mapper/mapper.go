package mapper

import (
	"errors"
	"reflect"

	"filestore-server/service/dbproxy/orm"
)

var funcs = map[string]interface{}{
	"/file/OnFileUploadFinished": orm.OnFileUploadFinished,
	"/file/GetFileMeta":          orm.GetFileMeta,
	"/file/GetFileMetaList":      orm.GetFileMetaList,
	"/file/UpdateFileLocation":   orm.UpdateFileLocation,

	"/user/UserSignup":  orm.UserSignup,
	"/user/UserSignin":  orm.UserSignin,
	"/user/UpdateToken": orm.UpdateToken,
	"/user/GetUserInfo": orm.GetUserInfo,
	"/user/UserExist":   orm.UserExist,

	"/ufile/OnUserFileUploadFinished": orm.OnUserFileUploadFinished,
	"/ufile/QueryUserFileMetas":       orm.QueryUserFileMetas,
	"/ufile/DeleteUserFile":           orm.DeleteUserFile,
	"/ufile/RenameFileName":           orm.RenameFileName,
	"/ufile/QueryUserFileMeta":        orm.QueryUserFileMeta,
}

func FuncCall(name string, params ...interface{}) (result []reflect.Value, err error) {
	if _, ok := funcs[name]; !ok {
		err = errors.New("函数名不存在.")
		return
	}

	// 通过反射可以动态调用对象的导出方法
	f := reflect.ValueOf(funcs[name])
	if len(params) != f.Type().NumIn() {
		err = errors.New("传入参数数量与被调用方法要求的数量不一致.")
		return
	}

	// 构造一个 Value的slice, 用作Call()方法的传入参数
	in := make([]reflect.Value, len(params))
	for k, param := range params {
		in[k] = reflect.ValueOf(param)
	}

	// 执行方法f, 并将方法结果赋值给result
	result = f.Call(in)
	return
}
