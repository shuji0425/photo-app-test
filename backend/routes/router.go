package routes

import (
	"backend/handlers"
	"backend/middlewares"

	"github.com/gin-gonic/gin"
)

// ルーターをセットアップ
func SetRouter() *gin.Engine {
	// Ginのルーター作成
	r := gin.Default()

	// CORSの設定
	middlewares.SetupCORS(r)

	// 画像を格納するディレクトリを静的ファイルサーバーとして公開
	r.Static("/images", "./uploads")

	// 画像のエンドポイント
	r.GET("/get/images", handlers.GetImageHandler)
	r.POST("/upload", handlers.UploadHandler)

	return r
}
