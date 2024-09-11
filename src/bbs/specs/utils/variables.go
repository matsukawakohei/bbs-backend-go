package utils

import (
	"bbs/internal/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var R *gin.Engine

var Db *gorm.DB

var User *model.User

var Token string

var ContentType = "application/json"

var Password = "password"
