package main

import (
	"backend/middlewares"
	"context"
	"fmt"
	"mime/multipart"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

const (
	endpoint        = "localhost:9000"
	accessKeyID     = "admin123"
	secretAccessKey = "admin123"
	bucketName      = "photo"
	useSSL          = false
)

func main() {
	// Ginのルーター作成
	r := gin.Default()
	// CORSの設定
	middlewares.SetupCORS(r)

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
	var uploadedURLs []string

	// ゴルーチンの完了を待つ
	var wg sync.WaitGroup
	// 最大5つまで同時処理
	sem := make(chan struct{}, 5)

	// アップロードされたファイルを保存
	for _, file := range files {
		wg.Add(1)
		go func(f *multipart.FileHeader) {
			defer wg.Done()

			// 制限
			sem <- struct{}{}

			// クラウドにアップロード
			openedFile, err := file.Open()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "ファイルを開けませんでした"})
				return
			}
			defer openedFile.Close()

			cloudURL, err := uploadToCloud(openedFile, file.Filename)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "クラウドへのアップロードに失敗しました"})
				return
			}

			// アップロードURLをスライスに追加
			uploadedURLs = append(uploadedURLs, cloudURL)

			// ゴルーチン終了後にチャンネルから受け取る
			<-sem
		}(file)
	}

	// 全てのゴルーチンを待つ
	wg.Wait()

	// 成功レスポンス
	c.JSON(http.StatusOK, gin.H{"message": "ファイルアップロードに成功しました", "cloudURLs": uploadedURLs})
}

// クラウドにファイルをアップロードする
func uploadToCloud(file multipart.File, fileName string) (string, error) {
	// クライアント作成
	cloudClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return "", err
	}

	// ファイル名の重複を避けるためタイムスタンプを付与
	uniqueFileName := fmt.Sprintf("%d-%s", time.Now().UnixNano(), fileName)

	// アップロード
	_, err = cloudClient.PutObject(
		context.Background(), bucketName, uniqueFileName, file, -1, minio.PutObjectOptions{ContentType: "image/jpeg"},
	)
	if err != nil {
		return "", nil
	}

	// クラウドのファイルURLを作成
	cloudURL := fmt.Sprintf("http://%s/%s/%s", endpoint, bucketName, uniqueFileName)
	return cloudURL, nil
}

// 写真を全て取得
func getImageHandler(c *gin.Context) {
	// クライアント作成
	cloudClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "クラウドの接続に失敗しました"})
		return
	}

	// バケット内のオブジェクト一覧を取得
	var imageList []map[string]string
	objectCh := cloudClient.ListObjects(context.Background(), bucketName, minio.ListObjectsOptions{})

	for object := range objectCh {
		if object.Err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "クラウドの画像取得に失敗しました"})
			return
		}

		// クラウド上の画像URLを作成
		imageList = append(imageList, map[string]string{
			"id":  object.Key,
			"url": fmt.Sprintf("http://%s/%s/%s", endpoint, bucketName, object.Key),
		})
	}

	// 画像リストを返却
	c.JSON(http.StatusOK, imageList)
}
