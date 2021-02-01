package view_models

type RegisterViewModel struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	Realm    string `json:"realm" validate:"required"`
}
