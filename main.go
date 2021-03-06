package main

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rabbitmeow/golang-simple-api/config"
	"github.com/rabbitmeow/golang-simple-api/controllers"
	"github.com/rabbitmeow/golang-simple-api/utils/db"
)

var conf = config.ReadConfig()
var key = conf.ServerKey

func middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		bearer := c.Request.Header.Get("Authorization")
		if bearer != "" {
			strSplit := strings.Split(bearer, " ")
			if strSplit[0] == "Bearer" && strSplit[1] != "" {
				if strSplit[1] == key {
					c.Next()
				} else {
					c.AbortWithStatus(401)
					c.JSON(401, gin.H{
						"status":  401,
						"message": "unauthorized",
					})
				}
			} else {
				c.AbortWithStatus(401)
				c.JSON(401, gin.H{
					"status":  401,
					"message": "unauthorized",
				})
			}
		} else {
			c.AbortWithStatus(401)
			c.JSON(401, gin.H{
				"status":  401,
				"message": "unauthorized",
			})
		}
	}
}

func main() {
	welcomeController := new(controllers.WelcomeController)
	nameController := new(controllers.NameController)
	peopleController := new(controllers.PeopleController)

	db.Init()
	peopleController.Create()
	if conf.Mode == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()
	router.Use(middleware())
	v1 := router.Group("/v1")
	{
		v1.GET("/", welcomeController.Welcome)
		v1.GET("/name", nameController.Name)
		v1.GET("/peoples", peopleController.FetchAll)
	}

	router.Run(":" + conf.Port)
}
