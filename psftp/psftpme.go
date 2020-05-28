package psftp

import (
	"bufio"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"time"
)

type PsFtpMeConnStatus int

var PsFtpMeStatus PsFtpMeConnStatus

const (
	Off PsFtpMeConnStatus = iota
	Disconnected
	Connecting
	Connected
)

func (status PsFtpMeConnStatus) String() string {
	return [...]string{
		"Off",
		"Disconnected",
		"Connecting",
		"Connected",
	}[status]
}

func StartPsFtpMe() {
	// Update Config Value
	*PsFtpMe = true

	// Update Share URI
	FtpURI = PublicFtpURI

	// Attempt to Connect to Proxy
	go func() {
		// Keep Trying Until Stopped
		for *PsFtpMe {
			// Close Previous Connection
			if PsFtpMeConn != nil {
				_ = PsFtpMeConn.Close()
			}

			// Update Network Status
			PsFtpMeStatus = Connecting

			// Try to Connect...
			var err error
			if PsFtpMeConn, err = net.Dial("tcp", *PsFtpMeAddress); err != nil {
				if *VeryVerbose {
					log.Println(err.Error())
				}
				goto RETRY
			} else {
				// Update Network Status
				PsFtpMeStatus = Connected

				// Reeeeeeaddddd!
				reader := bufio.NewReader(PsFtpMeConn)

				// Conversations Can Be Tough. Let's Be Attentive Listeners!
				var line string
				for {
					// Yum! What a Delicious Byte!
					var char byte
					if char, err = reader.ReadByte(); err != nil {
						log.Println(err.Error())
						goto RETRY
					}
					if char == '\r' {
						// Unsafely Assume \r Followed by \n
						if char, err = reader.ReadByte(); err != nil {
							log.Println(err.Error())
							goto RETRY
						}
					}
					if char == '\n' {
						// What's in a Name!?
						if line != "220 psftp.me" {
							log.Println("Received Unwelcome Greeting!")
						} else {
							// Well Shit! Hello! Let's Talk!
							if _, err = io.WriteString(PsFtpMeConn, "|"+User+Pass+Filename+"|"); err != nil {
								log.Println(err.Error())
							} else {
								// Size!
								if _, err = io.WriteString(PsFtpMeConn, strconv.FormatInt(ZipFileStat.Size(), 10)+"|"); err != nil {
									log.Println(err.Error())
									break
								}

								for {
									// Ping!
									if _, err = io.WriteString(PsFtpMeConn, "E"); err != nil {
										log.Println(err.Error())
										break
									}

									// Consume...
									if char, err = reader.ReadByte(); err != nil {
										log.Println(err.Error())
										goto RETRY
									}
									switch char {
									case 'E': // Pong!?
										if *VeryVerbose {
											log.Println("sent pong to server's ping")
										}
									case '^': // Upload!?
										// LET'S FUCKIN' DO THIS!
										if _, err = io.WriteString(PsFtpMeConn, "^"); err != nil {
											log.Println(err.Error())
											break
										}

										// Stream the Data
										if zipFile, err := os.Open(ZipFile); err != nil {
											log.Println(err.Error())
											break
										} else {
											buf := make([]byte, 64512) // Just Under 64k Blocks (Size of a Packet)
											for {
												if bytesRead, err := zipFile.Read(buf); err != nil {
													if err != io.EOF {
														log.Println(err.Error())
													}
													goto UploadFinished
												} else if bytesRead > 0 {
													totalBytesWritten := 0
													for totalBytesWritten < bytesRead {
														if bytesWritten, err := PsFtpMeConn.Write(buf[totalBytesWritten:bytesRead]); err != nil {
															log.Println(err.Error())
															goto UploadFinished
														} else {
															totalBytesWritten += bytesWritten
														}
													}
												}
											}
										UploadFinished:
										}
									default:
										log.Println(char)
										break
									}
								}
							}
						}
					} else {
						// Just Concatenatin' Some Incantations
						line += string(char)
					}
				}
			}

			// Let's Loop!
		RETRY:
			PsFtpMeStatus = Disconnected
		}
	}()
}

func StopPsFtpMe() {
	// Update Network Status
	PsFtpMeStatus = Off

	// Update Config Value
	*PsFtpMe = false

	// Update Share URI
	FtpURI = PrivateFtpURI

	// Close the Connection
	if PsFtpMeConn != nil {
		if err := PsFtpMeConn.Close(); err != nil {
			log.Println(err.Error())
		}
		PsFtpMeConn = nil
	}
}

func GetPsFtpMeTitle() string {
	// Return the Current Status, Too!
	var status string
	if *PsFtpMe {
		switch PsFtpMeStatus {
		case Connected: // ◆
			status = statusRunes[2]
		case Connecting: // ◢ ◣ ◤ ◥
			status = spinnerRunes[time.Now().Second()%len(spinnerRunes)]
		case Disconnected: // ◈
			status = statusRunes[1]
		}
	} else { // ◇
		status = statusRunes[0]
	}
	return status + " psftp.me (Share on the Internet)"
}
