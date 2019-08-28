package meta
//FileMeta:文件元信息结构
type FileMeta struct {
	FileSha1 string
	FileName string
	FileSize int64
	Location string
	UploadAt string
}


var fileMetas map[string]FileMeta
func init(){
	fileMetas=make(map[string]FileMeta)
}
//UpdateFileMeta:新增/更新文件元信息
func UpdateFileMeta(fmeta FileMeta){
	fileMetas[fmeta.FileSha1]=fmeta
}
//GetFileMeta:通过Sha1获取文件的元信息对象
func GetFileMeta(fileSha1 string)FileMeta{
	return fileMetas[fileSha1]
}

// RemoveFileMeta : 删除元信息
func RemoveFileMeta(fileSha1 string) {
	delete(fileMetas, fileSha1)
}

