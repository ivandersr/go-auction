package user_controller

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ivandersr/go-auction/config/rest_err"
	"github.com/ivandersr/go-auction/internal/usecase/user_usecase"
)

type UserController struct {
	userUseCase user_usecase.UserUsecaseInterface
}

func NewUserController(userUseCase user_usecase.UserUsecaseInterface) *UserController {
	return &UserController{
		userUseCase: userUseCase,
	}
}

func (u *UserController) FindUserById(c *gin.Context) {
	userId := c.Param("userId")

	if err := uuid.Validate(userId); err != nil {
		restErr := rest_err.NewBadRequestError("invalid fields", rest_err.Cause{
			Field:   "userId",
			Message: "invalid UUID value",
		})

		c.JSON(restErr.Code, restErr)
		return
	}

	userData, err := u.userUseCase.FindUserById(context.Background(), userId)

	if err != nil {
		restErr := rest_err.ConvertError(err)

		c.JSON(restErr.Code, restErr)
		return
	}

	c.JSON(http.StatusOK, userData)
}
