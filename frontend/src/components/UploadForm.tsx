import React, { useRef, useState } from "react";
import axios from "axios";

const UploadForm = () => {
  const [images, setImages] = useState<File[]>([]);
  const [message, setMessage] = useState("");
  const [isUploading, setIsUploading] = useState(false);
  const inputRef = useRef<HTMLInputElement>(null);
  const [uploadProgress, setUploadProgress] = useState(0);

  const handleFileChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    if (event.target.files) {
      const files = Array.from(event.target.files);
      setImages(files);
    }
  };

  // フォーム送信時のハンドラー
  const handleUpload = async (event: React.FormEvent) => {
    event.preventDefault();

    if (images.length === 0) {
      setMessage("画像を選択してください");
      return;
    }

    setMessage("")        // メッセージを初期化
    setIsUploading(true); // アップロード開始
    setUploadProgress(0); // 進捗リセット

    const formData = new FormData();
    images.forEach((image) => {
      formData.append("images", image);
    });

    try {
      const response = await axios.post(
        "http://localhost:8080/upload",
        formData,
        {
          onUploadProgress: (progressEvent) => {
            const percent = Math.round(
              (progressEvent.loaded * 100) / (progressEvent.total || 1)
            );
            setUploadProgress(Math.min(percent, 99));
          },
        }
      );

      if (response.status === 200) {
        setMessage("アップロード成功");
        setUploadProgress(100);
        setImages([]);
        // inputタグをからにする
        if (inputRef.current) {
          inputRef.current.value = "";
        }
      } else {
        setMessage("アップロードに失敗しました");
      }
    } catch (error) {
      setMessage("エラーが発生しました");
      console.error(error);
    } finally {
      setIsUploading(false);
    }
  };

  return (
    <div className="p-4 bg-gray-100 rounded-lg shadow-md max-w-md mx-auto">
      <h2 className="text-xl font-bold mb-4">画像アップロード</h2>
      <label className="block mb-2 text-sm font-semibold text-gray-700">
        画像を選択
      </label>
      <div className="flex items-center space-x-2">
        <input
          type="file"
          ref={inputRef}
          onChange={handleFileChange}
          accept="image/*"
          multiple
          className="hidden"
          disabled={isUploading}
          id="file-upload"
        />
        <label
          htmlFor="file-upload"
          className={`px-4 py-2 rounded bg-blue-500 text-white cursor-pointer hover:bg-blue-700 ${
            isUploading ? "bg-gray-400 cursor-not-allowed" : ""
          }`}
        >
          ファイルを選択
        </label>

        {/* 選択されたファイル名の表示 */}
        {images.length > 0 && (
          <span className="text-sm text-gray-600 ml-2">
            {images.length}枚の画像が選択されました
          </span>
        )}
      </div>

      <button
        onClick={handleUpload}
        disabled={isUploading}
        className={`mt-4 px-4 py-2 rounded ${
          isUploading
            ? "bg-gray-400 cursor-not-allowed"
            : "bg-blue-500 hover:bg-blue-700 text-white"
        }`}
      >
        {isUploading ? "アップロード中..." : "アップロード"}
      </button>

      {/* プログレスバー */}
      {isUploading && (
        <div className="mt-3">
          <div className="w-full bg-gray-300 h-4 rounded-lg relative overflow-hidden">
            <div
              className="bg-blue-500 h-full"
              style={{ width: `${uploadProgress}%` }}
            >
              <span className="absolute inset-0 flex justify-center items-center text-sm font-bold text-white z-10">
                {uploadProgress}%
              </span>
            </div>
          </div>
        </div>
      )}

      {message && (
        <p
          className={`mt-2 font-bold ${
            uploadProgress === 100 ? "text-green-500" : "text-red-500"
          }`}
        >
          {message}
        </p>
      )}
    </div>
  );
};

export default UploadForm;
