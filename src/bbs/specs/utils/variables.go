package utils

import (
	"bbs/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var R *gin.Engine

var Db *gorm.DB

var User *models.User

var Token string

var ContentType = "application/json"

var Password = "password"
