package middleware

import (
	"net/http"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/gogodjzhu/listen-tube/internal/app/auth"
)

type JWTMiddleware struct {
	*jwt.GinJWTMiddleware
	authService *auth.AuthService
}

func NewJWTMiddleware(authService *auth.AuthService) (*JWTMiddleware, error) {
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:            "test zone",
		SigningAlgorithm: "HS256",
		Key:              []byte("secret key"),
		Timeout:          time.Hour,
		MaxRefresh:       time.Hour,
		IdentityKey:      "listentube",
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*UserInfo); ok {
				return userinfoToClaims(v)
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return claimsToUserinfo(claims)
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			type login struct {
				Username string `form:"username" json:"username" binding:"required"`
				Password string `form:"password" json:"password" binding:"required"`
			}
			var loginVals login
			if err := c.ShouldBind(&loginVals); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			user, err := authService.Authenticate(loginVals.Username, loginVals.Password)
			if err != nil {
				return nil, jwt.ErrFailedAuthentication
			}
			return &UserInfo{
				UserName: user.Name,
			}, nil
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			// TODO: check if the user is authorized to perform the action
			return true
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		// specifiy where to get the token, format: "<source>:<name>", if there are multiple, use comma to separate
		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
		SendCookie:    true,
	})
	m := JWTMiddleware{
		GinJWTMiddleware: authMiddleware,
		authService:      authService,
	}
	return &m, err
}

func (m *JWTMiddleware) RegisterHandler(ctx *gin.Context) {
	type register struct {
		Username    string `form:"username" json:"username" binding:"required"`
		Password    string `form:"password" json:"password" binding:"required"`
		InvitedCode string `form:"invitedCode" json:"invitedCode" binding:"required"`
	}
	var req register
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.InvitedCode != "123456" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid invited code"})
		return
	}
	err := m.authService.Register(req.Username, req.Password)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "register success"})
}

func (m *JWTMiddleware) UserInfoHandler(ctx *gin.Context) {
	user := getCurrentUser(ctx)
	ctx.JSON(http.StatusOK, user)
}

type UserInfo struct {
	UserName string
}

func getCurrentUser(c *gin.Context) *UserInfo {
	claims := jwt.ExtractClaims(c)
	return claimsToUserinfo(claims)
}

func claimsToUserinfo(claims jwt.MapClaims) *UserInfo {
	return &UserInfo{
		UserName: claims["username"].(string),
	}
}

func userinfoToClaims(user *UserInfo) jwt.MapClaims {
	return jwt.MapClaims{
		"username": user.UserName,
	}
}
