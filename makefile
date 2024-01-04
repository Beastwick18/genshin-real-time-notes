.PHONY: icon windows
windows:
	env GOOS=windows GOARCH=amd64 go build -ldflags "-H=windowsgui"

icon:
	~/go/bin/2goarray FullData icon < icon/full.ico > icon/full.go
	~/go/bin/2goarray NotFullData icon < icon/not_full.ico > icon/notFull.go
