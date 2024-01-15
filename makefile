.PHONY: hsr_icon genshin_icon icon resin stamina all

all: resin stamina

resin:
	env GOOS=windows GOARCH=amd64 go build -ldflags "-H=windowsgui" -o resin.exe cmd/resin/main.go

stamina:
	env GOOS=windows GOARCH=amd64 go build -ldflags "-H=windowsgui" -o stamina.exe cmd/stamina/main.go

icon: hsr_icon genshin_icon

hsr_icon:
	${GOPATH}/bin/2goarray HsrFullData icon < assets/hsr/hsr_full.ico > pkg/icon/hsrFull.go
	${GOPATH}/bin/2goarray HsrNotFullData icon < assets/hsr/hsr_not_full.ico > pkg/icon/hsrNotFull.go
	${GOPATH}/bin/2goarray TrainingData icon < assets/hsr/daily.ico > pkg/icon/training.go
	${GOPATH}/bin/2goarray HsrExpeditionData icon < assets/hsr/expedition.ico > pkg/icon/hsrExpedition.go
	${GOPATH}/bin/2goarray EchoOfWarData icon < assets/hsr/echo.ico > pkg/icon/echo.go

genshin_icon:
	${GOPATH}/bin/2goarray FullData icon < assets/genshin/full.ico > pkg/icon/full.go
	${GOPATH}/bin/2goarray NotFullData icon < assets/genshin/not_full.ico > pkg/icon/notFull.go
	${GOPATH}/bin/2goarray CommissionData icon < assets/genshin/commission.ico > pkg/icon/commission.go
	${GOPATH}/bin/2goarray ExpeditionData icon < assets/genshin/expedition.ico > pkg/icon/expedition.go
	${GOPATH}/bin/2goarray RealmData icon < assets/genshin/realm.ico > pkg/icon/realm.go
	${GOPATH}/bin/2goarray WeeklyBossData icon < assets/genshin/domain.ico > pkg/icon/domain.go
