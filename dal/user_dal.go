package dal

import (
	"go-jwt-server/models"
	"go-jwt-server/types"
)

func QueryUser(db *types.DatabaseT, table *map[string]interface{}) (u *models.User, err error) {
	var user = models.User{}

	result := db.Where(*table).First(&user)

	return &user, result.Error
}

func FindUser(db *types.DatabaseT, id uint) (u *models.User, err error) {
	var user = models.User{}
	result := db.First(id, &user)

	return &user, result.Error
}

func AddUser(db *types.DatabaseT, user *models.User) (u *models.User, err error) {
	result := db.Create(&user)
	return user, result.Error
}

func UpdateUser(db *types.DatabaseT, id uint, user *models.User) (u *models.User, err error) {
	result := db.Updates(map[string]interface{}{
		"name":  user.Name,
		"email": user.Email,
	})
	return user, result.Error
}

func DeleteUser(db *types.DatabaseT, id uint) (u *models.User, err error) {
	var user = models.User{}

	db.First(&user)
	result := db.Delete(id)

	return &user, result.Error
}
