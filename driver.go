package main

import (
	"github.com/golobby/container"
	"go-jwt-server/logic"
	"go-jwt-server/models"
	"go-jwt-server/types"
	"goji.io"
	"goji.io/pat"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"net/http"
)

func main() {
	_ = container.Singleton(func() *types.DatabaseT {
		db, err := gorm.Open(sqlite.Open("db.sqlite"), &gorm.Config{})
		if err != nil {
			panic("failed to connect database")
		}

		return db
	})

	var db *types.DatabaseT
	err := container.Make(&db)

	if err != nil {
		panic("Failed to resolve db context")
	}

	err = db.AutoMigrate(&models.User{}, &models.Session{})
	if err != nil {
		panic("Migrations failed")
	}

	mux := goji.NewMux()
	mux.HandleFunc(pat.Post("/register"), logic.Register)
	mux.HandleFunc(pat.Post("/login"), logic.Login)
	mux.HandleFunc(pat.Post("/logout"), logic.Logout)

	_ = http.ListenAndServe("localhost:8000", mux)
}
