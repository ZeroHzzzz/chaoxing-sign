package models

type User struct {
	ID               int            `json:"id" gorm:"primaryKey"`
	Username         string         `json:"username" gorm:"unique;not null"`
	Password         string         `json:"password" gorm:"not null"`
	ChaoxingAccounts []ChaoxingUser `json:"chaoxing_account"`
}

type ChaoxingUser struct {
	ID    int    `json:"id" gorm:"primaryKey"`
	Phone string `json:"phone"`
	Pass  string `json:"pass"`
	Name  string `json:"name"`
}
