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
    <div className="p-2">
      <h2 className="text-2xl font-semibold mb-4">画像一覧</h2>
      <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 gap-4">
        {image ? image.map((image) => (
          <div key={image.id} className="bg-white shadow-lg rounded-lg overflow-hidden transform hover:scale-105 transition-transform duration-300">
            <img src={image.url} alt={`Image ${image.id}`} className="w-full h-full object-cover" loading="lazy" />
          </div>
        )) : (
          <p className="col-span-full text-center text-gray-500">画像を追加してください。</p>
        )}
      </div>
    </div>
  );
};

export default ImageList;