package auth

import (
	"errors"
	"net/http"

	"super-indo-api/internal/auth/domain"
	"super-indo-api/internal/auth/port"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	register port.RegisterUseCase
	login    port.LoginUseCase
}

func NewAuthHandler(register port.RegisterUseCase, login port.LoginUseCase) *AuthHandler {
	return &AuthHandler{register: register, login: login}
}

func (h *AuthHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/register", h.Register)
	rg.POST("/login", h.Login)
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Status:  http.StatusBadRequest,
			Message: "format JSON tidak valid: " + err.Error(),
		})
		return
	}

	user, err := domain.NewUser(req.Email, req.Password, req.Name)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	if err := h.register.Execute(c.Request.Context(), user); err != nil {
		if errors.Is(err, domain.ErrEmailAlreadyExists) {
			c.JSON(http.StatusConflict, Response{
				Status:  http.StatusConflict,
				Message: err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, Response{
			Status:  http.StatusInternalServerError,
			Message: "gagal mendaftarkan user",
		})
		return
	}

	c.JSON(http.StatusCreated, Response{
		Status:  http.StatusCreated,
		Message: "registrasi berhasil",
		Data:    toUserResponse(user),
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Status:  http.StatusBadRequest,
			Message: "format JSON tidak valid: " + err.Error(),
		})
		return
	}

	user, token, err := h.login.Execute(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, Response{
				Status:  http.StatusUnauthorized,
				Message: err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, Response{
			Status:  http.StatusInternalServerError,
			Message: "gagal login",
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Status:  http.StatusOK,
		Message: "login berhasil",
		Data: LoginResponse{
			Token: token,
			User:  toUserResponse(user),
		},
	})
}
