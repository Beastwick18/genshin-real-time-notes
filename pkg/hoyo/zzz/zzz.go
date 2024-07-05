package zzz

const BaseURL = "https://sg-act-nap-api.hoyolab.com/event/game_record_zzz/api/zzz/note"

const DailyURL = "https://sg-public-api.hoyolab.com/event/luna/os/sign"
const ActID = "e202406031448091"

// **Second** character in UID. First is always '1'
var Servers = map[byte]string{
	'0': "prod_gf_us",
	'3': "prod_gf_jp",
	'5': "prod_gf_eu",
	'7': "prod_gf_sg",
}

type ZzzEnergyProgress struct {
	Max     int `json:"max"`
	Current int `json:"current"`
}

type ZzzEnergy struct {
	Progress ZzzEnergyProgress `json:"progress"`
	Restore  int               `json:"restore"`
}

type ZzzVitality struct {
	Max     int `json:"max"`
	Current int `json:"current"`
}

type ZzzVhsSale struct {
	SaleState string `json:"sale_state"`
}

type ZzzData struct {
	Energy   ZzzEnergy   `json:"energy"`
	Vitality ZzzVitality `json:"vitality"`
	VhsSale  ZzzVhsSale  `json:"vhs_sale"`
	CardSign string      `json:"card_sign"`
}

type ZzzResponse struct {
	Retcode int     `json:"retcode"`
	Message string  `json:"message"`
	Data    ZzzData `json:"data"`
}

type ZzzDailyData struct {
	Code      string `json:"code"`
	RiskCode  int    `json:"risk_code"`
	GT        string `json:"gt"`
	Challenge string `json:"challenge"`
	Success   int    `json:"success"`
	IsRisk    bool   `json:"is_risk"`
}

type ZzzDailyResponse struct {
	Retcode int          `json:"retcode"`
	Message string       `json:"message"`
	Data    ZzzDailyData `json:"data"`
}
