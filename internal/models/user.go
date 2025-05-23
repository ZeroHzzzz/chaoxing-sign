package models

type User struct {
	ID       int    `json:"id" gorm:"primaryKey"`
	Email    string `json:"email" gorm:"unique;not null"`
	Username string `json:"username" gorm:"unique;not null"`
	Password string `json:"password" gorm:"not null"`
}

type ChaoxingAccount struct {
	ID        int    `json:"id" gorm:"primaryKey"`
	UserID    int    `json:"user_id" gorm:"not null;unique"`
	Phone     string `json:"phone" gorm:"not null"`
	Pass      string `json:"pass" gorm:"not null"`
	Name      string `json:"name"`
	CreatedAt int64  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt int64  `json:"updated_at" gorm:"autoUpdateTime"`
}
