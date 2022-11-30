package interfaces

import (
	"github.com/gin-gonic/gin"
	"github.com/mauryasaurav/timescale_database/models"
)

/* Declare All Interface Related to Users DB  */
type UserUseCase interface {
	CreateUser(ctx *gin.Context, req *models.UserCreateAndUpdate) (*models.UserResponse, error)
	GetUserBytesByFilter(ctx *gin.Context, name string) (*models.UserResponse, error)
}
