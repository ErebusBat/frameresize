package frameresize

import (
	"fmt"
	"image/jpeg"
	"log"
	"os"
	"path"
	"sync/atomic"

	"github.com/nfnt/resize"
)

// Helper function that handles the stats and synchronization aspects of resize
func (pf *Photoframe) resizeFiles() {
	for {

		file := <-pf.chImages
		pf.wgProcess.Add(1)

		resizeImage(file, pf.DestRoot, pf.Width, pf.Height)

		pf.wgProcess.Done()
		atomic.AddUint32(&pf.FilesResized, 1)
	}
}

func resizeImage(srcImage ImageInfo, dstPath string, width, height uint) {
	dstImage := fmt.Sprintf("%s_%dx%d.jpg", srcImage.Name, width, height)
	dstImage = path.Join(dstPath, dstImage)
	fmt.Printf("Resizing %s\n", dstImage)

	file, err := os.Open(srcImage.FullPath())
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
