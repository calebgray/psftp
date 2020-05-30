//go:generate go run main.go ../../assets/icon.ico ../../rsrc.syso

package main

import (
	"github.com/akavel/rsrc/rsrc"
	"log"
	"os"
)

func main() {
	log.SetFlags(log.LstdFlags)

	// Create a Binary Go File
	log.Println(rsrc.Embed(os.Args[2], "386", "", os.Args[1]).Error())
}
