package auth

import (
	"arch/ikeppu/github.com/configs"
	"arch/ikeppu/github.com/pkg/jwt"
	"arch/ikeppu/github.com/pkg/req"
	"arch/ikeppu/github.com/pkg/response"
	"net/http"
)

type AuthHandlerDeps struct {
	*configs.Config
	*AuthService
}
type AuthHandler struct {
	*configs.Config
	*AuthService
}

func NewAuthHandler(router *http.ServeMux, deps AuthHandlerDeps) {
	authHandler := &AuthHandler{
		Config:      deps.Config,
		AuthService: deps.AuthService,
	}
	router.HandleFunc("POST /auth/login", authHandler.Login)
	router.HandleFunc("POST /auth/register", authHandler.Register())
}

func (handler *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	body, err := req.HandleBody[LoginRequest](&w, r)

	if err != nil {
		return
	}

	email, err := handler.AuthService.Login(body.Email, body.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	token, err := jwt.NewJWT(handler.Config.Auth.Secret).Create(email)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := LoginResponse{
		Token: token,
	}

	response.Json(w, data, 200)
}

func (handler *AuthHandler) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[RegisterRequest](&w, r)

		if err != nil {
			return
		}

		email, err := handler.AuthService.Register(body.Email, body.Password, body.Name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		token, err := jwt.NewJWT(handler.Config.Auth.Secret).Create(email)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data := RegisterResponse{
			Token: token,
		}

		response.Json(w, data, 200)
	}
}
