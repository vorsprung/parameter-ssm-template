package sfill

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

var maxr int64 = 10
var yes = true

// Version of this package
const Version int = 1

// Para Parameter hook type
type Para struct {
	s *ssm.SSM
}

func attach() Para {
	me := Para{}
	me.setup()
	return me
}

func (p *Para) setup() {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-west-1")},
	)
	if err != nil {
		panic("yike")
	}
	p.s = ssm.New(sess)
}

// Getpath  take a base path for the AWS Parameter store and return a map of the key/values
// all keys have the base path removed so /stage/journals/root is given as
// journals/root
// TODO: there has to be a better way of structuring this
func (p Para) Getpath(path string) map[string]string {
	var base = regexp.MustCompile("^" + path + "/")
	rv := make(map[string]string)
	firsttime := true
	for walk, err := p.s.GetParametersByPath(&ssm.GetParametersByPathInput{Path: aws.String(path), MaxResults: &maxr,
		Recursive: &yes}); (walk.NextToken != nil || firsttime) && err == nil; {
		firsttime = false
		tok := walk.Parameters
		for _, v := range tok {
			kc := *v.Name
			k := base.ReplaceAllString(kc, "")
			rv[k] = *v.Value
		}
		walk, err = p.s.GetParametersByPath(&ssm.GetParametersByPathInput{Path: aws.String(path), MaxResults: &maxr,
			Recursive: &yes,
			NextToken: walk.NextToken})
	}
	return rv
}

//Parameterstorestringvalidate returns a matcher function
// which checks that strings are valid AWS Parameter Store keys
func Parameterstorestringvalidate() func(string) bool {
	matcher, err := regexp.Compile("^[a-zA-Z0-9.-_/]+$")
	if err != nil {
		exitErrorf("unexpected regexp problem %v", err)
	}
	return func(this string) bool {
		return matcher.MatchString(this)
	}
}

//Parameterbasevalidate ensure path parameters never have trailing /
func Parameterbasevalidate(base string) string {
	var lb = len(base)
	// remove a trailing / if one is present
	if base[lb-1:] == "/" {
		base = base[0:(lb - 1)]
	}
	return base
}

// take an array of "lines" with two white space seperated values
// split them, add the LH side to the base - this the Parameter name or key
// the RHS is the Parameter value
// all name/values are inserted into the Parameter store
// old values are overwritten
// Loader is implemented so that any invalid character (including space, tab and comma)
const find string = " \x08,"

//Lr  get string from start index of first character from "find" above
//like split()
func Lr(s string) []string {
	ix := strings.IndexAny(s, find) //first space,tab, comma
	return []string{string(s[:ix]), string(s[ix+1:])}
}

//Loader  get data from text file
func (p Para) Loader(kv []string, base string) {
	var ov = true
	var t = "String"
	m := Parameterstorestringvalidate()
	base = Parameterbasevalidate(base)
	for _, l := range kv {
		pair := Lr(l)
		middle := ""
		if pair[0][:1] != "/" {
			middle = "/"
		}
		var k = base + middle + pair[0]
		var v = pair[1]
		if !m(k) {
			fmt.Printf("invalid key  %s -> %s ", k, v)
		} else {

			_, err := p.s.PutParameter(&ssm.PutParameterInput{Name: aws.String(k), Overwrite: &ov, Type: &t, Value: aws.String(v)})
			//fmt.Printf(">> %s -> %s <<",k,v)
			if err != nil {
				fmt.Print("put parameter %s -> %s %v", k, v, err)
			}
		}
	}
}

//Delete a single key (and it's value) from the Parameter store
// intended for testing only
func (p Para) Delete(path string) error {
	_, e := p.s.DeleteParameter(&ssm.DeleteParameterInput{Name: aws.String(path)})
	return e
}

//Readflatfile load a file of lines into a string array
// skip blank lines and commented lines
func Readflatfile(filename string) []string {
	var rv []string
	f, err := os.Open(filename)
	if err != nil {
		exitErrorf("read error", err)
	}
	defer f.Close()
	r := bufio.NewReader(f)
	path, err := r.ReadString(10) // 0x0A separator = newline
	for err == nil {              // stop on eof
		if string(path[0]) != "#" && len(path) > 1 {
			rv = append(rv, path)
		}
		path, err = r.ReadString(10)
	}
	return rv
}

func main() {
	me := Para{}
	mymap := me.Getpath("/global/example/aaa")
	fmt.Print("result is\n")
	for k, v := range mymap {
		fmt.Printf("key[%s] value[%s]\n", k, v)
	}
}

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}
