package main

import (
	"fmt"
	"image/jpeg"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/nfnt/resize"
)

type stringSlice []string
type imageInfo struct {
	Path string
	Name string
}

var wgProcess sync.WaitGroup
var FilesResized uint32

func NewImageInfoFromFileInfo(path string, pi os.FileInfo) imageInfo {
	return imageInfo{
		Path: path,
		Name: pi.Name(),
	}
}
func (ii imageInfo) FullPath() string {
	return path.Join(ii.Path, ii.Name)
}

var IMAGE_EXT = &stringSlice{".JPG", ".JPEG", ".TIF", ".TIFF", ".PNG", ".GIF", ".BMP"}

func (s stringSlice) Contains(test string) bool {
	for _, a := range s {
		if a == test {
			return true
		}
	}
	return false
}

func scanDir(srcPath string, chImages chan imageInfo) (err error) {
	entries, err := ioutil.ReadDir(srcPath)
	if err != nil {
		return
	}
	// We are waiting
	wgProcess.Add(1)

	// Process entries
	for _, r := range entries {
		upper_name := strings.ToUpper(r.Name())
		if r.IsDir() {
			// Make sure we don't process hidden folders
			if !strings.HasPrefix(upper_name, ".") {
				newPath := path.Join(srcPath, r.Name())
				go scanDir(newPath, chImages)
			}
		} else {
			if IMAGE_EXT.Contains(filepath.Ext(upper_name)) {
				chImages <- NewImageInfoFromFileInfo(srcPath, r)
			}
		}
	}

	// Clean up
	wgProcess.Done()
	return
}

func printFiles(theChan chan imageInfo) {
	for {
		file := <-theChan
		fmt.Printf("%s\n", file.FullPath())
	}
}

// Helper function that handles the stats and synchronization aspects of resize
func resizeFiles(theChan chan imageInfo, dstPath string, width, height uint) {
	for {
		// Read from the channel, add our lock to wgProcess
		file := <-theChan
		wgProcess.Add(1)

		// Do the actual work
		resizeImage(file, dstPath, width, height)

		// Decrement our resize operation, increment our processing count
		wgProcess.Done()
		atomic.AddUint32(&FilesResized, 1)
	}
}

func resizeImage(srcImage imageInfo, dstPath string, width, height uint) {
	dstImage := fmt.Sprintf("%s_%dx%d.jpg", srcImage.Name, width, height)
	dstImage = path.Join(dstPath, dstImage)
	fmt.Printf("Resizing %s\n", dstImage)

	// open "test.jpg"
	file, err := os.Open(srcImage.FullPath())
	if err != nil {
		log.Fatal(err)
	}

	// decode jpeg into image.Image
	img, err := jpeg.Decode(file)
	if err != nil {
		log.Fatal(err)
	}
	file.Close()

	// resize to width 1000 using Lanczos resampling
	// and preserve aspect ratio
	resized := resize.Resize(0, height, img, resize.NearestNeighbor)

	out, err := os.Create(dstImage)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	// write new image to file
	jpeg.Encode(out, resized, nil)
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	// testVar, _ := os.Lstat("/data/void/audioBook/child_lee/07_persuader/012.mp3")
	src_path := "/data/Dropbox/Photos/Photostream/2015/2015-02-WDW"
	// src_path += "/photopass"
	// src_path += "/PhotoPass_20150208_54d6b3a4a029b_8"
	dst_path := "/tmp/photoframe"
	fmt.Println("Will output to", dst_path)

	FilesResized = 0
	chSrcImages := make(chan imageInfo)

	// go printFiles(chSrcImages)

	// Resize on as many cores as we have
	for x := 0; x < runtime.GOMAXPROCS(0); x++ {
		go resizeFiles(chSrcImages, dst_path, 640, 480)
	}

	// Start filling the queue
	go scanDir(src_path, chSrcImages)

	// Sleep to give enough time for scanDir() to set a lock on wgProcess
	time.Sleep(50 * time.Millisecond)

	// Wait for all existing operations to finish
	wgProcess.Wait()
	fmt.Printf("All Done, resized %d files\n", FilesResized)
}
