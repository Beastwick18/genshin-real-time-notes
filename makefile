.PHONY: icon windows resin stamina all

all: resin stamina

resin:
	env GOOS=windows GOARCH=amd64 go build -ldflags "-H=windowsgui" -o resin.exe cmd/resin/main.go

stamina:
	env GOOS=windows GOARCH=amd64 go build -ldflags "-H=windowsgui" -o stamina.exe cmd/stamina/main.go

icon:
	${GOPATH}/bin/2goarray FullData icon < assets/icons/full.ico > pkg/icon/full.go
	${GOPATH}/bin/2goarray NotFullData icon < assets/icons/not_full.ico > pkg/icon/notFull.go
	${GOPATH}/bin/2goarray HsrFullData icon < assets/icons/hsr_full.ico > pkg/icon/hsrFull.go
	${GOPATH}/bin/2goarray HsrNotFullData icon < assets/icons/hsr_not_full.ico > pkg/icon/hsrNotFull.go
