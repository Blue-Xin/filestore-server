package meta

// FileMeta：文件元信息结构
type FileMeta struct {
	FileSha1 string
	FileName string
	FileSize int64
	Location string
	UploadAt string
}

var fileMetas map[string]FileMeta

func init() {
	fileMetas = make(map[string]FileMeta)
}

// 增加或更新元信息
func UpdateFileMeta(meta FileMeta) {
	fileMetas[meta.FileSha1] = meta
}

// 通过sha1获取元文件信息
func GetFileMeta(fileSha1 string) FileMeta {
	meta := fileMetas[fileSha1]
	return meta
}
func RemoveFileMeta(fileHash string) {
	delete(fileMetas, fileHash)
}
