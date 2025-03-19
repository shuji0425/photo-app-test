package handlers

import (
	"backend/services"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 写真を全て取得
func GetImageHandler(c *gin.Context) {
	fmt.Println("GetImageHandler was called")

	imageList, err := services.GetImageList()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "クラウドの画像取得に失敗しました"})
		return
	}

	// 画像リストを返却
	c.JSON(http.StatusOK, imageList)
}
