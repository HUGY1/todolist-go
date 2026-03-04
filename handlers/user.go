package handlers

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"todolist/config"
	"todolist/models"
	"todolist/utils"

	"github.com/gin-gonic/gin"
)

func CreateUser(c *gin.Context) {
	log.Println("[CreateUser] 收到新增用户请求")

	var req models.CreateUser
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[CreateUser] 请求参数无效: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求参数无效，mobile,password为必填项",
			"error":   err.Error(),
		})
		return
	}

	// 检查用户名是否存在
	var count int64

	if req.Mobile != "" {
		config.GetDB().Model(&models.User{}).Where("mobile = ?", req.Mobile).Count(&count)
		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "手机号已被注册",
				"error":   errors.New("手机号已被注册 " + req.Mobile),
			})

			return
		}
	}

	// 加密密码
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		log.Printf("[CreateUser] HashPassword加密错误: %v", err)
		return
	}

	user := models.User{
		Mobile:   req.Mobile,
		UserName: req.Mobile,
		Password: hashedPassword,
	}

	// 插入数据库

	if err = config.GetDB().Create(&user).Error; err != nil {
		log.Printf("[CreateUser] 用户新建失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "新建用户失败",
			"error":   err.Error(),
		})
	}

	log.Printf("[CreateUser] 成功新建用户, userId=%d", user.UserID)
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    user,
	})

}

func UpdateUser(c *gin.Context) {
	log.Println("[UpdateUser] 收到更新用户请求")

	var req models.UpdateUser
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[UpdateUser] 请求参数无效: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求参数无效，userId为必填项",
			"error":   err.Error(),
		})
		return
	}

	var user models.User

	// 检查手机号是否存在
	result := config.GetDB().First(&user, req.UserID)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "用户不存在",
			"error":   errors.New("用户不存在:" + string(req.UserID)),
		})

		return
	}

	updates := make(map[string]interface{})

	if req.UserName != nil {
		userName := strings.TrimSpace(*req.UserName)

		var reqError error
		if userName == "" {
			reqError = errors.New("用户名不能为空")

			return
		}

		if (len(userName) < 3 || len(userName) > 20) && reqError == nil {

			reqError = errors.New("用户名长度必须在3-20个字符之间")
		}

		if reqError != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": reqError.Error(),
				"error":   reqError,
			})
			return
		}

		updates["user_name"] = req.UserName
	}

	if req.Mobile != nil {
		mobile := strings.TrimSpace(*req.Mobile)

		var reqError error
		if mobile != "" {
			if len(mobile) != 11 {
				reqError = errors.New("手机号格式不正确")
			}

			if reqError == nil {
				var count int64
				config.GetDB().Model(&models.User{}).Where("mobile = ? AND userId != ?", mobile, req.UserID).Count(&count)
				if count > 0 && reqError != nil {
					reqError = errors.New("手机号格式不正确")

				}

			}

			if reqError != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"success": false,
					"message": reqError.Error(),
					"error":   reqError,
				})
				return
			}
			updates["mobile"] = req.Mobile

		}

	}

	if req.Gender < 0 || req.Gender > 2 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "性别值无效",
			"error":   errors.New("性别值无效"),
		})
		return
	}
	updates["gender"] = req.Gender

	if req.Avatar != nil {
		avatar := strings.TrimSpace(*req.Avatar)
		updates["avatar"] = avatar
	}

	if len(updates) == 0 {
		log.Printf("[UpdateUser] 没有要更新的字段")
		c.JSON(http.StatusOK, gin.H{
			"success": true,
		})
		return
	}

	updateResult := config.GetDB().Model(&user).Updates(updates)
	if updateResult.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "更新用户失败",
			"error":   updateResult.Error.Error(),
		})
		return
	}

	config.GetDB().First(&user, req.UserID)

	log.Printf("[UpdateUser] 成功更新用户, userId=%d", user.UserID)
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    user,
	})

}
