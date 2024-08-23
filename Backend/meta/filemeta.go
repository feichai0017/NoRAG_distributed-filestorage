package meta

import (
	mydb "cloud_distributed_storage/database"
	"sync"
)

// FileMeta infomation struct
type FileMeta struct {
	FileSha1 string `json:"file_sha1"`
	FileName string `json:"file_name"`
	FileSize int64  `json:"file_size"`
	Location string `json:"location"`
	UploadAt string `json:"upload_at"`
}

var (
	mu        sync.RWMutex
	fileMetas map[string]FileMeta
)

func init() {
	fileMetas = make(map[string]FileMeta)
}

// UpdateFileMeta: increment/update meta info of file
func UpdateFileMeta(fmeta FileMeta) {
	mu.Lock()
	defer mu.Unlock()
	fileMetas[fmeta.FileSha1] = fmeta
}

// UpdateFileMetaDB:increment/update meta info of file to mysql
func UpdateFileMetaDB(fmeta FileMeta) bool {
	mu.Lock()
	defer mu.Unlock()
	return mydb.OnFileUploadFinished(fmeta.FileSha1, fmeta.FileName, fmeta.FileSize, fmeta.Location)
}

// GetFileMeta: find filemeta by sha1
func GetFileMeta(fileSha1 string) FileMeta {
	mu.RLock()
	defer mu.RUnlock()
	return fileMetas[fileSha1]
}

// GetFileMetaDB: get filemeta from mysql
func GetFileMetaDB(fileSha1 string) (FileMeta, error) {
	mu.RLock()
	defer mu.RUnlock()
	tfile, err := mydb.GetFileMeta(fileSha1)
	if err != nil {
		return FileMeta{}, err
	}
	fmeta := FileMeta{
		FileSha1: tfile.FileHash,
		FileName: tfile.FileName.String,
		FileSize: tfile.FileSize.Int64,
		Location: tfile.FileAddr.String,
	}
	return fmeta, nil
}

// GetLastFileMetas: returns the last 'count' file metas in order of upload time
func GetLastFileMetas(count int) []FileMeta {
	mu.RLock()
	defer mu.RUnlock()
	fMetaArray := make([]FileMeta, len(fileMetas))
	for _, v := range fileMetas {
		fMetaArray = append(fMetaArray, v)
	}

	ByTimestamp(fMetaArray)

	if count > len(fMetaArray) {
		return fMetaArray
	}

	return fMetaArray[:count]
}

// RemoveFileMeta: delete filemeta by filesha1
func RemoveFileMeta(fileSha1 string) {
	mu.Lock()
	defer mu.Unlock()
	delete(fileMetas, fileSha1)
}
