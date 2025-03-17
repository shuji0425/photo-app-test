import React, { useRef, useState } from "react";
import axios from "axios";

const UploadForm = () => {
  const [images, setImages] = useState<File[]>([]);
  const [message, setMessage] = useState("");
  const [isUploading, setIsUploading] = useState(false);
  const inputRef = useRef<HTMLInputElement>(null);

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

    // アップロード開始
    setIsUploading(true);

    const formData = new FormData();
    images.forEach((image) => {
      formData.append("images", image);
    });

    try {
      const response = await axios.post(
        "http://localhost:8080/upload",
        formData
      );

      if (response.status === 200) {
        setMessage("アップロード成功");
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
      <input
        type="file"
        ref={inputRef}
        onChange={handleFileChange}
        accept="image/*"
        multiple
        className="mb-2"
        disabled={isUploading}
      />
      <button
        onClick={handleUpload}
        disabled={isUploading}
        className={`px-4 py-2 rounded ${isUploading ? "bg-gray-400 cursor-not-allowed" : "bg-blue-500 hover:bg-blue-700 text-white"}`}
      >
        {isUploading ? "アップロード中..." : "アップロード"}
      </button>
      {message && <p className="mt-2 text-red-500">{message}</p>}
    </div>
  );
};

export default UploadForm;
