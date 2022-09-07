package main

import (
	"main/database"
	"main/docs"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lestrrat-go/jwx/jwa"
	ginjwt "github.com/sina-am/gin-jwt"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type AuthUser struct {
	Username string
	Password string
}

func main() {
	authMiddleware := ginjwt.JwtAuthentication{
		Authenticator: func(c *gin.Context) (interface{}, error) {
			form := AuthUser{}
			if err := c.ShouldBindJSON(&form); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"Error": "Bad request"})
				return nil, err
			}
			user, err := database.AuthenticateUser(form.Username, form.Password)
			if err != nil {
				c.JSON(http.StatusForbidden, gin.H{"Error": "Login failed"})
				return nil, err
			}
			return user.ID, nil
		},
		Authorizator: func(data interface{}) (interface{}, error) {
			return database.GetUserById(data.(uint))
		},
		SecretKey:   []byte("randomkey"),
		Algorithm:   jwa.HS512,
		IdentityKey: "user",
		TokenLookup: ginjwt.TokenLookup{From: ginjwt.Header, Name: "Authorization"},
	}

	database.InitDB()
	r := gin.Default()
	docs.SwaggerInfo.BasePath = "/api/v1"
	r.Use(authMiddleware.AuthenticationMiddleware())
	v1 := r.Group("/api/v1")
	v1.POST("/authenticate/", authMiddleware.LoginHandler)
	v1.GET("/providers/", GetProviders)
	v1.POST("/providers/", PostProviders)
	v1.POST("/sms/send/", SendSMS)
	v1.POST("/sms/bulk-send/")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	r.Run("localhost:8080")
}
