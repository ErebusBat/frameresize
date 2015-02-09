package frameresize

import (
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

var IMAGE_EXT = &StringSlice{".JPG", ".JPEG", ".TIF", ".TIFF", ".PNG", ".GIF", ".BMP"}

type Photoframe struct {
	Width        uint
	Height       uint
	DestRoot     string
	FilesResized uint32
	chImages     chan ImageInfo
	// sourcePath   string
	wgProcess sync.WaitGroup
}

func NewPhotoframe(dest_root string, width, height uint) *Photoframe {
	app := &Photoframe{
		DestRoot: dest_root,
		Width:    width,
		Height:   height,
	}

	app.FilesResized = 0
	app.chImages = make(chan ImageInfo)

	return app
}

func (pf *Photoframe) Process(src_path string) {
	// pf.sourcePath = src_path
	// Resize on as many cores as we have
	for x := 0; x < runtime.GOMAXPROCS(0); x++ {
		go pf.resizeFiles()
	}

	// Start filling the queue
	go pf.scanDir(src_path)

	// Sleep to give enough time for scanDir() to set a lock on wgProcess
	time.Sleep(50 * time.Millisecond)

	// Wait for all existing operations to finish
	pf.wgProcess.Wait()
}

func (pf *Photoframe) scanDir(srcPath string) (err error) {
	entries, err := ioutil.ReadDir(srcPath)
	if err != nil {
		return
	}
	// We are waiting
	pf.wgProcess.Add(1)

	// Process entries
	for _, r := range entries {
		upper_name := strings.ToUpper(r.Name())
		if r.IsDir() {
			// Make sure we don't process hidden folders
			if !strings.HasPrefix(upper_name, ".") {
				newPath := path.Join(srcPath, r.Name())
				go pf.scanDir(newPath)
			}
		} else {
			if IMAGE_EXT.Contains(filepath.Ext(upper_name)) {
				pf.chImages <- NewImageInfoFromFileInfo(srcPath, r)
			}
		}
	}

	// Clean up
	pf.wgProcess.Done()
	return
}

func (pf *Photoframe) NewFileName(image ImageInfo) string {
	// Hash output serves as a randomizer and name conflict resolution
	hashOutput := true

	dstImage := fmt.Sprintf("%s_%dx%d.jpg", image.Name, pf.Width, pf.Height)
	if hashOutput == true {
		hash := sha1.New()
		hash.Write([]byte(image.Path + "/"))
		hash.Write([]byte(dstImage))
		dstImage = fmt.Sprintf("%x.jpg", hash.Sum(nil))
	}
	return path.Join(pf.DestRoot, dstImage)
}
