import { useImages } from "../hooks/useImages";
import ImageCard from "./ImageCard";

const ImageList = () => {
  const { images, isLoading } = useImages();

  return (
    <div className="p-2">
      <h2 className="text-2xl font-semibold mb-4">画像一覧</h2>
      <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 gap-4">
        {isLoading
          ? Array.from({ length: 10 }).map((_, index) => (
              <div
                key={index}
                className="bg-gray-300 animate-pulse w-full"
                style={{ height: "300px" }} // CLS対策
              />
            ))
          : images.map((image, index) => (
              <ImageCard key={image.id} image={image} priority={index < 3} />
            ))}
      </div>
    </div>
  );
};

export default ImageList;
