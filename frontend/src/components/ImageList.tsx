import { Suspense } from "react";
import { useImages } from "../hooks/useImages";
import ImageCard from "./ImageCard";

const ImageList = () => {
  const { images } = useImages();

  return (
    <Suspense fallback={<p className="text-center">画像を読み込み中...</p>}>
      <div className="p-2">
        <h2 className="text-2xl font-semibold mb-4">画像一覧</h2>
        <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 gap-4">
          {images.length > 0 ? (
            images.map((image) => <ImageCard key={image.id} image={image} />)
          ) : (
            <p className="col-span-full text-center text-gray-500">
              画像を追加してください。
            </p>
          )}
        </div>
      </div>
    </Suspense>
  );
};

export default ImageList;
