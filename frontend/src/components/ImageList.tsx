import { useEffect, useState } from "react";

interface Image {
  id: string;
  url: string;
}

const ImageList = () => {
  const [image, setImage] = useState<Image[]>([]);

  const fetchImages = async () => {
    const response = await fetch("http://localhost:8080/get/images");
    const data = await response.json();
    setImage(data);
  };

  useEffect(() => {
    fetchImages();
  }, []);

  return (
    <div>
      <h2>画像一覧</h2>
      <div className="grid grid-cols-5 gap-4">
        {image ? image.map((image) => (
          <div key={image.id} className="border p-2">
            <img src={image.url} alt={`Image ${image.id}`} className="w-full h-auto" loading="lazy" />
          </div>
        )) : (
          <p>画像を追加してください。</p>
        )}
      </div>
    </div>
  );
};

export default ImageList;