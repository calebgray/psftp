//go:generate go run main.go ../../assets/template.md ../../USAGE.md Title=Usage Body=@../../USAGE.txt

package main

import (
	"github.com/Masterminds/sprig"
	"io/ioutil"
	"os"
	"strings"
	"text/template"
)

func getTemplateData() map[string]string {
	data := map[string]string{}
	for _, arg := range os.Args[3:] {
		// Parse: Key=Value
		splitter := strings.IndexByte(arg, '=')
		key := arg[0:splitter]
		value := arg[splitter+1:]

		// Replace @./file.name With File's Contents
		if value[0] == '@' {
			// Read File
			if val, err := ioutil.ReadFile(value[1:]); err != nil {
				panic(err.Error())
			} else {
				value = string(val)
			}
		}

		// Set Value
		data[key] = value
	}
	return data
}

func main() {
	if tmplRaw, err := ioutil.ReadFile(os.Args[1]); err != nil {
		panic(err.Error())
	} else if tmpl, err := template.New("usage").Funcs(sprig.GenericFuncMap()).Parse(string(tmplRaw)); err != nil {
		panic(err.Error())
	} else if tmplOut, err := os.Create(os.Args[2]); err != nil {
		panic(err.Error())
	} else if err = tmpl.Execute(tmplOut, getTemplateData()); err != nil {
		panic(err.Error())
	}
}
