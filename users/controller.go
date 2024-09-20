package users

import (
	"fmt"
	"go-api/model"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	Service UserService
}

func NewUserController(service UserService) UserController {
	return UserController{Service: service}
}

func generateMessage(action string, success bool) string {
	if success {
		return action + " successfully"
	}
	return "Failed to " + action
}

func sendSuccess(c *gin.Context, action string, data interface{}) {
	c.JSON(http.StatusOK, gin.H{"message": generateMessage(action, true), "data": data})
}

func sendError(c *gin.Context, status int, action string) {
	c.JSON(status, gin.H{"message": generateMessage(action, false)})
}

func (ctrl *UserController) GetUsers(c *gin.Context) {
	users, err := ctrl.Service.GetUsers()
	if err != nil {
		sendError(c, http.StatusInternalServerError, "retrieve users")
		return
	}
	sendSuccess(c, "retrieve users", users)
}

func (ctrl *UserController) GetUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		sendError(c, http.StatusBadRequest, "parse ID")
		return
	}
	user, err := ctrl.Service.GetUserByID(uint(id))
	if err != nil {
		sendError(c, http.StatusNotFound, "find user")
		return
	}
	sendSuccess(c, "retrieve user", user)
}

func (ctrl *UserController) CreateUser(c *gin.Context) {
	var user model.User
	if err := c.BindJSON(&user); err != nil {
		sendError(c, http.StatusBadRequest, "parse request")
		return
	}
	if err := ctrl.Service.CreateUser(&user); err != nil {
		sendError(c, http.StatusInternalServerError, "create user")
		return
	}
	sendSuccess(c, "create user", user)
}

// ทำ LOGIN
// UserLogin handles POST requests to login a user.
func (ctrl *UserController) UserLogin(c *gin.Context) {
	var user model.User
	if err := c.BindJSON(&user); err != nil {
		sendError(c, http.StatusBadRequest, "parse request")
		return
	}

	userFromdb, token, err := ctrl.Service.UserLogin(user)
	if err != nil {
		if strings.Contains(err.Error(), "username not found") {
			sendError(c, http.StatusUnauthorized, "invalid username")
		} else if strings.Contains(err.Error(), "incorrect password") {
			sendError(c, http.StatusUnauthorized, "incorrect password")
		} else {
			sendError(c, http.StatusInternalServerError, "login")
		}
		return
	}

	// Set cookie
	c.SetCookie(
		"token",
		fmt.Sprintf("Bearer %v", token),
		15*60, // 15 min
		"/",
		"localhost",
		false,
		false,
	)

	// Return user
	sendSuccess(c, "login", userFromdb)
}

// ////
func (ctrl *UserController) UpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		sendError(c, http.StatusBadRequest, "parse ID")
		return
	}
	var updatedUser model.User
	if err := c.BindJSON(&updatedUser); err != nil {
		sendError(c, http.StatusBadRequest, "parse request")
		return
	}
	updatedUser.ID = uint(id)
	if err := ctrl.Service.UpdateUser(updatedUser); err != nil {
		sendError(c, http.StatusInternalServerError, "update user")
		return
	}
	sendSuccess(c, "update user", updatedUser)
}

func (ctrl *UserController) PatchUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		sendError(c, http.StatusBadRequest, "parse ID")
		return
	}
	var updatedFields map[string]interface{}
	if err := c.BindJSON(&updatedFields); err != nil {
		sendError(c, http.StatusBadRequest, "parse request")
		return
	}
	user, err := ctrl.Service.GetUserByID(uint(id))
	if err != nil {
		sendError(c, http.StatusNotFound, "find user")
		return
	}
	if username, ok := updatedFields["username"].(string); ok {
		user.Username = username
	}
	if password, ok := updatedFields["password"].(string); ok {
		user.Password = password
	}
	if err := ctrl.Service.UpdateUser(user); err != nil {
		sendError(c, http.StatusInternalServerError, "update user")
		return
	}
	sendSuccess(c, "update user", user)
}

func (ctrl *UserController) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		sendError(c, http.StatusBadRequest, "parse ID")
		return
	}
	user := model.User{ID: uint(id)}
	if err := ctrl.Service.DeleteUser(user); err != nil {
		sendError(c, http.StatusInternalServerError, "delete user")
		return
	}
	sendSuccess(c, "delete user", nil)
}
