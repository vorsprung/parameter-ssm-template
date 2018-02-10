package loader

import (
	"fmt"
	"os"
	"sfill"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "takes two command line args, filename to read from and base path to write to\n")
		os.Exit(1)
	}
	me := sfill.Attach()
	filename := os.Args[1]
	list := sfill.Readflatfile(filename)
	base := os.Args[2]
	me.Loader(list, base)
}
