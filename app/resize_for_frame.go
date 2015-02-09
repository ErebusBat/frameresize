package main

import (
	"fmt"
	"runtime"

	. "github.com/ErebusBat/frameresize"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	app := NewPhotoframe(
		"/tmp/photoframe",
		640, 480,
	)

	fmt.Println("Will output to", app.DestRoot)

	src_path := "/data/Dropbox/Photos/Photostream/2015/2015-02-WDW"
	app.Process(src_path)

	fmt.Printf("All Done, resized %d files\n", app.FilesResized)
}
