package user_controller

import (
	"context"
	"net/http"

	"github.com/dpcamargo/distributed-auction-system-go/configuration/rest_err"
	"github.com/dpcamargo/distributed-auction-system-go/internal/usecase/user_usecase"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserController struct {
	UserUseCase user_usecase.UserUseCaseInterface
}

func NewUserController(userUsecase user_usecase.UserUseCaseInterface) *UserController {
	return &UserController{
		UserUseCase: userUsecase,
	}
}

func (u *UserController) FindUserById(c *gin.Context) {
	userId := c.Param("userId")

	if err := uuid.Validate(userId); err != nil {
		errRest := rest_err.NewBadRequestError("invalid fields", rest_err.Causes{
			Field:   "userId",
			Message: "invalid UUID value",
		})
		c.JSON(errRest.Code, errRest)
		return
	}

	userData, err := u.UserUseCase.FindUserById(context.Background(), userId)
	if err != nil {
		errRest := rest_err.ConvertError(err)
		c.JSON(errRest.Code, errRest)
		return
	}

	c.JSON(http.StatusOK, userData)
}
