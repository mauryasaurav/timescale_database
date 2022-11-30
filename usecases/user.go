package usecases

import (
	"math/rand"

	"github.com/jmoiron/sqlx"
	"github.com/mauryasaurav/timescale_database/interfaces"
	"github.com/mauryasaurav/timescale_database/models"
	"github.com/mauryasaurav/timescale_database/repository"

	"github.com/gin-gonic/gin"
)

type userUseCase struct {
	db   *sqlx.DB
	repo *repository.UserRepository
}

func NewUserUseCase(db *sqlx.DB, r *repository.UserRepository) interfaces.UserUseCase {
	return &userUseCase{
		db:   db,
		repo: r,
	}
}

func (u *userUseCase) CreateUser(ctx *gin.Context, user *models.UserCreateAndUpdate) (*models.UserResponse, error) {
	RandomInteger := rand.Intn(20002-1000) + 1000
	createUser := &models.UserCreateAndUpdate{
		UserName: user.UserName,
		Email:    user.Email,
		Pass:     user.Pass,
		ReqBytes: RandomInteger,
	}

	err := u.repo.Save(ctx, createUser)
	if err != nil {
		return nil, err
	}

	return &models.UserResponse{
		UserName: user.UserName,
		ReqBytes: RandomInteger,
	}, nil
}

func (u *userUseCase) GetUserBytesByFilter(ctx *gin.Context, name string) (*models.UserResponse, error) {
	user, err := u.repo.GetUserBytesByFilter(ctx, name)
	if err != nil {
		return nil, err
	}
	return user, nil
}
