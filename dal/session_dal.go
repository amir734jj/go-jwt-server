package dal

import (
	"go-jwt-server/models"
	"go-jwt-server/types"
)

func FindSession(db types.DatabaseT, id uint) models.Session {

	var session models.Session
	db.First(id, &session)

	return session
}

func AddSession(db types.DatabaseT, session models.Session) models.Session {
	db.Create(&session)

	return session
}

func UpdateSession(db types.DatabaseT, id uint, session models.Session) models.Session {
	db.Model(id).Updates("Name")

	return session
}

func DeleteSession(db types.DatabaseT, id uint) models.Session {
	var session models.Session
	db.First(&session)

	db.Delete(id)

	return session
}
