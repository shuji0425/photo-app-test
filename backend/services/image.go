package services

import (
	"bytes"
	"image"

	"github.com/chai2010/webp"
	"github.com/disintegration/imaging"
)

// 画像をリサイズしてWebpに変換
func ProcessImage(img image.Image) (*bytes.Buffer, error) {
	// 画像の幅または高さが1024pxより大きい場合にリサイズする
	imgBounds := img.Bounds()
	width := imgBounds.Dx()
	height := imgBounds.Dy()

	// リサイズ前に画像の幅を確認して希望通りなら飛ばす
	if width > 1024 || height > 1024 {
		// 画像をリサイズ（幅1024pxに縮小）
		img = imaging.Resize(img, 1024, 0, imaging.Linear)
	}

	// メモリ内のバッファを作成
	var buffer bytes.Buffer

	// WebPに変換して保存
	err := webp.Encode(&buffer, img, &webp.Options{Lossless: false, Quality: 80})
	if err != nil {
		return nil, err
	}

	return &buffer, nil
}
