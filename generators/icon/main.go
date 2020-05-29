//go:generate go run main.go ../../assets/icon.ico ../../icon.go Icon

package main

import (
	"io/ioutil"
	"os"
	"strconv"
)

func main() {
	// Create a Binary Go File
	if bytes, err := ioutil.ReadFile(os.Args[1]); err != nil {
		panic(err.Error())
	} else {
		src := "package main\n\nvar " + os.Args[3] + " = []byte{"
		for i, byte := range bytes {
			src += strconv.Itoa(int(byte))
			if i < len(bytes)-1 {
				src += ", "
			}
		}
		src += "}\n"
		ioutil.WriteFile(os.Args[2], []byte(src), 0777)
	}
}
