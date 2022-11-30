package handlers

import (
	"github.com/mauryasaurav/timescale_database/interfaces"
	"github.com/mauryasaurav/timescale_database/models"

	"net/http"

	"github.com/gin-gonic/gin"
)

type userHandler struct {
	userUseCase interfaces.UserUseCase
}

func NewUserHandler(e *gin.Engine, u interfaces.UserUseCase) {
	handler := userHandler{
		userUseCase: u,
	}
	e.POST("api/user", handler.CreateUserHandler)
	e.GET("api/user_summary", handler.UserSummaryHandler)

}

func (h *userHandler) CreateUserHandler(ctx *gin.Context) {
	req := new(models.UserCreateAndUpdate)
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.userUseCase.CreateUser(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusCreated, user)
}

func (h *userHandler) UserSummaryHandler(ctx *gin.Context) {
	name, _ := ctx.GetQuery("name")
	user, err := h.userUseCase.GetUserBytesByFilter(ctx, name)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusCreated, user)
}
