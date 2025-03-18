import { useState } from "react";
import Pica from "pica";

const pica = new Pica();

export const useImageCompressor = () => {
  const [compressionProgress, setCompressionProgress] = useState(0);

  const compressImage = async (
    file: File,
    format: "image/jpeg" | "image/webp" = "image/webp"
  ): Promise<File> => {
    return new Promise((resolve, reject) => {
      const img = new Image();
      img.src = URL.createObjectURL(file);
      img.onload = async () => {
        const canvas = document.createElement("canvas");
        const scale = Math.min(1024 / img.width, 1024 / img.height, 1);
        canvas.width = img.width * scale;
        canvas.height = img.height * scale;

        // 画像をリサイズ&圧縮
        const resizedCanvas = await pica.resize(img, canvas, {
          unsharpAmount: 80,
          unsharpRadius: 0.4,
          unsharpThreshold: 1.5,
        });

        // 圧縮率を更新
        setCompressionProgress(50);

        resizedCanvas.toBlob(
          (blob) => {
            if (!blob) return reject("画像圧縮に失敗しました");
            resolve(
              new File(
                [blob],
                "compressed." + (format === "image/webp" ? "webp" : "jpg"),
                { type: format }
              )
            );
            setCompressionProgress(100);
          },
          format,
          0.8 // 圧縮率
        );
      };
      img.onerror = () => reject("画像の読み込みに失敗しました");
    });
  };

  return { compressImage, compressionProgress };
};
