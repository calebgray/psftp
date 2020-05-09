package main

import (
	"crypto/sha256"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/mholt/archiver"
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
	"time"
)

const AppTitle = "P.S. FTP"

var AutoQuitTitle = map[bool]string {
	false: "[ ] Automatically Quit",
	true:  "[X] Automatically Quit",
}

var Verbose *bool
var AutoQuit *bool

var User string
var Pass string
var Filename string
var ZipFile string
var FtpURI string

var SortedFilenames []string

func onReady() {
	if systray.SetIcon(Icon) != nil {
		return
	}
	systray.SetTitle(AppTitle)
	systray.SetTooltip(AppTitle + " (" + strings.Join(SortedFilenames, ", ") + ")")

	for _, filename := range SortedFilenames {
		systray.AddMenuItem(filename, "", 0).Disable()
	}
	systray.AddSeparator()
	menuAutoQuit := systray.AddMenuItem(AutoQuitTitle[*AutoQuit], "Automatically quits P.S. FTP after the next successful download.", 0)
	systray.AddSeparator()
	menuCopy := systray.AddMenuItem("Copy to Clipboard", "Copies the link to your clipboard.", 0)
	systray.AddSeparator()
	menuQuit := systray.AddMenuItem("Quit", "P.S. FTP Never Quits!", 0)

	go func() {
		for {
			select {
			case <-menuAutoQuit.OnClickCh():
				*AutoQuit = !*AutoQuit
				menuAutoQuit.SetTitle(AutoQuitTitle[*AutoQuit])
				break
			case <-menuCopy.OnClickCh():
				clipboard.WriteAll(FtpURI)
				break
			case <-menuQuit.OnClickCh():
				systray.Quit()
				break
			}
		}
	}()
}

func onExit() {
	os.Exit(0)
}

func tryStartServer(hostname string, port int) {
	// Generate the FTP URI
	FtpURI = fmt.Sprint("ftp://" + User + ":" + Pass + "@" + hostname + ":" + strconv.Itoa(port) + "/" + Filename)
	clipboard.WriteAll(FtpURI)

	// (Attempt to) Start the Server
	_ = graval.NewFTPServer(&graval.FTPServerOpts{
		Factory:    &PSFTPDriverFactory{},
		ServerName: "Prier",
		Hostname:   hostname,
		Port:       port,
	}).ListenAndServe()
}

type Config struct {
	AutoQuit bool `json:"autoQuit"`
}

func main() {
	// Run the OS Shim!
	Shim()

	// Parse Command Line Flags
	flag.Usage = func() {
		_, _ = fmt.Fprintln(flag.CommandLine.Output(), "Usage: ", path.Base(os.Args[0]), " [-flags] <filenames...>")
		flag.PrintDefaults()
	}
	Verbose = flag.Bool("v", false, "whisper sweet nothings aloud")
	AutoQuit = flag.Bool("autoQuit", true, "automatically quit after a successful download")
	listeners := flag.String("listeners", "", "comma separated list of addresses (hostname or ip) to bind with an FTP server")
	blacklist := flag.String("blacklist", "^127.,^169.", "comma separated list of prefixes (starts with ^), matches, or suffixes (ends with $) to blacklist")
	tempDir := flag.String("tempDir", "/tmp", "temporary storage")
	minPort := flag.Int("minPort", 1024, "minimum value for the port")
	maxPort := flag.Int("maxPort", 65535, "maximum value for the port")
	port := flag.Int("port", 0, "override range and set a specific port")
	ipv4 := flag.Bool("ipv4", true, "listen to ipv4 addresses")
	ipv6 := flag.Bool("ipv6", false, "listen to ipv6 addresses")
	flag.Parse()

	// Config File?
	configFile, err := os.Open("psftp.json")
	if *Verbose && err != nil {
		// Error Reading Config...
		_, _ = fmt.Fprintln(flag.CommandLine.Output(), err)
	} else {
		// Read Config File
		data, _ := ioutil.ReadAll(configFile)
		var config Config
		json.Unmarshal(data, &config)
		configFile.Close()
		if config.AutoQuit {
			println("YES")
		} else {
			println("NO")
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
	ZipFile = path.Join(*tempDir, "psftp."+User+"."+Pass+"."+Filename)
	if *Verbose {
		_, _ = fmt.Fprintln(flag.CommandLine.Output(), "Temporary File: ", ZipFile)
	}
	_ = os.Remove(ZipFile)
	err = archiver.Archive(SortedFilenames, ZipFile)
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(ZipFile)

	// Bind to Public Addresses
	var addresses []string
	if len(*listeners) > 0 {
		addresses = strings.Split(*listeners, ",")
	} else {
		// Auto-Detect Addresses
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
	for _, address := range addresses {
		// Let's Filter Some Addresses...
		if strings.Contains(address, "::") {
			// IPv6 Disabled?
			if !*ipv6 {
				if *Verbose {
					_, _ = fmt.Fprint(flag.CommandLine.Output(), "Skipping IPv6: ", address, "\n")
				}
				goto NextAddress
			}
		} else if *ipv4 {
			// Blacklist
			for _, blacklisted := range strings.Split(*blacklist, ",") {
				blacklisted = strings.TrimSpace(blacklisted)
				if strings.HasPrefix(blacklisted, "^") && strings.HasPrefix(address, blacklisted[1:]) || strings.HasSuffix(blacklisted, "$") && strings.HasSuffix(address, blacklisted[:len(blacklisted)-1]) || strings.Contains(address, blacklisted) {
					if *Verbose {
						_, _ = fmt.Fprint(flag.CommandLine.Output(), "Skipping Blacklisted: ", address, "\n")
					}
					goto NextAddress
				}
			}
		} else {
			if *Verbose {
				_, _ = fmt.Fprint(flag.CommandLine.Output(), "Skipping Disabled: ", address, "\n")
			}
			goto NextAddress
		}

		// Fire Up the FTP Server
		if useRandomPort {
			// Generate a Random Port (between minPort and maxPort); Try to Start the Server on Random Port
			go func(hostname string) {
				for {
					tryStartServer(hostname, *minPort+rand.Intn(*maxPort-*minPort-1))
				}
			}(address)
		} else {
			// Try to Start the Server on a Port
			go func(hostname string, port int) {
				tryStartServer(hostname, port)
			}(address, *port)
		}
	NextAddress:
	}

	// Wait for Server(s) to Exit
	hour, _ := time.ParseDuration("1h")
	for {
		time.Sleep(hour)
	}
}
