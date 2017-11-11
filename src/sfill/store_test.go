
package sfill


import "testing"


func TestGet(t *testing.T) {  
    testpath := "/global/example/aaa"
    me:=Para{}
    me.setup()
    mymap := me.Getpath(testpath)
    if (mymap["bbb"] == "") {
       t.Errorf("cannot see example, result is ", mymap)
    }
    if (len(mymap) !=2 ) {
        t.Errorf("wrong number of values found, expected 2, was ",len(mymap))
    }
    
}

func TestReadFile(t *testing.T) {  
    testpath := "_kvtestdata.txt"
    strs := Readflatfile(testpath)
    if (len(strs) != 5) {
       t.Errorf("wrong number of values found, expected 5, was ",len(strs))
    }

}

func TestCRUD(t *testing.T) {
    var testlines []string=[]string{"/global/example/zzz/fish fingers",
                                    "/global/example/zzz/chip sauce",
                                    "/global/example/zzz/msg salt and vinegar"}

    base := "/testing"
    me:=Para{}
    me.setup()
    // create test
    me.Loader(testlines, base)
    // read test
    mymap := me.Getpath(base)
    // first k/v
    if (mymap["global/example/zzz/fish"] != "fingers") {
        t.Errorf("no fish fingers found ", mymap)
    }
    // k/v with spaces
    if (mymap["global/example/zzz/msg"] != "salt and vinegar") {
        t.Errorf("no salt and vinegar found ", mymap)
    }

    // update 
    var updateline []string=[]string{"/global/example/zzz/fish narwal"}
    me.Loader(updateline, base)
    mymap3 := me.Getpath(base)
    if (mymap3["global/example/zzz/fish"] != "narwal") {
        t.Errorf("lack of narwal update detected", mymap3)
    }
    // delete test
    me.Delete(base + "/global/example/zzz/fish")
    me.Delete(base + "/global/example/zzz/chip")
    me.Delete(base + "/global/example/zzz/msg")
    mymap2 := me.Getpath(base)
    if (mymap2["global/example/zzz/fish"] == "narwal") {
        t.Errorf("old narwal found ", mymap2)
    }
    
}
