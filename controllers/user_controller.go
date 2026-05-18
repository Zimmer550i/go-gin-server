package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"go-server/dto"
	"go-server/repositories"
	"go-server/services"
	"go-server/utils"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService *services.UserService
}

func NewUserController(userService *services.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

func (c *UserController) GetUsers(ctx *gin.Context) {
	users := c.userService.GetUsers()

	utils.Success(ctx, http.StatusOK, "Users fetched successfully", users)
}

func (c *UserController) CreateUser(ctx *gin.Context) {
	var req dto.CreateUserRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(ctx, "Invalid request body", err.Error())
		return
	}

	user := c.userService.CreateUser(req)

	utils.Success(ctx, http.StatusCreated, "User created successfully", user)
}

func (c *UserController) GetUserByID(ctx *gin.Context) {
	id, err := parseIDParam(ctx)
	if err != nil {
		utils.BadRequest(ctx, "Invalid user ID", err.Error())
		return
	}

	user, err := c.userService.GetUserByID(id)
	if err != nil {
		if errors.Is(err, repositories.ErrUserNotFound) {
			utils.NotFound(ctx, "User not found", err.Error())
			return
		}

		utils.Error(ctx, http.StatusInternalServerError, "Something went wrong", err.Error())
		return
	}

	utils.Success(ctx, http.StatusOK, "User fetched successfully", user)
}

func (c *UserController) DeleteUser(ctx *gin.Context) {
	id, err := parseIDParam(ctx)
	if err != nil {
		utils.BadRequest(ctx, "Invalid user ID", err.Error())
		return
	}

	user, err := c.userService.DeleteUser(id)
	if err != nil {
		if errors.Is(err, repositories.ErrUserNotFound) {
			utils.NotFound(ctx, "User not found", err.Error())
			return
		}

		utils.Error(ctx, http.StatusInternalServerError, "Something went wrong", err.Error())
		return
	}

	utils.Success(ctx, http.StatusOK, "User deleted successfully", user)
}

func parseIDParam(ctx *gin.Context) (int, error) {
	idParam := ctx.Param("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		return 0, err
	}

	if id <= 0 {
		return 0, errors.New("id must be greater than zero")
	}

	return id, nil
}