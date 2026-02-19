package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"todolist/config"
	"todolist/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type FileHandler struct {
	DB *gorm.DB
}

func NewFileHandler(db *gorm.DB) *FileHandler {
	return &FileHandler{DB: db}
}

// POST /upload
func Upload(c *gin.Context) {
	// 1. 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请上传文件", "detail": err.Error()})
		return
	}

	// 2. 创建存储目录
	uploadDir := "./uploads"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建目录失败"})
		return
	}

	// 3. 生成唯一文件名（防止重名覆盖）
	ext := filepath.Ext(file.Filename)
	newFileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	savePath := filepath.Join(uploadDir, newFileName)
	absPath, _ := filepath.Abs(savePath)
	fmt.Println("📄 文件保存到:", absPath)
	// 4. 保存文件到本地
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存文件失败"})
		return
	}

	// 5. 写入数据库
	fileRecord := models.File{
		FileName: file.Filename,
		FilePath: savePath,
		FileSize: file.Size,
		MimeType: file.Header.Get("Content-Type"),
	}
	if err := config.GetDB().Create(&fileRecord).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "数据库写入失败"})
		return
	}

	// 6. 返回 fileId
	c.JSON(http.StatusOK, gin.H{
		"message":  "上传成功",
		"fileId":   fileRecord.ID,
		"fileName": fileRecord.FileName,
	})
}
