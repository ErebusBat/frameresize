package main

import (
	"fmt"
	"os"
	"runtime"

	. "github.com/ErebusBat/frameresize"
)

func main() {
	if os.Getenv("GOMAXPROCS") == "" {
		fmt.Println("$GOMAXPROCS undefined... turning on AutoFAST(tm)")
		runtime.GOMAXPROCS(runtime.NumCPU())
	}

	app := NewPhotoframe(
		"/tmp/photoframe",
		640, 480,
	)

	fmt.Println("Will output to", app.DestRoot)

	src_path := "/data/Dropbox/Photos/Photostream/2015/2015-02-WDW"
	app.Process(src_path)

	fmt.Printf("All Done, resized %d files\n", app.FilesResized)
}
