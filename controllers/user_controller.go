package controllers

import (
	"context"
	"hollow/configs"
	"hollow/models"
	"hollow/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type UserController struct{}

// Register godoc
// @Summary      注册新用户
// @Description  注册一个新的用户账号
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user  body      models.User  true  "用户注册信息"
// @Success      201   {object}  RegisterResponse
// @Failure      400   {object}  ErrorResponse
// @Failure      500   {object}  ErrorResponse
// @Router       /api/auth/register [post]
func (uc *UserController) Register(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查邮箱是否已存在
	var existingUser models.User
	err := utils.DB.Collection("users").FindOne(context.Background(), bson.M{"email": user.Email}).Decode(&existingUser)
	if err != mongo.ErrNoDocuments {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email already exists"})
		return
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	user.Password = string(hashedPassword)
	user.CreatedAt = time.Now().Add(time.Hour * 8).Add(time.Hour * 8)
	user.UpdatedAt = time.Now().Add(time.Hour * 8).Add(time.Hour * 8)

	// 保存用户
	result, err := utils.DB.Collection("users").InsertOne(context.Background(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"id":      result.InsertedID,
	})
}

// Login godoc
// @Summary      用户登录
// @Description  使用邮箱和密码登录
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        login  body      LoginRequest  true  "登录信息"
// @Success      200    {object}  LoginResponse
// @Failure      400    {object}  ErrorResponse
// @Failure      401    {object}  ErrorResponse
// @Failure      500    {object}  ErrorResponse
// @Router       /api/auth/login [post]
func (uc *UserController) Login(c *gin.Context) {
	var loginData struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	err := utils.DB.Collection("users").FindOne(context.Background(), bson.M{"email": loginData.Email}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginData.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// 生成JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID.Hex(),
		"exp":     time.Now().Add(time.Hour * 24).Add(time.Hour * 8).Unix(),
	})

	tokenString, err := token.SignedString([]byte(configs.JWTSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
		"user": models.UserResponse{
			ID:    user.ID,
			Email: user.Email,
		},
	})
}

// LoginRequest represents the login request body
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email" example:"user@example.com"`
	Password string `json:"password" binding:"required" example:"password123"`
}

// LoginResponse represents the login response
type LoginResponse struct {
	Token string              `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	User  models.UserResponse `json:"user"`
}

// RegisterResponse represents the register response
type RegisterResponse struct {
	Message string `json:"message" example:"User registered successfully"`
	ID      string `json:"id" example:"507f1f77bcf86cd799439011"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error" example:"Invalid credentials"`
}
