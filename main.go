///go:generate go generate ./generators/usage
///go:generate go generate ./generators/template
//go:generate go generate ./generators/icon
//go:generate go generate ./generators/windows-icon

package main

import (
	"crypto/sha256"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/mholt/archiver/v3"
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
var ShowPsFtpMe *bool
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
		"showPsftpme":    *ShowPsFtpMe,
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

func tryStartServer(hostname string, port int) {
	// Generate the FTP URI's
	PrivateFtpURI = generateURI(hostname, port)
	psFtpMePortPosition := strings.LastIndexByte(*PsFtpMeAddress, ':')
	psFtpMePort, _ := strconv.Atoi((*PsFtpMeAddress)[psFtpMePortPosition+1:])
	PublicFtpURI = generateURI((*PsFtpMeAddress)[:psFtpMePortPosition], psFtpMePort)
	if *PsFtpMe {
		StartPsFtpMe()
	} else {
		StopPsFtpMe()
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

	// Clean Logs By Default
	log.SetFlags(0)

	// Parse Command Line Flags
	flag.Usage = func() {
		log.Println("Usage:", path.Base(os.Args[0]), "[-flags...] <filenames...>")
		flag.PrintDefaults()
	}
	Verbose = flag.Bool("v", false, "whisper sweet nothings aloud")
	VeryVerbose = flag.Bool("vv", false, "grab a loudspeaker while you're at it")
	AutoQuit = flag.Bool("autoQuit", true, "automatically quit after a successful download")
	ShowPsFtpMe = flag.Bool("showPsftpme", false, "show|true or hide|false psftp.me")
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

		// Debugging Logs
		log.SetFlags(log.LstdFlags | log.Lshortfile)
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
			*ShowPsFtpMe = evaluateConfigBool(config, "showPsftpme", *ShowPsFtpMe)
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
