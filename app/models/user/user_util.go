package user

import "gohub/pkg/database"

func IsEmailExist(email string) bool {
	var count int64
	database.DB.Model(User{}).Where("email = ?", email).Count(&count)
	return count > 0
}

func IsPhoneExist(phone string) bool {
	var count int64
	database.DB.Model(User{}).Where("phone = ?", phone).Count(&count)
	return count > 0
}

func GetByMulti(loginID string) (userModel User) {
	database.DB.Where("phone=?", loginID).Or("email=?", loginID).Or("name=?", loginID).First(&userModel)
	return
}

func GetByPhone(phone string) (userModel User) {
	database.DB.Where("phone=?", phone).First(&userModel)
	return
}

// GetByEmail 通过 Email 来获取用户
func GetByEmail(email string) (userModel User) {
	database.DB.Where("email = ?", email).First(&userModel)
	return
}

func Get(idstr string) (userModel User) {
	database.DB.Where("id", idstr).First(&userModel)
	return
}