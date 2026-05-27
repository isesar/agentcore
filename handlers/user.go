package handlers

import (
	"strconv"

	"agentcore/models"
	"agentcore/services"
	"github.com/gin-gonic/gin"
)

// GetUsers retrieves all users from the database
func GetUsers(c *gin.Context) {
	users, err := services.GetAllUsers()
	if err != nil {
		WriteInternalError(c, "Failed to retrieve users")
		return
	}
	c.JSON(200, users)
}

// GetUser retrieves a single user by ID
func GetUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		WriteValidationError(c, "Invalid user ID", nil)
		return
	}

	user, err := services.GetUserByID(id)
	if err != nil {
		WriteInternalError(c, "Failed to retrieve user")
		return
	}

	if user == nil {
		WriteNotFoundError(c, "User not found")
		return
	}

	c.JSON(200, user)
}

// CreateUser creates a new user
func CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		WriteValidationError(c, "Invalid request body", err.Error())
		return
	}

	createdUser, err := services.CreateUser(&user)
	if err != nil {
		WriteInternalError(c, "Failed to create user")
		return
	}

	c.JSON(201, createdUser)
}

// UpdateUser updates an existing user
func UpdateUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		WriteValidationError(c, "Invalid user ID", nil)
		return
	}

	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		WriteValidationError(c, "Invalid request body", err.Error())
		return
	}

	updatedUser, err := services.UpdateUser(id, &user)
	if err != nil {
		WriteInternalError(c, "Failed to update user")
		return
	}

	if updatedUser == nil {
		WriteNotFoundError(c, "User not found")
		return
	}

	c.JSON(200, updatedUser)
}

// DeleteUser deletes a user by ID
func DeleteUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		WriteValidationError(c, "Invalid user ID", nil)
		return
	}

	err = services.DeleteUser(id)
	if err != nil {
		WriteInternalError(c, "Failed to delete user")
		return
	}

	c.JSON(200, gin.H{"message": "User deleted successfully"})
}
