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

	// 画像アップロードのエンドポイント
	r.POST("/upload", uploadHnadler)

	// 8080ポート
	http.ListenAndServe(":8080", r)
}

// 写真をアップロードするハンドラー
func uploadHnadler(c *gin.Context) {
	// ファイルから画像を取得
	file, err := c.FormFile("photo")

	if err != nil {
		c.JSON(400, gin.H{"error": "ファイルを取得できませんでした"})
		return
	}

	// 保存先のパスを指定
	filePath := fmt.Sprintf("%s/%s", uploadDir, file.Filename)

	// ファイルを保存
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(500, gin.H{"error": "ファイルの保存に失敗しました"})
		return
	}

	// 成功レスポンス
	c.JSON(200, gin.H{"message": "ファイルアップロード成功", "filePath": filePath})
}
