package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
)

func createLib() {
	file := "params/foo.libsonnet"
	err := os.MkdirAll(filepath.Dir(file), 0755)
	if err != nil {
		log.Fatalln(err)
	}
	err = ioutil.WriteFile(file, []byte(`
local target = std.extVar('target');
local foo = target.cluster;
local bar = {
  attr1: 'unknown',
  attr2: 'unknown',
  att3: 'unknown',
} + (if std.objectHas(target, 'bar') then target.bar else {});

foo {
  bar: bar,
}
`), 0644)
	if err != nil {
		log.Fatalln(err)
	}
}

func createParams(count int) {
	for i := 0; i < count; i++ {
		file := fmt.Sprintf("params/p%d.libsonnet", i)
		err := os.MkdirAll(filepath.Dir(file), 0755)
		if err != nil {
			log.Fatalln(err)
		}
		err = ioutil.WriteFile(file, []byte(fmt.Sprintf(`
local foo = import 'foo.libsonnet';
// important - must call a func of foo whether or not that exists
// and refer to the attributes defined as defaults
foo.func(std.thisFile,{
	attr1: foo.bar.attr1,
	attr2: foo.bar.attr2,
	attr3: foo.bar.attr3,
})
`)), 0644)
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func createCaller(count int) {
	var b bytes.Buffer
	fmt.Fprintln(&b, `local params = {`)
	for i := 0; i < count; i++ {
		fmt.Fprintf(&b, "\tp%d: import 'params/p%d.libsonnet',\n", i, i)
	}
	fmt.Fprintln(&b, `};`)
	fmt.Fprintln(&b, "// important - the values of the map above must be referenced")
	fmt.Fprintln(&b, `std.foldl(function (prev,k) prev + { [k] : params[k] }, std.objectFields(params), {})`)
	err := ioutil.WriteFile("caller.jsonnet", []byte(b.String()), 0644)
	if err != nil {
		log.Fatalln(err)
	}
}

func main() {
	count := 20
	if len(os.Args) == 2 {
		i, err := strconv.Atoi(os.Args[1])
		if err != nil {
			log.Fatalln("bad count:", err)
		}
		count = i
	}
	cmd := exec.Command("rm", "-rf", "params/")
	err := cmd.Run()
	if err != nil {
		log.Fatalln(err)
	}
	createLib()
	createParams(count)
	createCaller(count)
}
