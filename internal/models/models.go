package models

import (
	"net/http"
)

type UserCookieType struct {
	// Name  string `json:"name"`
	// Pid   string `json:"pid"`
	// Refer string `json:"refer"`
	// Blank string `json:"_blank"`
	// T   bool   `json:"t"`

	// basic cookie
	UID string `json:"UID"`
	Vc3 string `json:"vc3"`
	Uid string `json:"_uid"`
	D   string `json:"_d"`
	Uf  string `json:"uf"`

	// user params
	Lv  string `json:"lv"`
	Fid string `json:"fid"`
}

func (uc *UserCookieType) ToCookies() []*http.Cookie {
	var cookies []*http.Cookie

	uidValue := uc.Uid
	if uidValue == "" {
		uidValue = uc.UID
	}

	cookies = append(cookies, &http.Cookie{Name: "uf", Value: uc.Uf})
	cookies = append(cookies, &http.Cookie{Name: "UID", Value: uidValue})
	cookies = append(cookies, &http.Cookie{Name: "_uid", Value: uidValue})
	cookies = append(cookies, &http.Cookie{Name: "_d", Value: uc.D})
	cookies = append(cookies, &http.Cookie{Name: "fid", Value: uc.Fid})
	cookies = append(cookies, &http.Cookie{Name: "vc3", Value: uc.Vc3})

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
