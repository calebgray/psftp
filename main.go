package main

import (
	"crypto/sha256"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/mholt/archiver/v3"
	"github.com/riftbit/go-systray"
	"github.com/yob/graval"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

const AppTitle = "P.S. FTP"

// TODO: Display All Addrs on Systray
// TODO: Data Views
// TODO: https://github.com/jessevdk/go-flags

var statusRunes = strings.Split("◇◈◆", "")
var spinnerRunes = strings.Split("◢◣◤◥", "")

var AutoQuitTitle = map[bool]string{
	false: statusRunes[0] + " Automatically Quit",
	true:  statusRunes[2] + " Automatically Quit",
}

var Verbose *bool
var VeryVerbose *bool
var AutoQuit *bool
var PsFtpMe *bool
var PsFtpMeAddress *string
var ConfigPath *string

var User string
var Pass string
var Filename string
var ZipFile string
var ZipFileStat os.FileInfo
var ZipFileReady = make(chan bool, 1)
var PublicFtpURI string
var PrivateFtpURI string
var FtpURI string

var PsFtpMeConn net.Conn

var SortedFilenames []string

func generateURI(hostname string, port int) string {
	// Put Together as Friendly of a URI as We're Able
	strPort := ""
	if port != 21 {
		strPort = ":" + strconv.Itoa(port)
	}
	return fmt.Sprint("ftp://" + User + ":" + Pass + "@" + hostname + strPort + "/" + Filename)
}

func saveConfig() {
	// Serialize...
	if configJson, err := json.MarshalIndent(map[string]interface{}{
		"autoQuit":       *AutoQuit,
		"psftpme":        *PsFtpMe,
		"psftpmeAddress": *PsFtpMeAddress,
	}, "", "\t"); err != nil {
		log.Println(err.Error())
	} else {
		// ... Write to Disk!
		if err = ioutil.WriteFile(*ConfigPath, configJson, 0644); err != nil {
			log.Println(err.Error())
		}
	}
}

var lastPsFtpMeTitle string

func refreshPsFtpMeTitle(menuItem *systray.MenuItem) {
	// Only Update the Menu When There's New Text
	currentPsFtpMeTitle := getPsFtpMeTitle()
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
	menuPsFtpMe := systray.AddMenuItem(getPsFtpMeTitle(), "Generates a disposable Internet accessible link!", 0)
	menuAutoQuit := systray.AddMenuItem(AutoQuitTitle[*AutoQuit], "Automatically quits P.S. FTP after the next successful download.", 0)
	systray.AddSeparator()

	// Copy to Clipboard
	menuCopy := systray.AddMenuItem("Copy to Clipboard", "Copies the link to your clipboard.", 0)
	systray.AddSeparator()

	// Don't Quit!
	menuQuit := systray.AddMenuItem("Quit", "P.S. FTP Never Quits!", 0)

	// Events-in-a-Thread
	go func() {
		for {
			select {
			case <-time.After(333 * time.Millisecond):
				// Refresh...
				refreshPsFtpMeTitle(menuPsFtpMe)
			case <-menuPsFtpMe.OnClickCh():
				// Toggle!
				if *PsFtpMe {
					stopPsFtpMe()
				} else {
					startPsFtpMe()
				}
				_ = clipboard.WriteAll(FtpURI)
				refreshPsFtpMeTitle(menuPsFtpMe)
				saveConfig()
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

func tryStartServer(hostname string, port int) {
	// Generate the FTP URI's
	PrivateFtpURI = generateURI(hostname, port)
	psFtpMePortPosition := strings.LastIndexByte(*PsFtpMeAddress, ':')
	psFtpMePort, _ := strconv.Atoi((*PsFtpMeAddress)[psFtpMePortPosition+1:])
	PublicFtpURI = generateURI((*PsFtpMeAddress)[:psFtpMePortPosition], psFtpMePort)
	if *PsFtpMe {
		startPsFtpMe()
	} else {
		stopPsFtpMe()
	}
	_ = clipboard.WriteAll(FtpURI)

	// (Attempt to) Start the Server
	_ = graval.NewFTPServer(&graval.FTPServerOpts{
		Factory:    &PSFTPDriverFactory{},
		ServerName: "psftp",
		Hostname:   hostname,
		Port:       port,
	}).ListenAndServe()
}

func evaluateConfigBool(config map[string]interface{}, name string, fallback bool) bool {
	if raw, exists := config[name]; exists {
		if value, ok := raw.(bool); ok {
			return value
		}
	}
	return fallback
}

func evaluateConfigString(config map[string]interface{}, name string, fallback string) string {
	if raw, exists := config[name]; exists {
		if value, ok := raw.(string); ok {
			return value
		}
	}
	return fallback
}

func main() {
	// Run the OS Shim!
	Shim()

	// Setup Logging
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Parse Command Line Flags
	flag.Usage = func() {
		log.Println("Usage: ", path.Base(os.Args[0]), " [-flags] <filenames...>")
		flag.PrintDefaults()
	}
	Verbose = flag.Bool("v", false, "whisper sweet nothings aloud")
	VeryVerbose = flag.Bool("vv", false, "grab a loudspeaker while you're at it")
	AutoQuit = flag.Bool("autoQuit", true, "automatically quit after a successful download")
	PsFtpMe = flag.Bool("psftpme", false, "generate a (temporary but) Internet accessible URI to your files")
	PsFtpMeAddress = flag.String("psftpmeAddress", "psftp.me:21", "point to a different psftp.me server")
	configDisabled := flag.Bool("noconfig", false, "disable the config file")
	ConfigPath = flag.String("config", "psftp.json", "path to config file")
	listeners := flag.String("listeners", "", "comma separated list of addresses (hostname or ip) to bind with an FTP server")
	blacklist := flag.String("blacklist", "^127.,^169.", "comma separated list of prefixes (starts with ^), matches, or suffixes (ends with $) to blacklist")
	tempDir := flag.String("tempDir", "/tmp", "temporary storage")
	minPort := flag.Int("minPort", 1024, "minimum value for the port")
	maxPort := flag.Int("maxPort", 65535, "maximum value for the port")
	port := flag.Int("port", 0, "override range and set a specific port")
	ipv4 := flag.Bool("ipv4", true, "listen to ipv4 addresses")
	ipv6 := flag.Bool("ipv6", false, "listen to ipv6 addresses")
	flag.Parse()

	// Pre-Process Flags/Commands
	if *VeryVerbose {
		*Verbose = true
	}

	// Config File?
	if !*configDisabled && len(*ConfigPath) > 0 {
		configFile, err := os.Open(*ConfigPath)
		if *Verbose && err != nil {
			// Error Reading Config...
			log.Println(err.Error())
		} else {
			// Read Config File
			data, _ := ioutil.ReadAll(configFile)
			var config map[string]interface{}
			_ = json.Unmarshal(data, &config)
			_ = configFile.Close()

			// Parsed Config Values Override Command Line Arguments
			*AutoQuit = evaluateConfigBool(config, "autoQuit", *AutoQuit)
			*PsFtpMe = evaluateConfigBool(config, "psftpme", *PsFtpMe)
			*PsFtpMeAddress = evaluateConfigString(config, "psftpmeAddress", *PsFtpMeAddress)
		}
	}

	// Halp!?
	if len(flag.Args()) == 0 {
		flag.Usage()
		os.Exit(0)
	}

	// Post-Process Command Line
	useRandomPort := *port == 0

	// Parse Command Line Arguments
	uniqueFilenames := map[string]bool{}
	for _, filename := range flag.Args() {
		_, err := os.Stat(filename)
		uniqueFilenames[filename] = err == nil
	}
	if len(uniqueFilenames) == 0 {
		os.Exit(-1)
	}
	for filename, _ := range uniqueFilenames {
		SortedFilenames = append(SortedFilenames, filename)
	}
	sort.Strings(SortedFilenames)
	User = fmt.Sprintf("%x", sha256.Sum256([]byte(strings.Join(SortedFilenames, "\n"))))
	Pass = User[21:42]
	Filename = User[42:] + ".zip"
	User = User[:21]

	// Create the Zip
	go func() {
		ZipFile = path.Join(*tempDir, "psftp."+User+"."+Pass+"."+Filename)
		if *VeryVerbose {
			log.Println("Temporary File: ", ZipFile)
		}
		_ = os.Remove(ZipFile)
		err := archiver.Archive(SortedFilenames, ZipFile)
		if err != nil {
			log.Fatal(err)
		}
		ZipFileStat, err = os.Stat(ZipFile)
		if err != nil {
			log.Fatal(err)
		}
		ZipFileReady <- true
	}()

	// Determine Bind Addresses
	var addresses []string
	if len(*listeners) > 0 {
		// Parse Command Line!
		addresses = strings.Split(*listeners, ",")
	} else {
		// Auto-Detect Addresses!
		addrs, err := net.InterfaceAddrs()
		if err != nil {
			addresses = []string{"::"}
		} else {
			for _, addr := range addrs {
				address := addr.String()
				address = address[:strings.LastIndex(address, "/")] // Hack Off the CIDR
				addresses = append(addresses, address)
			}
		}
	}

	// Time to Start the Server(s)
	var wg sync.WaitGroup
	for _, address := range addresses {
		// Let's Filter Some Addresses...
		if strings.Contains(address, "::") {
			// IPv6 Disabled?
			if !*ipv6 {
				if *VeryVerbose {
					log.Println("Skipping IPv6: ", address)
				}
				goto NextAddress
			}
		} else if *ipv4 {
			// Blacklist
			for _, blacklisted := range strings.Split(*blacklist, ",") {
				blacklisted = strings.TrimSpace(blacklisted)
				if strings.HasPrefix(blacklisted, "^") && strings.HasPrefix(address, blacklisted[1:]) || strings.HasSuffix(blacklisted, "$") && strings.HasSuffix(address, blacklisted[:len(blacklisted)-1]) || strings.Contains(address, blacklisted) {
					if *VeryVerbose {
						log.Println("Skipping Blacklisted: ", address)
					}
					goto NextAddress
				}
			}
		} else {
			// Everything Disabled...?
			if *VeryVerbose {
				log.Println("Skipping Disabled: ", address)
			}
			goto NextAddress
		}

		// Fire Up the FTP Server
		if useRandomPort {
			// Try to Start the Server on a Random Port
			go func(hostname string) {
				for {
					wg.Add(1)
					tryStartServer(hostname, *minPort+rand.Intn(*maxPort-*minPort-1))
					wg.Done()
				}
			}(address)
		} else {
			// Try to Start the Server on a Specific Port
			go func(hostname string, port int) {
				wg.Add(1)
				tryStartServer(hostname, port)
				wg.Done()
			}(address, *port)
		}
	NextAddress:
	}

	// Wait for Server(s) to Exit
	wg.Wait()
}
