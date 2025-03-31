package models

import (
	"net/http"
)

type UserCookieType struct {
	Name  string `json:"name"`
	UID   string `json:"UID"`
	Fid   string `json:"fid"`
	Pid   string `json:"pid"`
	Refer string `json:"refer"`
	Blank string `json:"_blank"`
	T     bool   `json:"t"`
	Vc3   string `json:"vc3"`
	Uid   string `json:"_uid"`
	D     string `json:"_d"`
	Uf    string `json:"uf"`
	Lv    string `json:"lv"`
}

func (uc *UserCookieType) ToCookies() []*http.Cookie {
	var cookies []*http.Cookie

	uidValue := uc.Uid
	if uidValue == "" {
		uidValue = uc.UID
	}

	if uc.Fid != "" {
		cookies = append(cookies, &http.Cookie{
			Name:  "fid",
			Value: uc.Fid,
		})
	}

	if uc.Uf != "" {
		cookies = append(cookies, &http.Cookie{
			Name:  "uf",
			Value: uc.Uf,
		})
	}

	if uc.D != "" {
		cookies = append(cookies, &http.Cookie{
			Name:  "_d",
			Value: uc.D,
		})
	}

	if uidValue != "" {
		cookies = append(cookies, &http.Cookie{
			Name:  "UID",
			Value: uidValue,
		})
	}

	// 检查并添加 vc3
	if uc.Vc3 != "" {
		cookies = append(cookies, &http.Cookie{
			Name:  "vc3",
			Value: uc.Vc3,
		})
	}

	return cookies
}

type CourseType struct {
	CourseID string `json:"courseId"`
	ClassID  string `json:"classId"`
}

type ActivityType struct {
	ActivityID string `json:"activityId"`
	Name       string `json:"name"`
	CourseID   string `json:"courseId"`
	ClassID    string `json:"classId"`
	OtherID    int    `json:"otherId"`
	IfPhoto    int    `json:"ifPhoto"`
	ChatID     string `json:"chatId"`
}
