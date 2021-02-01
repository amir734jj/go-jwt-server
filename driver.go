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

	authSub := goji.SubMux()
	authSub.HandleFunc(pat.Post("/register"), logic.Register)
	authSub.HandleFunc(pat.Post("/login"), logic.Login)
	authSub.HandleFunc(pat.Get("/logout"), logic.Logout)
	mux.Handle(pat.New("/authorize/*"), authSub)

	apiSub := goji.SubMux()
	apiSub.Use(logic.AuthorizeMiddleware)
	apiSub.HandleFunc(pat.Get("/amir"), logic.AuthorizedApi)
	mux.Handle(pat.New("/api/*"), apiSub)

	err = http.ListenAndServe("localhost:8000", mux)

	if err != nil {
		panic(err)
	}
}
