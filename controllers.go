package main

import (
	"log"
	"main/database"
	"main/entities"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// @BasePath /api/v1

// PingExample godoc
// @Summary Providers
// @Schemes
// @Description get providers
// @Tags example
// @Accept json
// @Produce json
// @Success 200 {string} Helloworld
// @Router /providers/ [get]
func GetProviders(c *gin.Context) {
	if c.GetBool("IsAuthenticated") {
		userId, err := c.Get("User")
		if !err {
			return
		}
		c.JSON(http.StatusOK, database.GetProvidersByUserId(userId.(uint)))
		log.Printf("Token: %d\n", userId.(uint))
	} else {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "permission denied"})
	}
}

func PostProviders(c *gin.Context) {
	provider := entities.Provider{}
	if err := c.ShouldBindJSON(&provider); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}
	database.CreateProvider(&provider)
	c.JSON(http.StatusCreated, gin.H{"message": provider})
}

func SendSMS(c *gin.Context) {
	providerId, ok := c.GetQuery("provider")
	if !ok {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "not found"})
		return
	}
	pID, err := strconv.ParseUint(providerId, 10, 32)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "not found"})
		return
	}

	database.GetProviderById(uint(pID))

}
