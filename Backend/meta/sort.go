package meta

import (
	"sort"
	"time"
)

// ByTimestamp sorts fileMetas by upload timestamp in descending order
func ByTimestamp(fileMetas []FileMeta) {
	sort.Slice(fileMetas, func(i, j int) bool {
		t1, _ := time.Parse(time.RFC3339, fileMetas[i].UploadAt)
		t2, _ := time.Parse(time.RFC3339, fileMetas[j].UploadAt)
		return t1.After(t2) // descending order
	})
}

// ByFileSize sorts fileMetas by file size in descending order
func ByFileSize(fileMetas []FileMeta) {
	sort.Slice(fileMetas, func(i, j int) bool {
		return fileMetas[i].FileSize > fileMetas[j].FileSize // descending order
	})
}

// ByFileName sorts fileMetas by file name in ascending order
func ByFileName(fileMetas []FileMeta) {
	sort.Slice(fileMetas, func(i, j int) bool {
		return fileMetas[i].FileName < fileMetas[j].FileName // ascending order
	})
}
