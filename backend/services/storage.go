package services

import (
	"backend/config"
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// MinIOクライアントを作成
func NewMinIOClient() (*minio.Client, error) {
	return minio.New(config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.AccessKeyID, config.SecretAccessKey, ""),
		Secure: config.UseSSL,
	})
}

// バケット内の画像一覧を取得
func GetImageList() ([]map[string]string, error) {
	fmt.Println("GetImageHandler was called")
	// クライアント作成
	cloudClient, err := NewMinIOClient()
	if err != nil {
		return nil, err
	}

	// バケット内のオブジェクト一覧を取得
	var imageList []map[string]string
	objectCh := cloudClient.ListObjects(context.Background(), config.BucketName, minio.ListObjectsOptions{})

	for object := range objectCh {
		if object.Err != nil {
			return nil, object.Err
		}

		// クラウド上の画像URLを作成
		imageList = append(imageList, map[string]string{
			"id":  object.Key,
			"url": fmt.Sprintf("http://%s/%s/%s", config.Endpoint, config.BucketName, object.Key),
		})
	}

	return imageList, nil
}

// クラウドにファイルをアップロードする
func UploadToCloud(buffer *bytes.Buffer, fileName string) (string, error) {
	// クライアント作成
	cloudClient, err := NewMinIOClient()
	if err != nil {
		return "", err
	}

	// ファイル名の重複を避けるためタイムスタンプを付与
	uniqueFileName := fmt.Sprintf("%d-%s", time.Now().UnixNano(), fileName)

	// アップロード
	_, err = cloudClient.PutObject(
		context.Background(), config.BucketName, uniqueFileName, buffer, int64(buffer.Len()),
		minio.PutObjectOptions{
			ContentType: "image/jpeg",
		},
	)
	if err != nil {
		return "", err
	}

	// クラウドのファイルURLを作成
	return fmt.Sprintf("http://%s/%s/%s", config.Endpoint, config.BucketName, uniqueFileName), nil
}
