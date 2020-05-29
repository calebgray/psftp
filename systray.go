//+build windows linux darwin

package main

import (
	"github.com/atotto/clipboard"
	"github.com/riftbit/go-systray"
	"os"
	"strings"
	"time"
)

var lastPsFtpMeTitle string

func refreshPsFtpMeTitle(menuItem *systray.MenuItem) {
	// Only Update the Menu When There's New Text
	currentPsFtpMeTitle := GetPsFtpMeTitle()
	if currentPsFtpMeTitle != lastPsFtpMeTitle {
		_ = menuItem.SetTitle(currentPsFtpMeTitle)
		lastPsFtpMeTitle = currentPsFtpMeTitle
	}
}

func systrayBegin() {
	// Build the System Tray
	if systray.SetIcon(Icon) != nil {
		return
	}
	_ = systray.SetTitle(AppTitle)
	_ = systray.SetTooltip(AppTitle + " (" + strings.Join(SortedFilenames, ", ") + ")")

	// Start with Files
	for _, filename := range SortedFilenames {
		_ = systray.AddMenuItem(filename, "", 0).Disable()
	}
	systray.AddSeparator()

	// psftp.me, Auto-Quit
	var menuPsFtpMe *systray.MenuItem
	if *ShowPsFtpMe {
		menuPsFtpMe = systray.AddMenuItem(GetPsFtpMeTitle(), "Generates a disposable Internet accessible link!", 0)
	}
	menuAutoQuit := systray.AddMenuItem(AutoQuitTitle[*AutoQuit], "Automatically quits P.S. FTP after the next successful download.", 0)
	systray.AddSeparator()

	// Copy to Clipboard
	menuCopy := systray.AddMenuItem("Copy to Clipboard", "Copies the link to your clipboard.", 0)
	systray.AddSeparator()

	// Don't Quit!
	menuQuit := systray.AddMenuItem("Quit", "P.S. FTP Never Quits!", 0)

	// Events-in-a-Thread
	go func() {
		// Dynamic Menus
		if *ShowPsFtpMe {
			for {
				select {
				case <-time.After(333 * time.Millisecond):
					// Refresh...
					refreshPsFtpMeTitle(menuPsFtpMe)
				case <-menuPsFtpMe.OnClickCh():
					// Toggle!
					if *PsFtpMe {
						StopPsFtpMe()
					} else {
						StartPsFtpMe()
					}
					_ = clipboard.WriteAll(FtpURI)
					refreshPsFtpMeTitle(menuPsFtpMe)
					saveConfig()
				}
			}
		}

		// Static Menus
		for {
			select {
			case <-menuAutoQuit.OnClickCh():
				// Auto-Quit!?
				*AutoQuit = !*AutoQuit
				_ = menuAutoQuit.SetTitle(AutoQuitTitle[*AutoQuit])
				saveConfig()
			case <-menuCopy.OnClickCh():
				// Copy!
				_ = clipboard.WriteAll(FtpURI)
			case <-menuQuit.OnClickCh():
				// Quit!
				_ = systray.Quit()
			}
		}
	}()
}

func systrayExit() {
	// Attempt to Clean Up
	_ = os.Remove(ZipFile)

	// Clean Exit
	os.Exit(0)
}
