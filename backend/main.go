package main

import (
	"backend/middlewares"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

const uploadDir = "./uploads"

func main() {
	// Ginのルーター作成
	r := gin.Default()
	// CORSの設定
	middlewares.SetupCORS(r)

	// アップロードディレクトリがないときは作成
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		os.Mkdir(uploadDir, os.ModePerm)
	}

	// 画像を格納するディレクトリを静的ファイルサーバーとして公開
	r.Static("/images", "./uploads")

	// 画像のエンドポイント
	r.GET("/get/images", getImageHandler)
	r.POST("/upload", uploadHnadler)

	// 8080ポート
	http.ListenAndServe(":8080", r)
}

// 写真をアップロードするハンドラー
func uploadHnadler(c *gin.Context) {
	// ファイルを取得
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ファイルの取得に失敗しました"})
		return
	}

	files := form.File["images"]

	// アップロードされたファイルを保存
	for _, file := range files {
		// 保存先のパスを指定
		filePath := fmt.Sprintf("%s/%s", uploadDir, file.Filename)
		// ファイルを保存
		if err := c.SaveUploadedFile(file, filePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ファイルの保存に失敗しました"})
			return
		}
	}

	// 成功レスポンス
	c.JSON(http.StatusOK, gin.H{"message": "ファイルアップロードに成功しました"})
}

// 写真を全て取得
func getImageHandler(c *gin.Context) {
	// 画像ファイルが保存されているディレクトリパス
	uploadDir := "./uploads"

	// ファイルをリストアップ
	files, err := os.ReadDir(uploadDir)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to read files"})
		return
	}

	// 画像ファイルのメタデータを格納
	var imageList []map[string]string
	for _, file := range files {
		if !file.IsDir() {
			// ファイルディレクトリでないときはファイルのURLを構築
			imageList = append(imageList, map[string]string{
				"id":  file.Name(),
				"url": "/images/" + file.Name(),
			})
		}
	}

	// 画像リストを返却
	c.JSON(http.StatusOK, imageList)
}
