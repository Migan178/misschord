package middlewares

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Migan178/misschord-backend/internal/configs"
	customErrors "github.com/Migan178/misschord-backend/internal/errors"
	"github.com/Migan178/misschord-backend/internal/models"
	"github.com/Migan178/misschord-backend/internal/repository"
	ginJWT "github.com/appleboy/gin-jwt/v3"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func NewAuthMiddleware() (*ginJWT.GinJWTMiddleware, error) {
	config := configs.GetConfig()

	return ginJWT.New(&ginJWT.GinJWTMiddleware{
		Realm:          "misschord-be",
		Key:            []byte(config.Backend.AuthKey),
		Timeout:        24 * time.Hour,
		MaxRefresh:     7 * time.Hour,
		IdentityKey:    "id",
		SendCookie:     true,
		SecureCookie:   false,
		CookieHTTPOnly: true,
		CookieName:     "token_jwt",
		CookieSameSite: http.SameSiteLaxMode,
		TokenLookup:    "header:Authorization, cookie:token_jwt",

		PayloadFunc:     payloadFunc,
		IdentityHandler: identityHandler,
		Authenticator:   authenticator,
		Authorizer:      authorizer,
		Unauthorized:    unauthorized,
	})
}

func payloadFunc(data any) jwt.MapClaims {
	if data, ok := data.(*models.UserToken); ok {
		return jwt.MapClaims{
			"id":    data.ID,
			"email": data.Email,
		}
	}

	return jwt.MapClaims{}
}

func identityHandler(c *gin.Context) any {
	claims := ginJWT.ExtractClaims(c)

	return claims["id"].(string)
}

func authenticator(c *gin.Context) (any, error) {
	var loginData models.LoginUserRequest

	if err := c.ShouldBindJSON(&loginData); err != nil {
		return nil, ginJWT.ErrMissingLoginValues
	}

	user, err := repository.GetDatabase().Users.GetByEmail(c.Request.Context(), loginData.Email)
	if err != nil {
		var dbErr *repository.DatabaseError
		if errors.As(err, &dbErr) {
			if dbErr.Code == repository.ErrorCodeAuthenticationFailed {
				return nil, ginJWT.ErrFailedAuthentication
			}
		}

		fmt.Printf("database err: %v\n", err)
		return nil, customErrors.ErrInternalServer
	}

	ok, err := repository.CheckPassword(loginData.Password, user.HashedPassword)
	if err != nil {
		return nil, ginJWT.ErrFailedAuthentication
	}

	if !ok {
		return nil, ginJWT.ErrFailedAuthentication
	}

	return &models.UserToken{
		ID:    strconv.Itoa(user.ID),
		Email: user.Email,
	}, nil
}

func authorizer(c *gin.Context, data any) bool {
	id, ok := data.(string)
	if !ok {
		return false
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		return false
	}

	if idInt < 0 {
		return false
	}

	return true
}

func unauthorized(c *gin.Context, code int, message string) {
	if message == customErrors.ErrInternalServer.Error() {
		c.JSON(http.StatusInternalServerError, customErrors.APIError{
			Code:    customErrors.ErrorCodeInternalError,
			Message: "an error occurred",
		})
		return
	}

	c.JSON(code, customErrors.APIError{
		Code:    customErrors.ErrorCodeAuthorizationError,
		Message: message,
	})
}
