import { useState } from "react";
import { Image } from "../types/image";

// サイレンダリング削減
const ImageCard = ({ image }: { image: Image }) => {
  const [loaded, setLoaded] = useState(false);

  return (
    <div className="bg-white shadow-lg rounded-lg overflow-hidden transform hover:scale-105 transition-transform duration-300">
      {/* プレースホルダー (Shimmer Effect) */}
      {!loaded && <div className="w-full h-48 bg-gray-300 animate-pulse" />}

      <img
        srcSet={`${image.url}?w=400&h=300&fit=crop 400w, ${image.url}?w=800&h=600&fit=crop 800w`}
        sizes="(max-width: 640px) 400px, (max-width: 1024px) 800px, 1200px"
        src={image.url}
        alt={`Image ${image.id}`}
        className={`w-full h-full object-cover transition-opacity duration-500 ${
          loaded ? "opacity-100" : "opacity-0"
        }`}
        loading="lazy"
        onLoad={() => setLoaded(true)}
      />
    </div>
  );
};

export default ImageCard;
