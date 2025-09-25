package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/your-org/custos/internal/application/dto"
	"github.com/your-org/custos/internal/application/usecase/auth"
	"github.com/your-org/custos/pkg/errors"
)

type AuthHandler struct {
	registerUC *auth.RegisterUseCase
	loginUC    *auth.LoginUseCase
}

func NewAuthHandler(registerUC *auth.RegisterUseCase, loginUC *auth.LoginUseCase) *AuthHandler {
	return &AuthHandler{
		registerUC: registerUC,
		loginUC:    loginUC,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, &dto.ErrorResponse{
			Code:    "INVALID_REQUEST",
			Message: "Invalid request format",
		})
		return
	}

	userInfo, err := h.registerUC.Execute(c.Request.Context(), &req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, &dto.SuccessResponse{
		Data: userInfo,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, &dto.ErrorResponse{
			Code:    "INVALID_REQUEST",
			Message: "Invalid request format",
		})
		return
	}

	loginResp, err := h.loginUC.Execute(c.Request.Context(), &req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, &dto.SuccessResponse{
		Data: loginResp,
	})
}

func (h *AuthHandler) handleError(c *gin.Context, err error) {
	if domainErr, ok := err.(*errors.DomainError); ok {
		statusCode := h.getStatusCodeFromError(domainErr.Code)
		c.JSON(statusCode, &dto.ErrorResponse{
			Code:    domainErr.Code,
			Message: domainErr.Message,
			Fields:  domainErr.Fields,
		})
		return
	}

	c.JSON(http.StatusInternalServerError, &dto.ErrorResponse{
		Code:    "INTERNAL_SERVER_ERROR",
		Message: "Internal server error",
	})
}

func (h *AuthHandler) getStatusCodeFromError(code string) int {
	switch code {
	case errors.CodeUserNotFound, errors.CodeInvalidCredentials:
		return http.StatusUnauthorized
	case errors.CodeUserAlreadyExists:
		return http.StatusConflict
	case errors.CodeInvalidPassword:
		return http.StatusBadRequest
	case errors.CodeTokenExpired, errors.CodeTokenInvalid:
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError
	}
}