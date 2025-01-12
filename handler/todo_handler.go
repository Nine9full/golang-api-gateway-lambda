package handler

import (
	"net/http"

	"github.com/Nine9full/workshop-deployment/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetTodos(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var todos []model.Todo

		if tx := db.Find(&todos); tx.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": tx.Error.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": todos})

	}
}

func CreateTodo(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		var todo model.Todo

		if err := c.BindJSON(&todo); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if tx := db.Save(&todo); tx.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": tx.Error.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": todo})
	}
}

func DeleteTodo(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		var todo model.Todo

		id := c.Param("id")

		if tx := db.Where("id = ?", id).Delete(&todo); tx.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": tx.Error.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": todo})
	}
}

func GetTodoByID(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		var todo model.Todo

		id := c.Param("id")

		if tx := db.Where("id = ?", id).First(&todo); tx.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": tx.Error.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": todo})
	}
}
