package user_controller

import (
	"github.com/gin-gonic/gin"
	"github.com/tiagoncardoso/fc-pge-auction-multithread/configuration/rest_err"
	"github.com/tiagoncardoso/fc-pge-auction-multithread/internal/infra/api/web/validation"
	"github.com/tiagoncardoso/fc-pge-auction-multithread/internal/usecase/user_usecase"
	"net/http"
)

type CreateUserController struct {
	createUserUseCase user_usecase.CreateUserUseCaseInterface
}

func NewCreateUserController(createUserUseCase user_usecase.CreateUserUseCaseInterface) *CreateUserController {
	return &CreateUserController{
		createUserUseCase: createUserUseCase,
	}
}

func (u *CreateUserController) CreateUser(c *gin.Context) {
	var userInputDTO struct {
		Name string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&userInputDTO); err != nil {
		restErr := validation.ValidateErr(err)

		c.JSON(restErr.Code, restErr)
		return
	}

	err := u.createUserUseCase.CreateUser(c, userInputDTO.Name)
	if err != nil {
		errRest := rest_err.ConvertError(err)
		c.JSON(errRest.Code, err)
		return
	}

	c.JSON(http.StatusCreated, "User created successfully")
}
