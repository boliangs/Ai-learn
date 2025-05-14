package controller

import (
	"net/http"

	"ai-interview/internal/model"
	"ai-interview/internal/service"
	"ai-interview/pkg/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// UserController 用户控制器
type UserController struct {
	userService *service.UserService
}

// NewUserController 创建新的用户控制器
func NewUserController() *UserController {
	return &UserController{
		userService: service.NewUserService(),
	}
}

// Register 用户注册
func (c *UserController) Register(ctx *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required,min=3,max=32"`
		Password string `json:"password" binding:"required,min=6,max=32"`
		Email    string `json:"email" binding:"required,email"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的输入"})
		return
	}

	// 创建用户
	user := &model.User{
		Username: input.Username,
		Password: input.Password,
		Email:    input.Email,
	}

	if err := c.userService.CreateUser(user); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "注册失败: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "注册成功",
		"data": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
	})
}

// Login 用户登录
func (c *UserController) Login(ctx *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的输入"})
		return
	}

	// 验证用户
	user, err := c.userService.VerifyUser(input.Username, input.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	// 生成token
	token, err := utils.GenerateToken(user.ID, user.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "生成token失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "登录成功",
		"data": gin.H{
			"token": token,
			"user": gin.H{
				"id":       user.ID,
				"username": user.Username,
				"email":    user.Email,
			},
		},
	})
}

// GetProfile 获取用户信息
func (c *UserController) GetProfile(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")
	user, err := c.userService.GetUser(userID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
	})
}

// UpdateProfile 更新用户信息
func (c *UserController) UpdateProfile(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")
	var input struct {
		Email    string `json:"email" binding:"omitempty,email"`
		Password string `json:"password" binding:"omitempty,min=6,max=32"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的输入"})
		return
	}

	// 更新用户信息
	updates := make(map[string]interface{})
	if input.Email != "" {
		updates["email"] = input.Email
	}
	if input.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "密码加密失败"})
			return
		}
		updates["password"] = string(hashedPassword)
	}

	if err := c.userService.UpdateUser(userID, updates); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "更新成功",
	})
}
