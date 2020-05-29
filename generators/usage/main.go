//go:generate go run main.go

package main

import (
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

func main() {
	// Back to Project Root
	if err := os.Chdir("../.."); err != nil {
		panic(err)
	}

	// Go?
	if goBin, err := exec.LookPath("go"); err != nil {
		panic(err)
	} else {
		// Get Usage Output
		cmd := exec.Command(goBin, "run", "-ldflags", "-linkmode=internal", ".")
		if usage, err := cmd.CombinedOutput(); err != nil {
			panic(err)
		} else {
			// Save Usage Output
			if err := ioutil.WriteFile("USAGE.txt", []byte(strings.TrimSpace(string(usage))), 0777); err != nil {
				panic(err)
			}
		}
	}
}
