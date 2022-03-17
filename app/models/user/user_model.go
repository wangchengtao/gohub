package user

import (
	"gohub/app/models"
	"gohub/pkg/database"
	"gohub/pkg/hash"
)

type User struct {
	models.BaseModel

	Name     string `json:"name,omitempty"`
	Email    string `json:"-"`
	Phone    string `json:"-"`
	Password string `json:"-"`

	models.CommonTimestampsField
}

func (u *User) Create() {
	database.DB.Create(&u)
}

func (u *User) ComparePassword(_password string) bool {
	return hash.BcryptCheck(_password, u.Password)
}

func (u *User) Save() (rowsAffected int64) {
	result := database.DB.Save(&u)

	return result.RowsAffected
}