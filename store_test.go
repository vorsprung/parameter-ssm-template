package sfill

import (
	"os"
	"testing"
)

func TestGet(t *testing.T) {
	testpath := "/global/example/aaa"
	if !(awsConfig()) {
		t.Skip("skipping test AWS not available")
	}

	me := Para{}
	me.setup()
	mymap := me.Getpath(testpath)
	if mymap["bbb"] == "" {
		t.Errorf("cannot see example, result is %s", mymap)
	}
	if len(mymap) != 2 {
		t.Errorf("wrong number of values found, expected 2, was %d", len(mymap))
	}

}

func TestReadFile(t *testing.T) {
	testpath := "_kvtestdata.txt"
	strs := Readflatfile(testpath)
	if len(strs) != 5 {
		t.Errorf("wrong number of values found, wanted 5, was %d", len(strs))
	}

}

func TestCRUD(t *testing.T) {
	var testlines = []string{"/global/example/zzz/fish fingers",
		"/global/example/zzz/chip sauce",
		"/global/example/zzz/msg salt and vinegar"}

	if !(awsConfig()) {
		t.Skip("skipping test AWS not available")
	}
	base := "/testing"
	me := Para{}
	me.setup()
	// create test
	me.Loader(testlines, base)
	// read test
	mymap := me.Getpath(base)
	// first k/v
	if mymap["global/example/zzz/fish"] != "fingers" {
		t.Error("no fish fingers found ", mymap)
	}
	// k/v with spaces
	if mymap["global/example/zzz/msg"] != "salt and vinegar" {
		t.Error("no salt and vinegar found ", mymap)
	}

	// update
	var updateline = []string{"/global/example/zzz/fish narwal"}
	me.Loader(updateline, base)
	mymap3 := me.Getpath(base)
	if mymap3["global/example/zzz/fish"] != "narwal" {
		t.Error("lack of narwal update detected", mymap3)
	}
	// delete test
	me.Delete(base + "/global/example/zzz/fish")
	me.Delete(base + "/global/example/zzz/chip")
	me.Delete(base + "/global/example/zzz/msg")
	mymap2 := me.Getpath(base)
	if mymap2["global/example/zzz/fish"] == "narwal" {
		t.Error("old narwal found ", mymap2)
	}

}

func awsConfig() bool {
	if _, err := os.Stat(os.Getenv("HOME") + "/.aws/config"); os.IsNotExist(err) {
		return false
	}
	return true
}
