package jwt

import (
	"net/http"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/gogodjzhu/listen-tube/internal/app/auth"
	"github.com/gogodjzhu/listen-tube/web/controller/middleware/interceptor"
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
		CookieSameSite: http.SameSiteNoneMode,
		// SendAuthorization: true,
		// SecureCookie:      true,
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
				UserName:   user.Name,
				UserCredit: user.Credit,
			}, nil
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			// TODO: check if the user is authorized to perform the action
			return true
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, interceptor.APIResponseDTO[string]{
				Code: 1,
				Msg:  message,
			})
		},
		LoginResponse: func(c *gin.Context, code int, message string, time time.Time) {
			c.JSON(code, interceptor.APIResponseDTO[string]{
				Code: 0,
				Msg:  "ok",
				Data: message,
			})
		},
		LogoutResponse: func(c *gin.Context, code int) {
			c.JSON(code, interceptor.APIResponseDTO[string]{
				Code: 0,
				Msg:  "ok",
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
	var req RegisterRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, interceptor.APIResponseDTO[string]{
			Code: 1,
			Msg:  err.Error(),
		})
		return
	}
	if req.InvitedCode != "123456" {
		ctx.JSON(http.StatusBadRequest, interceptor.APIResponseDTO[RegisterResult]{
			Code: 1,
			Msg:  "invalid invited code",
		})
		return
	}
	err := m.authService.Register(req.Username, req.Password)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, interceptor.APIResponseDTO[RegisterResult]{
			Code: 1,
			Msg:  err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, interceptor.APIResponseDTO[RegisterResult]{
		Code: 0,
		Msg:  "success!",
		Data: RegisterResult(req.Username),
	})
}

func (m *JWTMiddleware) UserInfoHandler(ctx *gin.Context) {
	user := GetCurrentUser(ctx)
	ctx.JSON(http.StatusOK, interceptor.APIResponseDTO[UserInfo]{
		Code: 0,
		Msg:  "ok",
		Data: *user,
	})
}

type UserInfo struct {
	UserName   string
	UserCredit string
}

func GetCurrentUser(c *gin.Context) *UserInfo {
	claims := jwt.ExtractClaims(c)
	return claimsToUserinfo(claims)
}

func claimsToUserinfo(claims jwt.MapClaims) *UserInfo {
	return &UserInfo{
		UserName:   claims["username"].(string),
		UserCredit: claims["usercredit"].(string),
	}
}

func userinfoToClaims(user *UserInfo) jwt.MapClaims {
	return jwt.MapClaims{
		"username":   user.UserName,
		"usercredit": user.UserCredit,
	}
}

type RegisterRequest struct {
	Username    string `form:"username" json:"username" binding:"required"`
	Password    string `form:"password" json:"password" binding:"required"`
	InvitedCode string `form:"invitedCode" json:"invitedCode" binding:"required"`
}

type RegisterResult string
