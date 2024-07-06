VERSION := v0.0.5

.PHONY: resin stamina charge all clean zip

ENV := CGO_ENABLED=1 GOOS=windows GOARCH=amd64
LDFLAGS := -ldflags "-H=windowsgui"

all: resin stamina charge

resin:
	${ENV} go build $(LDFLAGS) -o resin.exe cmd/resin/main.go

stamina:
	${ENV} go build $(LDFLAGS) -o stamina.exe cmd/stamina/main.go

charge:
	${ENV} go build $(LDFLAGS) -o charge.exe cmd/charge/main.go

zip: resin stamina
	zip "real-time-notes-$(VERSION)-x86_64.zip" resin.exe stamina.exe

clean:
	rm -rf resin*.exe
	rm -rf stamina*.exe
	rm -rf charge*.exe
	rm -rf login/*.exe.WebView2

