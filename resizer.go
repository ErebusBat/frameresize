package frameresize

import (
	"fmt"
	"image/jpeg"
	"log"
	"os"
	"sync/atomic"

	"github.com/nfnt/resize"
)

// Helper function that handles the stats and synchronization aspects of resize
func (pf *Photoframe) resizeFiles() {
	for {

		file := <-pf.chImages
		pf.wgProcess.Add(1)

		dstImage := pf.NewFileName(file)
		if fileExists(dstImage) {
			fmt.Printf("SKIPping %s => %s\n", file.Name, dstImage)
		} else {
			fmt.Printf("Resizing %s => %s\n", file.Name, dstImage)
			resizeImage(file.FullPath(), dstImage, pf.Width, pf.Height)
		}

		pf.wgProcess.Done()
		atomic.AddUint32(&pf.FilesResized, 1)
	}
}

func fileExists(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}
	return true
}

func resizeImage(srcImage string, dstImage string, width, height uint) {
	file, err := os.Open(srcImage)
	if err != nil {
		log.Fatal(err)
	}

	img, err := jpeg.Decode(file)
	if err != nil {
		log.Fatal(err)
	}
	file.Close()

	resized := resize.Resize(0, height, img, resize.NearestNeighbor)

	out, err := os.Create(dstImage)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	jpeg.Encode(out, resized, nil)
}
