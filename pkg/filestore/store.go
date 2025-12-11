package filestore

import (
	"io"
)

type FileStore interface {
	// 上传静态资源池资源
	Save(path, mime string, content io.ReadSeeker) (string, error)
}
