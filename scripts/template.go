package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"sfill"
)
// this is a very crude templating system using regexp to replace a token 
// assuming the base is set to /testing then
// and that there is a key /testing/aaa/bbb/ccc with a value mincebeef
// ${aaa/bbb/ccc} is replaced with mincebeef
//
func main() {

	if len(os.Args) != 3 {
		panic("must have two args, filename and basepath")
	}

    me:=sfill.Attach()
	substitutes := me.Getpath(os.Args[2])

	rname := os.Args[1] + ".templ"

    // check .templ file exists
    _,err := os.Stat(rname)
	if os.IsNotExist(err){
		fmt.Printf("%s.templ must exist", rname)
        os.Exit(2)
	}

    
	rsource, err := ioutil.ReadFile(rname)

	source := string(rsource)

	start := "${"
	finish := "}"
	for k, v := range substitutes {
		tofind, err := regexp.Compile(regexp.QuoteMeta(start + k + finish))
		if err != nil {
			panic(k)
		}
		source = tofind.ReplaceAllString(source, v)
	}

	fmt.Print(source)
}
