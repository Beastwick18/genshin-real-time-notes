package hsr

const BaseURL = "https://bbs-api-os.hoyolab.com/game_record/hkrpg/api/note"

var Servers = map[byte]string{
	'1': "prod_gf_cn",
	'2': "prod_gf_cn",
	'5': "prod_qd_cn",
	'6': "prod_official_usa",
	'7': "prod_official_eur",
	'8': "prod_official_asia",
	'9': "prod_official_cht",
}

type HsrExpedition struct {
	Avatars       []string `json:"avatars"`
	Status        string   `json:"status"`
	RemainingTime int      `json:"remaining_time"`
	Name          string   `json:"name"`
	ItemURL       string   `json:"item_url"`
}

type HsrData struct {
	CurrentStamina        int             `json:"current_stamina"`
	MaxStamina            int             `json:"max_stamina"`
	StaminaRecoveryTime   int             `json:"stamina_recover_time"`
	AcceptedExpiditionNum int             `json:"accepted_epedition_num"` // intentional typo. somewhere, a hoyo intern has messed up lol
	TotalExpeditionNum    int             `json:"total_expedition_num"`
	Expeditions           []HsrExpedition `json:"expeditions"`
	CurrentTrainScore     int             `json:"current_train_score"`
	MaxTrainScore         int             `json:"max_train_score"`
	CurrentRogueScore     int             `json:"current_rogue_score"`
	MaxRogueScore         int             `json:"max_rogue_score"`
	WeeklyCocoonCnt       int             `json:"weekly_cocoon_cnt"`
	WeeklyCocoonLimit     int             `json:"weekly_cocoon_limit"`
	CurrentReserveStamina int             `json:"current_reserve_stamina"`
	IsReserveStaminaFull  bool            `json:"is_reserve_stamina_full"`
}

type HsrResponse struct {
	Retcode int     `json:"retcode"`
	Message string  `json:"message"`
	Data    HsrData `json:"data"`
}
