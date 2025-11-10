package hash

import "github.com/gin-gonic/gin"

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// Success sends a successful JSON response with status code 200
func Success(c *gin.Context, data interface{}) {
	c.JSON(200, Response{
		Success: true,
		Data:    data,
	})
}

// Created sends a successful JSON response with status code 201
func Created(c *gin.Context, data interface{}) {
	c.JSON(201, Response{
		Success: true,
		Data:    data,
	})
}

// Error sends an error JSON response with the specified status code
func Error(c *gin.Context, statusCode int, message string, err error) {
	c.JSON(statusCode, Response{
		Success: false,
		Message: message,
		Error:   err.Error(),
	})
}

// BadRequest sends a 400 Bad Request response
func BadRequest(c *gin.Context, message string, err error) {
	Error(c, 400, message, err)
}

// Unauthorized sends a 401 Unauthorized response
func Unauthorized(c *gin.Context, message string, err error) {
	Error(c, 401, message, err)
}

// Forbidden sends a 403 Forbidden response
func Forbidden(c *gin.Context, message string, err error) {
	Error(c, 403, message, err)
}

// NotFound sends a 404 Not Found response
func NotFound(c *gin.Context, message string, err error) {
	Error(c, 404, message, err)
}

// InternalServerError sends a 500 Internal Server Error response
func InternalServerError(c *gin.Context, message string, err error) {
	Error(c, 500, message, err)
}
