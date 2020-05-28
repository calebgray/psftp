package psftp

import (
	"github.com/ffred/guitocons"
	"github.com/riftbit/go-systray"
)

func Shim() {
	// Attach the Console to a Windows GUI App
	_ = guitocons.Guitocons()

	// Add the Tray Icon
	go func() {
		println(systray.Run(systrayBegin, systrayExit))
	}()
}
