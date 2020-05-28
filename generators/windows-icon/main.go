//go:generate go run main.go ../../assets/icon.ico ../../rsrc.syso

package main

import (
	"github.com/akavel/rsrc/rsrc"
	"os"
)

func main() {
	// Create a Binary Go File
	_ = rsrc.Embed(os.Args[2], "386", "", os.Args[1])
}
