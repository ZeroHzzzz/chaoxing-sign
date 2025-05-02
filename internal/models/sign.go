package models

type CourseType struct {
	CourseID string `json:"courseId"`
	ClassID  string `json:"classId"`
}

type ActivityType struct {
	ActivityID           string `json:"activityId"`
	Name                 string `json:"name"`
	CourseID             string `json:"courseId"`
	ClassID              string `json:"classId"`
	OtherID              int    `json:"otherId"` // 这个是用来区分签到类型的
	IfPhoto              int    `json:"ifPhoto"`
	OpenPreventCheatFlag int    `json:"openPreventCheatFlag"` // 验证码
	ChatID               string `json:"chatId"`
}

type LocationType struct {
	Address   string `json:"address"`
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
	Altitude  string `json:"altitude"`
}

type SignConfigType struct {
	Locations []LocationType `json:"locations"`
	PicUrl    string         `json:"picUrl"` // 图片上传到网盘后的url
}
