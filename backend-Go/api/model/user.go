package model

const (
	RoleAdmin = "admin"
	RoleUser  = "user"
)

type User struct {
	ID           int64  `gorm:"primaryKey" json:"id"`
	Username     string `gorm:"not null" json:"username"`
	Email        string `gorm:"not null" json:"email"`
	Role         string `gorm:"not null;default:user" json:"role"`
	PasswordHash string `gorm:"column:password_hash;not null" json:"-"`
}

func (User) TableName() string {
	return "users"
}
