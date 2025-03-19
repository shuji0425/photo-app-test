package handlers

import (
	"backend/services"
	"image"
	"mime/multipart"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

// 写真をアップロードするハンドラー
func UploadHandler(c *gin.Context) {
	// ファイルを取得
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ファイルの取得に失敗しました"})
		return
	}

	files := form.File["images"]
	var uploadedURLs []string
	var mu sync.Mutex     // データの競合防止
	var wg sync.WaitGroup // ゴルーチンの完了を待つ

	// 同時処理の個数
	sem := make(chan struct{}, 10)

	// アップロードされたファイルを保存
	for _, file := range files {
		wg.Add(1)
		go func(f *multipart.FileHeader) {
			defer wg.Done()

			// 制限
			sem <- struct{}{}

			// 画像のデコード
			img, err := decodeImage(f)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "画像のデコードに失敗しました"})
				<-sem
				return
			}

			// 画像の変換
			buffer, err := services.ProcessImage(img)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "画像の変換に失敗しました"})
				<-sem
				return
			}

			// クラウドアップロード
			cloudURL, err := services.UploadToCloud(buffer, f.Filename)
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

// 画像をデコード
func decodeImage(fileHeader *multipart.FileHeader) (image.Image, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	return img, err
}
