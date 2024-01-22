package genshin

const BaseURL = "https://bbs-api-os.hoyolab.com/game_record/genshin/api/dailyNote"

const DailyURL = "https://sg-hk4e-api.hoyolab.com/event/sol/sign?lang=en-us"
const ActID = "e202102251931481"

var Servers = map[byte]string{
	'1': "cn_gf01",
	'2': "cn_gf01",
	'5': "cn_qd01",
	'6': "os_usa",
	'7': "os_euro",
	'8': "os_asia",
	'9': "os_cht",
}

type GenshinAttendanceRewards struct {
	Status   string `json:"status"`
	Progress int    `json:"progress"`
}

type GenshinTaskRewards struct {
	Status string `json:"status"`
}

type GenshinDailyTask struct {
	TotalNum                  int                        `json:"total_num"`
	FinishedNum               int                        `json:"finished_num"`
	IsExtraTaskRewardReceived bool                       `json:"is_extra_task_reward_received"`
	TaskRewards               []GenshinTaskRewards       `json:"task_rewards"`
	AttendanceRewards         []GenshinAttendanceRewards `json:"attendance_rewards"`
	AttendanceVisible         bool                       `json:"attendance_visible"`
}

type GenshinTransformerRecovery struct {
	Day     int  `json:"Day"`
	Hour    int  `json:"Hour"`
	Minute  int  `json:"Minute"`
	Second  int  `json:"Second"`
	Reached bool `json:"reached"`
}

type GenshinTransformer struct {
	Obtained     bool                       `json:"obtained"`
	RecoveryTime GenshinTransformerRecovery `json:"recovery_time"`
	Wiki         string                     `json:"wiki"`
	Noticed      bool                       `json:"noticed"`
	LatestJobID  int                        `json:"latest_job_id"`
}

type GenshinExpedition struct {
	AvatarSideIcon string `json:"avatar_side_icon"`
	Status         string `json:"status"`
	RemainedTime   string `json:"remained_time"`
}

type GenshinData struct {
	CurrentResin              int                    `json:"current_resin"`
	MaxResin                  int                    `json:"max_resin"`
	ResinRecoveryTime         string                 `json:"resin_recovery_time"`
	FinishedTaskNum           int                    `json:"finished_task_num"`
	TotalTaskNum              int                    `json:"total_task_num"`
	IsExtraTaskRewardReceived bool                   `json:"is_extra_task_reward_received"`
	RemainResinDiscountNum    int                    `json:"remain_resin_discount_num"`
	ResinDiscountNumLimit     int                    `json:"resin_discount_num_limit"`
	CurrentExpeditionNum      int                    `json:"current_expedition_num"`
	MaxExpeditionNum          int                    `json:"max_expedition_num"`
	Expeditions               []GenshinExpedition    `json:"expeditions"`
	CurrentHomeCoin           int                    `json:"current_home_coin"`
	MaxHomeCoin               int                    `json:"max_home_coin"`
	HomeCoinRecoveryTime      string                 `json:"home_coin_recovery_time"`
	CalendarURL               string                 `json:"calendar_url"`
	Transformer               string                 `json:"transformer"`
	DailyTask                 GenshinDailyTask       `json:"daily_task"`
	ArchonQuestProgress       map[string]interface{} `json:"archon_quest_progress"`
}

type GenshinResponse struct {
	Retcode int         `json:"retcode"`
	Message string      `json:"message"`
	Data    GenshinData `json:"data"`
}

type GenshinDailyResponse struct {
	Retcode int                    `json:"retcode"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data"`
}
