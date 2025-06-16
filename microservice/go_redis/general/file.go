package general

import (
	"os"
	"time"
)

type FileInfo struct {
	FileName string
	IsDir    bool
	ModTime  time.Time
}

func ReadDir(pathname string) ([]FileInfo, error) {

	f, err := os.Open(pathname)
	if err != nil {
		return nil, err
	}
	files, err := f.Readdir(0)

	if err != nil {
		return nil, err
	}

	var flist []FileInfo = []FileInfo{}
	for _, v := range files {
		var info FileInfo
		info.FileName = v.Name()
		info.IsDir = v.IsDir()
		info.ModTime = v.ModTime()

		flist = append(flist, info)
	}

	return flist, nil
}
