# ![P.S. FTP Logo](https://github.com/calebgray/psftp/raw/master/assets/icon.png) P.S. FTP
## Portable. Simple. FTP.
Drag and drop files/folders to `psftp` then paste the `ftp://` link in your clipboard wherever pleases you.

## GUI Legend
| Off   | Halfway | On    | Retrying |
| :---: | :---:   | :---: | :---:    |
| ◇     | ◈      | ◆     | ◢ ◣ ◤ ◥ |

## Features
* [X] Portable and Simple FTP Server
* [X] System/Notification Tray/Icon
* [X] Drag and Drop Multiple Files/Folders
* [X] Command Line Interface
* [X] Config File
* [X] Retries and Timeouts
* [X] Automatic Builds
* [X] Go Generate Build System
* [X] Versioning
* [ ] Auto-Update
* [ ] psftp.me Integration
* [ ] Release Candidate(s)
* [ ] Release First Stable Version
* [ ] TLS/SSL Support

## Development
```
# Clone!
git clone https://github.com/calebgray/psftp.git

# Option 1: Dockerfile
cd psftp
docker build -t psftp . && docker run --rm -it psftp

# Option 2: Ubuntu
cd psftp
sudo apt install -y build-essential curl git-all pkg-config libxxf86vm-dev libappindicator3-dev gcc-mingw-w64-x86-64 # Prerequisites
go generate .
go get -ldflags -linkmode=internal .
go build -ldflags -linkmode=internal .
./psftp

# Regenerate Assets
go generate . && go build -ldflags -linkmode=internal .
```

## Donations
[![By PayPal](https://github.com/calebgray/psftp/raw/master/assets/paypal.png)](https://www.paypal.com/cgi-bin/webscr?cmd=_s-xclick&hosted_button_id=AXXTUBFDC4DY2&source=url)

---

# Works in Progress

## psftp.me
Optionally, generate your very own Internet accessible link (limited resources; please support my development by donating). <3

Simply execute `psftp -psftpme` from the CLI, or enable it from the GUI.
