package main

import (
	"backend/middlewares"
	"bytes"
	"context"
	"fmt"
	"image"
	"mime/multipart"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/chai2010/webp"
	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// minioの設定
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
	var mu sync.Mutex // データの競合防止

	// ゴルーチンの完了を待つ
	var wg sync.WaitGroup
	// 同時処理の個数
	sem := make(chan struct{}, 10)

	// アップロードされたファイルを保存
	for _, file := range files {
		wg.Add(1)
		go func(f *multipart.FileHeader) {
			defer wg.Done()

			// 制限
			sem <- struct{}{}

			// 画像をリサイズしてアップロード
			cloudURL, err := processAndUpload(f)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "クラウドへのアップロードに失敗しました"})
				<-sem
				return
			}

			// ミューテックスをロックしてURLを追加
			mu.Lock()
			// アップロードURLをスライスに追加
			uploadedURLs = append(uploadedURLs, cloudURL)
			mu.Unlock()

			// ゴルーチン終了後にチャンネルから受け取る
			<-sem
		}(file)
	}

	// 全てのゴルーチンを待つ
	wg.Wait()

	// 成功レスポンス
	c.JSON(http.StatusOK, gin.H{"message": "ファイルアップロードに成功しました", "cloudURLs": uploadedURLs})
}

// 画像をリサイズしてWebpに変換
func processAndUpload(fileHeader *multipart.FileHeader) (string, error) {
	// ファイルを開く
	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	// 画像をデコード
	img, _, err := image.Decode(file)
	if err != nil {
		return "", err
	}

	// 画像の幅または高さが1024pxより大きい場合にリサイズする
	imgBounds := img.Bounds()
	width := imgBounds.Dx()
	height := imgBounds.Dy()

	// リサイズ前に画像の幅を確認して希望通りなら飛ばす
	if width > 1024 || height > 1024 {
		// 画像をリサイズ（幅1024pxに縮小）
		img = imaging.Resize(img, 1024, 0, imaging.Linear)
	}

	// 一時ファイルを作成
	tempFile, err := os.CreateTemp("", "resized-*.webp")
	if err != nil {
		return "", err
	}
	// 一時ファイルを削除
	defer os.Remove(tempFile.Name())

	// メモリ内のバッファを作成
	var buffer bytes.Buffer

	// WebPに変換して保存
	err = encodeWebP(&buffer, img)
	if err != nil {
		return "", err
	}

	// アップロード
	cloudURL, err := uploadToCloud(&buffer, fileHeader)
	if err != nil {
		return "", err
	}

	return cloudURL, nil
}

// WebPにエンコード
func encodeWebP(buffer *bytes.Buffer, img image.Image) error {
	return webp.Encode(buffer, img, &webp.Options{Lossless: false, Quality: 80})
}

// クラウドにファイルをアップロードする
func uploadToCloud(buffer *bytes.Buffer, fileHeader *multipart.FileHeader) (string, error) {
	// クライアント作成
	cloudClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return "", err
	}

	// ファイル名の重複を避けるためタイムスタンプを付与
	uniqueFileName := fmt.Sprintf("%d-%s", time.Now().UnixNano(), fileHeader.Filename)

	// マルチパートアップロードの最適化（5MBチャンク）
	partSize := 5 * 1024 * 1024

	// アップロード
	_, err = cloudClient.PutObject(
		context.Background(), bucketName, uniqueFileName, buffer, int64(buffer.Len()),
		minio.PutObjectOptions{
			ContentType: "image/jpeg",
			PartSize:    uint64(partSize),
		},
	)
	if err != nil {
		return "", err
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
