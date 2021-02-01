package dal

import (
	"go-jwt-server/models"
	"go-jwt-server/types"
)

func QuerySession(db *types.DatabaseT, table *map[string]interface{}) (u *models.Session, err error) {
	var session = models.Session{}

	result := db.Where(*table).First(&session)

	return &session, result.Error
}

func AddSession(db *types.DatabaseT, session *models.Session) *models.Session {
	db.Create(&session)

	return session
}

func DeleteSessions(db *types.DatabaseT, userId uint) (sessions []models.Session, err error) {
	cond := map[string]interface{}{
		"user_id": userId,
	}
	sessions = make([]models.Session, 10)
	db.Find(&sessions, cond)

	ids := make([]uint, 10)
	for _, session := range sessions {
		ids = append(ids, session.Id)
	}

	result := db.Delete(&sessions, ids)

	return sessions, result.Error
}
