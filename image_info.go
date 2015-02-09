package frameresize

import (
	"os"
	"path"
)

type ImageInfo struct {
	Path string
	Name string
}

func NewImageInfoFromFileInfo(path string, pi os.FileInfo) ImageInfo {
	return ImageInfo{
		Path: path,
		Name: pi.Name(),
	}
}
func (ii ImageInfo) FullPath() string {
	return path.Join(ii.Path, ii.Name)
}
