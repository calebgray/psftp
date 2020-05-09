
build:
	2goarray Icon main < icon.ico > icon.go
	rsrc -ico icon.ico
	go build -ldflags "-linkmode=internal -H=windowsgui"
