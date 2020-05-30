//go:generate go run main.go

package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

func main() {
	log.SetFlags(log.LstdFlags)

	// Back to Project Root
	if err := os.Chdir("../.."); err != nil {
		log.Panic(err.Error())
	}

	// Go?
	if goBin, err := exec.LookPath("go"); err != nil {
		log.Panic(err.Error())
	} else {
		// Get Usage Output
		cmd := exec.Command(goBin, "run", "-ldflags", "-linkmode=internal", ".")
		if usage, err := cmd.CombinedOutput(); err != nil {
			log.Panic(err.Error())
		} else {
			// Save Usage Output
			if err := ioutil.WriteFile("USAGE.txt", []byte(strings.TrimSpace(string(usage))), 0777); err != nil {
				log.Panic(err.Error())
			}
		}
	}
}
