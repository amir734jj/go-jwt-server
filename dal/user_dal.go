package dal

import (
	"crypto/sha256"
	b64 "encoding/base64"
	"go-jwt-server/models"
	"go-jwt-server/types"
)

const SALT = "SALT_HERE"

func FindUser(db *types.DatabaseT, id uint) (u *models.User, err error) {

	var user models.User
	result := db.First(id, &user)

	return &user, result.Error
}

func AddUser(db *types.DatabaseT, user *models.User) (u *models.User, err error) {
	h := sha256.New()

	h.Write([]byte(user.Password + SALT))
	user.Password = b64.StdEncoding.EncodeToString(h.Sum(nil))

	result := db.Create(&user)

	return user, result.Error
}

func UpdateUser(db *types.DatabaseT, id uint, user *models.User) (u *models.User, err error) {
	result := db.Model(id).Updates("Name")

	return user, result.Error
}

func DeleteUser(db *types.DatabaseT, id uint) (u *models.User, err error) {
	var user models.User
	db.First(&user)

	result := db.Delete(id)

	return &user, result.Error
}
