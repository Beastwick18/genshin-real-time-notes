VERSION := v0.0.5

.PHONY: hsr_icon genshin_icon icon resin stamina all clean zip

GCC := /usr/bin/x86_64-w64-mingw32-gcc
ENV := CGO_ENABLED=1 CC=${GCC} GOOS=windows GOARCH=amd64
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

icon: zzz_icon hsr_icon genshin_icon

zzz_icon:
	${GOPATH}/bin/2goarray ZzzFullData icon < assets/zzz/zzz_full.ico > pkg/icon/zzzFull.go
	${GOPATH}/bin/2goarray ZzzNotFullData icon < assets/zzz/zzz_not_full.ico > pkg/icon/zzzNotFull.go
	${GOPATH}/bin/2goarray ZzzEngagementData icon < assets/zzz/daily.ico > pkg/icon/zzzDaily.go
	${GOPATH}/bin/2goarray ZzzEngagementDoneData icon < assets/zzz/daily_done.ico > pkg/icon/zzzDailyDone.go
	${GOPATH}/bin/2goarray ZzzCheckInData icon < assets/zzz/zzz_checkin.ico > pkg/icon/zzzCheckIn.go
	${GOPATH}/bin/2goarray ZzzTicketData icon < assets/zzz/ticket.ico > pkg/icon/zzzTicket.go
	${GOPATH}/bin/2goarray ZzzTapeData icon < assets/zzz/tape.ico > pkg/icon/zzzTape.go

hsr_icon:
	${GOPATH}/bin/2goarray HsrFullData icon < assets/hsr/hsr_full.ico > pkg/icon/hsrFull.go
	${GOPATH}/bin/2goarray HsrNotFullData icon < assets/hsr/hsr_not_full.ico > pkg/icon/hsrNotFull.go
	${GOPATH}/bin/2goarray TrainingData icon < assets/hsr/daily.ico > pkg/icon/training.go
	${GOPATH}/bin/2goarray HsrExpeditionData icon < assets/hsr/expedition.ico > pkg/icon/hsrExpedition.go
	${GOPATH}/bin/2goarray EchoOfWarData icon < assets/hsr/echo.ico > pkg/icon/echo.go
	${GOPATH}/bin/2goarray HsrCheckInData icon < assets/hsr/checkinhsr.ico > pkg/icon/hsrCheckIn.go

genshin_icon:
	${GOPATH}/bin/2goarray FullData icon < assets/genshin/full.ico > pkg/icon/full.go
	${GOPATH}/bin/2goarray NotFullData icon < assets/genshin/not_full.ico > pkg/icon/notFull.go
	${GOPATH}/bin/2goarray CommissionData icon < assets/genshin/commission.ico > pkg/icon/commission.go
	${GOPATH}/bin/2goarray ExpeditionData icon < assets/genshin/expedition.ico > pkg/icon/expedition.go
	${GOPATH}/bin/2goarray RealmData icon < assets/genshin/realm.ico > pkg/icon/realm.go
	${GOPATH}/bin/2goarray WeeklyBossData icon < assets/genshin/domain.ico > pkg/icon/domain.go
	${GOPATH}/bin/2goarray GenshinCheckInData icon < assets/genshin/checkin.ico > pkg/icon/checkin.go
