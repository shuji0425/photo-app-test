import React, { useState } from "react";
import axios from "axios";

const UploadForm = () => {
  const [file, setFile] = useState<File | null>(null);
  const [message, setMessage] = useState("");

  const handleFileChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    if (event.target.files) {
      setFile(event.target.files[0]);
    }
  };

  const handleUpload = async () => {
    if (!file) {
      setMessage("ファイルを選択してください");
      return;
    }

    const formData = new FormData();
    formData.append("photo", file);

    try {
      const response = await axios.post("http://localhost:8080/upload", formData);

      if (response.status === 200) {
        setMessage("アップロード成功");
      } else {
        setMessage("アップロードに失敗しました");
      }
    } catch (error) {
      setMessage("エラーが発生しました");
      console.error(error);
    }
  };

  return (
    <div className="p-4 bg-gray-100 rounded-lg shadow-md max-w-md mx-auto">
      <h2 className="text-xl font-bold mb-4">画像アップロード</h2>
      <input type="file" onChange={handleFileChange} className="mb-2" />
      <button onClick={handleUpload} className="bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-700">
        アップロード
      </button>
      {message && <p className="mt-2 text-red-500">{message}</p>}
    </div>
  )
}

export default UploadForm